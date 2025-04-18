package services_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/robertvitoriano/go-encoder-microservice/application/services"
	"github.com/stretchr/testify/require"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env")
	}
	if err := deleteTestFiles(); err != nil {
		fmt.Println(err)
	}
}

func TestUploadManagerUpload(t *testing.T) {

	if err := deleteTestFiles(); err != nil {
		fmt.Println(err)
	}

	video, videoRepository := prepare()
	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = videoRepository

	err := videoService.Download("go-encodertest")

	require.Nil(t, err)

	err = videoService.Fragment()

	require.Nil(t, err)

	err = videoService.Encode()

	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()

	videoUpload.OutputBucket = "go-upload-teste"
	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + video.ID
	done := make(chan string)
	go videoUpload.ProcessUpload(50, done)

	uploadResult := <-done

	require.Equal(t, "Uploaded completed", uploadResult)
	if err := deleteTestFiles(); err != nil {
		fmt.Println(err)
	}

}
