package repositories

import (
	"fmt"

	"github.com/Codeflix-FullCycle/encoder/domain"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type JobsRepository interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobsRepositoryDb struct {
	db *gorm.DB
}

func NewJobsRepositoryDb(db *gorm.DB) *JobsRepositoryDb {
	return &JobsRepositoryDb{
		db: db,
	}
}

func (repo JobsRepositoryDb) Find(id string) (*domain.Job, error) {
	var job domain.Job
	repo.db.Preload("Video").First(&job, "id = ?", id)

	if job.ID == "" {
		return nil, fmt.Errorf("job does not exist")
	}
	return &job, nil
}

func (repo JobsRepositoryDb) Insert(job *domain.Job) (*domain.Job, error) {
	if job.ID == "" {
		job.ID = uuid.NewV4().String()
	}

	err := repo.db.Create(job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (repo JobsRepositoryDb) Update(job *domain.Job) (*domain.Job, error) {
	err := repo.db.Save(&job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}
