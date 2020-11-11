package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin"
)

func main() {

	// Example values (from sign_message.go)
	privateKey := "ac40784f09304a88b8db0cfcaff3a4e7d81b6b6ccdef68fad8ddf853ea1d1dce"
	signature := "H/sEz5QDQYkXCox9shPB4MMVAVUM/JzfbPHNpPRwNl+hMI2gxy3x7xs9Ed5ryuny5s2hY4Qxc5uirqjMyEEON6k="
	message := "This is the example message"

	rawKey, err := bitcoin.PrivateKeyFromString(privateKey)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Get an address from private key
	var address string
	address, err = bitcoin.GetAddressFromPrivateKey(rawKey, true)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Verify the signature
	if err = bitcoin.VerifyMessage(address, signature, message); err != nil {
		log.Fatalf("verify failed: %s", err.Error())
	} else {
		log.Println("verification passed")
	}
}
