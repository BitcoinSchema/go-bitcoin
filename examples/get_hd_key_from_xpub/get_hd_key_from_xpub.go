package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {
	// Start with an existing xPub
	xPub := "xpub661MyMwAqRbcH3WGvLjupmr43L1GVH3MP2WQWvdreDraBeFJy64Xxv4LLX9ZVWWz3ZjZkMuZtSsc9qH9JZR74bR4PWkmtEvP423r6DJR8kA"

	// Convert to a HD key
	key, err := bitcoin.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	log.Printf("converted key: %s private: %v", key.String(), key.IsPrivate())
}
