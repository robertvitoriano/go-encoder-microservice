package repositories

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/robertvitoriano/go-encoder-microservice/domain"
	uuid "github.com/satori/go.uuid"
)

type JobRepository interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepositoryDB struct {
	Connection *gorm.DB
}

func (repository *JobRepositoryDB) Insert(job *domain.Job) (*domain.Job, error) {
	if job.ID == "" {
		job.ID = uuid.NewV4().String()
	}

	err := repository.Connection.Create(job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (repository *JobRepositoryDB) Find(id string) (*domain.Job, error) {
	var job domain.Job

	if id == "" {
		return nil, fmt.Errorf("ID MUST BE GIVEN")
	}

	repository.Connection.Preload("Video").First(&job, "id = ?", id)

	if job.ID == "" {
		return nil, fmt.Errorf("JOB DOES NOT EXIST")
	}

	return &job, nil

}

func (repository *JobRepositoryDB) Update(job *domain.Job) (*domain.Job, error) {
	err := repository.Connection.Save(&job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}
