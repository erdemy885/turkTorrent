package torrent

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tf TorrentFile

func openJSONTorrent(t *testing.T) bencodeTorrent {
	file, err := os.Open("test-data/archlinux-2019.12.01-x86_64.iso.json")
	assert.Nil(t, err)
	defer file.Close()

	fileStats, err := file.Stat()
	assert.Nil(t, err)

	fileBuffer := make([]byte, fileStats.Size())
	bytesRead, err := file.Read(fileBuffer)
	assert.Nil(t, err)
	assert.Equal(t, fileStats.Size(), int64(bytesRead))

	var jsonTorrent bencodeTorrent
	json.Unmarshal(fileBuffer, &jsonTorrent)
	return jsonTorrent
}

func TestMain(m *testing.M) {
	var err error
	tf, err = Open("test-data/archlinux-2019.12.01-x86_64.iso.torrent")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

func TestTorrent(t *testing.T) {
	jsonTorrent := openJSONTorrent(t)

	assert.Equal(t, tf.Announce, jsonTorrent.Announce)
	assert.Equal(t, tf.Name, jsonTorrent.Info.Name)
	assert.Equal(t, tf.Length, jsonTorrent.Info.Length)
	assert.Equal(t, tf.PieceLength, jsonTorrent.Info.PieceLength)

	splitPieces := strings.Split(jsonTorrent.Info.Pieces, " ")
	jsonTorrent.Info.Pieces = strings.Join(splitPieces, "")
	pieces, err := hex.DecodeString(jsonTorrent.Info.Pieces)
	jsonTorrent.Info.Pieces = string(pieces)
	assert.Nil(t, err)

	infoHash, err := jsonTorrent.Info.hash()
	assert.Equal(t, tf.InfoHash, infoHash)
	assert.Nil(t, err)

	pieceHashes, err := jsonTorrent.Info.splitPieces()
	assert.Nil(t, err)
	for i, hash := range tf.PieceHashes {
		for j, hashByte := range hash {
			assert.Equal(t, fmt.Sprintf("%02X", hashByte), splitPieces[i*20+j])
			assert.Equal(t, hashByte, pieceHashes[i][j])
		}
	}

}
