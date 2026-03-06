package main

import (
	"bitTorrentClient/connect"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strings"
)

// to test: run 'caffeinate -i go run .'

const torrentFilename = "debian-13.3.0-amd64-netinst.iso.torrent"
const outputFilename = "downloaded_debian.iso"
const expectedHash = "c9f09d24b7e834e6834f2ffa565b33d6f1f540d04bd25c79ad9953bc79a8ac02"

func main() {

	f, err := os.Open(torrentFilename)
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

	err = os.WriteFile(outputFilename, data, 0644)
	if err != nil {
		log.Fatal("Could not write .iso file.")
	}

	log.Printf("%v Download complete %v\n", strings.Repeat("-", 25), strings.Repeat("-", 25))

	// using sha256 to test correctness since sha1 sums are deprecated on the debian site
	h := sha256.Sum256(data)
	actualHash := fmt.Sprintf("%x", h)
	if actualHash != expectedHash {
		log.Fatal("SHA256 sums don't match")
	}

	log.Println("Download successful!")
}
