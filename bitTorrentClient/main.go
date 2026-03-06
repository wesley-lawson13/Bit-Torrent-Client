package main

import (
	"bitTorrentClient/connect"
	"crypto/rand"
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
        cli, err := connect.NewClient(peer, [20]byte(peerId), tf)
        if err != nil {
            log.Printf("peer %v failed: %v\n", peer, err)
            continue
        }
        pw := connect.PieceWork{
            PieceIndex: 0,
            ExHash: tf.InfoHash,
            PieceLen: connect.CalculatePieceLength(tf, 0),
        }

        ret, err := connect.DownloadPiece(&cli, &pw)
        if err != nil {
            log.Fatalf("download for peer %v failed: %v", peer, err)
        }

        log.Printf("peer %v succeeded!: buffer length = %v\n", peer, len(ret))
        cli.Conn.Close()
        break
    }
}
