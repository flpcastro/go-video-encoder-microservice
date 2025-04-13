package repositories_test

import (
	"testing"
	"time"

	"github.com/flpcastro/go-video-encoder-microservice/application/repositories"
	"github.com/flpcastro/go-video-encoder-microservice/domain"
	"github.com/flpcastro/go-video-encoder-microservice/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestJobRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepository(db)
	repo.Insert(video)

	job, err := domain.NewJob(
		"output_path",
		"pending",
		video,
	)
	require.Nil(t, err)

	jobRepo := repositories.NewJobRepository(db)
	jobRepo.Insert(job)

	j, err := jobRepo.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, job.ID, j.ID)
	require.Equal(t, job.VideoID, j.VideoID)
}

func TestJobRepositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepository(db)
	repo.Insert(video)

	job, err := domain.NewJob(
		"output_path",
		"pending",
		video,
	)
	require.Nil(t, err)

	jobRepo := repositories.NewJobRepository(db)
	jobRepo.Insert(job)

	job.Status = "completed"

	jobRepo.Update(job)

	j, err := jobRepo.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, job.Status, j.Status)
}
