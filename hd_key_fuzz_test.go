package bitcoin

import (
	"testing"
)

// FuzzGenerateHDKeyFromString tests the GenerateHDKeyFromString function
// with various xPriv strings to ensure robust error handling
func FuzzGenerateHDKeyFromString(f *testing.F) {
	// Seed corpus with valid and invalid xPriv strings
	f.Add("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	f.Add("xprv9s21ZrQH143K3XJueaaswvbJ38UX3FhnXkcA7xF8kqeN62qEu116M1XnqaDpSE7SoKp8NxejVJG9dfpuvBC314VZNdB7W1kQN3Viwgkjr8L")
	f.Add("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	f.Add("")
	f.Add("0")
	f.Add("invalid")
	f.Add("xprv")
	f.Add("xpub661MyMwAqRbcH3WGvLjupmr43L1GVH3MP2WQWvdreDraBeFJy64Xxv4LLX9ZVWWz3ZjZkMuZtSsc9qH9JZR74bR4PWkmtEvP423r6DJR8kA")
	f.Add("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")

	f.Fuzz(func(_ *testing.T, xPriv string) {
		// The function should not panic regardless of input
		_, _ = GenerateHDKeyFromString(xPriv)
		// Test passes if no panic occurs
	})
}

// FuzzGetHDKeyFromExtendedPublicKey tests the GetHDKeyFromExtendedPublicKey
// function with various xPub strings to ensure robust error handling
func FuzzGetHDKeyFromExtendedPublicKey(f *testing.F) {
	// Seed corpus with valid and invalid xPub strings
	f.Add("xpub661MyMwAqRbcH3WGvLjupmr43L1GVH3MP2WQWvdreDraBeFJy64Xxv4LLX9ZVWWz3ZjZkMuZtSsc9qH9JZR74bR4PWkmtEvP423r6DJR8kA")
	f.Add("xpub661MyMwAqRbcGjhmJnvR198z2x9XnnDhz2yBtLuTdXQ2VBQj8eJ9RnxmXxKnRPhYy6nLsmabmUfVkbajvP7aZASrrnoZkzmwgyjiNskiefG")
	f.Add("")
	f.Add("0")
	f.Add("invalid")
	f.Add("xpub")
	f.Add("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	f.Add("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")

	f.Fuzz(func(_ *testing.T, xPub string) {
		// The function should not panic regardless of input
		_, _ = GetHDKeyFromExtendedPublicKey(xPub)
		// Test passes if no panic occurs
	})
}
