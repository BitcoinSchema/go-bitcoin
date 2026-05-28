// Package main demonstrates how to create a WIF (Wallet Import Format) string.
package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin/v3"
)

func main() {
	// Create a wif
	wifString, err := bitcoin.CreateWifString()
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("wif key: %s", wifString)
}
