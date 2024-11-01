package p2p

import (
	"bytes"
	"crypto/sha1"
	"log"
	"runtime"
	"time"

	"github.com/amha-mersha/GoTorrent/message"
	"github.com/amha-mersha/GoTorrent/peers"
)

const MaxBlockSize = 16384
const MaxBacklog = 5

type Torrent struct {
	Peers       []peers.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type pieceWork struct {
	index int
	hash  [20]byte
	len   int
}

type pieceProgress struct {
	index      int
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}

type pieceResult struct {
	index int
	buf   []byte
}

func (tor *Torrent) calcBoundForPiece(index int) (int, int) {
	start := index * tor.PieceLength
	end := min(start+tor.PieceLength, tor.Length)
	return start, end
}

func (tor *Torrent) Download() ([]byte, error) {
	// create a work queue and store the pieces inside it
	log.Println("Starting download for", tor.Name)
	workQueue := make(chan *pieceWork, len(tor.PieceHashes))

	for index, hash := range tor.PieceHashes {
		st, end := tor.calcBoundForPiece(index)
		workQueue <- &pieceWork{index, hash, end - st}
	}
	// create a result queue to accept the values and store
	resultQueue := make(chan *pieceResult)
	// for each piece pass a work to be downloaded
	log.Println("Downloading from", len(tor.Peers), "peers")
	for _, peer := range tor.Peers {
		go tor.downloadFromPeer(&peer, workQueue, resultQueue)
	}

	// collect the downloaded pieces
	buf := make([]byte, tor.Length)
	donePieces := 0
	for donePieces < len(tor.PieceHashes) {
		res := <-resultQueue
		begin, end := tor.calcBoundForPiece(res.index)
		copy(buf[begin:end], res.buf)
		donePieces++

		percent := float64(donePieces) / float64(len(tor.PieceHashes)) * 100
		numWorkers := runtime.NumGoroutine() - 1
		log.Printf("(%0.2f%%) Downloaded piece #%d from %d peers\n", percent, res.index, numWorkers)
	}
	close(workQueue)

	return buf, nil
}

func (state *pieceProgress) readMessage(msg *message.Message, client *peers.Client) error {
	if msg == nil {
		return nil
	}

	switch msg.ID {
	case message.MsgUnchoke:
		client.Choked = false
	case message.MsgChoke:
		client.Choked = true
	case message.MsgHave:
		index, err := msg.ParseHave()
		if err != nil {
			return err
		}
		client.Bitfield.SetPiece(index)
	case message.MsgPiece:
		n, err := msg.ParsePiece(state.index, state.buf)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil
}

func (tor *Torrent) downloadFromPeer(peer *peers.Peer, workQueue chan *pieceWork, resultQueue chan *pieceResult) {
	client, err := peers.NewClient(peer, tor.PeerID, tor.InfoHash)
	if err != nil {
		log.Printf("Could not handshake with %s. Disconnecting\n", peer.IP)
		return
	}
	defer client.Conn.Close()
	log.Printf("Completed handshake with %s\n", peer.IP)

	client.SendUnchock()
	client.SendInterested()

	// start downloading pieces from the work queue
	for work := range workQueue {
		if !client.Bitfield.HasPiece(work.index) {
			workQueue <- work
			continue
		}

		resp, err := tryToDownload(client, work)
		if err != nil {
			log.Println("Exiting", err)
			workQueue <- work
			return
		}

		if !integrityCheck(work, resp.buf) {
			log.Printf("Piece #%d failed integrity check\n", work.index)
			workQueue <- work
			continue
		}
		client.SendHave(work.index)
		resultQueue <- resp
	}
}

func tryToDownload(client *peers.Client, wk *pieceWork) (*pieceResult, error) {
	client.Conn.SetDeadline(time.Now().Add(time.Second * 30))
	defer client.Conn.SetDeadline(time.Time{})

	state := pieceProgress{
		buf:   make([]byte, wk.len),
		index: wk.index,
	}
	for state.downloaded < wk.len {
		if !client.Choked {
			for state.backlog < MaxBacklog && state.requested < wk.len {
				blockSize := min(MaxBlockSize, wk.len-state.requested)
				err := client.SendRequest(wk.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize
			}
		}
		msg, err := message.Read(client.Conn)
		if err != nil {
			return nil, err
		}
		err = state.readMessage(msg, client)
		if err != nil {
			return nil, err
		}
	}

	return &pieceResult{buf: state.buf, index: state.index}, nil
}

func integrityCheck(work *pieceWork, buf []byte) bool {
	hash := sha1.Sum(buf)
	return bytes.Equal(hash[:], work.hash[:])
}
