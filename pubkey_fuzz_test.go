package bitcoin

import (
	"testing"
)

// FuzzPubKeyFromString tests the PubKeyFromString function with various
// hex-encoded public key strings to ensure robust error handling
func FuzzPubKeyFromString(f *testing.F) {
	// Seed corpus with valid and invalid public keys
	// Compressed pubkey (33 bytes = 66 hex chars)
	f.Add("02ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3")
	f.Add("03ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3")
	// Uncompressed pubkey (65 bytes = 130 hex chars)
	f.Add("04ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3414a345b88b5e7e32c5f469f2c6e2e0e8c5d4a6f7e8c5d4a6f7e8c5d4a6f7e8c5d")
	f.Add("")
	f.Add("0")
	f.Add("00")
	f.Add("02")
	f.Add("03")
	f.Add("04")
	f.Add("invalid")
	f.Add("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")
	f.Add("02ea87") // Too short
	f.Add("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	f.Add("0000000000000000000000000000000000000000000000000000000000000000")

	f.Fuzz(func(_ *testing.T, pubKey string) {
		// The function should not panic regardless of input
		_, _ = PubKeyFromString(pubKey)
		// Test passes if no panic occurs
	})
}
