package main

import (
	"crypto/rand"
	"fmt"

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
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(peers)
}
