package connect

import (
    "net"
    "time"
    "errors"
)

func NewClient(p Peer, peerId [20]byte, tf TorrentFile) (net.Conn, error) {

    // create connection
    timeout := 3 * time.Second
    dialer := net.Dialer{
        Timeout: timeout,
    }

    conn, err := dialer.Dial("tcp", p.String()) // conn will be closed later in an external function
    if err != nil {
        return nil, errors.New("Error dialing peer in NewClient.")
    }

    // create handshake and serialize
    h := newHandshake(tf.InfoHash, peerId)
    payload := serialize(h)

    // send to connection
    _, err = conn.Write(payload)
    if err != nil {
        return nil, errors.New("Error writing serialized handshake to the peer.")
    }

    // read and deserialize the response
    hResp, err := deserialize(conn)
    if err != nil {
        return nil, err
    }

    if h.InfoHash != hResp.InfoHash {
        return nil, errors.New("Mismatched infoHash values.")
    }

    return conn, nil
}
