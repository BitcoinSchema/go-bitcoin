package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin"
	"github.com/bitcoinsv/bsvutil"
)

func main() {

	// Start with an HD key (we will make one for this example)
	hdKey, err := bitcoin.GenerateHDKey(bitcoin.SecureSeedLength)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get an address
	var rawAddress *bsvutil.LegacyAddressPubKeyHash
	if rawAddress, err = bitcoin.GetAddressFromHDKey(hdKey); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Success!
	log.Printf("got address: %s", rawAddress.String())
}
