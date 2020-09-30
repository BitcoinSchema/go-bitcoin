package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin"
)

func main() {
	// Create a private key for the example
	privateKey, err := bitcoin.CreatePrivateKeyString()
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Sign the message (returning a signature)
	var signature string
	if signature, err = bitcoin.SignMessage(privateKey, "This is the example message"); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Final signature for the given message
	log.Printf("private key: %s signature: %s", privateKey, signature)
}
