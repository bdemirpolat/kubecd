package application

import (
	"encoding/json"
	"github.com/bdemirpolat/kubecd/logger"
	mocks "github.com/bdemirpolat/kubecd/mocks/application"
	"github.com/bdemirpolat/kubecd/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestCreateHandler(t *testing.T) {
	_, err := logger.Init()
	if err != nil {
		assert.Nil(t, err)
	}

	applicationService := &mocks.ServiceInterface{}
	applicationService.On("CreateService", mock.AnythingOfType("models.ApplicationCreateReq")).Return(&models.ApplicationCreateRes{ID: 1}, nil)
	applicationHandler := NewHandler(applicationService)

	app := fiber.New()
	applicationHandler.RegisterHandlers(app)

	b, err := json.Marshal(models.ApplicationCreateReq{
		Name:        time.Now().String(),
		Interval:    2000,
		URL:         "http://xxx.com",
		ManifestDir: "x",
		Username:    "x",
		Token:       "x",
	})
	assert.Nil(t, err)
	body := strings.NewReader(string(b))
	req, err := http.NewRequest("POST", "/applications", body)
	req.Header.Add("Content-Type", "application/json")
	assert.Nil(t, err)
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestUpdateHandler(t *testing.T) {
	_, err := logger.Init()
	if err != nil {
		assert.Nil(t, err)
	}

	applicationService := &mocks.ServiceInterface{}
	applicationService.On("UpdateService", mock.AnythingOfType("models.ApplicationUpdateReq")).Return(&models.ApplicationUpdateRes{ID: 1}, nil)
	applicationHandler := NewHandler(applicationService)

	app := fiber.New()
	applicationHandler.RegisterHandlers(app)

	b, err := json.Marshal(models.ApplicationUpdateReq{
		Name:        time.Now().String(),
		Interval:    2000,
		URL:         "http://xxx.com",
		ManifestDir: "x",
		Username:    "x",
		Token:       "x",
	})
	assert.Nil(t, err)
	body := strings.NewReader(string(b))
	req, err := http.NewRequest("PUT", "/applications/1", body)
	req.Header.Add("Content-Type", "application/json")
	assert.Nil(t, err)
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestGetHandler(t *testing.T) {
	_, err := logger.Init()
	if err != nil {
		assert.Nil(t, err)
	}

	applicationService := &mocks.ServiceInterface{}
	applicationService.On("GetService", mock.AnythingOfType("models.ApplicationGetReq")).Return(&models.ApplicationGetRes{Data: &models.Application{}}, nil)
	applicationHandler := NewHandler(applicationService)

	app := fiber.New()
	applicationHandler.RegisterHandlers(app)

	req, err := http.NewRequest("GET", "/applications/1", nil)

	assert.Nil(t, err)
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestListHandler(t *testing.T) {
	_, err := logger.Init()
	if err != nil {
		assert.Nil(t, err)
	}

	applicationService := &mocks.ServiceInterface{}
	applicationService.On("ListService", mock.AnythingOfType("models.ApplicationListReq")).Return(&models.ApplicationListRes{Data: &[]models.Application{}}, nil)
	applicationHandler := NewHandler(applicationService)

	app := fiber.New()
	applicationHandler.RegisterHandlers(app)

	req, err := http.NewRequest("GET", "/applications", nil)

	assert.Nil(t, err)
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestDeleteHandler(t *testing.T) {
	_, err := logger.Init()
	if err != nil {
		assert.Nil(t, err)
	}

	applicationService := &mocks.ServiceInterface{}
	applicationService.On("DeleteService", mock.AnythingOfType("models.ApplicationDeleteReq")).Return(&models.ApplicationDeleteRes{}, nil)
	applicationHandler := NewHandler(applicationService)

	app := fiber.New()
	applicationHandler.RegisterHandlers(app)

	req, err := http.NewRequest("DELETE", "/applications/1", nil)

	assert.Nil(t, err)
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 204, res.StatusCode)
}
