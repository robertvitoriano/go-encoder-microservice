package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (v *VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {

	path := strings.Split(objectPath, os.Getenv("LOCAL_STORAGE_PATH")+"/")

	f, err := os.Open(objectPath)

	if err != nil {
		return err
	}

	defer f.Close()

	writerClient := client.Bucket(v.OutputBucket).Object(path[1]).NewWriter(ctx)

	writerClient.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	if _, err = io.Copy(writerClient, f); err != nil {
		return err
	}

	if err := writerClient.Close(); err != nil {
		return err
	}
	return nil

}

func (vu *VideoUpload) loadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			return err
		}

		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	pathIndexChannel := make(chan int, runtime.NumCPU())

	resultChannel := make(chan string)

	err := vu.loadPaths()

	if err != nil {
		return err
	}

	if uploadClient, ctx, err := getClientUpload(); err == nil {

		for proccess := 0; proccess < concurrency; proccess++ {
			go vu.uploadWorkder(pathIndexChannel, resultChannel, uploadClient, ctx)
		}

		go func() {
			for i := 0; i < len(vu.Paths); i++ {
				pathIndexChannel <- i
			}
		}()

		for result := range resultChannel {
			if result != "" {
				doneUpload <- result
				break
			}
		}

		close(pathIndexChannel)

		return nil
	}

	return err
}

func (vu *VideoUpload) uploadWorkder(pathIndexChannel chan int, resultChannel chan string, uploadClient *storage.Client, ctx context.Context) {
	for i := range pathIndexChannel {
		err := vu.UploadObject(vu.Paths[i], uploadClient, ctx)

		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[i])
			log.Printf("Error during the upload: %v. Error: %v", vu.Paths[i], err)
			resultChannel <- err.Error()
		}

		resultChannel <- ""

	}

	resultChannel <- "Uploaded completed"

}
