// Package main demonstrates how to get addresses for a derivation path from an HD key.
package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {
	// Start with an HD key (we will make one for this example)
	hdKey, err := bitcoin.GenerateHDKey(bitcoin.SecureSeedLength)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get the addresses for the given path
	var addresses []string
	addresses, err = bitcoin.GetAddressesForPath(hdKey, 2, true)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("address 1: %s address 2: %s", addresses[0], addresses[1])
}
