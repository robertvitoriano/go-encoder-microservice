package services

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/robertvitoriano/go-encoder-microservice/application/repositories"
	"github.com/robertvitoriano/go-encoder-microservice/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(v.Video.FilePath)

	reader, err := obj.NewReader(ctx)

	if err != nil {
		return err
	}
	defer reader.Close()

	log.Printf("Video %v has been stored", v.Video.ID)

	body, err := io.ReadAll(reader)

	if err != nil {
		return err
	}

	file, err := os.Create(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4")

	if err != nil {
		return err
	}

	_, err = file.Write(body)

	if err != nil {
		return err
	}

	return nil
}

func (v *VideoService) Fragment() error {
	err := os.Mkdir(os.Getenv("LOCAL_STORAGE_PATH")+"/"+v.Video.ID, os.ModePerm)

	if err != nil {
		return err
	}
	source := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func printOutput(output []byte) {
	if len(output) > 0 {
		log.Printf("======> Output: %s\n", string(output))
	}
}
