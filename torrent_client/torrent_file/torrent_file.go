package torrent_file

import (
	"bytes"
	"crypto/sha1"
	"os"

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
	Files       []file `bencode:"files"`  // multiple files
}

type benFile struct {
	Announce     string   `bencode:"announce"`
	AnnounceList []string `bencode:"announce-list"`
	Info         benInfo  `bencode:"info"`
}

type TorrentFile struct {
	InfoHash [20]byte
	Length   int
	Name     string
	Files    []file
	PieceLen int
	Pieces   [][20]byte
	Announce string
}

func DecodeTorrentFile(path string) (*benFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var torrentFile benFile
	err = bencode.Unmarshal(file, &torrentFile)

	if err != nil {
		return nil, err
	}
	return &torrentFile, nil
}

func hashInfo(info *benInfo) ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, info)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buf.Bytes()), nil
}

func (bf *benFile) toTorrentFile() *TorrentFile {
	infoHash, _ := hashInfo(&bf.Info)
	pieceLen := bf.Info.PieceLength
	pieces := bf.Info.Pieces
	pieceHashes := make([][20]byte, len(pieces)/20)
	for i := 0; i < len(pieces); i += 20 {
		var pieceHash [20]byte
		copy(pieceHash[:], pieces[i:i+20])
		pieceHashes[i/20] = pieceHash
	}
	return &TorrentFile{
		InfoHash: infoHash,
		Length:   bf.Info.Length,
		Name:     bf.Info.Name,
		Files:    bf.Info.Files,
		PieceLen: pieceLen,
		Pieces:   pieceHashes,
		Announce: bf.Announce,
	}
}
