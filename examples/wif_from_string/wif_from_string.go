package main

import (
	"log"

	"github.com/libsv/go-bk/wif"

	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {
	// Create a wif
	wifString, err := bitcoin.CreateWifString()
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Create a wif from a string
	var wifKey *wif.WIF
	wifKey, err = bitcoin.WifFromString(wifString)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("wif key: %s is also: %s", wifString, wifKey.String())
}
