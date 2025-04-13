package domain_test

import (
	"testing"
	"time"

	"github.com/flpcastro/go-video-encoder-microservice/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidadeIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()
	require.Error(t, err)
}

func TestVideoIdIsNotUUID(t *testing.T) {
	video := domain.NewVideo()
	video.ID = "123"
	video.ResourceID = "123"
	video.FilePath = "123"
	err := video.Validate()
	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.ResourceID = "123"
	video.FilePath = "123"
	video.CreatedAt = time.Now()
	err := video.Validate()
	require.Nil(t, err)
}
