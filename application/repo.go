package application

import (
	"github.com/bdemirpolat/kubecd/logger"
	"github.com/bdemirpolat/kubecd/models"
	"gorm.io/gorm"
)

// RepoInterface contains repo functions
type RepoInterface interface {
	Create(application *models.Application) (uint, error)
	Update(application *models.Application) (uint, error)
	Get(id uint) (*models.Application, error)
	List(page, limit int) (*[]models.Application, error)
	Delete(id uint) error
}

// RepoStruct
type RepoStruct struct {
	db *gorm.DB
}

// NewRepo returns new RepoStruct / RepoInterface with db
func NewRepo(db *gorm.DB) RepoInterface {
	return &RepoStruct{db: db}
}

// compile time proof
var _ RepoInterface = &RepoStruct{}

// Create creates new record with gorm.Create function
func (r *RepoStruct) Create(application *models.Application) (uint, error) {
	err := r.db.Create(application).Error
	if err != nil {
		return 0, err
	}
	err = Loop(r)
	if err != nil {
		logger.SugarLogger.Error(err)
		return 0, err
	}
	return application.ID, nil
}

// Update updates a record with gorm.Updates function
func (r *RepoStruct) Update(application *models.Application) (uint, error) {
	err := r.db.Updates(application).Error
	if err != nil {
		return 0, err
	}
	err = Loop(r)
	if err != nil {
		logger.SugarLogger.Error(err)
		return 0, err
	}
	return application.ID, nil
}

// Get finds first record with the given id and gorm.First function
func (r *RepoStruct) Get(id uint) (*models.Application, error) {
	a := &models.Application{}
	err := r.db.First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return a, nil
}

// List finds records with given page and limit parameters
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

// Delete record
func (r *RepoStruct) Delete(id uint) error {
	return r.db.Delete(&models.Application{}, id).Error
}
