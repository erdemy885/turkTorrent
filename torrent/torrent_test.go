package torrent

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

var tf TorrentFile

func TestMain(m *testing.M) {
	var err error
	tf, err = Open("test-data\\archlinux-2019.12.01-x86_64.iso.torrent")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

// TO-DO: compare outputs to json
func TestPieceHashes(t *testing.T) {
	fmt.Println(tf)
	for i, hashes := range tf.PieceHashes {
		fmt.Print(strconv.Itoa(i) + ":")
		for _, hash := range hashes {
			fmt.Printf(" %02X", hash)
		}
		fmt.Println()
	}
}

//TO-DO: write testcases for each field of TorrentFile
