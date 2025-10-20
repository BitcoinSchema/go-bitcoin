package main

import (
	"encoding/hex"
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
)

func main() {
	// This data will be encrypted / shared
	testString := "testing 1, 2, 3..."

	// User 1's private key
	privKey1, _ := bitcoin.CreatePrivateKey()

	// User 2's private key
	privKey2, _ := bitcoin.CreatePrivateKey()

	// User 1 encrypts using their private key and user 2's pubkey
	_, _, encryptedData, err := bitcoin.EncryptShared(privKey1, privKey2.PubKey(), []byte(testString))
	if err != nil {
		log.Fatalf("failed to encrypt data for sharing %s", err)
	}

	// Generate the shared key
	user2SharedPrivKey, _ := bitcoin.GenerateSharedKeyPair(privKey2, privKey1.PubKey())

	// User 2 can decrypt using the shared private key
	var decryptedTestData []byte
	decryptedTestData, err = bec.Decrypt(user2SharedPrivKey, encryptedData)
	if err != nil {
		log.Fatalf("failed to decrypt test data %s", err)
	}

	// Success
	log.Printf("test string: %s", testString)
	log.Printf("encrypted: %s", hex.EncodeToString(encryptedData))
	log.Printf("decrypted: %s", decryptedTestData)
}
