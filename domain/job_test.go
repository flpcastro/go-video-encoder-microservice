package domain_test

import (
	"testing"
	"time"

	"github.com/flpcastro/go-video-encoder-microservice/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.ResourceID = "resource-id"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob(
		"path",
		"converted",
		video,
	)
	require.NotNil(t, job)
	require.Nil(t, err)
}
