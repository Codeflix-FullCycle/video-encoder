package repositories

import (
	"fmt"

	"github.com/Codeflix-FullCycle/encoder/domain"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type VideosRepository interface {
	Find(id string) (*domain.Video, error)
	Insert(video *domain.Video) (*domain.Video, error)
}

type VideosRepositoryDb struct {
	Db *gorm.DB
}

func NewVideosRepositoryDb(db *gorm.DB) *VideosRepositoryDb {
	return &VideosRepositoryDb{
		Db: db,
	}
}

func (repo VideosRepositoryDb) Find(id string) (*domain.Video, error) {
	var video domain.Video
	repo.Db.Preload("Jobs").First(&video, "id = ?", id)

	if video.ID == "" {
		return nil, fmt.Errorf("video dos not exist")
	}
	return &video, nil
}

func (repo VideosRepositoryDb) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}

	err := repo.Db.Create(video).Error

	if err != nil {
		return nil, err
	}

	return video, nil
}
