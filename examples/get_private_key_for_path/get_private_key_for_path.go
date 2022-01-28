package main

import (
	"encoding/hex"
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/bitcoinsv/bsvd/bsvec"
)

func main() {

	// Start with an HD key (we will make one for this example)
	hdKey, err := bitcoin.GenerateHDKey(bitcoin.SecureSeedLength)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get a private key from a specific path (chain/num)
	var privateKey *bsvec.PrivateKey
	privateKey, err = bitcoin.GetPrivateKeyByPath(hdKey, 10, 2)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("private key: %s for chain/path: %d/%d", hex.EncodeToString(privateKey.Serialize()), 10, 2)
}
