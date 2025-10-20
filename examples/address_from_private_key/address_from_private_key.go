// Package main demonstrates how to derive a Bitcoin address from a private key.
package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {
	// Start with a private key (we will make one for this example)
	privateKey, err := bitcoin.CreatePrivateKey()
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get an address
	var address string
	if address, err = bitcoin.GetAddressFromPrivateKey(privateKey, true, true); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("found address: %s from private key: %s", address, privateKey)
}
