package db

import (
	"fmt"

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

func (postgres *PostgreSQLDB) AddTorrent(torrent domains.Torrent, file domains.TorrentFile) error {
	err := postgres.DB.Transaction(func(tx *gorm.DB) error {
		var existingTorrent domains.Torrent
		if err := tx.Where("info_hash = ?", torrent.InfoHash).First(&existingTorrent).Error; err == nil {
			return fmt.Errorf("torrent with InfoHash %s already exists", torrent.InfoHash)
		}

		if err := tx.Create(&torrent).Error; err != nil {
			return fmt.Errorf("failed to add torrent: %v", err)
		}

		file.TorrentID = torrent.ID

		if err := tx.Create(&file).Error; err != nil {
			return fmt.Errorf("failed to add torrent file: %v", err)
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
