package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/Codeflix-FullCycle/encoder/application/repositories"
	"github.com/Codeflix-FullCycle/encoder/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideosRepository
}

func NewVideoService() *VideoService {
	return &VideoService{}
}

func (vs *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)

	obj := bkt.Object(vs.Video.FilePath)
	r, err := obj.NewReader(ctx)

	if err != nil {
		return err
	}

	defer r.Close()
	body, err := io.ReadAll(r)

	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("localStoragePath") + "/" + vs.Video.ID + ".mp4")
	if err != nil {
		return err
	}
	_, err = f.Write(body)

	if err != nil {
		return err
	}

	defer f.Close()
	log.Printf("O video %v has been storage", vs.Video.ID)

	return nil
}

func (vs *VideoService) Fragment() error {
	err := os.Mkdir(os.Getenv("localStoragePath")+"/"+vs.Video.ID, os.ModePerm)

	if err != nil {
		return err
	}

	source := os.Getenv("localStoragePath") + "/" + vs.Video.ID + ".mp4"
	target := os.Getenv("localStoragePath") + "/" + vs.Video.ID + ".frag"

	cmd := exec.Command("mp4Fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	parseOutput(output)
	return nil
}

func (vs *VideoService) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+vs.Video.ID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+vs.Video.ID)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	parseOutput(output)

	return nil
}

func (vs *VideoService) Finish() error {

	err := os.Remove(os.Getenv("localStoragePath") + "/" + vs.Video.ID + ".mp4")
	if err != nil {
		log.Println("error removing mp4 ", vs.Video.ID, ".mp4")
		return err
	}

	err = os.Remove(os.Getenv("localStoragePath") + "/" + vs.Video.ID + ".frag")
	if err != nil {
		log.Println("error removing frag ", vs.Video.ID, ".frag")
		return err
	}

	err = os.RemoveAll(os.Getenv("localStoragePath") + "/" + vs.Video.ID)
	if err != nil {
		log.Println("error removing mp4 ", vs.Video.ID, ".mp4")
		return err
	}

	log.Println("files have been removed: ", vs.Video.ID)

	return nil

}

func parseOutput(output []byte) {
	if len(output) > 0 {
		fmt.Printf("====> Output: %s\n", string(output))
	}
}
