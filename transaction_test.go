package bitcoin

import (
	"fmt"
	"testing"
)

// TestTxFromHex will test the method TxFromHex()
func TestTxFromHex(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputHex      string
		expectedTxID  string
		expectedNil   bool
		expectedError bool
	}{
		{"", "", true, true},
		{"0", "", true, true},
		{"000", "", true, true},
		{"bad-hex", "", true, true},
		{"01000000012adda020db81f2155ebba69e7c841275517ebf91674268c32ff2f5c7e2853b2c010000006b483045022100872051ef0b6c47714130c12a067db4f38b988bfc22fe270731c2146f5229386b02207abf68bbf092ec03e2c616defcc4c868ad1fc3cdbffb34bcedfab391a1274f3e412102affe8c91d0a61235a3d07b1903476a2e2f7a90451b2ed592fea9937696a07077ffffffff02ed1a0000000000001976a91491b3753cf827f139d2dc654ce36f05331138ddb588acc9670300000000001976a914da036233873cc6489ff65a0185e207d243b5154888ac00000000", "64cd12102af20195d54a107e0ee5989ac5db3491893a0b9d42e24354732a22a5", false, false},
	}

	// Run tests
	for _, test := range tests {
		if rawTx, err := TxFromHex(test.inputHex); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.inputHex, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.inputHex)
		} else if rawTx == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was nil but not expected", t.Name(), test.inputHex)
		} else if rawTx != nil && test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was NOT nil but expected to be nil", t.Name(), test.inputHex)
		} else if rawTx != nil && rawTx.GetTxID() != test.expectedTxID {
			t.Errorf("%s Failed: [%s] inputted [%s] expected but failed comparison of txIDs, got: %s", t.Name(), test.inputHex, test.expectedTxID, rawTx.GetTxID())
		}
	}
}

// ExampleTxFromHex example using TxFromHex()
func ExampleTxFromHex() {
	tx, err := TxFromHex("01000000012adda020db81f2155ebba69e7c841275517ebf91674268c32ff2f5c7e2853b2c010000006b483045022100872051ef0b6c47714130c12a067db4f38b988bfc22fe270731c2146f5229386b02207abf68bbf092ec03e2c616defcc4c868ad1fc3cdbffb34bcedfab391a1274f3e412102affe8c91d0a61235a3d07b1903476a2e2f7a90451b2ed592fea9937696a07077ffffffff02ed1a0000000000001976a91491b3753cf827f139d2dc654ce36f05331138ddb588acc9670300000000001976a914da036233873cc6489ff65a0185e207d243b5154888ac00000000")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("txID: %s", tx.GetTxID())
	// Output:txID: 64cd12102af20195d54a107e0ee5989ac5db3491893a0b9d42e24354732a22a5
}

// BenchmarkTxFromHex benchmarks the method TxFromHex()
func BenchmarkTxFromHex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = TxFromHex("01000000012adda020db81f2155ebba69e7c841275517ebf91674268c32ff2f5c7e2853b2c010000006b483045022100872051ef0b6c47714130c12a067db4f38b988bfc22fe270731c2146f5229386b02207abf68bbf092ec03e2c616defcc4c868ad1fc3cdbffb34bcedfab391a1274f3e412102affe8c91d0a61235a3d07b1903476a2e2f7a90451b2ed592fea9937696a07077ffffffff02ed1a0000000000001976a91491b3753cf827f139d2dc654ce36f05331138ddb588acc9670300000000001976a914da036233873cc6489ff65a0185e207d243b5154888ac00000000")
	}
}

// TestCreateTx will test the method CreateTx()
func TestCreateTx(t *testing.T) {

	// Example from: https://github.com/libsv/libsv

	// Use a new UTXO
	utxo := &Utxo{
		TxID:      "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		Vout:      0,
		ScriptSig: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
		Satoshis:  1000,
	}

	// Add a pay-to address
	payTo := &PayToAddress{
		Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
		Satoshis: 500,
	}

	// Add some op return data
	opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
	opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

	// Generate the TX
	rawTx, err := CreateTx([]*Utxo{utxo}, []*PayToAddress{payTo}, []OpReturnData{opReturn1, opReturn2}, "L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Show the results
	t.Logf("created tx: %s", rawTx)
}
