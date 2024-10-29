package torrent_file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeTorrentFile(t *testing.T) {
	torrentFile, err := DecodeTorrentFile("./testdata/archlinux-2019.12.01-x86_64.iso.torrent")
	assert.Nil(t, err, "Expected no error when decoding torrent file")

	assert.NotNil(t, torrentFile, "Expected torrent file to be decoded")
	assert.Equal(t, "http://tracker.archlinux.org:6969/announce", torrentFile.Announce, "Expected announce URL to be http://tracker.archlinux.org:6969/announce")

	assert.Equal(t, "archlinux-2019.12.01-x86_64.iso", torrentFile.Info.Name, "Expected correct file name")
	assert.Equal(t, 670040064, torrentFile.Info.Length, "Expected correct file length")
	assert.Equal(t, 524288, torrentFile.Info.PieceLength, "Expected correct piece length")

	assert.Equal(t, 25560, len(torrentFile.Info.Pieces), "Expected correct length for pieces string")
}
