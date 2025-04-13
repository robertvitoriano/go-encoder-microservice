package services_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	if err := deleteTestFiles(); err != nil {
		fmt.Println(err)
	}
}
func deleteTestFiles() error {
	files, err := ioutil.ReadDir("/tmp")
	if err != nil {
		return fmt.Errorf("could not read /tmp directory: %v", err)
	}

	for _, file := range files {
		name := file.Name()
		if name == "bento4" || strings.Contains(name, "go-build") {
			continue
		}
		err := os.RemoveAll(filepath.Join("/tmp", name))
		if err != nil {
			return fmt.Errorf("could not remove %s: %v", name, err)
		}
	}

	return nil
}

func runTestWithCleanup(t *testing.T, testFunc func(t *testing.T)) {
	t.Helper()
	t.Cleanup(func() {
		if err := deleteTestFiles(); err != nil {
			fmt.Println(err)
		}
	})

	testFunc(t)
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
	runTestWithCleanup(t, func(t *testing.T) {
		video, videoRepository := prepare()
		videoService := services.NewVideoService()
		videoService.Video = video
		videoService.VideoRepository = videoRepository

		err := videoService.Download("go-encodertest")

		require.Nil(t, err)
	})

}

func TestVideoServiceFragmentation(t *testing.T) {
	runTestWithCleanup(t, func(t *testing.T) {

		video, videoRepository := prepare()
		videoService := services.NewVideoService()
		videoService.Video = video
		videoService.VideoRepository = videoRepository

		err := videoService.Download("go-encodertest")

		require.Nil(t, err)

		err = videoService.Fragment()

		require.Nil(t, err)

	})
}

func TestVideoServiceEncode(t *testing.T) {
	runTestWithCleanup(t, func(t *testing.T) {

		video, videoRepository := prepare()
		videoService := services.NewVideoService()
		videoService.Video = video
		videoService.VideoRepository = videoRepository

		err := videoService.Download("go-encodertest")

		require.Nil(t, err)

		err = videoService.Fragment()

		require.Nil(t, err)
	})
}
