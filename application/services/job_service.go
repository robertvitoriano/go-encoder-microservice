package services

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/robertvitoriano/go-encoder-microservice/application/repositories"
	"github.com/robertvitoriano/go-encoder-microservice/domain"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (j *JobService) changeStatus(status string) error {
	var err error
	fmt.Printf("STATUS UPDATED FROM %v to %v\n", j.Job.Status, status)
	j.Job.Status = status
	j.Job, err = j.JobRepository.Update(j.Job)

	if err != nil {
		return j.failJob(err)
	}
	return nil
}

func (j *JobService) failJob(err error) error {
	j.Job.Status = "FAILED"
	j.Job.Error = err.Error()

	_, err = j.JobRepository.Update(j.Job)

	if err != nil {
		return err
	}

	return err
}
func (j *JobService) Start() error {
	err := j.changeStatus("DOWNLOADING")

	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Download(os.Getenv("INPUT_BUCKET"))

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeStatus("FRAGMENTING")

	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Fragment()

	if err != nil {
		return j.failJob(err)

	}

	err = j.changeStatus("ENCODING")

	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Encode()

	if err != nil {
		return j.failJob(err)

	}

	j.performUpload()

	return nil

}

func (j *JobService) performUpload() error {
	err := j.changeStatus("UPLOADING")

	if err != nil {
		return j.failJob(err)

	}

	videoUpload := NewVideoUpload()

	videoUpload.OutputBucket = os.Getenv("OUTPUT_BUCKET")

	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + j.VideoService.Video.ID

	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))

	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	uploadResult := <-doneUpload

	if uploadResult != "Uploaded completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return nil
}
