package peers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/amha-mersha/GoTorrent/bitfield"
	"github.com/amha-mersha/GoTorrent/message"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

type Client struct {
	Peer     *Peer
	Conn     net.Conn
	Choked   bool
	Bitfield bitfield.BitField
	PeerID   [20]byte
	InfoHash [20]byte
}

func Unmarshal(peersStr []byte) ([]Peer, error) {
	const peerSize = 6
	numPeers := len(peersStr) / peerSize
	if len(peersStr)%peerSize != 0 {
		err := fmt.Errorf("Received malformed peers")
		return nil, err
	}
	peers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peersStr[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16([]byte(peersStr[offset+4 : offset+6]))
	}
	return peers, nil
}

func (p *Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

func NewClient(peer *Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}
	err = performHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}
	log.Println("Handshake successful with", peer.String())

	bf, err := recvBitfield(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	log.Println("Received bitfield from", peer.String())

	return &Client{
		Conn:     conn,
		Choked:   true,
		Peer:     peer,
		InfoHash: infoHash,
		PeerID:   peerID,
		Bitfield: bf,
	}, nil
}

func performHandshake(conn net.Conn, infoHash, peerID [20]byte) error {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})

	_, err := conn.Write(message.BuildHandshake(infoHash, peerID))
	if err != nil {
		return err
	}

	recvInfoHash, _, err := message.ParseHandshake(conn)
	if err != nil {
		return err
	}
	if !bytes.Equal(recvInfoHash[:], infoHash[:]) {
		return fmt.Errorf("info hash does not match")
	}
	return nil
}

func recvBitfield(conn net.Conn) ([]byte, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})
	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}
	if msg.ID != message.MsgBitfield {
		return nil, fmt.Errorf("Expected bitfield but got ID %d", msg.ID)
	}
	return msg.Payload, nil
}

func (client *Client) SendUnchock() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := client.Conn.Write(msg.Serialize())
	return err
}

func (client *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := client.Conn.Write(msg.Serialize())
	return err
}

func (client *Client) SendHave(index int) error {
	msg := message.BuildHaveMessage(index)
	_, err := client.Conn.Write(msg.Serialize())
	return err
}

func (client *Client) SendRequest(index, start, length int) error {
	msg := message.BuildRequest(index, start, length)
	_, err := client.Conn.Write(msg.Serialize())
	return err
}
