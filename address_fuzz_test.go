package bitcoin

import (
	"testing"
)

// FuzzSet58 tests the Set58 method with various inputs to ensure it handles
// malformed base58 strings correctly without panicking
func FuzzSet58(f *testing.F) {
	// Seed corpus with valid and edge-case addresses
	f.Add([]byte("1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi2"))
	f.Add([]byte("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"))
	f.Add([]byte("1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2"))
	f.Add([]byte(""))
	f.Add([]byte("0"))
	f.Add([]byte("1"))
	f.Add([]byte("invalid"))
	f.Add([]byte("!!!"))
	f.Add([]byte("1KCEAmV"))
	f.Add([]byte("000000000000000000000000000000"))

	f.Fuzz(func(_ *testing.T, data []byte) {
		var a A25
		// The function should not panic regardless of input
		_ = a.Set58(data)
		// Test passes if no panic occurs
	})
}

// FuzzValidA58 tests the ValidA58 function with various inputs to ensure
// it properly validates base58 addresses without panicking
func FuzzValidA58(f *testing.F) {
	// Seed corpus with valid and invalid addresses
	f.Add([]byte("1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi2"))
	f.Add([]byte("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"))
	f.Add([]byte("1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2"))
	f.Add([]byte("1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi")) // Invalid checksum
	f.Add([]byte(""))
	f.Add([]byte("0"))
	f.Add([]byte("1"))
	f.Add([]byte("invalid"))
	f.Add([]byte("1KCEAmV"))
	f.Add([]byte("000000000000000000000000000000"))
	f.Add([]byte("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"))

	f.Fuzz(func(_ *testing.T, data []byte) {
		// The function should not panic regardless of input
		// We don't care about the result, just that it doesn't crash
		_, _ = ValidA58(data)
		// Test passes if no panic occurs
	})
}

// FuzzGetAddressFromScript tests the GetAddressFromScript function with
// various hex-encoded scripts to ensure proper error handling
func FuzzGetAddressFromScript(f *testing.F) {
	// Seed corpus with valid and invalid scripts
	f.Add("76a914eb0bd5edba389198e73f8efabddfc61666969ff788ac") // Valid P2PKH script
	f.Add("76a914000000000000000000000000000000000000000088ac")
	f.Add("")
	f.Add("00")
	f.Add("76a914")
	f.Add("invalid")
	f.Add("0000000000")
	f.Add("ffffffffffffffffffffffffffffffff")
	f.Add("6a")   // OP_RETURN
	f.Add("a914") // P2SH prefix

	f.Fuzz(func(_ *testing.T, script string) {
		// The function should not panic regardless of input
		// We accept both success and error, just no panics
		_, _ = GetAddressFromScript(script)
		// Test passes if no panic occurs
	})
}

// FuzzGetAddressFromPubKeyString tests the GetAddressFromPubKeyString function
// with various hex-encoded public keys to ensure robust error handling
func FuzzGetAddressFromPubKeyString(f *testing.F) {
	// Seed corpus with valid and invalid public keys
	// Compressed pubkey (33 bytes)
	f.Add("02ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3", true, true)
	// Uncompressed pubkey (65 bytes)
	f.Add("04ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3414a345b88b5e7e32c5f469f2c6e2e0e8c5d4a6f7e8c5d4a6f7e8c5d4a6f7e8c", false, true)
	f.Add("", true, true)
	f.Add("00", true, true)
	f.Add("invalid", true, true)
	f.Add("0000000000000000000000000000000000000000000000000000000000000000", true, true)
	f.Add("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", true, true)

	f.Fuzz(func(_ *testing.T, pubKey string, compressed, mainnet bool) {
		// The function should not panic regardless of input
		_, _ = GetAddressFromPubKeyString(pubKey, compressed, mainnet)
		// Test passes if no panic occurs
	})
}
