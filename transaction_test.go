package bitcoin

import (
	"fmt"
	"testing"

	"github.com/libsv/libsv/transaction"
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

	// Private key (from wif)
	privateKey, err := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Generate the TX
	var rawTx *transaction.Transaction
	rawTx, err = CreateTx(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{opReturn1, opReturn2},
		privateKey,
	)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Show the results
	t.Logf("created tx: %s", rawTx.ToString())
}

func TestCreateEmptyTx(t *testing.T) {
	// Generate the TX
	var rawTx *transaction.Transaction
	rawTx, err := CreateTx(
		nil,
		nil,
		nil,
		nil,
	)
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

	// Private key (from wif)
	privateKey, err := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Generate the TX
	rawTx, err := CreateTx(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{opReturn1, opReturn2},
		privateKey,
	)
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

	// Private key (from wif)
	privateKey, _ := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")

	for i := 0; i < b.N; i++ {
		_, _ = CreateTx(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			privateKey,
		)
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
		{[]*Utxo{{
			TxID:      "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:      0,
			ScriptSig: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:  1000,
		}},
			[]*PayToAddress{{
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 1500,
			}},
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100bd31b3d9fbe18468086c0470e99f096e370f0c6ff41b6bb71f1a1d5c1b068ce302204f0c83d792a40337909b8b1bcea192722161f48dc475c653b7c352baa38eea6c412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff02f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
			true,
			true,
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

		// Private key (from wif)
		privateKey, err := WifToPrivateKey(test.inputWif)
		if err != nil && !test.expectedError {
			t.Fatalf("error occurred: %s", err.Error())
		}

		if rawTx, err := CreateTx(test.inputUtxos, test.inputAddresses, test.inputOpReturns, privateKey); err != nil && !test.expectedError {
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

// TestCreateTxUsingWif will test the method CreateTxUsingWif()
func TestCreateTxUsingWif(t *testing.T) {

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
	_, err := CreateTxUsingWif(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{opReturn1, opReturn2},
		"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
	)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Invalid wif
	_, err = CreateTxUsingWif(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{opReturn1, opReturn2},
		"",
	)
	if err == nil {
		t.Fatalf("error should have occurred")
	}
}

// ExampleCreateTxUsingWif example using CreateTxUsingWif()
func ExampleCreateTxUsingWif() {

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
	rawTx, err := CreateTxUsingWif(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{opReturn1, opReturn2},
		"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
	)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("rawTx: %s", rawTx.ToString())
	// Output:rawTx: 0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100eea3d606bd1627be6459a9de4860919225db74843d2fc7f4e7caa5e01f42c2d0022017978d9c6a0e934955a70e7dda71d68cb614f7dd89eb7b9d560aea761834ddd4412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff03f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000
}

// BenchmarkCreateTxUsingWif benchmarks the method CreateTxUsingWif()
func BenchmarkCreateTxUsingWif(b *testing.B) {
	// Use a new UTXO
	utxo := &Utxo{TxID: "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576", Vout: 0, ScriptSig: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac", Satoshis: 1000}

	// Add a pay-to address
	payTo := &PayToAddress{Address: "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL", Satoshis: 500}

	// Add some op return data
	opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
	opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

	for i := 0; i < b.N; i++ {
		_, _ = CreateTxUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
	}
}

// TestCalculateFeeForTx will test the method CalculateFeeForTx()
func TestCalculateFeeForTx(t *testing.T) {
	t.Parallel()

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
		Satoshis: 868,
	}

	// Add some op return data
	opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
	opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

	// Generate the TX
	rawTx, err := CreateTxUsingWif(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{opReturn1, opReturn2},
		"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
	)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	expectedFee := uint64(132)

	// Tx to calculate
	// t.Log(rawTx.ToString())
	// t.Log("tx size: ", len(rawTx.ToBytes()))

	// Calculate fee
	if fee := CalculateFeeForTx(rawTx, nil, nil); fee != expectedFee {
		t.Fatalf("expected fee: %d got: %d", expectedFee, fee)
	}
}

// TestCalculateFeeForTxVariousTxs will test the method CalculateFeeForTx()
func TestCalculateFeeForTxVariousTxs(t *testing.T) {

	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputHex          string
		inputStandardRate *FeeAmount
		inputDataRate     *FeeAmount
		expectedTxID      string
		expectedSatoshis  uint64
	}{
		{"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100e07b7661af4e4b521c012a146b25da2c7b9d606e9ceaae28fa73eb347ef6da6f0220527f0638a89ff11cbe53d5f8c4c2962484a370dcd9463a6330f45d31247c2512412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0364030000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000",
			nil,
			nil,
			"e75fa79ee5fbb589201f769c01835e14ca595b7bbfa0e602050a2a90cf28d129",
			132,
		},
		{"0200000001203d3a9d8e2ccfe2d6bb1bce6ad8a9c1251a58f9c788737b21e3e19588e89110010000006a47304402201e7fe22f20d02a5cd6978b21bc68aa31eb74530c6b75b47caef19d6f2a95f47802206b20de32c7398fa822397b49389a49224742a9d3c175a851376a25deb6428cb24121023bc787128d6296be6a534b32e90724413179545a0ec2720a7258de330fc54544ffffffff02000000000000000053006a4c4fe5b08fe881aae6b8b8e6888f3130e69c883234e697a5e6b7b1e59cb3e7babfe4b88be994a6e6a087e8b59befbc8ce586a0e5869b425356e5a596e98791e7ad89e4bda0e69da5e68c91e68898efbc817c0c0000000000001976a9146382643e30d2ea7e52eb07eac8767ed219c53e3a88ac00000000",
			nil,
			nil,
			"1affabe9b5adc3a6930a06002a447e834681004a5e6767a649d0371a806e7b1d",
			141,
		},
		{"0200000001a39b865acf5d100aa35341a63e4ba6dc7da101f3b740cd5372e25b0fb23306c5000000006a47304402207cdcab521641801cd0427501c0264073cfd11f2693bb9a109d00482c643f16b30220798801d1acc9dea23013f17cbd7cc0edec996abc1493f9b9a8ba61723b4ec01d41210250a932cf2543f8f35dba3fce46f53ca821409f4e233d9f8090a165f5ac42a0e8ffffffff014a060000000000001976a914f21ff89bd2259699d229ebc1f9c29ac3c3e0411888ac00000000",
			nil,
			nil,
			"175fc22ffbd76f1cdb7ec2c40474abedbb4bb6080e9c8d22736ffc2d48e85fd2",
			95,
		},
		{"010000000129cdf21be448b28b0f88ecc9317d9ea1a385a6928d855f55b2d62a3bbc8631cd000000006b483045022100a8b7d34e10e817647f753a9a1380606945bd9c255ebb87c7a00a701903abc6e20220284651051ef520eb5f13901c0f792d7eae2e75b9e85fff36258d9592839227614121029843434e00940e0829196042983e46fe801364b1a586f3d4160556684d63e8d1ffffffff030000000000000000fd48026a2231394878696756345179427633744870515663554551797131707a5a56646f4175744cff545742544348205457455443482054574554534820545745544348205457455443482054574554434820545745544348205457455443482054564554434820545745544348205457454943482054574554434820545745544348205457455453482054574554434820545745544348205457415443482054574554434820545745574348205457455443482054574554434820545745544548205457455443482054574554434820545745545348205457455443482054574554434820545745544f482054574554434820544d45544348205457455443482054574554454820545745544348205768617427732077697468696e205457455443483f203a500a746578742f706c61696e04746578741f7477657463685f7477746578745f313536383139363330383230372e747874017c223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540b7477646174615f6a736f6e046e756c6c0375726c046e756c6c07636f6d6d656e74046e756c6c076d625f757365720434363631047479706504706f73740974696d657374616d700f3638343131393038353536373938330361707006747765746368017c22313550636948473232534e4c514a584d6f53556157566937575371633768436676610d424954434f494e5f45434453412231414b48566959674247626d78693871694a6b4e766f484e654475396d334d665045102323636f6d70757465645f7369672323953f0000000000001976a9147346219a70740418422b19c0ecd696671d3ce17088acdbdf0500000000001976a914b4376ffbd47974340a0a0d940681e168c1d6ae0788ac00000000",
			nil,
			nil,
			"d612aed6c3f12756ea1d3e9a48eef2faf05eaf20ff6031911d8610332f8d3f9a",
			410,
		},
		{"01000000018558c697a0bb502f9aec70b27c48d3c87f0df28e16f9c8f43a43b1327ea69304010000006b483045022100e0114cf815c0d60fc72e7c8a4c39b03920ac03e62b5c98c68882c1235f2e0ac90220278319c67e79c4397ad989cb9be453acfbbb8e8a9f927c6f4509d9a1da4a17eb412103637720f48f059b854605da3d5a959d204e1e6cbd0d03ee91a38a067bdee90558ffffffff020000000000000000d96a2231394878696756345179427633744870515663554551797131707a5a56646f4175742c436865636b6f757420546f6e6963506f773a2068747470733a2f2f746e6370772e636f2f36313239373839330d746578742f6d61726b646f776e055554462d38017c223150755161374b36324d694b43747373534c4b79316b683536575755374d745552350353455403617070086d6574616c656e73047479706507636f6d6d656e740375726c2268747470733a2f2f6f66666572732e746f6e6963706f772e636f6d2f6f666665727304757365720444756465d3bb1100000000001976a914e0190be1c0ced1e92c26b3c6d5ccbbe853b57e7288ac00000000",
			nil,
			nil,
			"ef3fe744c7be4b81f2881f9c11433bc4905a032a0bf15bcda31f90c28ba1b356",
			209,
		},
		{"010000000190c5208bed05b4e54746bab5a5d4a5e4324e2999c180b7ea7105047f1f16b84a030000006a47304402202ca5d7a2cfb2388babb549b10c47ed20cdadafb845af835c7d5ff04e933ba1c102200a8b7289bbd3c0cc62172afe0006ba374ddd0132a7e4fb4e92ebcff5ce9217db412102812d641ff356c362815f8fc03bd061c2ae57e02f5dc3083f61785c0e5b198039ffffffff040000000000000000ad006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f696403323934106f666665725f73657373696f6e5f696440653638616631393439626131633239326131376431393635343638383234383663653635313830663465363439383235613561363532376634646132303761661f0c0000000000001976a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac31790000000000001976a9147c8ced9ee0f48192822a0148f27b5a1f24aa42d388ac53f74600000000001976a914852f89e9b05d6adfc905842ee0d301947d675df988ac00000000",
			nil,
			nil,
			"8785ca5f11795a38eb1f50f62562cb5e0335b283762fe8a2c7e96d5f7f79bb15",
			220,
		},
		{"010000000190c5208bed05b4e54746bab5a5d4a5e4324e2999c180b7ea7105047f1f16b84a030000006a47304402202ca5d7a2cfb2388babb549b10c47ed20cdadafb845af835c7d5ff04e933ba1c102200a8b7289bbd3c0cc62172afe0006ba374ddd0132a7e4fb4e92ebcff5ce9217db412102812d641ff356c362815f8fc03bd061c2ae57e02f5dc3083f61785c0e5b198039ffffffff040000000000000000ad006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f696403323934106f666665725f73657373696f6e5f696440653638616631393439626131633239326131376431393635343638383234383663653635313830663465363439383235613561363532376634646132303761661f0c0000000000001976a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac31790000000000001976a9147c8ced9ee0f48192822a0148f27b5a1f24aa42d388ac53f74600000000001976a914852f89e9b05d6adfc905842ee0d301947d675df988ac00000000",
			&FeeAmount{Bytes: DefaultRateBytes, Satoshis: 1000},
			&FeeAmount{Bytes: DefaultRateBytes, Satoshis: 1000},
			"8785ca5f11795a38eb1f50f62562cb5e0335b283762fe8a2c7e96d5f7f79bb15",
			441,
		},
		{"010000000190c5208bed05b4e54746bab5a5d4a5e4324e2999c180b7ea7105047f1f16b84a030000006a47304402202ca5d7a2cfb2388babb549b10c47ed20cdadafb845af835c7d5ff04e933ba1c102200a8b7289bbd3c0cc62172afe0006ba374ddd0132a7e4fb4e92ebcff5ce9217db412102812d641ff356c362815f8fc03bd061c2ae57e02f5dc3083f61785c0e5b198039ffffffff040000000000000000ad006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f696403323934106f666665725f73657373696f6e5f696440653638616631393439626131633239326131376431393635343638383234383663653635313830663465363439383235613561363532376634646132303761661f0c0000000000001976a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac31790000000000001976a9147c8ced9ee0f48192822a0148f27b5a1f24aa42d388ac53f74600000000001976a914852f89e9b05d6adfc905842ee0d301947d675df988ac00000000",
			&FeeAmount{Bytes: DefaultRateBytes, Satoshis: 250},
			&FeeAmount{Bytes: DefaultRateBytes, Satoshis: 250},
			"8785ca5f11795a38eb1f50f62562cb5e0335b283762fe8a2c7e96d5f7f79bb15",
			109,
		},
		{"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006a47304402200083bb297d53210cf9379b3f47de2eff38e6906e5982fbfeef9bf59778750f3e022046da020811e9a2d1e6db8da103d17598abc194125612be6b108d49cb60cbca95412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0000000000",
			&FeeAmount{Bytes: 157, Satoshis: 1},
			&FeeAmount{Bytes: 157, Satoshis: 1},
			"d3350a4ef4b2c72b23e5117979590d768e61f2102337e2ae956d152a80cd37ac",
			1,
		},
	}

	// Run tests
	var satoshis uint64
	for _, test := range tests {

		// Get the tx from hex string
		tx, err := TxFromHex(test.inputHex)
		if err != nil {
			t.Fatalf("error occurred: %s", err.Error())
		}

		// Test the function
		if satoshis = CalculateFeeForTx(tx, test.inputStandardRate, test.inputDataRate); satoshis != test.expectedSatoshis {
			t.Errorf("%s Failed: [%s] [%v] [%v] inputted [%d] expected but got: %d", t.Name(), test.inputHex, test.inputStandardRate, test.inputDataRate, test.expectedSatoshis, satoshis)
		} else if tx.GetTxID() != test.expectedTxID {
			t.Errorf("%s Failed: [%s] [%v] [%v] inputted [%s] expected but got: %s", t.Name(), test.inputHex, test.inputStandardRate, test.inputDataRate, test.expectedTxID, tx.GetTxID())
		}
	}
}

// TestA25_ComputeChecksum will test the method ComputeChecksum()
func TestA25_ComputeChecksum(t *testing.T) {
	utxo := &Utxo{
		TxID:      "",
		Vout:      0,
		ScriptSig: "",
		Satoshis:  1000,
	}
	tx, err := CreateTxUsingWif([]*Utxo{utxo}, nil, nil, "L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	t.Log(tx.ToString())
}

// ExampleCalculateFeeForTx example using CalculateFeeForTx()
func ExampleCalculateFeeForTx() {

	// Get the tx from hex string
	rawTx := "0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100e07b7661af4e4b521c012a146b25da2c7b9d606e9ceaae28fa73eb347ef6da6f0220527f0638a89ff11cbe53d5f8c4c2962484a370dcd9463a6330f45d31247c2512412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0364030000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000"
	tx, err := TxFromHex(rawTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Calculate the fee using default rates
	estimatedFee := CalculateFeeForTx(tx, nil, nil)

	fmt.Printf("tx id: %s estimated fee: %d satoshis", tx.GetTxID(), estimatedFee)
	// Output:tx id: e75fa79ee5fbb589201f769c01835e14ca595b7bbfa0e602050a2a90cf28d129 estimated fee: 132 satoshis
}

// BenchmarkCalculateFeeForTx benchmarks the method CalculateFeeForTx()
func BenchmarkCalculateFeeForTx(b *testing.B) {

	// Get the tx from hex string
	rawTx := "0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100e07b7661af4e4b521c012a146b25da2c7b9d606e9ceaae28fa73eb347ef6da6f0220527f0638a89ff11cbe53d5f8c4c2962484a370dcd9463a6330f45d31247c2512412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0364030000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000"
	tx, err := TxFromHex(rawTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	for i := 0; i < b.N; i++ {
		_ = CalculateFeeForTx(tx, nil, nil)
	}
}

// TestCalculateFeeForTxPanic tests for nil case in CalculateFeeForTx()
func TestCalculateFeeForTxPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()

	_ = CalculateFeeForTx(nil, nil, nil)
}

// TestCreateTxWithChange tests for nil case in CreateTxWithChange()
func TestCreateTxWithChange(t *testing.T) {

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

	// Private key (from wif)
	privateKey, err := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Generate the TX
	var rawTx *transaction.Transaction
	rawTx, err = CreateTxWithChange(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{opReturn1, opReturn2},
		"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
		nil,
		nil,
		privateKey,
	)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Show the results
	t.Logf("created tx: %s", rawTx.ToString())

	// Expected
	expectedFee := uint64(149)
	expectedChange := uint64(352)

	// Test the right fee
	fee := CalculateFeeForTx(rawTx, nil, nil)
	if fee != expectedFee {
		t.Fatalf("fee expected: %d vs %d", expectedFee, fee)
	}

	// Test that we got the right amount of change (satoshis)
	for _, out := range rawTx.GetOutputs() {
		if out.GetLockingScriptHexString() == "76a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac" {
			if out.Satoshis != expectedChange {
				t.Fatalf("incorrect change expected: %d vs %d", out.Satoshis, expectedChange)
			}
		}
	}
}

// TestCreateTxWithChangeErrors tests for nil case in CreateTxWithChange()
func TestCreateTxWithChangeErrors(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputUtxos         []*Utxo
		inputAddresses     []*PayToAddress
		inputOpReturns     []OpReturnData
		inputWif           string
		inputChangeAddress string
		inputStandardRate  *FeeAmount
		inputDataRate      *FeeAmount
		expectedRawTx      string
		expectedNil        bool
		expectedError      bool
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
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006a4730440220767b6b7483746c1b07e8f58dd953bdc9ade3696036fe093cc404660704d7b407022063858090101d990568af473834d63372c9d31cb678733f9e50a62b2f047c2aa5412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff03f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac72010000000000001976a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
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
				Satoshis: 1500,
			}},
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100bd31b3d9fbe18468086c0470e99f096e370f0c6ff41b6bb71f1a1d5c1b068ce302204f0c83d792a40337909b8b1bcea192722161f48dc475c653b7c352baa38eea6c412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff02f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
			true,
			true,
		},
		{nil,
			[]*PayToAddress{{
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 500,
			}},
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
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
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006a47304402203fa8a408ef7b523a7efab01e0016d4e4879b3e9a56418691cd441344e326fd2c022066997fde9d357a83dec6fc80722089f740f393d6ef949f966f6739da44282f44412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0278030000000000001976a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
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
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006a4730440220143c042ecbdb9296d84a87155dcaf558b657be627ab54ffe464028c050cf0da902205117bd7aa2f85d3cd2f98b504a52599e2ea69ea6afb6a640637c12ddcef8c2f0412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0188030000000000001976a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac00000000",
			false,
			false,
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
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
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
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
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
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 500,
			}},
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"",
			nil,
			nil,
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
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 1001,
			}},
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
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
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 950,
			}},
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"",
			true,
			true,
		},
	}

	// Run tests
	var rawTx *transaction.Transaction
	for _, test := range tests {
		privateKey, err := WifToPrivateKey(test.inputWif)
		if err != nil && !test.expectedError {
			t.Fatalf("error occurred: %s", err.Error())
		}

		if rawTx, err = CreateTxWithChange(test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputChangeAddress, test.inputStandardRate, test.inputDataRate, privateKey); err != nil && !test.expectedError {
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

// ExampleCreateTxWithChange example using CreateTxWithChange()
func ExampleCreateTxWithChange() {

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

	// Get private key from wif
	privateKey, err := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Generate the TX
	var rawTx *transaction.Transaction
	rawTx, err = CreateTxWithChange(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{opReturn1, opReturn2},
		"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
		nil,
		nil,
		privateKey,
	)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("rawTx: %s", rawTx.ToString())
	// Output:rawTx: 0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100beab95997a8b4b0e805aa16af1fed54a0ff80d6e45a330f71787795c394ff99d02207904b4930ccf4dae9e87d1f4b18be343f2fd73bb5500870d7194726919b5f6d8412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff04f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac60010000000000001976a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000
}

// BenchmarkCreateTxWithChange benchmarks the method CreateTxWithChange()
func BenchmarkCreateTxWithChange(b *testing.B) {
	// Use a new UTXO
	utxo := &Utxo{TxID: "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576", Vout: 0, ScriptSig: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac", Satoshis: 1000}

	// Add a pay-to address
	payTo := &PayToAddress{Address: "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL", Satoshis: 500}

	// Add some op return data
	opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
	opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

	// Get private key from wif
	privateKey, _ := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")

	for i := 0; i < b.N; i++ {
		_, _ = CreateTxWithChange(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			privateKey,
		)
	}
}

// TestCreateTxWithChangeUsingWif tests for nil case in CreateTxWithChangeUsingWif()
func TestCreateTxWithChangeUsingWif(t *testing.T) {

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
	rawTx, err := CreateTxWithChangeUsingWif(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{opReturn1, opReturn2},
		"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
		nil,
		nil,
		"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
	)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Show the results
	t.Logf("created tx: %s", rawTx.ToString())

	// Expected
	expectedFee := uint64(149)
	expectedChange := uint64(352)

	// Test the right fee
	fee := CalculateFeeForTx(rawTx, nil, nil)
	if fee != expectedFee {
		t.Fatalf("fee expected: %d vs %d", expectedFee, fee)
	}

	// Test that we got the right amount of change (satoshis)
	for _, out := range rawTx.GetOutputs() {
		if out.GetLockingScriptHexString() == "76a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac" {
			if out.Satoshis != expectedChange {
				t.Fatalf("incorrect change expected: %d vs %d", out.Satoshis, expectedChange)
			}
		}
	}

	// Invalid wif
	_, err = CreateTxWithChangeUsingWif(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{opReturn1, opReturn2},
		"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
		nil,
		nil,
		"",
	)
	if err == nil {
		t.Fatalf("error should have occurred")
	}
}

// ExampleCreateTxWithChangeUsingWif example using CreateTxWithChangeUsingWif()
func ExampleCreateTxWithChangeUsingWif() {

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
	rawTx, err := CreateTxWithChangeUsingWif(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{opReturn1, opReturn2},
		"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
		nil,
		nil,
		"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
	)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("rawTx: %s", rawTx.ToString())
	// Output:rawTx: 0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100beab95997a8b4b0e805aa16af1fed54a0ff80d6e45a330f71787795c394ff99d02207904b4930ccf4dae9e87d1f4b18be343f2fd73bb5500870d7194726919b5f6d8412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff04f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac60010000000000001976a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000
}

// BenchmarkCreateTxWithChangeUsingWif benchmarks the method CreateTxWithChangeUsingWif()
func BenchmarkCreateTxWithChangeUsingWif(b *testing.B) {
	// Use a new UTXO
	utxo := &Utxo{TxID: "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576", Vout: 0, ScriptSig: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac", Satoshis: 1000}

	// Add a pay-to address
	payTo := &PayToAddress{Address: "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL", Satoshis: 500}

	// Add some op return data
	opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
	opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

	for i := 0; i < b.N; i++ {
		_, _ = CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
	}
}
