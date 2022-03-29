package models

import "time"

// Application represent application config. Must create for each application.
type Application struct {
	ID                uint          `json:"id" gorm:"primaryKey"`
	Name              string        `json:"name" gorm:"unique;type:varchar(50)"`
	Interval          time.Duration `json:"interval"`
	URL               string        `json:"url" gorm:"type:varchar(255)"`
	ManifestDir       string        `json:"manifest_dir" gorm:"type:varchar(255)"`
	Username          string        `json:"username" gorm:"type:varchar(100)"`
	Token             string        `json:"token" gorm:"type:varchar(255)"`
	Head              string        `json:"head" gorm:"type:varchar(255)"`
	LastCheck         time.Time     `json:"last_check"`
	LastStatusMessage string        `json:"last_status_message" gorm:"type:varchar(255)"`
}

type ApplicationCreateReq struct {
	Name        string        `json:"name" validate:"required"`
	Interval    time.Duration `json:"interval" validate:"required"`
	URL         string        `json:"url" validate:"required,url"`
	ManifestDir string        `json:"manifest_dir" validate:"required"`
	Username    string        `json:"username" validate:"required"`
	Token       string        `json:"token" validate:"required"`
}

type ApplicationCreateRes struct {
	ID uint `json:"id"`
}

type ApplicationUpdateReq struct {
	ID          uint          `json:"id" validate:"omitempty"`
	Name        string        `json:"name" validate:"required"`
	Interval    time.Duration `json:"interval" validate:"required"`
	URL         string        `json:"url" validate:"required,url"`
	ManifestDir string        `json:"manifest_dir" validate:"required"`
	Username    string        `json:"username" validate:"required"`
	Token       string        `json:"token" validate:"required"`
}

type ApplicationUpdateRes struct {
	ID uint `json:"id"`
}

type ApplicationGetReq struct {
	ID uint `json:"id" validate:"required"`
}

type ApplicationGetRes struct {
	Data *Application `json:"data"`
}

type ApplicationListReq struct {
	Page  int `query:"page" validate:"required"`
	Limit int `query:"limit" validate:"required"`
}

type ApplicationListRes struct {
	Data *[]Application `json:"data"`
}

type ApplicationDeleteReq struct {
	ID uint `json:"id" validate:"required"`
}

type ApplicationDeleteRes struct {
}
