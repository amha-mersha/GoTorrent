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
	ID        string
	Name      string
	Announce  string
	CreatedAt time.Time
}
