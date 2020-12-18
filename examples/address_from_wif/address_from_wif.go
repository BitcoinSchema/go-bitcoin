package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin"
)

func main() {

	// Convert the wif into a private key
	privateKey, err := bitcoin.WifToPrivateKey("5KgHn2qiftW5LQgCYFtkbrLYB1FuvisDtacax8NCvumw3UTKdcP")
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get an address
	var address string
	if address, err = bitcoin.GetAddressFromPrivateKey(privateKey, true); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("found address: %s from private key: %s", address, privateKey)
}
