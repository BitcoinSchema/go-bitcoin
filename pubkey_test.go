package bitcoin

import (
	"fmt"
	"testing"
)

// TestPubKeyFromPrivateKey will test the method PubKeyFromPrivateKey()
func TestPubKeyFromPrivateKey(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputKey       string
		expectedPubKey string
		expectedError  bool
	}{
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", "031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f", false},
		{"0", "", true},
		{"", "", true},
	}

	// Run tests
	for _, test := range tests {
		if pubKey, err := PubKeyFromPrivateKey(test.inputKey); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.inputKey, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.inputKey)
		} else if pubKey != test.expectedPubKey {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, but got: %s", t.Name(), test.inputKey, test.expectedPubKey, pubKey)
		}
	}
}

// ExamplePubKeyFromPrivateKey example using PubKeyFromPrivateKey()
func ExamplePubKeyFromPrivateKey() {
	pubKey, err := PubKeyFromPrivateKey("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("pubkey generated: %s", pubKey)
	// Output:pubkey generated: 031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f
}

// BenchmarkPubKeyFromPrivateKey benchmarks the method PubKeyFromPrivateKey()
func BenchmarkPubKeyFromPrivateKey(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for i := 0; i < b.N; i++ {
		_, _ = PubKeyFromPrivateKey(key)
	}
}