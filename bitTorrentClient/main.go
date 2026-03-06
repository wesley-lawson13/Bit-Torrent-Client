package main

import (
	"bitTorrentClient/connect"
	"crypto/rand"
	"log"
	"os"
)

func main() {

	f, err := os.Open("debian-13.3.0-amd64-netinst.iso.torrent")
	if err != nil {
		log.Fatal(err)
	}

	tf, err := connect.Open(f)
	if err != nil {
		log.Fatal(err)
	}

	peerId := make([]byte, 20)
	if _, err := rand.Read(peerId); err != nil {
		log.Fatalf("Error generating the peer Id: %v\n", err)
	}

	peers, err := tf.GetPeers([20]byte(peerId), 6881)
	if err != nil {
		log.Fatal(err)
	}

    data, err := connect.Download(tf, peers, [20]byte(peerId))
    if err != nil {
        log.Fatalf("Error downloading file: %v\n", err)
    }

    err = os.WriteFile("downloaded_debian.iso", data, 0644)
    if err != nil {
        log.Fatal("Could not write .iso file.")
    }
}
