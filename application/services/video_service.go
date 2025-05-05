package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/flpcastro/go-video-encoder-microservice/application/repositories"
	"github.com/flpcastro/go-video-encoder-microservice/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService(video *domain.Video, videoRepository repositories.VideoRepository) VideoService {
	return VideoService{
		Video:           video,
		VideoRepository: videoRepository,
	}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()

	c, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bucket := c.Bucket(bucketName)
	obj := bucket.Object(v.Video.FilePath)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	videoPath := fmt.Sprintf("%s/%s.mp4", os.Getenv("LOCALSTORAGE_PATH"), v.Video.ID)
	f, err := os.Create(videoPath)
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Println("Video downloaded successfully")
	return nil
}

func (v *VideoService) Fragment() error {
	fragmentPath := fmt.Sprintf("%s/%s", os.Getenv("LOCALSTORAGE_PATH"), v.Video.ID)
	err := os.Mkdir(fragmentPath, os.ModePerm)
	if err != nil {
		return err
	}

	src := fmt.Sprintf("%s/%s.mp4", os.Getenv("LOCALSTORAGE_PATH"), v.Video.ID)
	target := fmt.Sprintf("%s/%s.frag", os.Getenv("LOCALSTORAGE_PATH"), v.Video.ID)

	cmd := exec.Command("mp4fragment", src, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)
	return nil
}

func (v *VideoService) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, fmt.Sprintf("%s/%s.frag", os.Getenv("LOCALSTORAGE_PATH"), v.Video.ID))
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, fmt.Sprintf("%s/%s", os.Getenv("LOCALSTORAGE_PATH"), v.Video.ID))
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)
	return nil
}

func (v *VideoService) Finish() error {
	err := os.Remove(fmt.Sprintf("%s/%s.mp4", os.Getenv("LOCALSTORAGE_PATH"), v.Video.ID))
	if err != nil {
		log.Printf("Error removing mp4 %s: %s", v.Video.ID, err.Error())
		return err
	}

	err = os.Remove(fmt.Sprintf("%s/%s.frag", os.Getenv("LOCALSTORAGE_PATH"), v.Video.ID))
	if err != nil {
		log.Printf("Error removing frag %s: %s", v.Video.ID, err.Error())
		return err
	}

	err = os.RemoveAll(fmt.Sprintf("%s/%s", os.Getenv("LOCALSTORAGE_PATH"), v.Video.ID))
	if err != nil {
		log.Printf("Error removing directory %s: %s", v.Video.ID, err.Error())
		return err
	}

	log.Printf("Video %s removed successfully", v.Video.ID)
	return nil
}

func (v *VideoService) InsertVideo() error {
	_, err := v.VideoRepository.Insert(v.Video)
	if err != nil {
		return err
	}

	return nil
}

func printOutput(output []byte) {
	if len(output) > 0 {
		log.Printf("====> Output: %s\n", string(output))
	}
}
