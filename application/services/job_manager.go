package services

import (
	"fmt"

	"github.com/robertvitoriano/go-encoder-microservice/domain"
	"github.com/robertvitoriano/go-encoder-microservice/framework/utils"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, resultChannel chan JobWorkerResult, jobService JobService, workerId int) {

	for message := range messageChannel {
		err := utils.IsJson(string(message.Body))

		if err != nil {
			fmt.Println("BODY IS NOT JSON")
		}
	}

}
