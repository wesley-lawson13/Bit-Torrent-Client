package connect

import (
	"encoding/binary"
    "log"
	"errors"
	"time"
    "crypto/sha1"
)

type pieceWork struct {
	PieceIndex int
	ExpectedHash     [hashLen]byte
	PieceLen   int
}

type pieceResult struct {
	PieceIndex int
	Data       []byte
}

const blockSize = 1 << 14 // 2^14, block size
const deadline = 30 * time.Second

func calculatePieceLength(tf TorrentFile, index int) int {

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

func checkIntegrity(pw *pieceWork, buf []byte) error {

    h := [hashLen]byte(sha1.Sum(buf))
    if h != pw.ExpectedHash {
        return errors.New("SHA1 verification failed: computed peer hash does not match expected hash.")
    }
    return nil
}

func downloadPiece(cli *Client, pw *pieceWork) ([]byte, error) {

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

func (tf TorrentFile) pushAllPieces(wq chan *pieceWork) {

    for i, pieceHash := range tf.PieceHashes {
        pw := pieceWork{
            PieceIndex: i,
            ExpectedHash: pieceHash,
            PieceLen: calculatePieceLength(tf, i),
        }

        wq <- &pw
    }
}

func Download(tf TorrentFile, peers []Peer, peerId [20]byte) ([]byte, error) {

    numPieces := len(tf.PieceHashes)
    workQueue, results := make(chan *pieceWork, numPieces), make(chan *pieceResult, numPieces)

    tf.pushAllPieces(workQueue)

    for _, peer := range peers {

        go func(peer Peer) {

            cli, err := NewClient(peer, peerId, tf)
            if err != nil {
                return
            }
            defer cli.Conn.Close()

            for pw := range workQueue {

                buf, err := downloadPiece(&cli, pw)
                if err != nil {
                    workQueue <- pw
                    return
                }

                err = checkIntegrity(pw, buf)
                if err != nil {
                    workQueue <- pw
                    return
                }

                pr := pieceResult{
                    PieceIndex: pw.PieceIndex,
                    Data: buf,
                }
                results <- &pr
            }

        }(peer)
    }

    buf, offset := make([]byte, tf.Length), 0

    piecesRecieved := 0
    for pr := range results {

        offset = pr.PieceIndex * tf.PieceLength
        copy(buf[offset:offset+len(pr.Data)], pr.Data)
        log.Printf("Piece at index %v download complete -- offset: %v, data size: %v\n", pr.PieceIndex, offset, len(pr.Data))

        // increment pieces recieved and close when all pieces have been recieved
        piecesRecieved++
        if piecesRecieved == numPieces {
            break
        }
    }

    if len(buf) != tf.Length {
        return buf, errors.New("Download incomplete: not all pieces were recieved.")
    }

    return buf, nil
}
