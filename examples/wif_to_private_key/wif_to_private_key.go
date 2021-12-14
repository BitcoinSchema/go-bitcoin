package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {

	// Convert the wif into a private key
	privateKey, err := bitcoin.WifToPrivateKeyString("5KgHn2qiftW5LQgCYFtkbrLYB1FuvisDtacax8NCvumw3UTKdcP")
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("private key: %s", privateKey)
}
