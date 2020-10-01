package bitcoin

import (
	"fmt"
	"testing"
)

// TestGenerateHDKey will test the method GenerateHDKey()
func TestGenerateHDKey(t *testing.T) {

	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputSeed     uint8
		expectedNil   bool
		expectedError bool
	}{
		{0, false, false},
		{1, true, true},
		{15, true, true},
		{65, true, true},
		{RecommendedSeedLength, false, false},
		{SecureSeedLength, false, false},
	}

	// Run tests
	for _, test := range tests {
		if hdKey, err := GenerateHDKey(test.inputSeed); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%d] inputted and error not expected but got: %s", t.Name(), test.inputSeed, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%d] inputted and error was expected", t.Name(), test.inputSeed)
		} else if hdKey == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%d] inputted and was nil but not expected", t.Name(), test.inputSeed)
		} else if hdKey != nil && test.expectedNil {
			t.Errorf("%s Failed: [%d] inputted and was NOT nil but expected to be nil", t.Name(), test.inputSeed)
		}
	}
}

// ExampleGenerateHDKey example using GenerateHDKey()
func ExampleGenerateHDKey() {
	hdKey, err := GenerateHDKey(SecureSeedLength)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	// Cannot show the private/public key since they change each time
	fmt.Printf("created HD key successfully! (length: %d)", len(hdKey.String()))

	// Output:created HD key successfully! (length: 111)
}

// BenchmarkGenerateHDKey benchmarks the method GenerateHDKey()
func BenchmarkGenerateHDKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateHDKey(RecommendedSeedLength)
	}
}

// BenchmarkGenerateHDKeySecure benchmarks the method GenerateHDKey()
func BenchmarkGenerateHDKeySecure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateHDKey(SecureSeedLength)
	}
}

// TestGenerateHDKeyPair will test the method GenerateHDKeyPair()
func TestGenerateHDKeyPair(t *testing.T) {

	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputSeed     uint8
		expectedError bool
	}{
		{0, false},
		{1, true},
		{15, true},
		{65, true},
		{RecommendedSeedLength, false},
		{SecureSeedLength, false},
	}

	// Run tests
	for _, test := range tests {
		if privateKey, publicKey, err := GenerateHDKeyPair(test.inputSeed); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%d] inputted and error not expected but got: %s", t.Name(), test.inputSeed, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%d] inputted and error was expected", t.Name(), test.inputSeed)
		} else if err == nil && len(privateKey) == 0 {
			t.Errorf("%s Failed: [%d] inputted and private key was empty", t.Name(), test.inputSeed)
		} else if err == nil && len(publicKey) == 0 {
			t.Errorf("%s Failed: [%d] inputted and pubic key was empty", t.Name(), test.inputSeed)
		}
	}
}

// ExampleGenerateHDKeyPair example using GenerateHDKeyPair()
func ExampleGenerateHDKeyPair() {
	xPrivateKey, xPublicKey, err := GenerateHDKeyPair(SecureSeedLength)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	// Cannot show the private/public key since they change each time
	fmt.Printf("created HD key successfully! (xPrivateKey length: %d) (xPublicKey length: %d)", len(xPrivateKey), len(xPublicKey))

	// Output:created HD key successfully! (xPrivateKey length: 111) (xPublicKey length: 111)
}

// BenchmarkGenerateHDKeyPair benchmarks the method GenerateHDKeyPair()
func BenchmarkGenerateHDKeyPair(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = GenerateHDKeyPair(RecommendedSeedLength)
	}
}

// BenchmarkGenerateHDKeyPairSecure benchmarks the method GenerateHDKeyPair()
func BenchmarkGenerateHDKeyPairSecure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = GenerateHDKeyPair(SecureSeedLength)
	}
}
