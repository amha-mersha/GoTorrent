package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLDB struct {
	DB *gorm.DB
}

func InitPostgreSQL(connection_url string) (*gorm.DB, error) {
	PostgreSQLDB, err := gorm.Open(postgres.Open(connection_url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return PostgreSQLDB, nil
}
