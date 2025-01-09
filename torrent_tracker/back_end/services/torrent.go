package services

import (
	"time"

	"github.com/amha-mersha/GoTorrent/db"
	"github.com/amha-mersha/GoTorrent/domains"
)

type Service struct {
	torrentDB *db.PostgreSQLDB
	peerDB    *db.RedisInst
}

func NewService(torrentDB *db.PostgreSQLDB, peerDB *db.RedisInst) *Service {
	return &Service{
		torrentDB: torrentDB,
		peerDB:    peerDB,
	}
}

func (s *Service) Announce(info_hash, peer_id, ip string, port, uploaded, downloaded, left int) ([]domains.Peer, error) {
	_, err := s.torrentDB.GetTorrent(info_hash)
	if err != nil {
		return nil, err
	}

	peer := domains.Peer{
		InfoHash: info_hash,
		PeerID:   peer_id,
		IP:       ip,
		Port:     port,
		LastSeen: time.Now(),
	}

	peers, err := s.peerDB.GetPeers(info_hash)
	if err != nil {
		return nil, err
	}

	err = s.peerDB.AddPeer(info_hash, peer)
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func (s *Service) UploadTorrent() {

}
