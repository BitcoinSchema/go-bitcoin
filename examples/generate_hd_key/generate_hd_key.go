package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin"
)

func main() {
	xPrivateKey, xPublicKey, err := bitcoin.GenerateHDKeyPair(bitcoin.SecureSeedLength)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("xPrivateKey: %s \n xPublicKey: %s", xPrivateKey, xPublicKey)
}
