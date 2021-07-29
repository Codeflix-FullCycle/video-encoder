package services_test

import (
	"log"
	"testing"
	"time"

	"github.com/Codeflix-FullCycle/encoder/application/repositories"
	"github.com/Codeflix-FullCycle/encoder/application/services"
	"github.com/Codeflix-FullCycle/encoder/domain"
	"github.com/Codeflix-FullCycle/encoder/framework/database"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalln(err)
	}
}

func prepare() (*domain.Video, repositories.VideosRepository) {
	db := database.NewDbTest()
	defer db.Close()

	videoRepository := repositories.NewVideosRepositoryDb(db)

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.ResourceId = uuid.NewV4().String()
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()
	video.FilePath = "sports - you are the right one (legendado).mp4"

	return video, videoRepository
}

func TestDownload(t *testing.T) {

	video, repo := prepare()

	videoService := services.NewVideoService()

	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("code-flix-tests")

	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	err = videoService.Finish()
	require.Nil(t, err)
}
