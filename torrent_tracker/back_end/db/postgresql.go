package db

import (
	"github.com/amha-mersha/GoTorrent/domains"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLDB struct {
	DB *gorm.DB
}

func InitPostgreSQL(connectionURL string) (*PostgreSQLDB, error) {
	db, err := gorm.Open(postgres.Open(connectionURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PostgreSQLDB{DB: db}, nil
}

func (postgres *PostgreSQLDB) GetTorrent(id string) (*domains.Torrent, error) {
	torrent := &domains.Torrent{}
	err := postgres.DB.First(torrent, "id = ?", id).Error
	return torrent, err
}
