package main

import (
	"log"

	"github.com/libsv/go-bt/v2/bscript"

	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {
	// Start with an HD key (we will make one for this example)
	hdKey, err := bitcoin.GenerateHDKey(bitcoin.SecureSeedLength)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get an address
	var rawAddress *bscript.Address
	if rawAddress, err = bitcoin.GetAddressFromHDKey(hdKey, true); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("got address: %s", rawAddress.AddressString)
}
