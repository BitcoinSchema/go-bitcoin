package bitcoin

import (
	"testing"
)

// FuzzPubKeyFromSignature tests the PubKeyFromSignature function with
// various signature and message combinations to ensure robust error handling
func FuzzPubKeyFromSignature(f *testing.F) {
	// Seed corpus with valid and invalid base64 signatures
	f.Add("IG8OFjRpLGYe3GAF3xsYN97qTmd+ZhfNGhm8YGVCHUJyKoNzOxm8Yh3oC3KCiNjXQo5u4DLQ3xI3OhqUz9MO4Qc=", "test message")
	f.Add("", "")
	f.Add("invalid", "test")
	f.Add("AAAA", "test")
	f.Add("", "test message")
	f.Add("IG8OFjRpLGYe3GAF3xsYN97qTmd+ZhfNGhm8YGVCHUJyKoNzOxm8Yh3oC3KCiNjXQo5u4DLQ3xI3OhqUz9MO4Qc=", "")
	f.Add("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", "message")

	f.Fuzz(func(_ *testing.T, sig, data string) {
		// The function should not panic regardless of input
		_, _, _ = PubKeyFromSignature(sig, data)
		// Test passes if no panic occurs
	})
}

// FuzzVerifyMessage tests the VerifyMessage function with various
// address, signature, and message combinations
func FuzzVerifyMessage(f *testing.F) {
	// Seed corpus with valid and invalid combinations
	f.Add("1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi2", "IG8OFjRpLGYe3GAF3xsYN97qTmd+ZhfNGhm8YGVCHUJyKoNzOxm8Yh3oC3KCiNjXQo5u4DLQ3xI3OhqUz9MO4Qc=", "test message", true)
	f.Add("", "", "", true)
	f.Add("invalid", "invalid", "test", false)
	f.Add("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "AAAA", "test", true)
	f.Add("", "IG8OFjRpLGYe3GAF3xsYN97qTmd+ZhfNGhm8YGVCHUJyKoNzOxm8Yh3oC3KCiNjXQo5u4DLQ3xI3OhqUz9MO4Qc=", "test", true)

	f.Fuzz(func(_ *testing.T, address, sig, data string, mainnet bool) {
		// The function should not panic regardless of input
		// It should return an error for invalid inputs
		_ = VerifyMessage(address, sig, data, mainnet)
		// Test passes if no panic occurs
	})
}

// FuzzVerifyMessageDER tests the VerifyMessageDER function with various
// hash, public key, and signature combinations
func FuzzVerifyMessageDER(f *testing.F) {
	// Seed corpus with valid and invalid combinations
	validHash := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
	}
	zeroHash := []byte{}

	f.Add(validHash, "02ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3", "3045022100...")
	f.Add(zeroHash, "", "")
	f.Add(validHash, "invalid", "invalid")
	f.Add(zeroHash, "02ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3", "00")

	f.Fuzz(func(_ *testing.T, hashBytes []byte, pubKey, signature string) {
		// Convert []byte to [32]byte
		var hash [32]byte
		if len(hashBytes) >= 32 {
			copy(hash[:], hashBytes[:32])
		} else {
			copy(hash[:], hashBytes)
		}

		// The function should not panic regardless of input
		_, _ = VerifyMessageDER(hash, pubKey, signature)
		// Test passes if no panic occurs
	})
}
