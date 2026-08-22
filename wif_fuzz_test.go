package bitcoin

import (
	"testing"
)

// FuzzDecodeWIF tests DecodeWIF with various strings to ensure it never panics
// regardless of input.
func FuzzDecodeWIF(f *testing.F) {
	// Seed corpus with valid (compressed/uncompressed) and invalid WIFs
	f.Add(testWIF)
	f.Add(testUncompressedWIF)
	f.Add(testCompressedWIF)
	f.Add("")
	f.Add("0")
	f.Add("invalid")
	f.Add("5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1E") // truncated
	f.Add("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")

	f.Fuzz(func(_ *testing.T, wif string) {
		// The function should not panic regardless of input
		_, _ = DecodeWIF(wif)
	})
}
