package main

import (
	"encoding/hex"
	"log"

	"github.com/bitcoinschema/go-bitcoin"
)

func main() {

	// Start with a private key (keep this a secret)
	privateKey, err := bitcoin.CreatePrivateKey()
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}
	log.Println("private key (used for encryption): ", hex.EncodeToString(privateKey.Serialize()))

	// Encrypt
	var data string
	if data, err = bitcoin.EncryptWithPrivateKey(privateKey, `{"some":"data"}`); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Encrypted data
	log.Println("encrypted: ", data)

	// Decrypt the data
	var decrypted string
	decrypted, err = bitcoin.DecryptWithPrivateKey(privateKey, data)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}
	log.Println("decrypted: ", decrypted)

}
