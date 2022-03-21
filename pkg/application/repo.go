package application

import (
	"github.com/bdemirpolat/kubecd/pkg/models"
	"gorm.io/gorm"
)

type RepoInterface interface {
	Create(application *models.Application) (uint, error)
	Update(application *models.Application) (uint, error)
	Get(id uint) (*models.Application, error)
	List(page, limit int) (*[]models.Application, error)
}

type RepoStruct struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) RepoInterface {
	return &RepoStruct{db: db}
}

var _ RepoInterface = &RepoStruct{}

func (r *RepoStruct) Create(application *models.Application) (uint, error) {
	return application.ID, r.db.Create(application).Error
}

func (r *RepoStruct) Update(application *models.Application) (uint, error) {
	return application.ID, r.db.Updates(application).Error
}

func (r *RepoStruct) Get(id uint) (*models.Application, error) {
	a := &models.Application{}
	err := r.db.First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *RepoStruct) List(page, limit int) (*[]models.Application, error) {
	var applications []models.Application
	q := r.db.Model(&models.Application{})

	if page == 0 {
		page = 1
	}

	if page > 1 {
		q = q.Offset(page * limit)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}
	return &applications, q.Find(&applications).Error
}
