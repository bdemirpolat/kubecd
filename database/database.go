package database

import (
	"github.com/bdemirpolat/kubecd/models"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const defaultDBfile = "/etc/kubecd/kubecd.db"

func Init() (*gorm.DB, error) {
	dbFile := defaultDBfile
	if dbFileFromEnv := os.Getenv("KUBECD_DBFILE"); dbFileFromEnv != "" {
		dbFile = dbFileFromEnv
	}
	var err error
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//db = db.Debug()
	err = db.AutoMigrate(&models.Application{})
	if err != nil {
		return nil, err
	}
	return db, err
}
