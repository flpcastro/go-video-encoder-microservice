package repositories

import (
	"github.com/flpcastro/go-video-encoder-microservice/domain"
	"github.com/jinzhu/gorm"
)

type JobRepository interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepositoryDB struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepositoryDB {
	return &JobRepositoryDB{
		db: db,
	}
}

func (jr *JobRepositoryDB) Insert(job *domain.Job) (*domain.Job, error) {
	err := jr.db.Create(job).Error
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (jr *JobRepositoryDB) Find(id string) (*domain.Job, error) {
	var job domain.Job
	err := jr.db.Preload("Video").First(&job, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	if job.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	return &job, nil
}

func (jr *JobRepositoryDB) Update(job *domain.Job) (*domain.Job, error) {
	err := jr.db.Save(&job).Error
	if err != nil {
		return nil, err
	}

	return job, nil
}
