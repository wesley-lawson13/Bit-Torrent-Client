package connect

import (
	"encoding/binary"
	"errors"
	"net"
	"net/url"
	"strconv"
)

// Peer encodes connection information for a peer
type Peer struct {
	IP   net.IP
	Port uint16
}

// Unmarshal parses peer IP addresses and ports from a buffer
func Unmarshal(peersBin []byte) ([]Peer, error) {

	const peerSize = 6 // 4 for IP, 2 for port
	numPeers := len(peersBin) / peerSize
	if len(peersBin)%peerSize != 0 {
		err := errors.New("Received malformed peers")
		return nil, err
	}

	peers := make([]Peer, numPeers)
	for i := range numPeers {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4 : offset+6])
	}

	return peers, nil
}

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {

	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	params.Set("info_hash", string(t.InfoHash[:]))
	params.Set("event", "started")

	base.RawQuery = params.Encode()
	return base.String(), nil
}
