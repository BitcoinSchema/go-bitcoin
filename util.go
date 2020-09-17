package bitcoin

import "encoding/hex"

func HexDecode(str string) []byte {
	b, _ := hex.DecodeString(str)
	return b
}
