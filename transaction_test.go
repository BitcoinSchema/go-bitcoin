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
	t.Logf("created tx: %s", rawTx.ToString())
}

// ExampleCreateTx example using CreateTx()
func ExampleCreateTx() {

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
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("rawTx: %s", rawTx.ToString())
	// Output:rawTx: 0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100eea3d606bd1627be6459a9de4860919225db74843d2fc7f4e7caa5e01f42c2d0022017978d9c6a0e934955a70e7dda71d68cb614f7dd89eb7b9d560aea761834ddd4412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff03f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000
}

// BenchmarkCreateTx benchmarks the method CreateTx()
func BenchmarkCreateTx(b *testing.B) {
	// Use a new UTXO
	utxo := &Utxo{TxID: "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576", Vout: 0, ScriptSig: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac", Satoshis: 1000}

	// Add a pay-to address
	payTo := &PayToAddress{Address: "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL", Satoshis: 500}

	// Add some op return data
	opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
	opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

	for i := 0; i < b.N; i++ {
		_, _ = CreateTx([]*Utxo{utxo}, []*PayToAddress{payTo}, []OpReturnData{opReturn1, opReturn2}, "L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	}
}

// TestCreateTxErrors will test the method CreateTx()
func TestCreateTxErrors(t *testing.T) {

	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputUtxos     []*Utxo
		inputAddresses []*PayToAddress
		inputOpReturns []OpReturnData
		inputWif       string
		expectedRawTx  string
		expectedNil    bool
		expectedError  bool
	}{
		{[]*Utxo{{
			TxID:      "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:      0,
			ScriptSig: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:  1000,
		}},
			[]*PayToAddress{{
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 500,
			}},
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100bd31b3d9fbe18468086c0470e99f096e370f0c6ff41b6bb71f1a1d5c1b068ce302204f0c83d792a40337909b8b1bcea192722161f48dc475c653b7c352baa38eea6c412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff02f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
			false,
			false,
		},
		{nil,
			[]*PayToAddress{{
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 500,
			}},
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"",
			true,
			true,
		},
		{[]*Utxo{{
			TxID:      "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:      0,
			ScriptSig: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:  1000,
		}}, nil,
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006a47304402205ba1a246371bf8db3fb6dfa75e1edaa18b6b86dc1775dc3f2aa3c38f22803ccc022057850f794ebf78e542228d301420d4ec896c30a2bc009b7e55c66120f6c5a57a412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0100000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
			false,
			false,
		},
		{[]*Utxo{{
			TxID:      "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:      0,
			ScriptSig: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:  1000,
		}}, nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006a47304402200083bb297d53210cf9379b3f47de2eff38e6906e5982fbfeef9bf59778750f3e022046da020811e9a2d1e6db8da103d17598abc194125612be6b108d49cb60cbca95412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0000000000",
			false,
			false,
		},
		{[]*Utxo{{
			TxID:      "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:      0,
			ScriptSig: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:  1000,
		}},
			[]*PayToAddress{{
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 500,
			}},
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"",
			"",
			true,
			true,
		},
		{[]*Utxo{{
			TxID:      "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:      0,
			ScriptSig: "invalid-script",
			Satoshis:  1000,
		}},
			[]*PayToAddress{{
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 500,
			}},
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"",
			true,
			true,
		},
		{[]*Utxo{{
			TxID:      "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:      0,
			ScriptSig: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:  1000,
		}},
			[]*PayToAddress{{
				Address:  "invalid-address",
				Satoshis: 500,
			}},
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"",
			true,
			true,
		},
	}

	// Run tests
	for _, test := range tests {
		if rawTx, err := CreateTx(test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputWif); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] [%v] [%v] [%s] inputted and error not expected but got: %s", t.Name(), test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputWif, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] [%v] [%v] [%s] inputted and error was expected", t.Name(), test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputWif)
		} else if rawTx == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] [%v] [%v] [%s] inputted and nil was not expected", t.Name(), test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputWif)
		} else if rawTx != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] [%v] [%v] [%s] inputted and nil was expected", t.Name(), test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputWif)
		} else if rawTx != nil && rawTx.ToString() != test.expectedRawTx {
			t.Errorf("%s Failed: [%v] [%v] [%v] [%s] inputted [%s] expected but failed comparison of scripts, got: %s", t.Name(), test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputWif, test.expectedRawTx, rawTx.ToString())
		}
	}
}
