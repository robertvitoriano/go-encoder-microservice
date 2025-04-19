package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	fmt.Println("READER", reader)

	body, err := ioutil.ReadAll(reader)

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

func (v *VideoService) Encode() error {

	commandArgs := []string{}
	commandArgs = append(commandArgs, os.Getenv("LOCAL_STORAGE_PATH")+"/"+v.Video.ID+".frag")
	commandArgs = append(commandArgs, "--use-segment-timeline")
	commandArgs = append(commandArgs, "-o")
	commandArgs = append(commandArgs, os.Getenv("LOCAL_STORAGE_PATH")+"/"+v.Video.ID)
	commandArgs = append(commandArgs, "-f")
	commandArgs = append(commandArgs, "--exec-dir")
	commandArgs = append(commandArgs, "/opt/bento4/bin/mp4dash")

	command := exec.Command("mp4dash", commandArgs...)

	output, err := command.CombinedOutput()

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

func DeleteTestFiles() error {
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

func (v *VideoService) Finish() error {
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

func (v *VideoService) InsertVideo() error {
	_, err := v.VideoRepository.Insert(v.Video)

	if err != nil {
		return err
	}

	return nil
}
