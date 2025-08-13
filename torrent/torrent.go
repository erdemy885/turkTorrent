package torrent

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/erdemy885/turkTorrent/p2p"
	"github.com/jackpal/bencode-go"
)

const Port uint16 = 6881

type bencodeInfo struct {
	Pieces      string `bencode:"pieces" json:"pieces"`
	PieceLength int    `bencode:"piece length" json:"piece length"`
	Length      int    `bencode:"length" json:"length"`
	Name        string `bencode:"name" json:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce" json:"announce"`
	Info     bencodeInfo `bencode:"info" json:"info"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (bto *bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}

	pieces, err := bto.Info.splitPieces()
	if err != nil {
		return TorrentFile{}, err
	}

	torFile := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieces,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}
	return torFile, nil
}

func (info *bencodeInfo) hash() ([20]byte, error) {
	var buffer bytes.Buffer
	if err := bencode.Marshal(&buffer, *info); err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buffer.Bytes()), nil
}

func (info *bencodeInfo) splitPieces() ([][20]byte, error) {
	if len(info.Pieces)%20 != 0 {
		return nil, fmt.Errorf("invalid piece array of size %v", len(info.Pieces))
	}
	numPieces := len(info.Pieces) / 20
	result := make([][20]byte, numPieces)

	for i := range numPieces {
		copy(result[i][:], info.Pieces[i*20:(i+1)*20])
	}

	return result, nil
}

func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()

	bto := bencodeTorrent{}
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return TorrentFile{}, err
	}

	return bto.toTorrentFile()
}

func (tf *TorrentFile) DownloadToFile(path string) error {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])
	if err != nil {
		return err
	}

	peers, err := tf.GetPeers(peerID, Port)
	if err != nil {
		return err
	}

	torrent := p2p.Torrent{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    tf.InfoHash,
		PieceHashes: tf.PieceHashes,
		PieceLength: tf.PieceLength,
		Length:      tf.Length,
		Name:        tf.Name,
	}

	buf, err := torrent.Download()
	if err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
