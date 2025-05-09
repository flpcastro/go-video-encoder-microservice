package main

import (
	"log"
	"os"
	"strconv"

	"github.com/flpcastro/go-video-encoder-microservice/application/services"
	"github.com/flpcastro/go-video-encoder-microservice/framework/database"
	"github.com/flpcastro/go-video-encoder-microservice/framework/queue"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var db database.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	autoMigrateDB, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("Error loading variable AUTO_MIGRATE_DB: %v", err)
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalf("Error loading variable DEBUG: %v", err)
	}

	db.AutoMigrateDb = autoMigrateDB
	db.Debug = debug
	db.DsnTest = os.Getenv("DSN_TEST")
	db.Dsn = os.Getenv("DSN")
	db.DbTypeTest = os.Getenv("DB_TYPE_TEST")
	db.DbType = os.Getenv("DB_TYPE")
	db.Env = os.Getenv("ENV")
}

func main() {
	messageChannel := make(chan amqp.Delivery)
	jobReturnChannel := make(chan services.JobWorkerResult)

	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer dbConnection.Close()

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	jobManager := services.NewJobManager(
		dbConnection,
		rabbitMQ,
		jobReturnChannel,
		messageChannel,
	)
	jobManager.Start(ch)
}
