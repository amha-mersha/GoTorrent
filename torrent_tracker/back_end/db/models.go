package db

import (
	"github.com/amha-mersha/GoTorrent/domains"
)

func (sqlite *SQLiteDB) MigrateModels() error {
	return sqlite.DB.AutoMigrate(&domains.Peer{})
}

func (postgres *PostgreSQLDB) MigrateModels() error {
	return postgres.DB.AutoMigrate(&domains.Torrent{})
}
