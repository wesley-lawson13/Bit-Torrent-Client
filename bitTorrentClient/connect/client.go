package connect

import (
	"errors"
	"net"
	"time"
)

type Client struct {
	Conn     net.Conn
	Bitfield Bitfield
	PeerId   [20]byte
	Choked   bool
}

func (c *Client) read() (*Message, error) {

	m, err := readMessage(c.Conn)
	if err != nil {
		return nil, err
	}

	// keep alive message
	if m == nil {
		return nil, nil
	}

	return m, nil
}

func (c *Client) send(m *Message) error {

	payload := m.serialize()
	_, err := c.Conn.Write(payload)
	if err != nil {
		return errors.New("Error writing message payload to peer.")
	}

	return nil
}

func newClient(p Peer, peerId [20]byte, tf TorrentFile) (Client, error) {

	// create connection
	timeout := 3 * time.Second
	dialer := net.Dialer{
		Timeout: timeout,
	}

	conn, err := dialer.Dial("tcp", p.String()) // conn will be closed later in an external function
	if err != nil {
		return Client{}, errors.New("Error dialing peer in NewClient.")
	}

	// create handshake and serialize
	h := newHandshake(tf.InfoHash, peerId)
	payload := h.serialize()

	// send to connection
	_, err = conn.Write(payload)
	if err != nil {
		return Client{}, errors.New("Error writing serialized handshake to the peer.")
	}

	// read and deserialize the handshake response
	hResp, err := deserialize(conn)
	if err != nil {
		return Client{}, err
	}

	if h.InfoHash != hResp.InfoHash {
		return Client{}, errors.New("Mismatched infoHash values.")
	}

	// read first message from the peer - should be a bitfield message
	cli := Client{
		Conn:     conn,
		Bitfield: nil,
		PeerId:   peerId,
		Choked:   true,
	}

	m, err := cli.read()
	if err != nil {
		return cli, err
	}

	if m != nil && m.MessageId != MsgBitfield {
		return cli, errors.New("Expected bitfield message, got unexpected messageId.")
	}

	// set bitfield
	if m != nil {
		cli.Bitfield = m.Payload
	}

	err = cli.send(&Message{MessageId: MsgInterested})
	if err != nil {
		return cli, err
	}

	return cli, nil
}
