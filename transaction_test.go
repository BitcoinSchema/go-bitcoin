package bitcoin

import "testing"

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
	opReturns := &OpReturnData{Data: "This is the example data!"}

	// Generate the TX
	rawTx, err := CreateTx([]*Utxo{utxo}, []*PayToAddress{payTo}, []*OpReturnData{opReturns}, "L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Show the results
	t.Logf("created tx: %s", rawTx)
}
