package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
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
