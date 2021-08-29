package services

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/Codeflix-FullCycle/encoder/application/repositories"
	"github.com/Codeflix-FullCycle/encoder/domain"
	"github.com/Codeflix-FullCycle/encoder/framework/queue"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db              *gorm.DB
	Domain          domain.Job
	MessageChannel  chan amqp.Delivery
	JobReturnResult chan JobWorkerResult
	RabbitMQ        *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func (j *JobManager) NewJobManager(
	db *gorm.DB,
	RabbitMQ *queue.RabbitMQ,
	messageChannel chan amqp.Delivery,
	JobReturnResult chan JobWorkerResult) JobManager {
	return JobManager{
		Db:              db,
		Domain:          domain.Job{},
		MessageChannel:  messageChannel,
		JobReturnResult: JobReturnResult,
		RabbitMQ:        RabbitMQ,
	}
}

func (j *JobManager) Start(ch *amqp.Channel) {
	videoService := *NewVideoService()
	videoService.VideoRepository = repositories.NewVideosRepositoryDb(j.Db)

	jobService := NewJobService(
		repositories.NewJobsRepositoryDb(j.Db),
		videoService)

	concurrence, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))

	if err != nil {
		log.Fatalf("error loading var: CONCURRENCY_WORKERS.")
	}

	for qtdProcess := 0; qtdProcess < concurrence; qtdProcess++ {
		go JobWorker(j.MessageChannel, j.JobReturnResult, jobService, j.Domain, qtdProcess)
	}

	for jobResult := range j.JobReturnResult {
		if jobResult.Error != nil {
			err = j.checkParseErrors(jobResult)
		} else {
			err = j.notifySuccess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(true)
		}
	}
}

func (j *JobManager) notify(jobJson []byte) error {
	err := j.RabbitMQ.Notify(
		string(jobJson),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)

	return err
}

func (j *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID: %v. Error during the job: %v with video: %v. Error: %v",
			jobResult.Message.DeliveryTag, jobResult.Job.ID, jobResult.Job.Video.ID, jobResult.Error.Error())
	} else {
		log.Printf("MessageID: %v. Error parsing message: %v", jobResult.Message.DeliveryTag, jobResult.Error)
	}

	notification := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Error.Error(),
	}

	jsMessage, err := json.Marshal(notification)

	if err != nil {
		log.Print(err)
	}

	err = j.notify(jsMessage)

	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) notifySuccess(jobResult JobWorkerResult, ch *amqp.Channel) error {
	jobJson, err := json.Marshal(jobResult)

	if err != nil {
		return err
	}

	err = j.notify(jobJson)

	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)

	if err != nil {
		return err
	}

	return nil
}
