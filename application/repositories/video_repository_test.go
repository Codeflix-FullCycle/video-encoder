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

func TestVideoRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	videoRepository := repositories.NewVideosRepositoryDb(db)

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.ResourceId = uuid.NewV4().String()
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()
	video.FilePath = "path"

	videoInserted, err := videoRepository.Insert(video)

	require.Empty(t, err)
	require.NotEmpty(t, videoInserted.ID)

	v, err := videoRepository.Find(video.ID)

	require.Empty(t, err)
	require.NotEmpty(t, v.ID)
	require.Equal(t, video.ID, v.ID)
}
