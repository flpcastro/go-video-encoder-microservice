package services_test

import (
	"log"
	"testing"
	"time"

	"github.com/flpcastro/go-video-encoder-microservice/application/repositories"
	"github.com/flpcastro/go-video-encoder-microservice/application/services"
	"github.com/flpcastro/go-video-encoder-microservice/domain"
	"github.com/flpcastro/go-video-encoder-microservice/framework/database"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func prepare() (*domain.Video, *repositories.VideoRepositoryDB) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "test.mp4"
	video.CreatedAt = time.Now()

	r := repositories.NewVideoRepository(db)

	return video, r
}

func TestVideoServiceDownload(t *testing.T) {
	v, r := prepare()
	videoService := services.NewVideoService(v, r)
	err := videoService.Download("test-bucket")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	err = videoService.Finish()
	require.Nil(t, err)
}
