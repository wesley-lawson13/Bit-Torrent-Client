package connect

import (
	"net"
	"time"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func NewHandshake(peer Peer, peerId [20]byte) (*Handshake, error) {

	var h *Handshake
	conn, err := net.DialTimeout("tcp", peer.IP.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	return h, nil
}
