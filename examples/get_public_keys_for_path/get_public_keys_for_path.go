// Package main demonstrates getting public keys for a derivation path.
package main

import (
	"encoding/hex"
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
)

func main() {
	// Start with an HD key (we will make one for this example)
	hdKey, err := bitcoin.GenerateHDKey(bitcoin.SecureSeedLength)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get keys by path (example showing 5 sets of keys)
	var pubKeys []*bec.PublicKey
	for i := uint32(1); i <= 5; i++ {
		if pubKeys, err = bitcoin.GetPublicKeysForPath(hdKey, i); err != nil {
			log.Fatalf("error occurred: %s", err.Error())
		}
		for index, key := range pubKeys {
			log.Printf("#%d found at m/%d/%d key: %s", i, index, i, hex.EncodeToString(key.SerialiseCompressed()))
		}
	}
}
