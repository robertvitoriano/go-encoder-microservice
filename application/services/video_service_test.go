package services_test

import (
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/robertvitoriano/go-encoder-microservice/application/repositories"
	"github.com/robertvitoriano/go-encoder-microservice/application/services"
	"github.com/robertvitoriano/go-encoder-microservice/domain"
	"github.com/robertvitoriano/go-encoder-microservice/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env")
	}
}
func prepare() (*domain.Video, repositories.VideoRepository) {
	db := database.NewDbTest()

	defer db.Close()

	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.FilePath = "video_teste.mp4"
	video.CreatedAt = time.Now()

	videoRepository := repositories.VideoRepositoryDB{Connection: db}

	return video, &videoRepository

}

func TestVideoServiceDownload(t *testing.T) {
	video, videoRepository := prepare()
	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = videoRepository

	err := videoService.Download("go-encodertest")

	require.Nil(t, err)

}

func TestVideoServiceFragmentation(t *testing.T) {
	video, videoRepository := prepare()
	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = videoRepository

	err := videoService.Download("go-encodertest")

	require.Nil(t, err)

	err = videoService.Fragment()

	require.Nil(t, err)
}
