package bitcoin

import (
	"testing"
)

// FuzzDecryptWithPrivateKeyString tests the DecryptWithPrivateKeyString function
// with various private key and encrypted data combinations
func FuzzDecryptWithPrivateKeyString(f *testing.F) {
	// Seed corpus with valid and invalid combinations
	validPrivKey := "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"
	f.Add(validPrivKey, "")
	f.Add(validPrivKey, "00")
	f.Add(validPrivKey, "invalid")
	f.Add("", "")
	f.Add("", "00")
	f.Add("invalid", "invalid")
	f.Add("0", "0")
	f.Add(validPrivKey, "ffffffffffffffffffffffffffffffff")
	f.Add("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", "00")

	f.Fuzz(func(_ *testing.T, privateKey, encryptedData string) {
		// The function should not panic regardless of input
		// It should handle invalid inputs gracefully with errors
		_, _ = DecryptWithPrivateKeyString(privateKey, encryptedData)
		// Test passes if no panic occurs
	})
}
