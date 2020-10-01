package bitcoin

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvutil/hdkeychain"
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

// TestGetPrivateKeyByPath will test the method GetPrivateKeyByPath()
func TestGetPrivateKeyByPath(t *testing.T) {

	t.Parallel()

	// Generate a valid key
	validKey, err := GenerateHDKey(RecommendedSeedLength)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Max depth key
	var maxKey *hdkeychain.ExtendedKey
	maxKey, err = GetHDKeyByPath(validKey, 1<<9, 1<<9)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Test depth limit
	for i := 0; i < 1<<8-1; i++ {
		maxKey, err = GetHDKeyByPath(maxKey, uint32(i), uint32(i))
		if err != nil {
			t.Log("hit the depth limit on HD key")
			break
		}
	}

	// Create the list of tests
	var tests = []struct {
		inputHDKey    *hdkeychain.ExtendedKey
		inputChain    uint32
		inputNum      uint32
		expectedNil   bool
		expectedError bool
	}{
		{nil, 0, 0, true, true},
		{validKey, 0, 0, false, false},
		{validKey, 10, 10, false, false},
		{validKey, 100, 100, false, false},
		{validKey, 2 ^ 31 + 1, 2 ^ 32 - 1, false, false},
		{validKey, 1 << 8, 1 << 8, false, false},
		{validKey, 1 << 9, 1 << 9, false, false},
		{validKey, 1 << 10, 1 << 10, false, false},
		{validKey, 1 << 11, 1 << 11, false, false},
		{validKey, 1 << 12, 1 << 12, false, false},
		{validKey, 1 << 16, 1 << 16, false, false},
		{validKey, 1<<32 - 1, 1<<32 - 1, false, false},
	}

	// Run tests
	for _, test := range tests {
		if privateKey, err := GetPrivateKeyByPath(test.inputHDKey, test.inputChain, test.inputNum); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] [%d] [%d] inputted and error not expected but got: %s", t.Name(), test.inputHDKey, test.inputChain, test.inputNum, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] [%d] [%d] inputted and error was expected", t.Name(), test.inputHDKey, test.inputChain, test.inputNum)
		} else if privateKey == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] [%d] [%d] inputted and was nil but not expected", t.Name(), test.inputHDKey, test.inputChain, test.inputNum)
		} else if privateKey != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] [%d] [%d] inputted and was NOT nil but expected to be nil", t.Name(), test.inputHDKey, test.inputChain, test.inputNum)
		} else if privateKey != nil && len(hex.EncodeToString(privateKey.Serialize())) == 0 {
			t.Errorf("%s Failed: [%v] [%d] [%d] inputted and should not be empty", t.Name(), test.inputHDKey, test.inputChain, test.inputNum)
		}
	}
}

// ExampleGetPrivateKeyByPath example using GetPrivateKeyByPath()
func ExampleGetPrivateKeyByPath() {

	hdKey, err := GenerateHDKey(SecureSeedLength)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get a private key at the path
	var privateKey *bsvec.PrivateKey
	privateKey, err = GetPrivateKeyByPath(hdKey, 0, 1)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("private key (%d) found at path %d/%d", len(privateKey.Serialize()), 0, 1)
	// Output:private key (32) found at path 0/1
}

// BenchmarkGetPrivateKeyByPath benchmarks the method GetPrivateKeyByPath()
func BenchmarkGetPrivateKeyByPath(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for i := 0; i < b.N; i++ {
		_, _ = GetPrivateKeyByPath(hdKey, 0, 1)
	}
}

// TestGetHDKeyByPath will test the method GetHDKeyByPath()
func TestGetHDKeyByPath(t *testing.T) {

	t.Parallel()

	// Generate a valid key
	validKey, err := GenerateHDKey(RecommendedSeedLength)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Max depth key
	var maxKey *hdkeychain.ExtendedKey
	maxKey, err = GetHDKeyByPath(validKey, 1<<9, 1<<9)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Test depth limit
	for i := 0; i < 1<<8-1; i++ {
		maxKey, err = GetHDKeyByPath(maxKey, uint32(i), uint32(i))
		if err != nil {
			t.Log("hit the depth limit on HD key")
			break
		}
	}

	// Create the list of tests
	var tests = []struct {
		inputHDKey    *hdkeychain.ExtendedKey
		inputChain    uint32
		inputNum      uint32
		expectedNil   bool
		expectedError bool
	}{
		{nil, 0, 0, true, true},
		{validKey, 0, 0, false, false},
		{validKey, 10, 10, false, false},
		{validKey, 100, 100, false, false},
		{validKey, 2 ^ 31 + 1, 2 ^ 32 - 1, false, false},
		{validKey, 1 << 8, 1 << 8, false, false},
		{validKey, 1 << 9, 1 << 9, false, false},
		{validKey, 1 << 10, 1 << 10, false, false},
		{validKey, 1 << 11, 1 << 11, false, false},
		{validKey, 1 << 12, 1 << 12, false, false},
		{validKey, 1 << 16, 1 << 16, false, false},
		{validKey, 1<<32 - 1, 1<<32 - 1, false, false},
	}

	// Run tests
	for _, test := range tests {
		if hdKey, err := GetHDKeyByPath(test.inputHDKey, test.inputChain, test.inputNum); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] [%d] [%d] inputted and error not expected but got: %s", t.Name(), test.inputHDKey, test.inputChain, test.inputNum, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] [%d] [%d] inputted and error was expected", t.Name(), test.inputHDKey, test.inputChain, test.inputNum)
		} else if hdKey == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] [%d] [%d] inputted and was nil but not expected", t.Name(), test.inputHDKey, test.inputChain, test.inputNum)
		} else if hdKey != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] [%d] [%d] inputted and was NOT nil but expected to be nil", t.Name(), test.inputHDKey, test.inputChain, test.inputNum)
		} else if hdKey != nil && len(hdKey.String()) == 0 {
			t.Errorf("%s Failed: [%v] [%d] [%d] inputted and should not be empty", t.Name(), test.inputHDKey, test.inputChain, test.inputNum)
		}
	}
}

// ExampleGetPrivateKeyByPath example using GetHDKeyByPath()
func ExampleGetHDKeyByPath() {

	hdKey, err := GenerateHDKey(SecureSeedLength)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get a child key
	var childKey *hdkeychain.ExtendedKey
	childKey, err = GetHDKeyByPath(hdKey, 0, 1)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("hd key (%d) found at path %d/%d", len(childKey.String()), 0, 1)
	// Output:hd key (111) found at path 0/1
}

// BenchmarkGetPrivateKeyByPath benchmarks the method GetHDKeyByPath()
func BenchmarkGetHDKeyByPath(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for i := 0; i < b.N; i++ {
		_, _ = GetHDKeyByPath(hdKey, 0, 1)
	}
}
