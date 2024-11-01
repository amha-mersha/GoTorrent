package torrent_file

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"os"

	"github.com/amha-mersha/GoTorrent/p2p"
	"github.com/jackpal/bencode-go"
)

type file struct {
	Length int    `bencode:"length"`
	Path   string `bencode:"path"`
}

type benInfo struct {
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"` // single file
}

type benFile struct {
	Announce string  `bencode:"announce"`
	Info     benInfo `bencode:"info"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func DecodeTorrentFile(path string) (*TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var torrentFile benFile
	err = bencode.Unmarshal(file, &torrentFile)
	log.Printf("Announce: %s\n", torrentFile.Announce)
	log.Printf("Name: %s\n", torrentFile.Info.Name)
	log.Printf("Piece length: %d\n", torrentFile.Info.PieceLength)

	if err != nil {
		return nil, err
	}
	return torrentFile.toTorrentFile(), nil
}

func hashInfo(info *benInfo) ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *info)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buf.Bytes()), nil
}

func (bf *benFile) toTorrentFile() *TorrentFile {
	infoHash, err := hashInfo(&bf.Info)

	if err != nil {
		return nil
	}
	pieceHashes, err := bf.Info.splitPieceHashes()
	if err != nil {
		return nil
	}
	log.Println("Received", len(pieceHashes), "piece hashes")
	log.Println("Info hash:", infoHash)
	return &TorrentFile{
		Announce:    bf.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bf.Info.PieceLength,
		Length:      bf.Info.Length,
		Name:        bf.Info.Name,
	}
}

func (tf *TorrentFile) DownloadToFiles(path string) error {
	peerID, err := generatePeerID()
	if err != nil {
		return err
	}

	peers, err := tf.requestPeers(Port, peerID)
	if err != nil {
		return err
	}
	log.Println("Received", len(peers), "peers")

	torrent := p2p.Torrent{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    tf.InfoHash,
		PieceHashes: tf.PieceHashes,
		PieceLength: tf.PieceLength,
		Length:      tf.Length,
		Name:        tf.Name,
	}
	resultBuff, err := torrent.Download()
	if err != nil {
		return err
	}
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = outFile.Write(resultBuff)
	if err != nil {
		return err
	}
	return nil
}

func (i *benInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}
