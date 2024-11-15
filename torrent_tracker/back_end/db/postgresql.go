package db

import "gorm.io/gorm"

var PostgreSQLDB *gorm.DB

func InitPostgreSQL(connection_url stirng) error {
	PostgreSQLDB, err := gorm.Open(postgres.Open(connection_url), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
