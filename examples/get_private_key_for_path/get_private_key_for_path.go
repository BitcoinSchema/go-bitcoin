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

	// Get a private key from a specific path (chain/num)
	var privateKey *bec.PrivateKey
	privateKey, err = bitcoin.GetPrivateKeyByPath(hdKey, 10, 2)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("private key: %s for chain/path: %d/%d", hex.EncodeToString(privateKey.Serialise()), 10, 2)
}
