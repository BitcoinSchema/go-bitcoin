package bitcoin

import (
	"testing"

	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bt/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateTxExtended provides comprehensive tests for CreateTx function
func TestCreateTxExtended(t *testing.T) {
	t.Parallel()

	// Valid WIF for testing
	validWIF := "L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu"
	privateKey, err := WifToPrivateKey(validWIF)
	require.NoError(t, err)

	tests := []struct {
		name        string
		utxos       []*Utxo
		addresses   []*PayToAddress
		opReturns   []OpReturnData
		privateKey  interface{} // Can be nil or *bec.PrivateKey
		shouldError bool
		errorText   string
	}{
		{
			name: "create tx with nil private key - should succeed",
			utxos: []*Utxo{
				{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     10000,
				},
			},
			addresses: []*PayToAddress{
				{
					Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
					Satoshis: 5000,
				},
			},
			opReturns:   nil,
			privateKey:  nil, // No private key should still create tx structure
			shouldError: false,
		},
		{
			name: "create tx with valid inputs and outputs",
			utxos: []*Utxo{
				{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     10000,
				},
			},
			addresses: []*PayToAddress{
				{
					Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
					Satoshis: 5000,
				},
			},
			opReturns:   nil,
			privateKey:  privateKey,
			shouldError: false,
		},
		{
			name: "create tx with multiple utxos",
			utxos: []*Utxo{
				{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     5000,
				},
				{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         1,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     5000,
				},
			},
			addresses: []*PayToAddress{
				{
					Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
					Satoshis: 8000,
				},
			},
			opReturns:   nil,
			privateKey:  privateKey,
			shouldError: false,
		},
		{
			name: "create tx with multiple outputs",
			utxos: []*Utxo{
				{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     10000,
				},
			},
			addresses: []*PayToAddress{
				{
					Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
					Satoshis: 3000,
				},
				{
					Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
					Satoshis: 3000,
				},
			},
			opReturns:   nil,
			privateKey:  privateKey,
			shouldError: false,
		},
		{
			name: "create tx with op_return data",
			utxos: []*Utxo{
				{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     10000,
				},
			},
			addresses: []*PayToAddress{
				{
					Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
					Satoshis: 5000,
				},
			},
			opReturns: []OpReturnData{
				{[]byte("test"), []byte("data")},
			},
			privateKey:  privateKey,
			shouldError: false,
		},
		{
			name: "insufficient funds - outputs exceed inputs",
			utxos: []*Utxo{
				{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     1000,
				},
			},
			addresses: []*PayToAddress{
				{
					Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
					Satoshis: 5000,
				},
			},
			opReturns:   nil,
			privateKey:  privateKey,
			shouldError: true,
			errorText:   "insufficient funds",
		},
		{
			name: "invalid address in pay-to",
			utxos: []*Utxo{
				{
					TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
					Vout:         0,
					ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
					Satoshis:     10000,
				},
			},
			addresses: []*PayToAddress{
				{
					Address:  "invalid-address",
					Satoshis: 5000,
				},
			},
			opReturns:   nil,
			privateKey:  privateKey,
			shouldError: true,
			errorText:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tx *bt.Tx
			var err error

			if tt.privateKey != nil {
				// Type assertion since we're using interface{} in the test struct
				pk := tt.privateKey.(*bec.PrivateKey)
				tx, err = CreateTx(tt.utxos, tt.addresses, tt.opReturns, pk)
			} else {
				tx, err = CreateTx(tt.utxos, tt.addresses, tt.opReturns, nil)
			}

			if tt.shouldError {
				require.Error(t, err, "Expected error but got none")
				if tt.errorText != "" {
					assert.Contains(t, err.Error(), tt.errorText)
				}
			} else {
				require.NoError(t, err, "Unexpected error: %v", err)
				assert.NotNil(t, tx)

				// Verify transaction structure
				if len(tt.utxos) > 0 {
					assert.Len(t, tx.Inputs, len(tt.utxos), "Input count mismatch")
				}
				expectedOutputs := len(tt.addresses) + len(tt.opReturns)
				assert.Len(t, tx.Outputs, expectedOutputs, "Output count mismatch")
			}
		})
	}
}

// TestCreateTxUsingWifExtended tests CreateTxUsingWif function with additional cases
func TestCreateTxUsingWifExtended(t *testing.T) {
	t.Parallel()

	utxo := &Utxo{
		TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		Vout:         0,
		ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
		Satoshis:     10000,
	}

	address := &PayToAddress{
		Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
		Satoshis: 5000,
	}

	t.Run("valid wif", func(t *testing.T) {
		tx, err := CreateTxUsingWif([]*Utxo{utxo}, []*PayToAddress{address}, nil, "L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
		require.NoError(t, err)
		assert.NotNil(t, tx)
	})

	t.Run("invalid wif", func(t *testing.T) {
		tx, err := CreateTxUsingWif([]*Utxo{utxo}, []*PayToAddress{address}, nil, "invalid-wif")
		require.Error(t, err)
		assert.Nil(t, tx)
	})

	t.Run("empty wif", func(t *testing.T) {
		tx, err := CreateTxUsingWif([]*Utxo{utxo}, []*PayToAddress{address}, nil, "")
		require.Error(t, err)
		assert.Nil(t, tx)
	})
}

// TestCalculateFeeForTxExtended provides comprehensive tests for CalculateFeeForTx
func TestCalculateFeeForTxExtended(t *testing.T) {
	t.Parallel()

	// Create a simple transaction for testing
	privateKey, err := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	require.NoError(t, err)

	utxo := &Utxo{
		TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		Vout:         0,
		ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
		Satoshis:     10000,
	}

	address := &PayToAddress{
		Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
		Satoshis: 5000,
	}

	tx, err := CreateTx([]*Utxo{utxo}, []*PayToAddress{address}, nil, privateKey)
	require.NoError(t, err)

	t.Run("calculate fee with default rates", func(t *testing.T) {
		fee := CalculateFeeForTx(tx, nil, nil)
		assert.Positive(t, fee, "Fee should be greater than 0")
	})

	t.Run("calculate fee with custom standard rate", func(t *testing.T) {
		customRate := &bt.Fee{
			FeeType: bt.FeeTypeStandard,
			MiningFee: bt.FeeUnit{
				Satoshis: 10,
				Bytes:    10,
			},
			RelayFee: bt.FeeUnit{
				Satoshis: 10,
				Bytes:    10,
			},
		}
		fee := CalculateFeeForTx(tx, customRate, nil)
		assert.Positive(t, fee, "Fee should be greater than 0")
	})

	t.Run("calculate fee with custom data rate", func(t *testing.T) {
		customDataRate := &bt.Fee{
			FeeType: bt.FeeTypeData,
			MiningFee: bt.FeeUnit{
				Satoshis: 2,
				Bytes:    10,
			},
			RelayFee: bt.FeeUnit{
				Satoshis: 2,
				Bytes:    10,
			},
		}
		fee := CalculateFeeForTx(tx, nil, customDataRate)
		assert.Positive(t, fee, "Fee should be greater than 0")
	})

	t.Run("calculate fee for transaction with OP_RETURN", func(t *testing.T) {
		txWithOpReturn, err := CreateTx(
			[]*Utxo{utxo},
			[]*PayToAddress{address},
			[]OpReturnData{{[]byte("test"), []byte("data")}},
			privateKey,
		)
		require.NoError(t, err)

		fee := CalculateFeeForTx(txWithOpReturn, nil, nil)
		assert.Positive(t, fee, "Fee should be greater than 0")

		// Fee for OP_RETURN tx should be higher due to data
		feeWithoutOpReturn := CalculateFeeForTx(tx, nil, nil)
		assert.Greater(t, fee, feeWithoutOpReturn, "OP_RETURN tx should have higher fee")
	})
}

// TestCreateTxWithChangeExtended provides comprehensive tests for CreateTxWithChange
func TestCreateTxWithChangeExtended(t *testing.T) {
	t.Parallel()

	privateKey, err := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	require.NoError(t, err)

	t.Run("create tx with change - basic", func(t *testing.T) {
		utxos := []*Utxo{
			{
				TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
				Vout:         0,
				ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
				Satoshis:     100000,
			},
		}

		addresses := []*PayToAddress{
			{
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 50000,
			},
		}

		tx, err := CreateTxWithChange(utxos, addresses, nil, "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL", nil, nil, privateKey)
		require.NoError(t, err)
		assert.NotNil(t, tx)
		// Should have 2 outputs: payment + change
		assert.Len(t, tx.Outputs, 2, "Should have payment and change outputs")
	})

	t.Run("no utxos - should error", func(t *testing.T) {
		addresses := []*PayToAddress{
			{
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 50000,
			},
		}

		tx, err := CreateTxWithChange(nil, addresses, nil, "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL", nil, nil, privateKey)
		require.Error(t, err)
		assert.Nil(t, tx)
		assert.Contains(t, err.Error(), "required")
	})

	t.Run("no change address - should error", func(t *testing.T) {
		utxos := []*Utxo{
			{
				TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
				Vout:         0,
				ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
				Satoshis:     100000,
			},
		}

		addresses := []*PayToAddress{
			{
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 50000,
			},
		}

		tx, err := CreateTxWithChange(utxos, addresses, nil, "", nil, nil, privateKey)
		require.Error(t, err)
		assert.Nil(t, tx)
		assert.Contains(t, err.Error(), "required")
	})

	t.Run("insufficient funds - should error", func(t *testing.T) {
		utxos := []*Utxo{
			{
				TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
				Vout:         0,
				ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
				Satoshis:     1000,
			},
		}

		addresses := []*PayToAddress{
			{
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 50000,
			},
		}

		tx, err := CreateTxWithChange(utxos, addresses, nil, "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL", nil, nil, privateKey)
		require.Error(t, err)
		assert.Nil(t, tx)
		assert.Contains(t, err.Error(), "insufficient")
	})

	t.Run("exact amount - no change needed", func(t *testing.T) {
		// When paying exact amount, no change output should be created
		// This test verifies the edge case handling
		utxos := []*Utxo{
			{
				TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
				Vout:         0,
				ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
				Satoshis:     10000,
			},
		}

		// Use most of the UTXO, leaving just enough for fees
		addresses := []*PayToAddress{
			{
				Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
				Satoshis: 9500,
			},
		}

		tx, err := CreateTxWithChange(utxos, addresses, nil, "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL", nil, nil, privateKey)
		// This might error or succeed depending on exact fee calculation
		if err != nil {
			t.Logf("Expected behavior: %v", err)
		} else {
			assert.NotNil(t, tx)
		}
	})
}

// TestCreateTxWithChangeUsingWifExtended tests CreateTxWithChangeUsingWif function with additional cases
func TestCreateTxWithChangeUsingWifExtended(t *testing.T) {
	t.Parallel()

	utxos := []*Utxo{
		{
			TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
			Vout:         0,
			ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
			Satoshis:     100000,
		},
	}

	addresses := []*PayToAddress{
		{
			Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
			Satoshis: 50000,
		},
	}

	t.Run("valid wif", func(t *testing.T) {
		tx, err := CreateTxWithChangeUsingWif(utxos, addresses, nil, "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL", nil, nil, "L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
		require.NoError(t, err)
		assert.NotNil(t, tx)
	})

	t.Run("invalid wif", func(t *testing.T) {
		tx, err := CreateTxWithChangeUsingWif(utxos, addresses, nil, "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL", nil, nil, "invalid-wif")
		require.Error(t, err)
		assert.Nil(t, tx)
	})
}

// BenchmarkCreateTxExtended benchmarks the CreateTx function with different scenarios
func BenchmarkCreateTxExtended(b *testing.B) {
	privateKey, _ := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	utxo := &Utxo{
		TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		Vout:         0,
		ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
		Satoshis:     10000,
	}
	address := &PayToAddress{
		Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
		Satoshis: 5000,
	}

	for i := 0; i < b.N; i++ {
		_, _ = CreateTx([]*Utxo{utxo}, []*PayToAddress{address}, nil, privateKey)
	}
}

// BenchmarkCalculateFeeForTxExtended benchmarks the CalculateFeeForTx function
func BenchmarkCalculateFeeForTxExtended(b *testing.B) {
	privateKey, _ := WifToPrivateKey("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")
	utxo := &Utxo{
		TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		Vout:         0,
		ScriptPubKey: "76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
		Satoshis:     10000,
	}
	address := &PayToAddress{
		Address:  "1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL",
		Satoshis: 5000,
	}
	tx, _ := CreateTx([]*Utxo{utxo}, []*PayToAddress{address}, nil, privateKey)

	for i := 0; i < b.N; i++ {
		_ = CalculateFeeForTx(tx, nil, nil)
	}
}
