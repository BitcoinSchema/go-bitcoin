package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {
	// Start with a private key (keep this a secret)
	privateKey, err := bitcoin.PrivateKeyFromString("b7a1f94ac7be8ed369421c3afe4eae548f10b96435e9c94e35590b85404a5ae4")
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	encryptedData := "00443645948b35d031859c200cd5c73e02ca0020985b837a7c1659660924ded90b0afb94c0de0b77a602fda965d5de3d94677e9200208593a7d1f7cc64023403b90c4562c0cb1cec5cb2849e7fbc5b1fe01b8570f1663bc4bca2e548981e355fb252168b48bc3c7a302a6da2c4d06e8f4900685b7bf9c9530b2b3b7f486d78d43eab21284545"

	// Encrypted data
	log.Println("encrypted: ", encryptedData)

	// Decrypt the data
	var decrypted string
	decrypted, err = bitcoin.DecryptWithPrivateKey(privateKey, encryptedData)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}
	log.Println("decrypted: ", decrypted)
}
