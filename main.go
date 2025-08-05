package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/erdemy885/turkTorrent/handshake"
	"github.com/erdemy885/turkTorrent/torrent"
)

func main() {
	tf, err := torrent.Open("debian-12.11.0-amd64-netinst.iso.torrent")
	if err != nil {
		fmt.Println(err)
		return
	}

	var id [20]byte
	rand.Read(id[:])

	peers, err := tf.GetPeers(id, 8888)
	hs := handshake.Handshake{
		Pstr:     "BitTorrent protocol",
		PeerID:   id,
		InfoHash: tf.InfoHash,
	}
	for _, peer := range peers {
		connection, err := peer.Connect()
		if err != nil {
			fmt.Println(err)
		} else {
			connection.Write(hs.Serialize())
			connection.SetReadDeadline(time.Now().Add(3 * time.Second))
			response, err := handshake.Read(connection)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(response)
			}
		}

	}
}
