package bitcoin

import (
	"fmt"
	"testing"

	"github.com/libsv/go-bt/v2"
	"github.com/stretchr/testify/assert"
)

// TestTxFromHex will test the method TxFromHex()
func TestTxFromHex(t *testing.T) {
	t.Parallel()

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

	for _, test := range tests {
		if rawTx, err := TxFromHex(test.inputHex); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.inputHex, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error was expected", t.Name(), test.inputHex)
		} else if rawTx == nil && !test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and was nil but not expected", t.Name(), test.inputHex)
		} else if rawTx != nil && test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and was NOT nil but expected to be nil", t.Name(), test.inputHex)
		} else if rawTx != nil && rawTx.TxID() != test.expectedTxID {
			t.Fatalf("%s Failed: [%s] inputted [%s] expected but failed comparison of txIDs, got: %s", t.Name(), test.inputHex, test.expectedTxID, rawTx.TxID())
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
	fmt.Printf("txID: %s", tx.TxID())
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
	t.Parallel()

	t.Run("basic valid tx", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 500,
		}

		opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
		opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

		privateKey, err := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
		assert.NoError(t, err)
		assert.NotNil(t, privateKey)

		var tx *bt.Tx
		tx, err = CreateTx(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			privateKey,
		)
		assert.NoError(t, err)
		assert.NotNil(t, tx)
		assert.Equal(t,
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100eea3d606bd1627be6459a9de4860919225db74843d2fc7f4e7caa5e01f42c2d0022017978d9c6a0e934955a70e7dda71d68cb614f7dd89eb7b9d560aea761834ddd4412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff03f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000",
			tx.String(),
		)
	})

	t.Run("empty tx", func(t *testing.T) {
		tx, err := CreateTx(
			nil,
			nil,
			nil,
			nil,
		)
		assert.NoError(t, err)
		assert.NotNil(t, tx)
		assert.Equal(t, "01000000000000000000", tx.String())
	})
}

// ExampleCreateTx example using CreateTx()
func ExampleCreateTx() {

	// Use a new UTXO
	utxo := &Utxo{
		TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		Vout:         0,
		ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
		Satoshis:     1000,
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

	fmt.Printf("rawTx: %s", rawTx.String())
	// Output:rawTx: 0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100eea3d606bd1627be6459a9de4860919225db74843d2fc7f4e7caa5e01f42c2d0022017978d9c6a0e934955a70e7dda71d68cb614f7dd89eb7b9d560aea761834ddd4412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff03f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000
}

// BenchmarkCreateTx benchmarks the method CreateTx()
func BenchmarkCreateTx(b *testing.B) {
	// Use a new UTXO
	utxo := &Utxo{TxID: "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576", Vout: 0, ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac", Satoshis: 1000}

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
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
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
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
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
			"010000000002f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
			false,
			false,
		},
		{[]*Utxo{{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}}, nil,
			[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006a47304402205ba1a246371bf8db3fb6dfa75e1edaa18b6b86dc1775dc3f2aa3c38f22803ccc022057850f794ebf78e542228d301420d4ec896c30a2bc009b7e55c66120f6c5a57a412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0100000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
			false,
			false,
		},
		{[]*Utxo{{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}}, nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006a47304402200083bb297d53210cf9379b3f47de2eff38e6906e5982fbfeef9bf59778750f3e022046da020811e9a2d1e6db8da103d17598abc194125612be6b108d49cb60cbca95412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0000000000",
			false,
			false,
		},
		{[]*Utxo{{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "invalid-script",
			Satoshis:     1000,
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
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
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

	var rawTx *bt.Tx
	for _, test := range tests {

		// Private key (from wif)
		privateKey, err := WifToPrivateKey(test.inputWif)
		if err != nil && !test.expectedError {
			t.Fatalf("error occurred: %s", err.Error())
		}

		if rawTx, err = CreateTx(test.inputUtxos, test.inputAddresses, test.inputOpReturns, privateKey); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%v] [%v] [%v] [%s] inputted and error not expected but got: %s", t.Name(), test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputWif, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%v] [%v] [%v] [%s] inputted and error was expected", t.Name(), test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputWif)
		} else if rawTx == nil && !test.expectedNil {
			t.Fatalf("%s Failed: [%v] [%v] [%v] [%s] inputted and nil was not expected", t.Name(), test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputWif)
		} else if rawTx != nil && test.expectedNil {
			t.Fatalf("%s Failed: [%v] [%v] [%v] [%s] inputted and nil was expected", t.Name(), test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputWif)
		} else if rawTx != nil && rawTx.String() != test.expectedRawTx {
			t.Fatalf("%s Failed: [%v] [%v] [%v] [%s] inputted [%s] expected but failed comparison of scripts, got: %s", t.Name(), test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputWif, test.expectedRawTx, rawTx.String())
		}
	}
}

// TestCreateTxUsingWif will test the method CreateTxUsingWif()
func TestCreateTxUsingWif(t *testing.T) {
	t.Parallel()

	t.Run("valid tx", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 500,
		}

		opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
		opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

		tx, err := CreateTxUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.NoError(t, err)
		assert.NotNil(t, tx)
	})

	t.Run("invalid wif", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 500,
		}

		opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
		opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

		tx, err := CreateTxUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			"",
		)
		assert.Error(t, err)
		assert.Nil(t, tx)
	})
}

// ExampleCreateTxUsingWif example using CreateTxUsingWif()
func ExampleCreateTxUsingWif() {

	// Use a new UTXO
	utxo := &Utxo{
		TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		Vout:         0,
		ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
		Satoshis:     1000,
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

	fmt.Printf("rawTx: %s", rawTx.String())
	// Output:rawTx: 0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100eea3d606bd1627be6459a9de4860919225db74843d2fc7f4e7caa5e01f42c2d0022017978d9c6a0e934955a70e7dda71d68cb614f7dd89eb7b9d560aea761834ddd4412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff03f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000
}

// BenchmarkCreateTxUsingWif benchmarks the method CreateTxUsingWif()
func BenchmarkCreateTxUsingWif(b *testing.B) {
	// Use a new UTXO
	utxo := &Utxo{TxID: "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576", Vout: 0, ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac", Satoshis: 1000}

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

	t.Run("valid tx", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 868,
		}

		opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
		opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

		rawTx, err := CreateTxUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.NoError(t, err)

		// Calculate fee
		assert.Equal(t, uint64(132), CalculateFeeForTx(rawTx, nil, nil))
	})

	t.Run("panic on tx nil", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = CalculateFeeForTx(nil, nil, nil)
		})
	})
}

// TestCalculateFeeForTxVariousTxs will test the method CalculateFeeForTx()
func TestCalculateFeeForTxVariousTxs(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		name              string
		inputHex          string
		inputStandardRate *bt.Fee
		inputDataRate     *bt.Fee
		expectedTxID      string
		expectedSatoshis  uint64
	}{
		{
			"tx-132",
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100e07b7661af4e4b521c012a146b25da2c7b9d606e9ceaae28fa73eb347ef6da6f0220527f0638a89ff11cbe53d5f8c4c2962484a370dcd9463a6330f45d31247c2512412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0364030000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000",
			nil,
			nil,
			"e75fa79ee5fbb589201f769c01835e14ca595b7bbfa0e602050a2a90cf28d129",
			132,
		},
		{
			"tx-141",
			"0200000001203d3a9d8e2ccfe2d6bb1bce6ad8a9c1251a58f9c788737b21e3e19588e89110010000006a47304402201e7fe22f20d02a5cd6978b21bc68aa31eb74530c6b75b47caef19d6f2a95f47802206b20de32c7398fa822397b49389a49224742a9d3c175a851376a25deb6428cb24121023bc787128d6296be6a534b32e90724413179545a0ec2720a7258de330fc54544ffffffff02000000000000000053006a4c4fe5b08fe881aae6b8b8e6888f3130e69c883234e697a5e6b7b1e59cb3e7babfe4b88be994a6e6a087e8b59befbc8ce586a0e5869b425356e5a596e98791e7ad89e4bda0e69da5e68c91e68898efbc817c0c0000000000001976a9146382643e30d2ea7e52eb07eac8767ed219c53e3a88ac00000000",
			nil,
			nil,
			"1affabe9b5adc3a6930a06002a447e834681004a5e6767a649d0371a806e7b1d",
			141,
		},
		{
			"tx-95",
			"0200000001a39b865acf5d100aa35341a63e4ba6dc7da101f3b740cd5372e25b0fb23306c5000000006a47304402207cdcab521641801cd0427501c0264073cfd11f2693bb9a109d00482c643f16b30220798801d1acc9dea23013f17cbd7cc0edec996abc1493f9b9a8ba61723b4ec01d41210250a932cf2543f8f35dba3fce46f53ca821409f4e233d9f8090a165f5ac42a0e8ffffffff014a060000000000001976a914f21ff89bd2259699d229ebc1f9c29ac3c3e0411888ac00000000",
			nil,
			nil,
			"175fc22ffbd76f1cdb7ec2c40474abedbb4bb6080e9c8d22736ffc2d48e85fd2",
			95,
		},
		{
			"tx-410",
			"010000000129cdf21be448b28b0f88ecc9317d9ea1a385a6928d855f55b2d62a3bbc8631cd000000006b483045022100a8b7d34e10e817647f753a9a1380606945bd9c255ebb87c7a00a701903abc6e20220284651051ef520eb5f13901c0f792d7eae2e75b9e85fff36258d9592839227614121029843434e00940e0829196042983e46fe801364b1a586f3d4160556684d63e8d1ffffffff030000000000000000fd48026a2231394878696756345179427633744870515663554551797131707a5a56646f4175744cff545742544348205457455443482054574554534820545745544348205457455443482054574554434820545745544348205457455443482054564554434820545745544348205457454943482054574554434820545745544348205457455453482054574554434820545745544348205457415443482054574554434820545745574348205457455443482054574554434820545745544548205457455443482054574554434820545745545348205457455443482054574554434820545745544f482054574554434820544d45544348205457455443482054574554454820545745544348205768617427732077697468696e205457455443483f203a500a746578742f706c61696e04746578741f7477657463685f7477746578745f313536383139363330383230372e747874017c223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540b7477646174615f6a736f6e046e756c6c0375726c046e756c6c07636f6d6d656e74046e756c6c076d625f757365720434363631047479706504706f73740974696d657374616d700f3638343131393038353536373938330361707006747765746368017c22313550636948473232534e4c514a584d6f53556157566937575371633768436676610d424954434f494e5f45434453412231414b48566959674247626d78693871694a6b4e766f484e654475396d334d665045102323636f6d70757465645f7369672323953f0000000000001976a9147346219a70740418422b19c0ecd696671d3ce17088acdbdf0500000000001976a914b4376ffbd47974340a0a0d940681e168c1d6ae0788ac00000000",
			nil,
			nil,
			"d612aed6c3f12756ea1d3e9a48eef2faf05eaf20ff6031911d8610332f8d3f9a",
			410,
		},
		{
			"tx-209",
			"01000000018558c697a0bb502f9aec70b27c48d3c87f0df28e16f9c8f43a43b1327ea69304010000006b483045022100e0114cf815c0d60fc72e7c8a4c39b03920ac03e62b5c98c68882c1235f2e0ac90220278319c67e79c4397ad989cb9be453acfbbb8e8a9f927c6f4509d9a1da4a17eb412103637720f48f059b854605da3d5a959d204e1e6cbd0d03ee91a38a067bdee90558ffffffff020000000000000000d96a2231394878696756345179427633744870515663554551797131707a5a56646f4175742c436865636b6f757420546f6e6963506f773a2068747470733a2f2f746e6370772e636f2f36313239373839330d746578742f6d61726b646f776e055554462d38017c223150755161374b36324d694b43747373534c4b79316b683536575755374d745552350353455403617070086d6574616c656e73047479706507636f6d6d656e740375726c2268747470733a2f2f6f66666572732e746f6e6963706f772e636f6d2f6f666665727304757365720444756465d3bb1100000000001976a914e0190be1c0ced1e92c26b3c6d5ccbbe853b57e7288ac00000000",
			nil,
			nil,
			"ef3fe744c7be4b81f2881f9c11433bc4905a032a0bf15bcda31f90c28ba1b356",
			209,
		},
		{
			"tx-220-1",
			"010000000190c5208bed05b4e54746bab5a5d4a5e4324e2999c180b7ea7105047f1f16b84a030000006a47304402202ca5d7a2cfb2388babb549b10c47ed20cdadafb845af835c7d5ff04e933ba1c102200a8b7289bbd3c0cc62172afe0006ba374ddd0132a7e4fb4e92ebcff5ce9217db412102812d641ff356c362815f8fc03bd061c2ae57e02f5dc3083f61785c0e5b198039ffffffff040000000000000000ad006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f696403323934106f666665725f73657373696f6e5f696440653638616631393439626131633239326131376431393635343638383234383663653635313830663465363439383235613561363532376634646132303761661f0c0000000000001976a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac31790000000000001976a9147c8ced9ee0f48192822a0148f27b5a1f24aa42d388ac53f74600000000001976a914852f89e9b05d6adfc905842ee0d301947d675df988ac00000000",
			nil,
			nil,
			"8785ca5f11795a38eb1f50f62562cb5e0335b283762fe8a2c7e96d5f7f79bb15",
			220,
		},
		{
			"tx-220-2",
			"01000000014179d984eb4738acddf290d1cbcb6ec0b44d945dc00020a9337738a4b5d4d0c5020000006b483045022100acc8301fb2f9bf70089e03f2a28ff8a384c0525e94a8ee5bb1ccb9b1f044bb91022012b5fb696bb92b2208f901e8bd43264a161cf97f72024cc1809224eee5fd7bab412102ad7b9a78da643560d0ffa411d27070d2468eccfcbb5b39ec30d69724a181944cffffffff045d020000000000001976a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac9e170000000000001976a9141531448d45a7985230d2ad467acc22f635f0d10988ac985b8800000000001976a914cf500ad6b6852a13017d60080e58522716e33b9088ac0000000000000000ac006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f6964023239106f666665725f73657373696f6e5f6964403361303930636439393361346564386239626261303438643030653134666466303931633266616134333735396236366434363136396634613066313463343900000000",
			nil,
			nil,
			"32050b05bf48440bbb3e7c4fa79aa6cb36a69e2921965780334ee168de41f676",
			220,
		},
		{
			"tx-220-3",
			"0100000001462c8567bdd70e5f4d36204c239eb1149e15f778925d5e4d36c3a0e7862bf887020000006b483045022100c55116cad7dd3c6bfba1ac4b696ebf60792032c4189c2c30735cfcad300a1c8f02201db5aa8fc38db419b559ff01a07a6dfaac1d5d7b940e89b398fa92bf387b1899412102ad7b9a78da643560d0ffa411d27070d2468eccfcbb5b39ec30d69724a181944cffffffff045c020000000000001976a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac9b170000000000001976a9141531448d45a7985230d2ad467acc22f635f0d10988acf3258800000000001976a914cf500ad6b6852a13017d60080e58522716e33b9088ac0000000000000000ac006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f6964023239106f666665725f73657373696f6e5f6964403634353530346139663966383533393862323434363339336237393066303265666236343738366661666533326535333737373062653064646437663532323700000000",
			nil,
			nil,
			"4232216b35ebaa9301e386010eaaba05902be7bd4db60e2f588ae9d0e3e583e9",
			220,
		},
		{
			"tx-220-4",
			"01000000014c0b196f630c02005b9c17069deef60dff492bdfd0a54a140aaaaf3769ab3d09020000006a47304402207b9acc2cb029617b625fe2a7fc4553b4b70d9f9ec87afc5aa5289866841cbf3702203047baf62881a098cb24bfa85e7f20c37513a8cbbd27efb4b30ad3519118ed7d4121028ae3c8bf25ae420e3e451be1b40f7edb8ff64e0ce32a2398b1d9374c1d50733cffffffff045c020000000000001976a914806773d8c6f90c544b98755e9a80915488357bbe88ac93170000000000001976a914a16150864dc3041666ab055f154da6c35335b6fb88acc2dd7500000000001976a914d1fc52081c4cf54454eeed6982d6b2977c40196088ac0000000000000000ad006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f696403313532106f666665725f73657373696f6e5f6964406637343736653061343763333863616664663061616164336564643539363831306235616539396130633737396637643035383438346330323933326238393500000000",
			nil,
			nil,
			"7a868c8529e598195c6c2cbf30006cb09fd55e9f3d3d809a38f191a586a3ce82",
			220,
		},
		{
			"tx-1000",
			"010000000190c5208bed05b4e54746bab5a5d4a5e4324e2999c180b7ea7105047f1f16b84a030000006a47304402202ca5d7a2cfb2388babb549b10c47ed20cdadafb845af835c7d5ff04e933ba1c102200a8b7289bbd3c0cc62172afe0006ba374ddd0132a7e4fb4e92ebcff5ce9217db412102812d641ff356c362815f8fc03bd061c2ae57e02f5dc3083f61785c0e5b198039ffffffff040000000000000000ad006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f696403323934106f666665725f73657373696f6e5f696440653638616631393439626131633239326131376431393635343638383234383663653635313830663465363439383235613561363532376634646132303761661f0c0000000000001976a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac31790000000000001976a9147c8ced9ee0f48192822a0148f27b5a1f24aa42d388ac53f74600000000001976a914852f89e9b05d6adfc905842ee0d301947d675df988ac00000000",
			&bt.Fee{FeeType: bt.FeeTypeStandard, MiningFee: bt.FeeUnit{
				Satoshis: 1000,
				Bytes:    1000,
			}, RelayFee: bt.FeeUnit{
				Satoshis: 1000,
				Bytes:    1000,
			}},
			&bt.Fee{FeeType: bt.FeeTypeData, MiningFee: bt.FeeUnit{
				Satoshis: 1000,
				Bytes:    1000,
			}, RelayFee: bt.FeeUnit{
				Satoshis: 1000,
				Bytes:    1000,
			}},
			"8785ca5f11795a38eb1f50f62562cb5e0335b283762fe8a2c7e96d5f7f79bb15",
			441,
		},
		{
			"tx-250",
			"010000000190c5208bed05b4e54746bab5a5d4a5e4324e2999c180b7ea7105047f1f16b84a030000006a47304402202ca5d7a2cfb2388babb549b10c47ed20cdadafb845af835c7d5ff04e933ba1c102200a8b7289bbd3c0cc62172afe0006ba374ddd0132a7e4fb4e92ebcff5ce9217db412102812d641ff356c362815f8fc03bd061c2ae57e02f5dc3083f61785c0e5b198039ffffffff040000000000000000ad006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f696403323934106f666665725f73657373696f6e5f696440653638616631393439626131633239326131376431393635343638383234383663653635313830663465363439383235613561363532376634646132303761661f0c0000000000001976a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac31790000000000001976a9147c8ced9ee0f48192822a0148f27b5a1f24aa42d388ac53f74600000000001976a914852f89e9b05d6adfc905842ee0d301947d675df988ac00000000",
			&bt.Fee{FeeType: bt.FeeTypeStandard, MiningFee: bt.FeeUnit{
				Satoshis: 250,
				Bytes:    1000,
			}, RelayFee: bt.FeeUnit{
				Satoshis: 250,
				Bytes:    1000,
			}},
			&bt.Fee{FeeType: bt.FeeTypeData, MiningFee: bt.FeeUnit{
				Satoshis: 250,
				Bytes:    1000,
			}, RelayFee: bt.FeeUnit{
				Satoshis: 250,
				Bytes:    1000,
			}},
			"8785ca5f11795a38eb1f50f62562cb5e0335b283762fe8a2c7e96d5f7f79bb15",
			109,
		},
		{
			"tx-157-1",
			"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006a47304402200083bb297d53210cf9379b3f47de2eff38e6906e5982fbfeef9bf59778750f3e022046da020811e9a2d1e6db8da103d17598abc194125612be6b108d49cb60cbca95412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0000000000",
			&bt.Fee{FeeType: bt.FeeTypeStandard, MiningFee: bt.FeeUnit{
				Satoshis: 1,
				Bytes:    157,
			}, RelayFee: bt.FeeUnit{
				Satoshis: 1,
				Bytes:    157,
			}},
			&bt.Fee{FeeType: bt.FeeTypeData, MiningFee: bt.FeeUnit{
				Satoshis: 1,
				Bytes:    157,
			}, RelayFee: bt.FeeUnit{
				Satoshis: 1,
				Bytes:    157,
			}},
			"d3350a4ef4b2c72b23e5117979590d768e61f2102337e2ae956d152a80cd37ac",
			1,
		},
	}

	var satoshis uint64
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tx, err := TxFromHex(test.inputHex)
			assert.NoError(t, err)
			assert.NotNil(t, tx)

			satoshis = CalculateFeeForTx(tx, test.inputStandardRate, test.inputDataRate)
			assert.Equal(t, test.expectedSatoshis, satoshis, satoshis)
			assert.Equal(t, test.expectedTxID, tx.TxID())
		})
	}
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

	fmt.Printf("tx id: %s estimated fee: %d satoshis", tx.TxID(), estimatedFee)
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

// TestCreateTxWithChange tests for nil case in CreateTxWithChange()
func TestCreateTxWithChange(t *testing.T) {
	t.Parallel()

	t.Run("basic tx", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 500,
		}

		opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
		opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

		privateKey, err := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
		assert.NoError(t, err)
		assert.NotNil(t, privateKey)

		var rawTx *bt.Tx
		rawTx, err = CreateTxWithChange(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			privateKey,
		)
		assert.NoError(t, err)

		// Test the right fee
		assert.Equal(t, uint64(149), CalculateFeeForTx(rawTx, nil, nil))

		// Test that we got the right amount of change (satoshis)
		for _, out := range rawTx.Outputs {
			if out.LockingScriptHexString() == "76a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac" {
				assert.Equal(t, uint64(351), out.Satoshis)
			}
		}
	})

	t.Run("valid txs", func(t *testing.T) {
		var tests = []struct {
			name               string
			inputUtxos         []*Utxo
			inputAddresses     []*PayToAddress
			inputOpReturns     []OpReturnData
			inputWif           string
			inputChangeAddress string
			inputStandardRate  *bt.Fee
			inputDataRate      *bt.Fee
			expectedRawTx      string
		}{
			{
				"simple tx - 1000/500",
				[]*Utxo{{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     1000,
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
				"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100b95aff403574aba31b1786e5f5ddb3c57356a13e6207b66babb16a6d851d7cfe02200d0f570f619e4c05b5b7213ce673f46549f9a7ee95814f3ec0cc0233fd54c85e412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff03f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac71010000000000001976a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
			},
			{
				"simple tx with small op return",
				[]*Utxo{{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     1000,
				}}, nil,
				[]OpReturnData{{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}},
				"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
				"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
				nil,
				nil,
				"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b48304502210092b9d4a913d7103e8770bdeb45528b6aed02126fafdfc64c728243714eac250002205fa4963a90fd69f6ae4cbf02c2ecc2367d585eb616244fb7c85adf4fef468a21412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0277030000000000001976a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
			},
			{
				"no pay-to, all goes to change",
				[]*Utxo{{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     1000,
				}}, nil,
				nil,
				"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
				"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
				nil,
				nil,
				"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006a4730440220029595a3bc3e94f92b1d08faf298ba938cb7c9824393789a1f3afc5ea5e172a802204d224492ab440c8180b45f39a25e270b0fb0d2fc74d8362db52e3ff960d9dc5d412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0187030000000000001976a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac00000000",
			},
			{
				"fee is removed from pay-to",
				[]*Utxo{{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     1000,
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
				"0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100fdf1c12b1512db7ced357358a333f3abf0f7b2d0b1a62c8e727d07627d702d5502207446779354f63e785bb8a46c9b2fb9a6da8a6d5503ee90ff35bb4c0774d38e62412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0276030000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
			},
		}

		var rawTx *bt.Tx
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				privateKey, err := WifToPrivateKey(test.inputWif)
				assert.NoError(t, err)
				assert.NotNil(t, privateKey)

				rawTx, err = CreateTxWithChange(
					test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputChangeAddress,
					test.inputStandardRate, test.inputDataRate, privateKey,
				)
				assert.NoError(t, err)
				assert.NotNil(t, rawTx)
				assert.Equal(t, test.expectedRawTx, rawTx.String())
			})
		}
	})

	t.Run("invalid txs", func(t *testing.T) {
		var tests = []struct {
			name               string
			inputUtxos         []*Utxo
			inputAddresses     []*PayToAddress
			inputOpReturns     []OpReturnData
			inputWif           string
			inputChangeAddress string
			inputStandardRate  *bt.Fee
			inputDataRate      *bt.Fee
			expectedRawTx      string
		}{
			{
				"tx-2",
				[]*Utxo{{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     1000,
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
			},
			{
				"tx-3",
				nil,
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
			},
			{
				"tx-6",
				[]*Utxo{{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "invalid-script",
					Satoshis:     1000,
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
			},
			{
				"tx-7",
				[]*Utxo{{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     1000,
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
			},
			{
				"tx-7",
				[]*Utxo{{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     1000,
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
			},
			{
				"tx-8",
				[]*Utxo{{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     1000,
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
			},
		}

		var rawTx *bt.Tx
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				privateKey, err := WifToPrivateKey(test.inputWif)
				assert.NoError(t, err)
				assert.NotNil(t, privateKey)

				rawTx, err = CreateTxWithChange(
					test.inputUtxos, test.inputAddresses, test.inputOpReturns, test.inputChangeAddress,
					test.inputStandardRate, test.inputDataRate, privateKey,
				)
				assert.Error(t, err)
				assert.Nil(t, rawTx)
			})
		}
	})
}

// ExampleCreateTxWithChange example using CreateTxWithChange()
func ExampleCreateTxWithChange() {

	// Use a new UTXO
	utxo := &Utxo{
		TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		Vout:         0,
		ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
		Satoshis:     1000,
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
	var rawTx *bt.Tx
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

	fmt.Printf("rawTx: %s", rawTx.String())
	// Output:rawTx: 0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100c1dc4a4c5f26a404ff3618013dc63777d0790d0b0ea3371c67ee3f1bb5126c3e02206be8c841918215337f9b6a6a6040bd058596f2bab5c8b8cb27f849b1474b9e4c412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff04f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac5f010000000000001976a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000
}

// BenchmarkCreateTxWithChange benchmarks the method CreateTxWithChange()
func BenchmarkCreateTxWithChange(b *testing.B) {
	// Use a new UTXO
	utxo := &Utxo{TxID: "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576", Vout: 0, ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac", Satoshis: 1000}

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
	t.Parallel()

	t.Run("valid tx", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 500,
		}

		opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
		opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.NoError(t, err)
		assert.NotNil(t, rawTx)

		// Test the right fee
		assert.Equal(t, uint64(149), CalculateFeeForTx(rawTx, nil, nil))

		// Test that we got the right amount of change (satoshis)
		for _, out := range rawTx.Outputs {
			if out.LockingScriptHexString() == "76a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac" {
				assert.Equal(t, uint64(351), out.Satoshis)
			}
		}
	})

	t.Run("valid tx - multiple inputs - send all", func(t *testing.T) {

		opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
		opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{{
				TxID:         "9e88ca8eec0845e9e864c024bc5e6711e670932c9c7d929f9fccdb2c440ae28e",
				Vout:         0,
				ScriptPubKey: "76a9147824dec00be2c45dad83c9b5e9f5d7ef05ba3cf988ac",
				Satoshis:     5689,
			}, {
				TxID:         "4e25b077d4cbb955b5a215feb53f963cf04688ff1777b9bea097c7ddbdf7ea42",
				Vout:         0,
				ScriptPubKey: "76a9147824dec00be2c45dad83c9b5e9f5d7ef05ba3cf988ac",
				Satoshis:     5689,
			}},
			[]*PayToAddress{{
				Address:  "1DE7mZ9g3zmzWLM729fneui7eypfX8BBfC",
				Satoshis: 5689 + 5689,
			}},
			[]OpReturnData{opReturn1, opReturn2},
			"1BxGFoRPSFgYxoAStEncL6HuELqPkV3JVj",
			nil,
			nil,
			"5JXAjNX7cbiWvmkdnj1EnTKPChauttKAJibXLm8tqWtDhXrRbKz",
		)
		assert.NoError(t, err)
		assert.NotNil(t, rawTx)

		// Test the right fee
		assert.Equal(t, uint64(206), CalculateFeeForTx(rawTx, nil, nil))

		assert.Equal(t, "01000000028ee20a442cdbcc9f9f927d9c2c9370e611675ebc24c064e8e94508ec8eca889e000000006b483045022100b79c96d45cb4bb9f2d1887ede3aced860c3d6d8fec169f710ec8a350d16004ad02205413acf2322aa9dbb5ea2a5166af8704ee354b941f4c04b60486e105c5c7ac344121034aaeabc056f33fd960d1e43fc8a0672723af02f275e54c31381af66a334634caffffffff42eaf7bdddc797a0beb97717ff8846f03c963fb5fe15a2b555b9cbd477b0254e000000006b483045022100f1cd918ebbaa7962d700c975cc793bf1351c4a064bf32685da40ac30194f13ea022062ad3e73b2cc96ed0e911f693d0dd844b6eff955f8d7e8f522594882a076f65c4121034aaeabc056f33fd960d1e43fc8a0672723af02f275e54c31381af66a334634caffffffff03a42b0000000000001976a914861c91132b67aec6d1bc111f13523cced19c9f2188ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000", rawTx.String())
	})

	t.Run("send entire utxo amount", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 1000,
		}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			nil,
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.NoError(t, err)
		assert.NotNil(t, rawTx)

		// Test the right fee
		assert.Equal(t, uint64(95), CalculateFeeForTx(rawTx, nil, nil))

		// Test that we got the right amount of change (satoshis)
		for _, out := range rawTx.Outputs {
			assert.Equal(t, uint64(904), out.Satoshis)
			assert.Equal(t, "76a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac", out.LockingScriptHexString())
		}
	})

	t.Run("send entire utxo amount - multiple pay to addresses", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{
				{
					Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
					Satoshis: 250,
				},
				{
					Address:  "1vr42cDUPWd1vbBSmRTaP4RG3kjXSCCM4",
					Satoshis: 250,
				},
				{
					Address:  "1D59KQRed6Nw32q5cLkvKwcqn3cJccLPx7",
					Satoshis: 250,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 250,
				},
			},
			nil,
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.NoError(t, err)
		assert.NotNil(t, rawTx)

		// Test the right fee
		assert.Equal(t, uint64(146), CalculateFeeForTx(rawTx, nil, nil))

		// Test that we got the right amount of change (satoshis)
		for _, out := range rawTx.Outputs {
			if out.LockingScriptHexString() == "76a9147dca96f8f7f4c4b400c46663b86a669bfb7e73c188ac" {
				assert.Equal(t, uint64(103), out.Satoshis)
			} else {
				assert.Equal(t, uint64(250), out.Satoshis)
			}
		}
	})

	t.Run("send almost entire utxo amount - 5 sat diff", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 995,
		}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			nil,
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.NoError(t, err)
		assert.NotNil(t, rawTx)

		// Test the right fee
		assert.Equal(t, uint64(95), CalculateFeeForTx(rawTx, nil, nil))
		assert.Equal(t, "0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006a47304402206cb30774f7cd99db0a713ce164577ae94bbb3fde02f5e3b5afafe354458f4afc02204caa9804b352172d2b4260f5c99edb2091b6c4b5f9fbce12280f76d6371c9abc412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0188030000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000", rawTx.String())

		// Test that we got the right amount of change (satoshis)
		for _, out := range rawTx.Outputs {
			if out.LockingScriptHexString() == "76a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac" {
				assert.Equal(t, uint64(904), out.Satoshis)
				assert.Equal(t, "76a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac", out.LockingScriptHexString())
			}
		}
	})

	t.Run("send almost entire utxo amount - 1 sat diff", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 999,
		}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			nil,
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.NoError(t, err)
		assert.NotNil(t, rawTx)

		// Test the right fee
		assert.Equal(t, uint64(96), CalculateFeeForTx(rawTx, nil, nil))
		assert.Equal(t, "0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100b7311ea368f243bae53913c8fd17a3a5774703fbb49ac1191d3331a1b2ef031d0220059f392cbe75eaa6d43adf45a6822c86c06c50e819a2352e5596f53049470f77412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0187030000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000", rawTx.String())

		// Test that we got the right amount of change (satoshis)
		for _, out := range rawTx.Outputs {
			assert.Equal(t, uint64(903), out.Satoshis)
			assert.Equal(t, "76a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac", out.LockingScriptHexString())
		}
	})

	t.Run("send more than utxos provided - pay to", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 1001,
		}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			nil,
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.Error(t, err)
		assert.Nil(t, rawTx)
	})

	t.Run("send entire utxo - with op return", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 1000,
		}

		opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
		opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.NoError(t, err)
		assert.NotNil(t, rawTx)

		assert.Equal(t, uint64(132), CalculateFeeForTx(rawTx, nil, nil))
		assert.Equal(t, "0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100b8d6cf4362122a00fcc50195b6ecd75e08abf6dc427f1b75039c6355195a265202200c6eaf5c45c192b1fb1e6384b10d2f526d73edd996f03c43e48fe43b17e9869f412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0363030000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000", rawTx.String())
	})

	t.Run("send entire utxo amount - last pay-to does not cover fee", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{
				{
					Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
					Satoshis: 950,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
			},
			nil,
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.NoError(t, err)
		assert.NotNil(t, rawTx)

		// Test the right fee
		assert.Equal(t, uint64(113), CalculateFeeForTx(rawTx, nil, nil))
		assert.Equal(t, "0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100ed46aac84c052c4b64c42782e2d09d0ebf679ba23c99cc22922c78d1ba185bc102204814d6d3d470e6c9ab46b4270ff98962dabd9df853d3856e73943433739faddd412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0245030000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac32000000000000001976a9147dca96f8f7f4c4b400c46663b86a669bfb7e73c188ac00000000", rawTx.String())

		// Test that we got the right amount of change (satoshis)
		for _, out := range rawTx.Outputs {
			if out.LockingScriptHexString() == "76a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac" {
				assert.Equal(t, uint64(837), out.Satoshis)
			} else {
				assert.Equal(t, uint64(50), out.Satoshis)
			}
		}
	})

	t.Run("send entire utxo - cannot cover fee in any pay-to", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
				{
					Address:  "1CU8AAJoPTvLCph2mnpXarExQ1rKdVZum5",
					Satoshis: 50,
				},
			},
			nil,
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.Error(t, err)
		assert.Nil(t, rawTx)
	})

	t.Run("send entire utxo using data, leaving 1-2 sats", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 50,
		}

		// fee = 998
		str := "043ce70e56f8f37ce5e00f992114594f24b45a49824321f9c9fe309cc3bc4112d546e88cdf95694278aaab92177b0e81a2d4bc4a1f8e18d2c28dae36b1f5a32be86f1c77191ef534d66a9e376472ea12d92e7b4e0f42874310d48189ad9d1c84bb9f3d7a84402565744d5721f67686baa55a3c2ba2865effa0c911e1565fa89d1dfdfa17ff42180435bd16b55076743d23dafc5af9fee20a3cddeb28b0dd86f058b8c49fa0249a50253331ee71c4455292222ec3f33450b95c0de42e68baeec7b5090862c41a19e903c2cc19df83ca3504b11f3e7142b063ed498226a0644d60524796ffbd16cfd50a7ce3dc81aaf4179035c46d0be00d200b1e90ea46b986045f24784a992bb323d7fbd9d694566ba566029f65a7edab65ad71bbffc49530219a5baa945ba59d3ec13a08dab88663694f733bf96e6ca07c5efc5e653bc0b12277812ba6d4f1c818384149b789b15583916a217c1b6bd971f628c70e660680edc8530296b6373f8db1d3bf2fd61c909775545fcc1d127d3e047d6740075398ff791a4eba34301efa928b6a3ace86af39ec4a8a7b88404c764a8038df81230be6e02ceb8808165b0a06e509fe94b8a644a14b5e7ee639f735c0125bbc64a52405dffb23f76e4211ea2fb86e9eadb5d130ee363998c7e1857097cdbe8f9c7b0143970402892ffe333de6da83190f6e6aabc21baacfb0924e3690db3d8ac282232d5be7c51f25409ceadf6a5f1f271844abf1fe5952af06a709ef2e3a722a54accd0291b926e2c29d7997c2e55989fa3d6ac2a72c47bf20b0d7d0d57f0ce358ed6dd6ded640d048108b92d6bac761734a98df9ae669f0de10fcd3f431ba0d7880e773d48c4eda955b66a7a7f6ffd0e7d151a45599de1173885754e8aa2c71aa9451a4cde2956f598d0db8013abc3d41b5acc1a99c6fe90218a565d55c81ca655ac1e196979f6b3aed7a6b643545cefdcfd57f03a24d3ee9e89a2acc2ed9234b628baef03f7c50d9d2e4ddc207a3377a911fd4b170483bd81b62958701fdfeb1dbc4b405d14437b60e25025350452531b90df95a21e9cbdc5093936392b80d68fbfe0066706857a9cfe6eea6944c2a548b4a0bb3560fc8ff0328d686c8deda686edce610eb8aadd06fbb4c8a54d8e9359c0507f7d2390b5139874198f6d66251123efae66606a99fad1c8676305c6f3590be8f9243942766714b59e66e4b8acbcd0c93e988b888c8243af505e52e01397290b7db6b773c65949bfe008b15e582"

		opReturn1 := OpReturnData{[]byte(str)}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1},
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.NoError(t, err)
		assert.NotNil(t, rawTx)
		assert.Equal(t, "0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100f1ab55c2f9fe4afb436ce18d097176eaca8a3d37a38966f606903fa1711f2bbe02202cedc7d854318c15e91fbdd779d427ef78e643e2dd9efa81ba5e245fad385ae9412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0201000000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac0000000000000000fd0107006a4dfc0630343363653730653536663866333763653565303066393932313134353934663234623435613439383234333231663963396665333039636333626334313132643534366538386364663935363934323738616161623932313737623065383161326434626334613166386531386432633238646165333662316635613332626538366631633737313931656635333464363661396533373634373265613132643932653762346530663432383734333130643438313839616439643163383462623966336437613834343032353635373434643537323166363736383662616135356133633262613238363565666661306339313165313536356661383964316466646661313766663432313830343335626431366235353037363734336432336461666335616639666565323061336364646562323862306464383666303538623863343966613032343961353032353333333165653731633434353532393232323265633366333334353062393563306465343265363862616565633762353039303836326334316131396539303363326363313964663833636133353034623131663365373134326230363365643439383232366130363434643630353234373936666662643136636664353061376365336463383161616634313739303335633436643062653030643230306231653930656134366239383630343566323437383461393932626233323364376662643964363934353636626135363630323966363561376564616236356164373162626666633439353330323139613562616139343562613539643365633133613038646162383836363336393466373333626639366536636130376335656663356536353362633062313232373738313262613664346631633831383338343134396237383962313535383339313661323137633162366264393731663632386337306536363036383065646338353330323936623633373366386462316433626632666436316339303937373535343566636331643132376433653034376436373430303735333938666637393161346562613334333031656661393238623661336163653836616633396563346138613762383834303463373634613830333864663831323330626536653032636562383830383136356230613036653530396665393462386136343461313462356537656536333966373335633031323562626336346135323430356466666232336637366534323131656132666238366539656164623564313330656533363339393863376531383537303937636462653866396337623031343339373034303238393266666533333364653664613833313930663665366161626332316261616366623039323465333639306462336438616332383232333264356265376335316632353430396365616466366135663166323731383434616266316665353935326166303661373039656632653361373232613534616363643032393162393236653263323964373939376332653535393839666133643661633261373263343762663230623064376430643537663063653335386564366464366465643634306430343831303862393264366261633736313733346139386466396165363639663064653130666364336634333162613064373838306537373364343863346564613935356236366137613766366666643065376431353161343535393964653131373338383537353465386161326337316161393435316134636465323935366635393864306462383031336162633364343162356163633161393963366665393032313861353635643535633831636136353561633165313936393739663662336165643761366236343335343563656664636664353766303361323464336565396538396132616363326564393233346236323862616566303366376335306439643265346464633230376133333737613931316664346231373034383362643831623632393538373031666466656231646263346234303564313434333762363065323530323533353034353235333162393064663935613231653963626463353039333933363339326238306436386662666530303636373036383537613963666536656561363934346332613534386234613062623335363066633866663033323864363836633864656461363836656463653631306562386161646430366662623463386135346438653933353963303530376637643233393062353133393837343139386636643636323531313233656661653636363036613939666164316338363736333035633666333539306265386639323433393432373636373134623539653636653462386163626364306339336539383862383838633832343361663530356535326530313339373239306237646236623737336336353934396266653030386231356535383200000000", rawTx.String())
	})

	t.Run("send entire utxo using data, too much data", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 50,
		}

		// fee = 1000+
		str := "000000000043ce70e56f8f37ce5e00f992114594f24b45a49824321f9c9fe309cc3bc4112d546e88cdf95694278aaab92177b0e81a2d4bc4a1f8e18d2c28dae36b1f5a32be86f1c77191ef534d66a9e376472ea12d92e7b4e0f42874310d48189ad9d1c84bb9f3d7a84402565744d5721f67686baa55a3c2ba2865effa0c911e1565fa89d1dfdfa17ff42180435bd16b55076743d23dafc5af9fee20a3cddeb28b0dd86f058b8c49fa0249a50253331ee71c4455292222ec3f33450b95c0de42e68baeec7b5090862c41a19e903c2cc19df83ca3504b11f3e7142b063ed498226a0644d60524796ffbd16cfd50a7ce3dc81aaf4179035c46d0be00d200b1e90ea46b986045f24784a992bb323d7fbd9d694566ba566029f65a7edab65ad71bbffc49530219a5baa945ba59d3ec13a08dab88663694f733bf96e6ca07c5efc5e653bc0b12277812ba6d4f1c818384149b789b15583916a217c1b6bd971f628c70e660680edc8530296b6373f8db1d3bf2fd61c909775545fcc1d127d3e047d6740075398ff791a4eba34301efa928b6a3ace86af39ec4a8a7b88404c764a8038df81230be6e02ceb8808165b0a06e509fe94b8a644a14b5e7ee639f735c0125bbc64a52405dffb23f76e4211ea2fb86e9eadb5d130ee363998c7e1857097cdbe8f9c7b0143970402892ffe333de6da83190f6e6aabc21baacfb0924e3690db3d8ac282232d5be7c51f25409ceadf6a5f1f271844abf1fe5952af06a709ef2e3a722a54accd0291b926e2c29d7997c2e55989fa3d6ac2a72c47bf20b0d7d0d57f0ce358ed6dd6ded640d048108b92d6bac761734a98df9ae669f0de10fcd3f431ba0d7880e773d48c4eda955b66a7a7f6ffd0e7d151a45599de1173885754e8aa2c71aa9451a4cde2956f598d0db8013abc3d41b5acc1a99c6fe90218a565d55c81ca655ac1e196979f6b3aed7a6b643545cefdcfd57f03a24d3ee9e89a2acc2ed9234b628baef03f7c50d9d2e4ddc207a3377a911fd4b170483bd81b62958701fdfeb1dbc4b405d14437b60e25025350452531b90df95a21e9cbdc5093936392b80d68fbfe0066706857a9cfe6eea6944c2a548b4a0bb3560fc8ff0328d686c8deda686edce610eb8aadd06fbb4c8a54d8e9359c0507f7d2390b5139874198f6d66251123efae66606a99fad1c8676305c6f3590be8f9243942766714b59e66e4b8acbcd0c93e988b888c8243af505e52e01397290b7db6b773c65949bfe008b15e582"

		opReturn1 := OpReturnData{[]byte(str)}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1},
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu",
		)
		assert.Error(t, err)
		assert.Nil(t, rawTx)
	})

	t.Run("invalid wif", func(t *testing.T) {
		utxo := &Utxo{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     1000,
		}

		payTo := &PayToAddress{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 500,
		}

		opReturn1 := OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}
		opReturn2 := OpReturnData{[]byte("prefix2"), []byte("more example data")}

		rawTx, err := CreateTxWithChangeUsingWif(
			[]*Utxo{utxo},
			[]*PayToAddress{payTo},
			[]OpReturnData{opReturn1, opReturn2},
			"1KQG5AY9GrPt3b5xrFqVh2C3YEhzSdu4kc",
			nil,
			nil,
			"",
		)
		assert.Error(t, err)
		assert.Nil(t, rawTx)
	})
}

// ExampleCreateTxWithChangeUsingWif example using CreateTxWithChangeUsingWif()
func ExampleCreateTxWithChangeUsingWif() {

	// Use a new UTXO
	utxo := &Utxo{
		TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		Vout:         0,
		ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
		Satoshis:     1000,
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

	fmt.Printf("rawTx: %s", rawTx.String())
	// Output:rawTx: 0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100c1dc4a4c5f26a404ff3618013dc63777d0790d0b0ea3371c67ee3f1bb5126c3e02206be8c841918215337f9b6a6a6040bd058596f2bab5c8b8cb27f849b1474b9e4c412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff04f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac5f010000000000001976a914c9d8699bdea34b131e737447b50a8b1af0b040bf88ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000
}

// BenchmarkCreateTxWithChangeUsingWif benchmarks the method CreateTxWithChangeUsingWif()
func BenchmarkCreateTxWithChangeUsingWif(b *testing.B) {
	// Use a new UTXO
	utxo := &Utxo{TxID: "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576", Vout: 0, ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac", Satoshis: 1000}

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
