package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/robertvitoriano/go-encoder-microservice/application/services"
	"github.com/robertvitoriano/go-encoder-microservice/framework/database"
	"github.com/robertvitoriano/go-encoder-microservice/framework/queue"
	"github.com/streadway/amqp"
)

var db database.Dabase

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	automigrateDB, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("Error parsing automigrate: %v", err)
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalf("Error parsing debug: %v", err)
	}
	db.AutoMigrateDb = automigrateDB
	db.Debug = debug
	db.DsnTest = os.Getenv("DSN_TEST")
	db.Dsn = os.Getenv("DSN")
	db.DbTypeTest = os.Getenv("DB_TYPE_TEST")
	db.DbType = os.Getenv("DB_TYPE")
	db.Env = os.Getenv("ENV")
}

func main() {
	messageChannel := make(chan amqp.Delivery)

	jobResultChannel := make(chan services.JobWorkerResult)

	dbConnection, err := db.Connect()

	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	defer dbConnection.Close()

	rabbitMQ := queue.NewRabbitMQ()

	rabbitMQChannel := rabbitMQ.Connect()

	defer rabbitMQChannel.Close()

	rabbitMQ.Consume(messageChannel)

	jobManager := services.NewJobManager(dbConnection, rabbitMQ, jobResultChannel, messageChannel)

	fmt.Println("Waiting from messages...")

	jobManager.Start(rabbitMQChannel)

}
