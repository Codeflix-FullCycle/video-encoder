package services

import (
	"github.com/Codeflix-FullCycle/encoder/domain"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChannel chan JobWorkerResult, jobService JobService, workerId string) {
	// for message := range messageChannel {
	// pegar essa mensagem do body
	// validar o json
	// validar o video
	// salvar o video no bd
	// startar o meu jobService
	// }
}
