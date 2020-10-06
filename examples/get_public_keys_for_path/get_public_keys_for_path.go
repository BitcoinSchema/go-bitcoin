package main

import (
	"encoding/hex"
	"log"

	"github.com/bitcoinschema/go-bitcoin"
	"github.com/bitcoinsv/bsvd/bsvec"
)

func main() {
	// Start with an HD key (we will make one for this example)
	hdKey, err := bitcoin.GenerateHDKey(bitcoin.SecureSeedLength)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get keys by path (example showing 5 sets of keys)
	var pubKeys []*bsvec.PublicKey
	for i := 1; i <= 5; i++ {
		if pubKeys, err = bitcoin.GetPublicKeysForPath(hdKey, uint32(i)); err != nil {
			log.Fatalf("error occurred: %s", err.Error())
		}
		for index, key := range pubKeys {
			log.Printf("#%d found at m/%d/%d key: %s", i, index, i, hex.EncodeToString(key.SerializeCompressed()))
		}
	}
}
