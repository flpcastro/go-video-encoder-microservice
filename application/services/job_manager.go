package services

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/flpcastro/go-video-encoder-microservice/application/repositories"
	"github.com/flpcastro/go-video-encoder-microservice/domain"
	"github.com/flpcastro/go-video-encoder-microservice/framework/queue"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	DB               *gorm.DB
	Domain           *domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
	Mutex            *sync.Mutex
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(
	db *gorm.DB,
	rabbitMQ *queue.RabbitMQ,
	jobReturnChannel chan JobWorkerResult,
	messageChannel chan amqp.Delivery,
) *JobManager {
	return &JobManager{
		DB:               db,
		Domain:           &domain.Job{},
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChannel,
		RabbitMQ:         rabbitMQ,
	}
}

func (jm *JobManager) Start(ch *amqp.Channel) {
	videoService := NewVideoService(jm.Domain.Video, repositories.NewVideoRepository(jm.DB))

	jobService := JobService{
		JobRepository: repositories.NewJobRepository(jm.DB),
		VideoService:  videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))
	if err != nil {
		log.Fatalf("Error loading variable CONCURRENCY_WORKERS: %v", err)
	}

	for qtdProcesses := range concurrency {
		go JobWorker(jm.MessageChannel, jm.JobReturnChannel, jobService, jm.Domain, qtdProcesses)
	}

	for jobResult := range jm.JobReturnChannel {
		if jobResult.Err != nil {
			err = jm.checkParseError(jobResult)
		} else {
			err = jm.notifySuccess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}
}

func (jm *JobManager) checkParseError(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID %d Error with Job: %s, with Video: %v. Error: %v", jobResult.Message.DeliveryTag, jobResult.Job.ID, jobResult.Job.Video.ID, jobResult.Err.Error())
	} else {
		log.Printf("MessageID %d Error parsing message: %s", jobResult.Message.DeliveryTag, jobResult.Err.Error())
	}

	errorMsg := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Err.Error(),
	}

	jobJSON, err := json.Marshal(errorMsg)
	if err != nil {
		log.Printf("Error marshalling error message: %v", err)
	}

	err = jm.notify(jobJSON)
	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)
	if err != nil {
		return err
	}

	return nil
}

func (jm *JobManager) notify(jobJSON []byte) error {
	err := jm.RabbitMQ.Notify(
		string(jobJSON),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)
	if err != nil {
		return err
	}

	return nil
}

func (jm *JobManager) notifySuccess(jobResult JobWorkerResult, ch *amqp.Channel) error {
	jm.Mutex.Lock()
	jobJSON, err := json.Marshal(jobResult.Job)
	jm.Mutex.Unlock()
	if err != nil {
		return err
	}

	err = jm.notify(jobJSON)
	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)
	if err != nil {
		return err
	}

	return nil
}
