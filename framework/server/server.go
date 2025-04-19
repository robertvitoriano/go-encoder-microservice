package main

import (
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
		log.Fatalf("Error loggin .env file")
	}

	automigrateDB, _ := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))

	debug, _ := strconv.ParseBool("DEBUG")

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
		log.Fatalf("Error connectioni to DB")
	}

	defer dbConnection.Close()

	rabbitMQ := queue.NewRabbitMQ()

	rabbitMQChannel := rabbitMQ.Connect()

	defer rabbitMQChannel.Close()

	rabbitMQ.Consume(messageChannel)

	jobManager := services.NewJobManager(dbConnection, rabbitMQ, jobResultChannel, messageChannel)

	jobManager.Start(rabbitMQChannel)

}
