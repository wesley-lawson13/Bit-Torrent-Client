package connect

import (
    "io"
    "errors"
)

type Handshake struct {
    InfoHash [20]byte
    PeerId [20]byte
}

const pStr = "BitTorrent protocol"
const handshakeLen = 68

func newHandshake(infoHash [20]byte, peerId [20]byte) *Handshake {

    h := Handshake{
        InfoHash: infoHash,
        PeerId: peerId,
    }

    return &h
}

func serialize(h *Handshake) []byte {

    buf := make([]byte, handshakeLen)
    buf[0] = byte(len(pStr))
    
    index := 1
    index += copy(buf[index:], pStr)
    index += copy(buf[index:], make([]byte, 8))
    index += copy(buf[index:], h.InfoHash[:])
    index += copy(buf[index:], h.PeerId[:])

    return buf
}

func deserialize(r io.Reader) (*Handshake, error) {

    var h Handshake

    buf := make([]byte, handshakeLen)
    _, err := io.ReadFull(r, buf)
    if err != nil {
        return nil, errors.New("Error deserializing handshake from reader.")
    }

    if string(buf[1:20]) != pStr {
        return nil, errors.New("Mismatched protocol types: deserialized handshake was not a BitTorrent response.")
    }

    _ = copy(h.InfoHash[:], buf[handshakeLen-40:handshakeLen-20])
    _ = copy(h.PeerId[:], buf[handshakeLen-20:handshakeLen])

    return &h, nil
}
