// Package main demonstrates how to derive a private key from an HD key using a derivation path.
package main

import (
	"encoding/hex"
	"log"

	"github.com/bitcoinschema/go-bitcoin/v3"
	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
)

func main() {
	// Start with an HD key (we will make one for this example)
	hdKey, err := bitcoin.GenerateHDKey(bitcoin.SecureSeedLength)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get a private key from a specific path (chain/num)
	var privateKey *ec.PrivateKey
	privateKey, err = bitcoin.GetPrivateKeyByPath(hdKey, 10, 2)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("private key: %s for chain/path: %d/%d", hex.EncodeToString(privateKey.Serialize()), 10, 2)
}
