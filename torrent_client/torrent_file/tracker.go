package torrent_file

import (
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
	q := base.Query()
	q.Set("info_hash", string(tf.InfoHash[:]))
	q.Set("peer_id", string(peerId[:]))
	q.Set("port", strconv.Itoa(int(port)))
	q.Set("left", strconv.Itoa(tf.Length))

	base.RawQuery = q.Encode()
	return base.String(), nil
}

func (tf *TorrentFile) requestPeers(port uint16, peerId [20]byte) ([]peers.Peer, error) {
	trackerUrl, err := tf.buildTrackerUrl(peerId, port)
	conn := http.Client{Timeout: 15 * time.Second}
	resp, err := conn.Get(trackerUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var trackerReponse peerResponse
	err = bencode.Unmarshal(resp.Body, &trackerReponse)
	if err != nil {
		return nil, err
	}
	return peers.Unmarshal([]byte(trackerReponse.Peers))
}
