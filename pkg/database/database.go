package database

import (
	"github.com/bdemirpolat/kubecd/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

const defaultDBfile = "/etc/kubecd/kubecd.db"

func Init() (*gorm.DB, error) {
	dbFile := defaultDBfile
	if dbFileFromEnv := os.Getenv("KUBECD_DBFILE"); dbFileFromEnv != "" {
		dbFile = dbFileFromEnv
	}
	var err error
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	db = db.Debug()
	db.AutoMigrate(&models.Application{})
	return db, err
}
