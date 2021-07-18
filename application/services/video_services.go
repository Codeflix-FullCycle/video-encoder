package services

import (
	"context"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/Codeflix-FullCycle/encoder/application/repositories"
	"github.com/Codeflix-FullCycle/encoder/domain"
)

type VideoService struct {
	video           *domain.Video
	videoRepository repositories.VideosRepository
}

func (vs *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)

	obj := bkt.Object(vs.video.FilePath)
	r, err := obj.NewReader(ctx)

	if err != nil {
		return err
	}

	defer r.Close()
	body, err := io.ReadAll(r)

	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("localStoragePath") + "/" + vs.video.ID + ".mp4")
	if err != nil {
		return err
	}
	_, err = f.Write(body)

	if err != nil {
		return err
	}

	defer f.Close()
	log.Printf("O video %v has been storage", vs.video.ID)

	return nil
}
