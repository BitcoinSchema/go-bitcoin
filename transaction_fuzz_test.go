package bitcoin

import (
	"testing"
)

// FuzzTxFromHex tests the TxFromHex function with various hex-encoded
// transaction strings to ensure robust error handling
func FuzzTxFromHex(f *testing.F) {
	// Seed corpus with valid and invalid transaction hex strings
	// Valid transaction
	f.Add("0100000001b7b0650a7c3a1bd4716369783476348b59f5404784970192cec1996e869505760000000000ffffffff0100000000000000001976a914eb0bd5edba389198e73f8efabddfc61666969ff788ac00000000")
	f.Add("")
	f.Add("00")
	f.Add("0100")
	f.Add("invalid")
	f.Add("zzzzzzzzzzzzzzzzzzzzzzzz")
	f.Add("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	f.Add("0000000000000000000000000000000000000000000000000000000000000000")
	// Truncated transaction
	f.Add("010000000")
	// Very long hex string
	f.Add("01000000010000000000000000000000000000000000000000000000000000000000000000")

	f.Fuzz(func(_ *testing.T, txHex string) {
		// The function should not panic regardless of input
		_, _ = TxFromHex(txHex)
		// Test passes if no panic occurs
	})
}
