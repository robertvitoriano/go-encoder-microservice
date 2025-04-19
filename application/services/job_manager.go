package services

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/robertvitoriano/go-encoder-microservice/application/repositories"
	"github.com/robertvitoriano/go-encoder-microservice/domain"
	"github.com/robertvitoriano/go-encoder-microservice/framework/queue"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db               *gorm.DB
	Domain           domain.Job
	MessageChannel   chan amqp.Delivery
	jobResultChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"messaage"`
	Error   string `json:"error"`
}

func NewJobManager(db *gorm.DB, rabbitMQ *queue.RabbitMQ, jobResultChannel chan JobWorkerResult, messageChannel chan amqp.Delivery) *JobManager {

	return &JobManager{
		Db:               db,
		jobResultChannel: jobResultChannel,
		Domain:           domain.Job{},
		MessageChannel:   messageChannel,
		RabbitMQ:         rabbitMQ,
	}
}

func (j *JobManager) Start(rabbitMQChannel *amqp.Channel) {
	videoService := NewVideoService()
	videoService.VideoRepository = &repositories.VideoRepositoryDB{Connection: j.Db}

	jobService := JobService{
		JobRepository: &repositories.JobRepositoryDB{Connection: j.Db},
		VideoService:  videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))

	if err != nil {
		log.Fatalf("Error loading var: CONCURRENCY_WORKERS")
	}
	for workerId := 0; workerId < concurrency; workerId++ {
		go JobWorker(j.MessageChannel, j.jobResultChannel, jobService, j.Domain, workerId)
	}

	for jobResult := range j.jobResultChannel {
		if jobResult.Error != nil {
			err = j.checkParseErrors(jobResult)
		} else {
			err = j.notifySuccess(jobResult, rabbitMQChannel)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}
}

func (j *JobManager) notify(jobJSON []byte) error {
	err := j.RabbitMQ.Notify(
		string(jobJSON),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"))

	if err != nil {
		return err
	}

	return nil
}
func (j *JobManager) notifySuccess(jobResult JobWorkerResult, rabbitMQChannel *amqp.Channel) error {
	jobJSONError, err := json.Marshal(jobResult.Job)

	if err != nil {
		return err
	}
	err = j.notify(jobJSONError)

	if err != nil {
		return err
	}
	err = jobResult.Message.Ack(false)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID: %v. Error during the job: %v with video: %v. Error: %v",
			jobResult.Message.DeliveryTag, jobResult.Job.ID, jobResult.Job.Video.ID, jobResult.Error.Error())
	} else {
		log.Printf("MessageID: %v. Error parsing message: %v", jobResult.Message.DeliveryTag, jobResult.Error)
	}

	errorMessage := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Error.Error(),
	}

	jobJSON, err := json.Marshal(errorMessage)

	if err != nil {
		return err
	}

	err = j.notify(jobJSON)

	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)

	if err != nil {
		return err
	}

	return nil
}
