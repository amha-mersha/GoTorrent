package torrent_file

import (
	"os"
	"testing"

	"github.com/jackpal/bencode-go"
	"github.com/stretchr/testify/assert"
)

func TestDecodeTorrentFile(t *testing.T) {
	torrentFile, err := DecodeTorrentFile("./testdata/archlinux-2019.12.01-x86_64.iso.torrent")

	assert.NoError(t, err, "Expected no error when decoding torrent file")
	assert.NotNil(t, torrentFile, "Expected torrent file to be decoded")

	assert.Equal(t, "http://tracker.archlinux.org:6969/announce", torrentFile.Announce, "Expected announce URL to match")
	assert.Equal(t, "archlinux-2019.12.01-x86_64.iso", torrentFile.Name, "Expected correct file name")
	assert.Equal(t, 670040064, torrentFile.Length, "Expected correct file length")
	assert.Equal(t, 524288, torrentFile.PieceLen, "Expected correct piece length")

	expectedPieceCount := 1278
	assert.Equal(t, expectedPieceCount, len(torrentFile.Pieces), "Expected correct length for pieces array")
}

func TestToTorrentFile(t *testing.T) {
	tf, err := DecodeTorrentFile("./testdata/archlinux-2019.12.01-x86_64.iso.torrent")
	assert.NoError(t, err, "Expected no error when decoding torrent file")

	file, err := os.Open("./testdata/archlinux-2019.12.01-x86_64.iso.torrent")
	assert.NoError(t, err, "Expected no error when opening torrent file.")
	defer file.Close()

	var torrentFile benFile
	err = bencode.Unmarshal(file, &torrentFile)

	assert.Equal(t, tf.Announce, torrentFile.Announce, "Expected announce URL to match")
	assert.Equal(t, tf.Name, torrentFile.Info.Name, "Expected file name to match")
	assert.Equal(t, tf.Length, torrentFile.Info.Length, "Expected file length to match")
	assert.Equal(t, tf.PieceLen, torrentFile.Info.PieceLength, "Expected piece length to match")

	pieceHashes := make([][20]byte, len(torrentFile.Info.Pieces)/20)
	for i := 0; i < len(torrentFile.Info.Pieces); i += 20 {
		var pieceHash [20]byte
		copy(pieceHash[:], torrentFile.Info.Pieces[i:i+20])
		pieceHashes[i/20] = pieceHash
	}
	assert.Equal(t, pieceHashes, tf.Pieces, "Expected piece hashes to match")
}
