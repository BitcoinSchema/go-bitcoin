// Package main demonstrates how to generate an HD key pair.
package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {
	xPrivateKey, xPublicKey, err := bitcoin.GenerateHDKeyPair(bitcoin.SecureSeedLength)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("xPrivateKey: %s \n xPublicKey: %s", xPrivateKey, xPublicKey)
}
