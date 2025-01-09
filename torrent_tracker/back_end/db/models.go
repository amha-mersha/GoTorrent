package db

import (
	"fmt"

	"github.com/amha-mersha/GoTorrent/domains"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domains.Torrent{}, &domains.TorrentFile{}); err != nil {
		return fmt.Errorf("failed to migrate models: %v", err)
	}
	return nil
}
