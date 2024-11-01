package torrent_file

import (
	"crypto/rand"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/amha-mersha/GoTorrent/peers"
	"github.com/jackpal/bencode-go"
)

const Port uint16 = 6881

type peerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (tf *TorrentFile) buildTrackerUrl(peerId [20]byte, port uint16) (string, error) {
	base, err := url.Parse(tf.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(tf.InfoHash[:])},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(tf.Length)},
	}
	base.RawQuery = params.Encode()
	log.Println("Tracker URL:", base.String())
	return base.String(), nil
}

func (tf *TorrentFile) requestPeers(port uint16, peerId [20]byte) ([]peers.Peer, error) {
	trackerUrl, err := tf.buildTrackerUrl(peerId, port)
	if err != nil {
		return nil, err
	}
	conn := &http.Client{Timeout: 15 * time.Second}
	resp, err := conn.Get(trackerUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	trackerReponse := peerResponse{}
	err = bencode.Unmarshal(resp.Body, &trackerReponse)
	if err != nil {
		return nil, err
	}
	log.Printf("Tracker response: %+v\n", trackerReponse)
	return peers.Unmarshal([]byte(trackerReponse.Peers))
}

func generatePeerID() ([20]byte, error) {
	var peerID [20]byte
	copy(peerID[:8], []byte("-GT0001-"))
	_, err := rand.Read(peerID[8:])
	if err != nil {
		return [20]byte{}, err
	}
	return peerID, nil
}
