package bitcoin

import "encoding/hex"

// HexDecode returns a decoded hex string without handling errors
// todo: why ignore the error? (@mrz)
func HexDecode(str string) []byte {
	b, _ := hex.DecodeString(str)
	return b
}
