package db

import (
	"github.com/amha-mersha/GoTorrent/domains"
)

func (postgres *PostgreSQLDB) MigrateModels() error {
	return postgres.DB.AutoMigrate(&domains.Torrent{})
}
