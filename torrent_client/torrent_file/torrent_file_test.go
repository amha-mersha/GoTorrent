package torrent_file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeTorrentFile(t *testing.T) {
	torrentFile, err := DecodeTorrentFile("./testdata/sample.torrent")
	assert.NoError(t, err, "Expected no error when decoding torrent file")

	assert.Equal(t, "http://tracker.example.com:8080/", torrentFile.Announce, "Expected announce URL to be http://tracker.example.com:8080/")

	assert.Equal(t, "testfile", torrentFile.Info.Name, "Expected correct file name")
	assert.Equal(t, 12345, torrentFile.Info.Length, "Expected correct file length")
	assert.Equal(t, 16384, torrentFile.Info.PieceLength, "Expected correct piece length")

	assert.Equal(t, 20, len(torrentFile.Info.Pieces), "Expected correct length for pieces string")
}
