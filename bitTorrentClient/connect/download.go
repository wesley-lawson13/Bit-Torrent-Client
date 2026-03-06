package connect

import (
	"encoding/binary"
	"errors"
	"time"
)

type PieceWork struct {
	PieceIndex int
	ExHash     [hashLen]byte
	PieceLen   int
}

type pieceResult struct {
	PieceIndex int
	Data       []byte
}

const blockSize = 1 << 14 // 2^14, block size
const deadline = 30 * time.Second

func CalculatePieceLength(tf TorrentFile, index int) int {

	uniPieces, overflow := tf.Length/tf.PieceLength, (tf.Length%tf.PieceLength != 0)
	totalPieces := uniPieces
	if overflow {
		totalPieces++
	}

	// At the last (smaller) piece
	if index == totalPieces-1 && overflow {
		return tf.Length % tf.PieceLength
	}

	return tf.PieceLength
}

func buildRequestPayload(index, offset, length uint32) []byte {

	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf[0:4], index)
	binary.BigEndian.PutUint32(buf[4:8], offset)
	binary.BigEndian.PutUint32(buf[8:12], length)

	return buf
}

func DownloadPiece(cli *Client, pw *PieceWork) ([]byte, error) {

	buf := make([]byte, pw.PieceLen)

	err := cli.Conn.SetDeadline(time.Now().Add(deadline))
	if err != nil {
		return buf, errors.New("Failed to set a deadline on the connection.")
	}

	bytesRequested, bytesReceived := 0, 0

	for bytesReceived < pw.PieceLen {

		if !cli.Choked && bytesRequested < pw.PieceLen {

			msg := Message{
				MessageId: MsgRequest,
				Payload:   buildRequestPayload(uint32(pw.PieceIndex), uint32(bytesRequested), uint32(min(blockSize, pw.PieceLen-bytesRequested))),
			}

			err = cli.Send(&msg)
			if err != nil {
				return buf, err
			}

			bytesRequested += blockSize
		}

		m, err := cli.Read()
		if err != nil {
			return buf, err
		}

		// keep alive case
		if m == nil {
			continue
		}

		switch m.MessageId {
		case MsgUnchoke:
			cli.Choked = false
		case MsgChoke:
			cli.Choked = true
		case MsgHave:
			index := binary.BigEndian.Uint32(m.Payload)
			cli.Bitfield.setPiece(int(index))
		case MsgPiece:
			offset := binary.BigEndian.Uint32(m.Payload[4:8])
			copy(buf[offset:], m.Payload[8:])
			bytesReceived += len(m.Payload) - 8
		}
	}
	return buf, nil
}
