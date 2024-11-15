package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var SQLiteDB *gorm.DB

func InitSQLite() error {
	SQLiteDB, err := gorm.Open(sqlite.Open("peers.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
