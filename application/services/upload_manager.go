package services

import (
	"context"
	"fmt"
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

func (vu *VideoUpload) UploadObject(
	ctx context.Context,
	objectPath string,
	client *storage.Client,
) error {
	path := strings.Split(objectPath, fmt.Sprintf("%s/", os.Getenv("LOCALSTORAGE_PATH")))

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	_, err = io.Copy(wc, f)
	if err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, uploadDone chan string) error {
	in := make(chan int, runtime.NumCPU())
	returnCh := make(chan string)

	err := vu.loadPaths()
	if err != nil {
		return err
	}

	ctx, uploadClient, err := getClientUpload()
	if err != nil {
		return err
	}

	for range concurrency {
		go vu.uploadWorker(ctx, in, returnCh, uploadClient)
	}

	go func() {
		for x := range vu.Paths {
			in <- x
		}
		close(in)
	}()

	for r := range returnCh {
		if r != "" {
			uploadDone <- r
			break
		}
	}

	return nil
}

func (vu *VideoUpload) uploadWorker(ctx context.Context, in chan int, returnCh chan string, uploadClient *storage.Client) {
	for x := range in {
		err := vu.UploadObject(ctx, vu.Paths[x], uploadClient)
		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("Error during the upload: %v. %v", vu.Paths[x], err)
			returnCh <- err.Error()
		}

		returnCh <- ""
	}

	returnCh <- "Upload Completed"
}

func (vu *VideoUpload) loadPaths() error {
	err := filepath.Walk(
		vu.VideoPath,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				vu.Paths = append(vu.Paths, path)
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func getClientUpload() (context.Context, *storage.Client, error) {
	ctx := context.Background()

	c, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return ctx, c, nil
}
