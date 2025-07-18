package torrent

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"net/http"
	"testing"

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

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)

	trackerResponse := trackerResponse{}
	err = bencode.Unmarshal(bytes.NewReader(body), &trackerResponse)
	assert.Nil(t, err)
	fmt.Printf("%+v", trackerResponse)
}
