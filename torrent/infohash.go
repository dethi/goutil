package torrent

import (
	"crypto/sha1"
	"fmt"

	"github.com/zeebo/bencode"
)

// Compute the info hash of a torrent represented as a byte array.
func InfoHash(data []byte) string {
	var torrent struct {
		Info interface{} `bencode:"info"`
	}

	if err := bencode.DecodeBytes(data, &torrent); err != nil {
		return ""
	}

	binfo, err := bencode.EncodeBytes(torrent.Info)
	if err != nil {
		return ""
	}
	hash := sha1.Sum(binfo)
	return fmt.Sprintf("%x", hash)
}
