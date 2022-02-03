package bitcoin

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/libsv/go-bk/bec"
	"github.com/stretchr/testify/assert"
)

// TestPubKeyFromPrivateKeyString will test the method PubKeyFromPrivateKeyString()
func TestPubKeyFromPrivateKeyString(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		inputKey       string
		expectedPubKey string
		expectedError  bool
	}{
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", "031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f", false},
		{"0", "", true},
		{"", "", true},
	}

	for _, test := range tests {
		if pubKey, err := PubKeyFromPrivateKeyString(test.inputKey, true); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.inputKey, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error was expected", t.Name(), test.inputKey)
		} else if pubKey != test.expectedPubKey {
			t.Fatalf("%s Failed: [%s] inputted and [%s] expected, but got: %s", t.Name(), test.inputKey, test.expectedPubKey, pubKey)
		}
	}
}

// ExamplePubKeyFromPrivateKeyString example using PubKeyFromPrivateKeyString()
func ExamplePubKeyFromPrivateKeyString() {
	pubKey, err := PubKeyFromPrivateKeyString("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", true)
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
	for i := 0; i < b.N; i++ {
		_, _ = PubKeyFromPrivateKeyString(key, true)
	}
}

// TestPubKeyFromPrivateKey will test the method PubKeyFromPrivateKey()
func TestPubKeyFromPrivateKey(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	assert.NoError(t, err)
	assert.NotNil(t, priv)

	var tests = []struct {
		inputKey       *bec.PrivateKey
		expectedPubKey string
		expectedError  bool
	}{
		{priv, "031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f", false},
	}

	for _, test := range tests {
		if pubKey := PubKeyFromPrivateKey(test.inputKey, true); pubKey != test.expectedPubKey {
			t.Fatalf("%s Failed: [%v] inputted and [%s] expected, but got: %s", t.Name(), test.inputKey, test.expectedPubKey, pubKey)
		}
	}
}

// TestPubKeyFromPrivateKeyPanic tests for nil case in PubKeyFromPrivateKey()
func TestPubKeyFromPrivateKeyPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		pubKey := PubKeyFromPrivateKey(nil, true)
		assert.NotEqual(t, 0, len(pubKey))
	})
}

// ExamplePubKeyFromPrivateKey example using PubKeyFromPrivateKey()
func ExamplePubKeyFromPrivateKey() {
	privateKey, err := PrivateKeyFromString("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
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
	for i := 0; i < b.N; i++ {
		_ = PubKeyFromPrivateKey(key, true)
	}
}

// TestPubKeyFromString will test the method PubKeyFromString()
func TestPubKeyFromString(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		inputKey       string
		expectedPubKey string
		expectedNil    bool
		expectedError  bool
	}{
		{"", "", true, true},
		{"0", "", true, true},
		{"00000", "", true, true},
		{"031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f", "031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f", false, false},
	}

	for _, test := range tests {
		if pubKey, err := PubKeyFromString(test.inputKey); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.inputKey, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error was expected", t.Name(), test.inputKey)
		} else if pubKey != nil && test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and nil was expected", t.Name(), test.inputKey)
		} else if pubKey == nil && !test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and nil was NOT expected", t.Name(), test.inputKey)
		} else if pubKey != nil && hex.EncodeToString(pubKey.SerialiseCompressed()) != test.expectedPubKey {
			t.Fatalf("%s Failed: [%s] inputted and [%s] expected, but got: %s", t.Name(), test.inputKey, test.expectedPubKey, hex.EncodeToString(pubKey.SerialiseCompressed()))
		}
	}
}

// ExamplePubKeyFromString example using PubKeyFromString()
func ExamplePubKeyFromString() {
	pubKey, err := PubKeyFromString("031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("pubkey from string: %s", hex.EncodeToString(pubKey.SerialiseCompressed()))
	// Output:pubkey from string: 031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f
}

// BenchmarkPubKeyFromString benchmarks the method PubKeyFromString()
func BenchmarkPubKeyFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = PubKeyFromString("031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f")
	}
}
