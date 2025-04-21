package services_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/flpcastro/go-video-encoder-microservice/application/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {
	v, r := prepare()
	videoService := services.NewVideoService(v, r)
	err := videoService.Download("test-bucket")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "test-bucket"

	videoUpload.VideoPath = fmt.Sprintf("%s/%s", os.Getenv("LOCALSTORAGE_PATH"), v.ID)

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(50, doneUpload)

	result := <-doneUpload
	require.Equal(t, result, "Upload Completed")

	err = videoService.Finish()
	require.Nil(t, err)
}
