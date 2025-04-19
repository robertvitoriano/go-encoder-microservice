package services

import (
	"encoding/json"
	"os"
	"time"

	"github.com/robertvitoriano/go-encoder-microservice/domain"
	"github.com/robertvitoriano/go-encoder-microservice/framework/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, resultChannel chan JobWorkerResult, jobService JobService, job domain.Job, workerId int) {

	for message := range messageChannel {
		err := utils.IsJson(string(message.Body))

		if err != nil {
			resultChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = json.Unmarshal([]byte(message.Body), &jobService.VideoService.Video)

		if err != nil {
			resultChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.VideoService.Video.ID = uuid.NewV4().String()

		err = jobService.VideoService.Video.Validate()
		if err != nil {
			resultChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.InsertVideo()

		if err != nil {
			resultChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		job = domain.Job{
			Video:            jobService.VideoService.Video,
			OutputBucketPath: os.Getenv("OUTPUT_BUCKET"),
			ID:               uuid.NewV4().String(),
			Status:           "STARTING",
			CreatedAt:        time.Now(),
		}

		_, err = jobService.JobRepository.Insert(&job)

		if err != nil {
			resultChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.Job = &job

		err = jobService.Start()

		if err != nil {
			resultChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		resultChannel <- returnJobResult(job, message, nil)

	}

}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult {
	result := JobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   err,
	}

	return result
}
