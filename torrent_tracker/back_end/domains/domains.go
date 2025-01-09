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
	ID           string      `gorm:"primaryKey;size:36;not null"` // Unique identifier (UUID)
	InfoHash     string      `gorm:"size:40;unique;not null"`     // Unique hash to avoid duplicates
	Title        string      `gorm:"size:255;not null"`           // Title of the material
	MaterialType string      `gorm:"size:50;not null"`            // Type of material (e.g., PDF, Video, etc.)
	Subject      string      `gorm:"size:100;not null"`           // Subject or category of the material
	Description  string      `gorm:"type:text"`                   // Description of the material
	Size         int64       `gorm:"not null"`                    // Total size of the torrent (bytes)
	CreatedAt    time.Time   `gorm:"autoCreateTime"`              // Creation timestamp
	PublishYear  int         `gorm:"size:4"`                      // Year of publication of the material
	TorrentFile  TorrentFile `gorm:"foreignKey:TorrentID"`        // Associated torrent file
}

type TorrentFile struct {
	ID        string    `gorm:"primaryKey;size:36;not null"` // Unique identifier for the file (UUID)
	TorrentID string    `gorm:"size:36;not null"`            // Foreign key to the Torrent table
	FileName  string    `gorm:"size:255;not null"`           // Name of the file
	FileType  string    `gorm:"size:50;not null"`            // File type (e.g., .mp4, .pdf, etc.)
	FileSize  int64     `gorm:"not null"`                    // File size in bytes
	CreatedAt time.Time `gorm:"autoCreateTime"`              // Timestamp
}
