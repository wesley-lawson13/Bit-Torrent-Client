package connect

import (
    "net/url"
    "errors"
    "strconv"
)

func (tf TorrentFile) buildTrackerUrl(peerId [20]byte, port uint16) (string, error) {

    u, err := url.Parse(tf.Announce)
    if err != nil {
        return "", errors.New("Error parsing the torrentFile Announce string.")
    }

    params := url.Values{
        "peer_id": []string{string(peerId[:])},
        "port": []string{strconv.Itoa(int(port))},
        "uploaded": []string{"0"},
        "downloaded": []string{"0"},
        "left": []string{strconv.Itoa(tf.Length)},
        "compact": []string{"1"},
    }

    u.RawQuery = params.Encode() + "&info_hash=" + url.QueryEscape(string(tf.InfoHash[:]))
    return u.String(), nil
}
