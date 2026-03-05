package main

import (
	"bitTorrentClient/connect"
	"crypto/rand"
	"fmt"
	"log"
	"os"
)

func main() {

	f, err := os.Open("connect/debian-13.3.0-amd64-netinst.iso.torrent")
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

	for _, peer := range peers {
        conn, err := connect.NewClient(peer, [20]byte(peerId), tf)
        if err != nil {
            // Use printf here because most clients won't connect
            log.Printf("Error connecting to the peer: %v\n", err)
            continue
        } 
        fmt.Printf("peer %v succeeded!\n", peer.String())
        defer conn.Close()
	}
}
