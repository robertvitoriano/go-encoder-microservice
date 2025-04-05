package repositories_test

import (
	"testing"
	"time"

	"github.com/robertvitoriano/go-encoder-microservice/application/repositories"
	"github.com/robertvitoriano/go-encoder-microservice/domain"
	"github.com/robertvitoriano/go-encoder-microservice/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestJobRepositoryInsert(t *testing.T) {

	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.CreatedAt = time.Now()
	video.FilePath = "path"

	repo := repositories.VideoRepositoryDB{Connection: db}
	_, err := repo.Insert(video)

	require.Nil(t, err)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDB{Connection: db}

	require.Nil(t, err)

	job, err = repoJob.Insert(job)

	require.Nil(t, err)

	require.NotEmpty(t, job.ID)

	createdJob, err := repoJob.Find(job.ID)

	require.Nil(t, err)
	require.Equal(t, job.ID, createdJob.ID)

	require.Equal(t, job.VideoID, createdJob.VideoID)

}
