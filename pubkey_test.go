package bitcoin

import (
	"encoding/hex"
	"fmt"
	"testing"

	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPubKeyFromPrivateKeyString will test the method PubKeyFromPrivateKeyString()
func TestPubKeyFromPrivateKeyString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputKey       string
		expectedPubKey string
		compressed     bool
		expectedError  bool
	}{
		{"compressed", testPrivateKeyHex, testPubKeyCompressed, true, false},
		{"uncompressed", testPrivateKeyHex, testPubKeyUncompressed, false, false},
		{caseSingleZero, "0", "", true, true},
		{caseEmpty, "", "", true, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			pubKey, err := PubKeyFromPrivateKeyString(test.inputKey, test.compressed)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.expectedPubKey, pubKey)
		})
	}
}

// ExamplePubKeyFromPrivateKeyString example using PubKeyFromPrivateKeyString()
func ExamplePubKeyFromPrivateKeyString() {
	pubKey, err := PubKeyFromPrivateKeyString(testPrivateKeyHex, true)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("pubkey generated: %s", pubKey)
	// Output:pubkey generated: 031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f
}

// BenchmarkPubKeyFromPrivateKeyString benchmarks the method PubKeyFromPrivateKeyString()
func BenchmarkPubKeyFromPrivateKeyString(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for b.Loop() {
		_, _ = PubKeyFromPrivateKeyString(key, true)
	}
}

// TestPubKeyFromPrivateKey will test the method PubKeyFromPrivateKey()
func TestPubKeyFromPrivateKey(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(testPrivateKeyHex)
	require.NoError(t, err)
	assert.NotNil(t, priv)

	tests := []struct {
		name           string
		inputKey       *ec.PrivateKey
		expectedPubKey string
		expectedError  bool
	}{
		{"compressed", priv, testPubKeyCompressed, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			pubKey := PubKeyFromPrivateKey(test.inputKey, true)
			assert.Equal(t, test.expectedPubKey, pubKey)
		})
	}
}

// TestPubKeyFromPrivateKeyPanic tests for nil case in PubKeyFromPrivateKey()
func TestPubKeyFromPrivateKeyPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		pubKey := PubKeyFromPrivateKey(nil, true)
		assert.NotEmpty(t, pubKey)
	})
}

// ExamplePubKeyFromPrivateKey example using PubKeyFromPrivateKey()
func ExamplePubKeyFromPrivateKey() {
	privateKey, err := PrivateKeyFromString(testPrivateKeyHex)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	pubKey := PubKeyFromPrivateKey(privateKey, true)
	fmt.Printf("pubkey generated: %s", pubKey)
	// Output:pubkey generated: 031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f
}

// BenchmarkPubKeyFromPrivateKey benchmarks the method PubKeyFromPrivateKey()
func BenchmarkPubKeyFromPrivateKey(b *testing.B) {
	key, _ := CreatePrivateKey()
	for b.Loop() {
		_ = PubKeyFromPrivateKey(key, true)
	}
}

// TestPubKeyFromString will test the method PubKeyFromString()
func TestPubKeyFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputKey       string
		expectedPubKey string
		expectedNil    bool
		expectedError  bool
	}{
		{caseEmpty, "", "", true, true},
		{caseSingleZero, "0", "", true, true},
		{"short zeros", "00000", "", true, true},
		{"valid compressed", testPubKeyCompressed, testPubKeyCompressed, false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			pubKey, err := PubKeyFromString(test.inputKey)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, pubKey)
			} else {
				require.NotNil(t, pubKey)
				assert.Equal(t, test.expectedPubKey, hex.EncodeToString(pubKey.Compressed()))
			}
		})
	}
}

// ExamplePubKeyFromString example using PubKeyFromString()
func ExamplePubKeyFromString() {
	pubKey, err := PubKeyFromString(testPubKeyCompressed)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("pubkey from string: %s", hex.EncodeToString(pubKey.Compressed()))
	// Output:pubkey from string: 031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f
}

// BenchmarkPubKeyFromString benchmarks the method PubKeyFromString()
func BenchmarkPubKeyFromString(b *testing.B) {
	for b.Loop() {
		_, _ = PubKeyFromString(testPubKeyCompressed)
	}
}
