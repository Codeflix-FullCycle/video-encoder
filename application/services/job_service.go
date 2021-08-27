package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/Codeflix-FullCycle/encoder/application/repositories"
	"github.com/Codeflix-FullCycle/encoder/domain"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobsRepository
	videoService  VideoService
}

func (j *JobService) Start() error {
	err := j.changeJobStaus("DOWNLOAD")

	if err != nil {
		return j.failJob(err)
	}

	err = j.videoService.Download(os.Getenv("inputBucketName"))

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStaus("FRAGMENTING")

	if err != nil {
		return j.failJob(err)
	}

	err = j.videoService.Fragment()

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStaus("ENCODING")

	if err != nil {
		return j.failJob(err)
	}

	err = j.videoService.Encode()

	if err != nil {
		return j.failJob(err)
	}

	err = j.performUpload()

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStaus("FINISH")

	if err != nil {
		return j.failJob(err)
	}

	err = j.videoService.Finish()

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStaus("COMPLETED")

	if err != nil {
		return j.failJob(err)
	}
	return nil
}

func (j *JobService) performUpload() error {

	err := j.changeJobStaus("UPLOADING")

	if err != nil {
		return j.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("inputBucketName")
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + j.videoService.Video.ID
	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	if uploadResult := <-doneUpload; uploadResult != "upload completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return nil
}

func (j *JobService) changeJobStaus(status string) error {
	var err error
	j.Job.Status = status

	j.Job, err = j.JobRepository.Update(j.Job)

	if err != nil {
		return err
	}
	return nil
}

func (j *JobService) failJob(error error) error {
	j.Job.Error = error.Error()

	j.Job.Status = "FAILED"

	if _, err := j.JobRepository.Update(j.Job); err != nil {
		return err
	}

	return nil
}
