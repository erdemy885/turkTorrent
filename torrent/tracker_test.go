package torrent

import (
	"crypto/rand"
	"net/http"
	"testing"

	"github.com/erdemy885/turkTorrent/peers"
	"github.com/jackpal/bencode-go"
	"github.com/stretchr/testify/assert"
)

func TestTracker(t *testing.T) {
	tf, err := Open("test-data/debian-12.11.0-amd64-netinst.iso.torrent")
	assert.Nil(t, err)

	var peerID [20]byte
	rand.Read(peerID[:])

	trackerURL, err := tf.buildTrackerURL(peerID, 6881)
	assert.Nil(t, err)

	resp, err := http.Get(trackerURL)
	assert.Nil(t, err)
	defer resp.Body.Close()

	trackerResponse := trackerResponse{}
	err = bencode.Unmarshal(resp.Body, &trackerResponse)
	assert.Nil(t, err)

	prs, err := tf.getPeers(peerID, 6881)
	assert.Nil(t, err)

	prs2, err := peers.Unmarshal([]byte(trackerResponse.Peers))
	assert.Nil(t, err)

	assert.Equal(t, prs, prs2)
}
