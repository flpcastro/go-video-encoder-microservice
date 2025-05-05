package services

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/flpcastro/go-video-encoder-microservice/domain"
	"github.com/flpcastro/go-video-encoder-microservice/framework/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     *domain.Job
	Message *amqp.Delivery
	Err     error
}

var (
	mutex = &sync.Mutex{}
)

func JobWorker(
	messageChannel chan amqp.Delivery,
	returnChannel chan JobWorkerResult,
	jobService JobService,
	job *domain.Job,
	workerID int,
) {
	for message := range messageChannel {
		err := utils.IsJSON(string(message.Body))
		if err != nil {
			returnChannel <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		mutex.Lock()
		err = json.Unmarshal(message.Body, &jobService.VideoService.Video)
		if err != nil {
			returnChannel <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		jobService.VideoService.Video.ID = uuid.NewV4().String()
		mutex.Unlock()

		err = jobService.VideoService.Video.Validate()
		if err != nil {
			returnChannel <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		mutex.Lock()
		err = jobService.VideoService.InsertVideo()
		mutex.Unlock()
		if err != nil {
			returnChannel <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		job.Video = jobService.VideoService.Video
		job.OutputBucketPath = os.Getenv("OUTPUT_BUCKET_NAME")
		job.ID = uuid.NewV4().String()
		job.Status = "STARTING"
		job.CreatedAt = time.Now()

		mutex.Lock()
		_, err = jobService.JobRepository.Insert(job)
		mutex.Unlock()
		if err != nil {
			returnChannel <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		jobService.Job = job
		err = jobService.Start()
		if err != nil {
			returnChannel <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		returnChannel <- returnJobResult(job, &message, nil)
	}
}

func returnJobResult(
	job *domain.Job,
	message *amqp.Delivery,
	err error,
) JobWorkerResult {
	return JobWorkerResult{
		Job:     job,
		Message: message,
		Err:     err,
	}
}
