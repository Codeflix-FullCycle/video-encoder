package repositories_test

import (
	"testing"
	"time"

	"github.com/Codeflix-FullCycle/encoder/application/repositories"
	"github.com/Codeflix-FullCycle/encoder/domain"
	"github.com/Codeflix-FullCycle/encoder/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestJobRepositoryInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.ResourceId = uuid.NewV4().String()
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()
	video.FilePath = "path"

	videoRepository := repositories.NewVideosRepositoryDb(db)
	videoRepository.Insert(video)

	job, err := domain.NewJob("output_path", "pending", video)
	require.Nil(t, err)

	jobRepo := repositories.NewJobsRepositoryDb(db)
	jobRepo.Insert(job)

	j, err := jobRepo.Find(job.ID)

	require.Nil(t, err)
	require.NotEmpty(t, j.ID)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, video.ID, j.VideoID)

}

func TestJobRepositoryUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.ResourceId = uuid.NewV4().String()
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()
	video.FilePath = "path"

	videoRepository := repositories.NewVideosRepositoryDb(db)
	videoRepository.Insert(video)

	job, err := domain.NewJob("output_path", "pending", video)
	require.Nil(t, err)

	jobRepo := repositories.NewJobsRepositoryDb(db)
	jobRepo.Insert(job)

	job.Status = "complete"
	jobRepo.Update(job)
	j, err := jobRepo.Find(job.ID)

	require.Nil(t, err)
	require.NotEmpty(t, j.ID)
	require.Equal(t, j.Status, job.Status)
}
