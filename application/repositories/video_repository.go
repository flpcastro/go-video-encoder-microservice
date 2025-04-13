package repositories

import (
	"github.com/flpcastro/go-video-encoder-microservice/domain"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type VideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDB struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepositoryDB {
	return &VideoRepositoryDB{
		db: db,
	}
}

func (vr *VideoRepositoryDB) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}

	err := vr.db.Create(video).Error
	if err != nil {
		return nil, err
	}

	return video, nil
}

func (vr *VideoRepositoryDB) Find(id string) (*domain.Video, error) {
	var video domain.Video
	err := vr.db.Preload("Jobs").First(&video, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	if video.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	return &video, nil
}
