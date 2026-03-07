# BitTorrent Client in Go

A fully functional BitTorrent client written in Go that implements the BitTorrent peer wire protocol from scratch to download Debian. 

The client successfully downloads files by communicating with real trackers, performing a TCP handshake with live peers, and assembling verified pieces concurrently across multiple connections. 

This project taught me how BitTorrent works in depth, but also helped me further my understanding of important Go programming principles like goroutines, channels, and detailed error checking and validation. 

After completing this project, I feel more confident in my skills as a programmer and my knowledge of networking principles and their modern implementations.

---

## Features

- Tracker HTTP communication to discover peers
- TCP peer handshaking using the BitTorrent peer wire protocol
- Bitfield parsing to determine which pieces each peer holds
- Concurrent piece downloading across multiple peers using goroutines
- SHA-1 integrity verification for every downloaded piece
- SHA-256 final file verification against official checksums
- Automatic piece requeueing on download or integrity failure

---

## Project Structure

```
bitTorrentClient/
├── main.go
├── debian-13.3.0-amd64-netinst.iso.torrent # The .torrent file I used for testing
└── connect/
    ├── torrentFile.go   # Torrent file parsing and struct definitions
    ├── peers.go         # Peer struct and compact peer list parsing
    ├── tracker.go       # Tracker HTTP communication and peer discovery
    ├── handshake.go     # BitTorrent handshake serialization and deserialization
    ├── message.go       # Peer wire protocol message types and serialization
    ├── bitfield.go      # Bitfield parsing and piece availability tracking
    ├── client.go        # TCP client, connection management, and messaging
    └── download.go      # Piece downloading, integrity verification, and orchestration
```

---

## How It Works

**1. Parse the torrent file**

The client reads a `.torrent` file and decodes the bencoded data into a `TorrentFile` one-dimensional struct containing the tracker URL, file length, piece length, and SHA-1 hashes for every piece.

**2. Discover peers**

The client sends an HTTP GET request to the tracker URL with the required query parameters including the info hash, peer ID, and port. The tracker responds with a compact binary list of peers which the client parses into IP and port pairs.

**3. Handshake with peers**

For each peer, the client opens a TCP connection and exchanges the 68-byte BitTorrent handshake message. Peers that respond with a mismatched info hash are rejected.

**4. Exchange messages**

After a successful handshake, the client reads the peer's bitfield to determine which pieces they hold, then sends an Interested message. The peer wire protocol messages handled include Choke, Unchoke, Have, Bitfield, Request, and Piece.

**5. Download pieces concurrently**

A work queue channel is populated with all pieces. One goroutine is spawned per peer, each pulling pieces off the queue, sending block requests in 16KB increments, and assembling the responses into a complete piece buffer.

**6. Verify and assemble**

Every downloaded piece is SHA-1 hashed and compared against the expected hash from the torrent file. Pieces that fail verification are requeued. Successfully verified pieces are written to the correct byte offset in the final output buffer.

---

## Usage

```bash
caffeinate -i go run . 
```

*Note: 'caffeinate -i' will ensure the program runs until completion. It is recommended since downloading can take ~10-25 minutes.* 

---

## Future Plans

**Bencode parsing** currently uses the `github.com/jackpal/bencode-go` third-party library. I'm currently working on building a custom Bencode parser so that the client works with code entirely built my me.

---
