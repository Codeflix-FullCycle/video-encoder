package domain_test

import (
	"testing"
	"time"

	"github.com/Codeflix-FullCycle/encoder/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)
}

func TestVideoIdIsAUuid(t *testing.T) {
	video := domain.NewVideo()

	video.ID = "id_test"
	video.FilePath = "path"
	video.ResourceId = "resource"
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()
	err := video.Validate()
	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.ResourceId = "resource"
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()
	err := video.Validate()

	require.Nil(t, err)

}
