package application

import (
	"errors"

	"github.com/bdemirpolat/kubecd/pkg/models"
)

// ServiceInterface contains service functions
type ServiceInterface interface {
	CreateService(req models.ApplicationCreateReq) (*models.ApplicationCreateRes, error)
	UpdateService(req models.ApplicationUpdateReq) (*models.ApplicationUpdateRes, error)
	GetService(req models.ApplicationGetReq) (*models.ApplicationGetRes, error)
	ListService(req models.ApplicationListReq) (*models.ApplicationListRes, error)
}

// ServiceStruct
type ServiceStruct struct {
	repo RepoInterface
}

// NewService returns new ServiceStruct / ServiceInterface with repo
func NewService(repo RepoInterface) ServiceInterface {
	return &ServiceStruct{repo: repo}
}

// compile time proof
var _ ServiceInterface = &ServiceStruct{}

// CreateService business logic layer for create operation
func (s *ServiceStruct) CreateService(req models.ApplicationCreateReq) (*models.ApplicationCreateRes, error) {
	application := &models.Application{
		Name:        req.Name,
		Interval:    req.Interval,
		URL:         req.URL,
		ManifestDir: req.ManifestDir,
		Username:    req.Username,
		Token:       req.Token,
	}
	id, err := s.repo.Create(application)
	if err != nil {
		return nil, err
	}

	err = Loop(s.repo)
	if err != nil {
		return nil, err
	}
	return &models.ApplicationCreateRes{ID: id}, nil
}

// UpdateService business logic layer for update operation
func (s *ServiceStruct) UpdateService(req models.ApplicationUpdateReq) (*models.ApplicationUpdateRes, error) {
	checkExists, err := s.repo.Get(req.ID)
	if err != nil && checkExists == nil {
		return nil, errors.New("record not found")
	}

	application := &models.Application{
		ID:          req.ID,
		Name:        req.Name,
		Interval:    req.Interval,
		URL:         req.URL,
		ManifestDir: req.ManifestDir,
		Username:    req.Username,
		Token:       req.Token,
	}
	id, err := s.repo.Update(application)
	if err != nil {
		return nil, err
	}

	err = Loop(s.repo)
	if err != nil {
		return nil, err
	}
	return &models.ApplicationUpdateRes{ID: id}, nil
}

// GetService business logic layer for get operation
func (s *ServiceStruct) GetService(req models.ApplicationGetReq) (*models.ApplicationGetRes, error) {
	application, err := s.repo.Get(req.ID)
	if err != nil {
		return nil, err
	}
	return &models.ApplicationGetRes{Data: application}, nil
}

// ListService business logic layer for list operation
func (s *ServiceStruct) ListService(req models.ApplicationListReq) (*models.ApplicationListRes, error) {
	applications, err := s.repo.List(req.Page, req.Limit)
	if err != nil {
		return nil, err
	}
	return &models.ApplicationListRes{Data: applications}, nil
}
