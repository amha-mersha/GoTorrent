package db

import (
	"fmt"
	"log"

	"github.com/amha-mersha/GoTorrent/domains"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteDB struct {
	DB *gorm.DB
}

func InitSQLite() (*SQLiteDB, error) {
	gormDB, err := gorm.Open(sqlite.Open("peers.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	query := `CREATE TABLE IF NOT EXISTS peers (
		info_hash TEXT NOT NULL,
		peer_id TEXT NOT NULL,
		ip TEXT NOT NULL,
		port INTEGER NOT NULL,
		last_seen DATETIME NOT NULL,
		PRIMARY KEY (info_hash, peer_id)`
	if err := gormDB.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}
	log.Println("SQLite database initialized and table created.")

	return &SQLiteDB{DB: gormDB}, nil
}

func (sqlite *SQLiteDB) AddOrUpdate(infoHash, peerID, ip string, port int) error {
	query := `
	INSERT INTO peers (info_hash, peer_id, ip, port, last_seen)
	VALUES (?, ?, ?, ?, datetime('now'))
	ON CONFLICT(info_hash, peer_id) DO UPDATE SET
		ip = excluded.ip,
		port = excluded.port,
		last_seen = CURRENT_TIMESTAMP;
	`
	if err := sqlite.DB.Exec(query, infoHash, peerID, ip, port).Error; err != nil {
		return err
	}
	return nil
}

func (sqlite *SQLiteDB) GetPeers(infoHash string) ([]domains.Peer, error) {
	var peers []domains.Peer

	query := `
	SELECT info_hash, peer_id, ip, port, last_seen
	FROM peers
	WHERE info_hash = ?
	AND last_seen > datetime('now', '-30 minutes'); -- Filter out stale peers
	`

	if err := sqlite.DB.Raw(query, infoHash).Scan(&peers).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve peers: %v", err)
	}

	return peers, nil
}
