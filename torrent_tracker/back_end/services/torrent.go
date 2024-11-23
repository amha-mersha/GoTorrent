package services

import (
	"github.com/amha-mersha/GoTorrent/db"
	"github.com/amha-mersha/GoTorrent/domains"
)

type Service struct {
	torrentDB *db.PostgreSQLDB
	peerDB    *db.SQLiteDB
}

func NewService(torrentDB *db.PostgreSQLDB, peerDB *db.SQLiteDB) *Service {
	return &Service{
		torrentDB: torrentDB,
		peerDB:    peerDB,
	}
}

func (s *Service) Announce(info_hash, peer_id, ip string, port, uploaded, downloaded, left int) ([]domains.Peer, error) {
	err := s.peerDB.AddOrUpdate(info_hash, peer_id, ip, port)
	if err != nil {
		return nil, err
	}

	peers, err := s.peerDB.GetPeers(info_hash)
	if err != nil {
		return nil, err
	}
	return peers, nil
}
