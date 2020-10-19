package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin"
)

func main() {
	// Start with an HD key (we will make one for this example)
	hdKey, err := bitcoin.GenerateHDKey(bitcoin.SecureSeedLength)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get the extended public key (xPub)
	var xPub string
	xPub, err = bitcoin.GetExtendedPublicKey(hdKey)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	log.Printf("xPub: %s", xPub)
}
