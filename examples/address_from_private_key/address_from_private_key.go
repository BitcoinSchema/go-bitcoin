package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin"
)

func main() {

	// Start with a private key (we will make one for this example)
	privateKey, err := bitcoin.CreatePrivateKeyString()
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get an address
	var address string
	if address, err = bitcoin.AddressFromPrivateKey(privateKey); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("found address: %s from private key: %s", address, privateKey)
}
