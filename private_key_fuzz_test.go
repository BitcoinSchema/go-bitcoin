package bitcoin

import (
	"testing"
)

// FuzzPrivateKeyFromString tests the PrivateKeyFromString function with
// various hex-encoded private key strings to ensure robust error handling
func FuzzPrivateKeyFromString(f *testing.F) {
	// Seed corpus with valid and invalid private keys
	f.Add("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	f.Add("E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F2B75")
	f.Add("0000000000000000000000000000000000000000000000000000000000000001")
	f.Add("")
	f.Add("0")
	f.Add("00")
	f.Add("invalid")
	f.Add("1234567")
	f.Add("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")
	f.Add("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036414") // Invalid (too short by 1 char)

	f.Fuzz(func(_ *testing.T, privKey string) {
		// The function should not panic regardless of input
		_, _ = PrivateKeyFromString(privKey)
		// Test passes if no panic occurs
	})
}

// FuzzWifToPrivateKey tests the WifToPrivateKey function with various
// WIF strings to ensure proper error handling and no panics
func FuzzWifToPrivateKey(f *testing.F) {
	// Seed corpus with valid and invalid WIF strings
	f.Add("5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei")
	f.Add("5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAbuatmU")
	f.Add("5KgHn2qiftW5LQgCYFtkbrLYB1FuvisDtacax8NCvumw3UTKdcP")
	f.Add("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	f.Add("")
	f.Add("0")
	f.Add("invalid")
	f.Add("5")
	f.Add("L")
	f.Add("5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1E") // Invalid checksum

	f.Fuzz(func(_ *testing.T, wif string) {
		// The function should not panic regardless of input
		_, _ = WifToPrivateKey(wif)
		// Test passes if no panic occurs
	})
}

// FuzzWifToPrivateKeyString tests the WifToPrivateKeyString function with
// various WIF strings to ensure robust error handling
func FuzzWifToPrivateKeyString(f *testing.F) {
	// Seed corpus with valid and invalid WIF strings
	f.Add("5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei")
	f.Add("5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAbuatmU")
	f.Add("5KgHn2qiftW5LQgCYFtkbrLYB1FuvisDtacax8NCvumw3UTKdcP")
	f.Add("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	f.Add("")
	f.Add("0")
	f.Add("invalid")
	f.Add("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")

	f.Fuzz(func(_ *testing.T, wif string) {
		// The function should not panic regardless of input
		_, _ = WifToPrivateKeyString(wif)
		// Test passes if no panic occurs
	})
}

// FuzzPrivateKeyToWifString tests the PrivateKeyToWifString function with
// various hex-encoded private keys to ensure robust error handling
func FuzzPrivateKeyToWifString(f *testing.F) {
	// Seed corpus with valid and invalid private keys
	f.Add("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	f.Add("E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F2B75")
	f.Add("0000000000000000000000000000000000000000000000000000000000000001")
	f.Add("000000")
	f.Add("6D792070726976617465206B6579")
	f.Add("")
	f.Add("0")
	f.Add("invalid")
	f.Add("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")

	f.Fuzz(func(_ *testing.T, privKey string) {
		// The function should not panic regardless of input
		_, _ = PrivateKeyToWifString(privKey)
		// Test passes if no panic occurs
	})
}
