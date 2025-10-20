package bitcoin

import (
	"testing"
)

// FuzzTxFromHex tests the TxFromHex function with various hex-encoded
// transaction strings to ensure robust error handling.
//
// Limitation: Some extremely malformed inputs with invalid varint values may cause
// fatal OOM errors in the underlying libsv/go-bt library that cannot be recovered.
// These represent bugs in the underlying library, not in this wrapper function.
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

	f.Fuzz(func(t *testing.T, txHex string) {
		// Skip inputs that are too long to prevent timeouts
		if len(txHex) > 2_000_000 {
			t.Skip("input too long")
		}

		// The function should not panic for most inputs
		// Some malformed inputs may cause OOM errors in the underlying library
		_, _ = TxFromHex(txHex)
		// Test passes if no panic occurs
	})
}
