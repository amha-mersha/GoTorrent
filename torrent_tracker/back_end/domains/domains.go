package domains

import "time"

type Peer struct {
	InfoHash string
	PeerID   string
	IP       string
	Port     int
	LastSeen time.Time
}

type Torrent struct {
	ID        string    `gorm:"primaryKey;size:36;not null"`
	Name      string    `gorm:"size:255;not null"`
	Announce  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
