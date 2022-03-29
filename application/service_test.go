package application

import (
	"github.com/bdemirpolat/kubecd/logger"
	mocks "github.com/bdemirpolat/kubecd/mocks/application"
	"github.com/bdemirpolat/kubecd/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateService(t *testing.T) {
	_, err := logger.Init()
	if err != nil {
		assert.Nil(t, err)
	}

	applicationRepo := &mocks.RepoInterface{}
	applicationRepo.On("Create", mock.AnythingOfType("*models.Application")).Return(uint(1), nil)
	applicationService := NewService(applicationRepo)

	req := models.ApplicationCreateReq{
		Name:        "test",
		Interval:    1,
		URL:         "http://xxx.com",
		ManifestDir: "x",
		Username:    "x",
		Token:       "x",
	}

	_, err = applicationService.CreateService(req)
	assert.Nil(t, err)
}

func TestUpdateService(t *testing.T) {
	_, err := logger.Init()
	if err != nil {
		assert.Nil(t, err)
	}

	applicationRepo := &mocks.RepoInterface{}
	applicationRepo.On("Update", mock.AnythingOfType("*models.Application")).Return(uint(1), nil)
	applicationRepo.On("Get", mock.AnythingOfType("uint")).Return(&models.Application{}, nil)
	applicationService := NewService(applicationRepo)

	req := models.ApplicationUpdateReq{
		Name:        "test",
		Interval:    1,
		URL:         "http://xxx.com",
		ManifestDir: "x",
		Username:    "x",
		Token:       "x",
	}

	_, err = applicationService.UpdateService(req)
	assert.Nil(t, err)
}

func TestGetService(t *testing.T) {
	_, err := logger.Init()
	if err != nil {
		assert.Nil(t, err)
	}

	applicationRepo := &mocks.RepoInterface{}
	applicationRepo.On("Get", mock.AnythingOfType("uint")).Return(&models.Application{}, nil)
	applicationService := NewService(applicationRepo)

	req := models.ApplicationGetReq{
		ID: 1,
	}

	_, err = applicationService.GetService(req)
	assert.Nil(t, err)
}

func TestListService(t *testing.T) {
	_, err := logger.Init()
	if err != nil {
		assert.Nil(t, err)
	}

	applicationRepo := &mocks.RepoInterface{}
	applicationRepo.On("List", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(&[]models.Application{}, nil)
	applicationService := NewService(applicationRepo)

	req := models.ApplicationListReq{}

	_, err = applicationService.ListService(req)
	assert.Nil(t, err)
}

func TestDeleteService(t *testing.T) {
	_, err := logger.Init()
	if err != nil {
		assert.Nil(t, err)
	}

	applicationRepo := &mocks.RepoInterface{}
	applicationRepo.On("Delete", mock.AnythingOfType("uint")).Return(nil)
	applicationService := NewService(applicationRepo)

	req := models.ApplicationDeleteReq{ID: 1}

	_, err = applicationService.DeleteService(req)
	assert.Nil(t, err)
}
