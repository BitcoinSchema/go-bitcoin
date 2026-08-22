package bitcoin

import (
	"testing"
)

// FuzzTxFromHex tests TxFromHex with various raw hex strings to ensure it never
// panics regardless of input.
func FuzzTxFromHex(f *testing.F) {
	// Seed corpus with valid and invalid raw transaction hex
	f.Add("01000000012adda020db81f2155ebba69e7c841275517ebf91674268c32ff2f5c7e2853b2c010000006b483045022100872051ef0b6c47714130c12a067db4f38b988bfc22fe270731c2146f5229386b02207abf68bbf092ec03e2c616defcc4c868ad1fc3cdbffb34bcedfab391a1274f3e412102affe8c91d0a61235a3d07b1903476a2e2f7a90451b2ed592fea9937696a07077ffffffff02ed1a0000000000001976a91491b3753cf827f139d2dc654ce36f05331138ddb588acc9670300000000001976a914da036233873cc6489ff65a0185e207d243b5154888ac00000000")
	f.Add("")
	f.Add("0")
	f.Add("00")
	f.Add("bad-hex")
	f.Add("deadbeef")
	f.Add("ffffffffffffffffffffffffffffffff")

	f.Fuzz(func(_ *testing.T, rawHex string) {
		// The function should not panic regardless of input
		_, _ = TxFromHex(rawHex)
	})
}
