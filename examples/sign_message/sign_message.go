package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {
	// Create a private key  (we will make one for this example)
	privateKey, err := bitcoin.CreatePrivateKeyString()
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Sign the message (returning a signature)

	// Note: If your signature references a compressed key,
	// the address you provide to verify must also come from a compressed key
	var signature string
	if signature, err = bitcoin.SignMessage(privateKey, "This is the example message", false); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Final signature for the given message
	log.Printf("private key: %s signature: %s", privateKey, signature)
}
