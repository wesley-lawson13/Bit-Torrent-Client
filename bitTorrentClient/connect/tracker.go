package connect

import (
	"errors"
	"github.com/jackpal/bencode-go"
	"net/http"
	"net/url"
	"strconv"
)

type bencodeTrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (tf TorrentFile) buildTrackerUrl(peerId [20]byte, port uint16) (string, error) {

	u, err := url.Parse(tf.Announce)
	if err != nil {
		return "", errors.New("Error parsing the torrentFile Announce string.")
	}

	params := url.Values{
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(tf.Length)},
		"compact":    []string{"1"},
	}

	u.RawQuery = params.Encode() + "&info_hash=" + url.QueryEscape(string(tf.InfoHash[:]))
	return u.String(), nil
}

func (tf TorrentFile) GetPeers(peerId [20]byte, port uint16) ([]Peer, error) {

	u, err := tf.buildTrackerUrl(peerId, port)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u)
	if err != nil {
		return nil, errors.New("Error making GET request to the tracker URL.")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unexpected status code when requesting the tracker URL.")
	}

	btr := bencodeTrackerResponse{}
	unmErr := bencode.Unmarshal(resp.Body, &btr)
	if unmErr != nil {
		return nil, unmErr
	}

	peers, err := parsePeers([]byte(btr.Peers))
	if err != nil {
		return nil, err
	}

	return peers, nil
}
