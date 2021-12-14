package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {

	// Use a new UTXO
	utxo := &bitcoin.Utxo{
		TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		Vout:         0,
		ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
		Satoshis:     1000,
	}

	// Add a pay-to address
	payTo := &bitcoin.PayToAddress{
		Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
		Satoshis: 500,
	}

	// Add some op return data
	opReturn1 := bitcoin.OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
	opReturn2 := bitcoin.OpReturnData{[]byte("prefix2"), []byte("more example data")}

	// Generate the TX
	rawTx, err := bitcoin.CreateTxUsingWif(
		[]*bitcoin.Utxo{utxo},
		[]*bitcoin.PayToAddress{payTo},
		[]bitcoin.OpReturnData{opReturn1, opReturn2},
		"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
	)
	if err != nil {
		log.Printf("error occurred: %s", err.Error())
		return
	}

	// Success!
	log.Printf("rawTx: %s", rawTx.ToString())
}
