package main

import (
	"log"
	"os"

	"github.com/amha-mersha/GoTorrent/torrent_file"
)

func main() {
	inPath := os.Args[1]
	outPath := os.Args[2]

	tf, err := torrent_file.DecodeTorrentFile(inPath)
	if err != nil {
		log.Fatal(err)
	}

	err = tf.DownloadToFiles(outPath)
	if err != nil {
		log.Fatal(err)
	}
}
