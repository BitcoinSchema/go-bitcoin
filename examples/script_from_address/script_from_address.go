package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin"
)

func main() {
	// Start with a private key
	privateKey, err := bitcoin.CreatePrivateKeyString()
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get an address
	var address string
	if address, err = bitcoin.AddressFromPrivateKey(privateKey); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get the script
	var script string
	if script, err = bitcoin.ScriptFromAddress(address); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("generated script: %s from address: %s", script, address)
}
