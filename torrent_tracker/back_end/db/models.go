package db

import (
	"time"
)

type Peer struct {
	ID         string `gorm:"primaryKey"`
	TorrentID  string `gorm:"index"`
	IP         string
	Port       int
	LastActive time.Time
}

type Torrent struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Announce  string
	CreatedAt time.Time
}

func MigrateModels() error {
	if err := SQLiteDB.AutoMigrate(&Peer{}); err != nil {
		return err
	}
	if err := PostgreSQLDB.AutoMigrate(&Torrent{}); err != nil {
		return err
	}
	return nil
}
