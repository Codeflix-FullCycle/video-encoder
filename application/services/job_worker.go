package services

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Codeflix-FullCycle/encoder/domain"
	"github.com/Codeflix-FullCycle/encoder/framework/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChannel chan JobWorkerResult, jobService JobService, job domain.Job, workerId int) {
	for message := range messageChannel {

		if err := utils.IsJson(string(message.Body)); err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		if err := json.Unmarshal(message.Body, &jobService.VideoService.Video); err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.VideoService.Video.ID = uuid.NewV4().String()

		if err := jobService.VideoService.Video.Validate(); err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		if err := jobService.VideoService.Insert(); err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		job.ID = uuid.NewV4().String()
		job.OutputBucketPath = os.Getenv("outputBucketName")
		job.Status = "STARTING"
		job.Video = jobService.VideoService.Video
		job.CreatedAt = time.Now()

		if _, err := jobService.JobRepository.Insert(&job); err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.Job = &job

		if err := jobService.Start(); err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		returnChannel <- returnJobResult(job, message, nil)
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
