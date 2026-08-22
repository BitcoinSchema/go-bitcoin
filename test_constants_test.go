package bitcoin

const (
	testPrivateKeyHex = "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"
	testWIF           = "L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu"
	testTxID          = "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576"
	testScriptPubKey  = "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac"
	testAddress       = "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL"
	testAddress2      = "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5"

	// testChangeAddress is the change address reused across the transaction tests.
	testChangeAddress = "1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc"

	// testUncompressedWIF / testCompressedWIF are the mainnet WIF encodings of
	// testPrivateKeyHex (uncompressed leads with '5', compressed with 'K'/'L').
	testUncompressedWIF = "5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei"
	testCompressedWIF   = "Kz32CUDnArL4eZrGM5NDMhJ5FrduV2MnumwEUePN3TP8AwSRRFvQ"

	// testPubKeyCompressed / testPubKeyUncompressed are the public key of
	// testPrivateKeyHex in hex (compressed = 33 bytes, uncompressed = 65 bytes).
	testPubKeyCompressed   = "031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f"
	testPubKeyUncompressed = "041b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f36e7ef720509250313fcf1b4c5af0dc7c5efa126efe2c3b7008e6f1487c61f31"

	// Frequently reused table-test case names (kept as constants for goconst).
	caseEmpty      = "empty"
	caseSingleZero = "single zero"
)
