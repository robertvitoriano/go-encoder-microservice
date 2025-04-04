package domain_test

import (
	"testing"
	"time"

	"github.com/robertvitoriano/go-encoder-microservice/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()

	err := video.Validate()

	require.Error(t, err)
}

func TestVideoIdIsNotAUuid(t *testing.T) {
	video := domain.NewVideo()

	video.ID = "ABC"
	video.ResourceId = "A"
	video.FilePath = "PATH"
	video.CreatedAt = time.Now()

	err := video.Validate()

	require.Error(t, err)
}
func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.ResourceId = "A"
	video.FilePath = "PATH"
	video.CreatedAt = time.Now()

	err := video.Validate()

	require.Nil(t, err)
}
