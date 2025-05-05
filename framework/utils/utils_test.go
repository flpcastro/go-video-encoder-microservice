package utils_test

import (
	"testing"

	"github.com/flpcastro/go-video-encoder-microservice/framework/utils"
	"github.com/stretchr/testify/require"
)

func TestIsJSON(t *testing.T) {
	json := `
		{
			"id": "525b5c4b-0f2d-4a3e-8b1c-5a7f3e6d9f2f",
			"file_path": "video.mp4",
			"status": "pending"
		}
	`

	err := utils.IsJSON(json)
	require.Nil(t, err)

	json = `fakejson`
	err = utils.IsJSON(json)
	require.Error(t, err)
}
