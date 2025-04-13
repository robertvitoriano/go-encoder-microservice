package repositories

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/robertvitoriano/go-encoder-microservice/domain"
	uuid "github.com/satori/go.uuid"
)

type VideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDB struct {
	Connection *gorm.DB
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &VideoRepositoryDB{Connection: db}
}

func (repository *VideoRepositoryDB) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}

	err := repository.Connection.Create(video).Error

	if err != nil {
		return nil, err
	}

	return video, nil
}

func (repository *VideoRepositoryDB) Find(id string) (*domain.Video, error) {
	var video domain.Video

	if id == "" {
		return nil, fmt.Errorf("ID MUST BE GIVEN")
	}

	repository.Connection.Preload("Jobs").First(&video, "id = ?", id)

	if video.ID == "" {
		return nil, fmt.Errorf("VIDEO DOES NOT EXIST")
	}

	return &video, nil

}
