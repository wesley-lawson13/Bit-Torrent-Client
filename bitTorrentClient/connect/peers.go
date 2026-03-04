package connect

import (
	"encoding/binary"
	"errors"
	"net/netip"
)

const peerLen = 6

type Peer struct {
	IP   netip.Addr
	Port uint16
}

func parsePeers(rawPeers []byte) ([]Peer, error) {

	peerSliceLen := len(rawPeers)
	if peerSliceLen%peerLen != 0 {
		return nil, errors.New("Malformed peers package: raw binary data is not a multiple of 6.")
	}

	var ret []Peer
	for i := 0; i < peerSliceLen; i += peerLen {
		var p Peer
		p.IP = netip.AddrFrom4([4]byte(rawPeers[i : i+4]))
		p.Port = binary.BigEndian.Uint16(rawPeers[i+4 : i+6])

		ret = append(ret, p)
	}

	return ret, nil
}
