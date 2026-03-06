package connect

import (
	"encoding/binary"
	"errors"
	"io"
)

type Message struct {
	MessageId MessageId
	Payload   []byte
}

type MessageId uint8

const (
	MsgChoke         MessageId = 0
	MsgUnchoke       MessageId = 1
	MsgInterested    MessageId = 2
	MsgNotInterested MessageId = 3
	MsgHave          MessageId = 4
	MsgBitfield      MessageId = 5
	MsgRequest       MessageId = 6
	MsgPiece         MessageId = 7
	MsgCancel        MessageId = 8
)

func (m *Message) serialize() []byte {

	// first check the special keep alive case
	if m == nil {
		return make([]byte, 4)
	}

	length := uint32(1 + len(m.Payload))

	buf := make([]byte, 5+len(m.Payload))
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.MessageId)
	copy(buf[5:], m.Payload[:])
	return buf
}

func readMessage(r io.Reader) (*Message, error) {

	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, errors.New("Error reading message length.")
	}

	length := binary.BigEndian.Uint32(lengthBuf)

	// Checking special keep alive case
	if length == 0 {
		return nil, nil
	}

	// message length is too large
	if length > 1<<17 {
		return nil, errors.New("Message length too large.")
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, errors.New("Error reading message payload and messageId.")
	}

	var m Message
	m.MessageId = MessageId(buf[0])
	m.Payload = buf[1:]

	return &m, nil
}
