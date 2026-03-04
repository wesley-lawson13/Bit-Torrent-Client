package connect

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"github.com/jackpal/bencode-go"
	"io"
)

// Using bencode parser for now.
type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

// global hashLen variable
const hashLen = 20

type TorrentFile struct {
	Announce    string
	InfoHash    [hashLen]byte
	PieceHashes [][hashLen]byte
	PieceLength int
	Length      int
	Name        string
}

// Open parses a torrent file
func Open(r io.Reader) (TorrentFile, error) {
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return TorrentFile{}, errors.New("Error unmarashaling the .torrent file to a bencodeTorrent struct.")
	}
	return bto.toTorrentFile()
}

// parses the individual hash pieces and puts them into an array
func (info *bencodeInfo) parsePieces() ([][hashLen]byte, error) {

	raw := []byte(info.Pieces)
	if len(raw)%hashLen != 0 {
		return nil, errors.New("Hash length is malformed.")
	}

	numPieces := len(raw) / hashLen
	pieceHashes := make([][hashLen]byte, numPieces)
	for i := range numPieces {
		copy(pieceHashes[i][:], raw[i*hashLen:(i+1)*hashLen])
	}

	return pieceHashes, nil
}

// Computes the info hash variable for simplified (one-dimensional) Torrent struct
func (info *bencodeInfo) computeHash() ([hashLen]byte, error) {

	var h [hashLen]byte
	var buf bytes.Buffer

	err := bencode.Marshal(&buf, *info)
	if err != nil {
		return h, errors.New("Unable to Marshal info data.")
	}

	h = sha1.Sum(buf.Bytes())
	return h, nil
}

// converting bencodeTorrent to a more basic TorrentFile
func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	var tf TorrentFile

	// initialize simple conversions
	tf.Announce, tf.Name, tf.Length, tf.PieceLength = bto.Announce, bto.Info.Name, bto.Info.Length, bto.Info.PieceLength

	// parse the individual pieces into an array of Byte arrays
	infoHash, err := bto.Info.computeHash()
	if err != nil {
		return tf, err
	}

	// computes the SHA-1 hash of the entire bencoded info dict
	hashPieces, err := bto.Info.parsePieces()
	if err != nil {
		return tf, err
	}

	// initialize infoHash, hashPieces
	tf.InfoHash, tf.PieceHashes = infoHash, hashPieces

	return tf, nil
}
