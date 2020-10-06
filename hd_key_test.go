package bitcoin

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvutil"
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
	// todo: make a better test (after 126 maxKey is now nil)
	for i := 0; i < 1<<8-1; i++ {
		maxKey, err = GetHDKeyByPath(maxKey, uint32(i), uint32(i))
		if i >= 126 && err == nil {
			t.Fatalf("expected to hit depth limit on HD key index: %d", i)
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

// ExampleGetHDKeyByPath example using GetHDKeyByPath()
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

// BenchmarkGetHDKeyByPath benchmarks the method GetHDKeyByPath()
func BenchmarkGetHDKeyByPath(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for i := 0; i < b.N; i++ {
		_, _ = GetHDKeyByPath(hdKey, 0, 1)
	}
}

// TestGetHDKeyByPath will test the method GetHDKeyByPath()
func TestGetHDChild(t *testing.T) {

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
	// todo: make a better test (after 126 maxKey is now nil)
	for i := 0; i < 1<<8-1; i++ {
		maxKey, err = GetHDKeyChild(maxKey, uint32(i))
		if i < 126 && err != nil {
			t.Fatalf("error occurred: %s", err.Error())
		}
		// TODO: make this better rather than grabbing the child twice. This is
		// basically a copy of the GetHDKeyByPath test
		maxKey, err = GetHDKeyChild(maxKey, uint32(i))
		if i >= 126 && err == nil {
			t.Fatalf("expected to hit depth limit on HD key index: %d", i)
		}
	}

	// Create the list of tests
	var tests = []struct {
		inputHDKey    *hdkeychain.ExtendedKey
		inputNum      uint32
		expectedNil   bool
		expectedError bool
	}{
		{nil, 0, true, true},
		{validKey, 0, false, false},
		{validKey, 10, false, false},
		{validKey, 100, false, false},
		{validKey, 2 ^ 31 + 1, false, false},
		{validKey, 1 << 8, false, false},
		{validKey, 1 << 9, false, false},
		{validKey, 1 << 10, false, false},
		{validKey, 1 << 11, false, false},
		{validKey, 1 << 12, false, false},
		{validKey, 1 << 16, false, false},
		{validKey, 1<<32 - 1, false, false},
	}

	// Run tests
	for _, test := range tests {
		if hdKey, err := GetHDKeyChild(test.inputHDKey, test.inputNum); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] [%d] inputted and error not expected but got: %s", t.Name(), test.inputHDKey, test.inputNum, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] [%d] inputted and error was expected", t.Name(), test.inputHDKey, test.inputNum)
		} else if hdKey == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] [%d] inputted and was nil but not expected", t.Name(), test.inputHDKey, test.inputNum)
		} else if hdKey != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] [%d] inputted and was NOT nil but expected to be nil", t.Name(), test.inputHDKey, test.inputNum)
		} else if hdKey != nil && len(hdKey.String()) == 0 {
			t.Errorf("%s Failed: [%v] [%d] inputted and should not be empty", t.Name(), test.inputHDKey, test.inputNum)
		}
	}
}

// ExampleGetHDKeyChild example using GetHDKeyChild()
func ExampleGetHDKeyChild() {

	hdKey, err := GenerateHDKey(SecureSeedLength)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get a child key
	var childKey *hdkeychain.ExtendedKey
	childKey, err = GetHDKeyChild(hdKey, 0)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("hd key (%d) found at path %d", len(childKey.String()), 0)
	// Output:hd key (111) found at path 0
}

// BenchmarkGetHDKeyChild benchmarks the method GetHDKeyChild()
func BenchmarkGetHDKeyChild(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for i := 0; i < b.N; i++ {
		_, _ = GetHDKeyChild(hdKey, 0)
	}
}

// TestGenerateHDKeyFromString will test the method GenerateHDKeyFromString()
func TestGenerateHDKeyFromString(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		input         string
		expectedNil   bool
		expectedError bool
	}{
		{"", true, true},
		{"0", true, true},
		{"1234567", true, true},
		{"xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE", false, false},
		{"xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUv", true, true},
		{"xprv9s21ZrQH143K3XJueaaswvbJ38UX3FhnXkcA7xF8kqeN62qEu116M1XnqaDpSE7SoKp8NxejVJG9dfpuvBC314VZNdB7W1kQN3Viwgkjr8L", false, false},
	}

	// Run tests
	for _, test := range tests {
		if hdKey, err := GenerateHDKeyFromString(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if hdKey == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was nil but not expected", t.Name(), test.input)
		} else if hdKey != nil && test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		} else if hdKey != nil && hdKey.String() != test.input {
			t.Errorf("%s Failed: [%s] inputted [%s] expected but got: %s", t.Name(), test.input, test.input, hdKey.String())
		}
	}
}

// ExampleGenerateHDKeyFromString example using GenerateHDKeyFromString()
func ExampleGenerateHDKeyFromString() {

	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("hd key generated from: %s", hdKey.String())
	// Output:hd key generated from: xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE
}

// BenchmarkGenerateHDKeyFromString benchmarks the method GenerateHDKeyFromString()
func BenchmarkGenerateHDKeyFromString(b *testing.B) {
	xPriv, _, _ := GenerateHDKeyPair(SecureSeedLength)
	for i := 0; i < b.N; i++ {
		_, _ = GenerateHDKeyFromString(xPriv)
	}
}

// TestGetPrivateKeyFromHDKey will test the method GetPrivateKeyFromHDKey()
func TestGetPrivateKeyFromHDKey(t *testing.T) {
	t.Parallel()

	validHdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Create the list of tests
	var tests = []struct {
		input         *hdkeychain.ExtendedKey
		expectedKey   string
		expectedNil   bool
		expectedError bool
	}{
		{nil, "", true, true},
		{new(hdkeychain.ExtendedKey), "", true, true},
		{validHdKey, "8511f5e1e35ab748e7639aa68666df71857866af13fda1d081d5917948a6cd34", false, false},
	}

	// Run tests
	for _, test := range tests {
		if privateKey, err := GetPrivateKeyFromHDKey(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), test.input)
		} else if privateKey == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and was nil but not expected", t.Name(), test.input)
		} else if privateKey != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		} else if privateKey != nil && hex.EncodeToString(privateKey.Serialize()) != test.expectedKey {
			t.Errorf("%s Failed: [%v] inputted [%s] expected but got: %s", t.Name(), test.input, test.expectedKey, hex.EncodeToString(privateKey.Serialize()))
		}
	}
}

// ExampleGetPrivateKeyFromHDKey example using GetPrivateKeyFromHDKey()
func ExampleGetPrivateKeyFromHDKey() {

	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var privateKey *bsvec.PrivateKey
	if privateKey, err = GetPrivateKeyFromHDKey(hdKey); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("private key: %s", hex.EncodeToString(privateKey.Serialize()))
	// Output:private key: 0ccf07f2cbe10dbe6f6034b7efbf62fc83cac3d44f49d67aa22ac8893d294e7a
}

// BenchmarkGetPrivateKeyFromHDKey benchmarks the method GetPrivateKeyFromHDKey()
func BenchmarkGetPrivateKeyFromHDKey(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for i := 0; i < b.N; i++ {
		_, _ = GetPrivateKeyFromHDKey(hdKey)
	}
}

// TestGetPublicKeyFromHDKey will test the method GetPublicKeyFromHDKey()
func TestGetPublicKeyFromHDKey(t *testing.T) {
	t.Parallel()

	validHdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Create the list of tests
	var tests = []struct {
		input         *hdkeychain.ExtendedKey
		expectedKey   string
		expectedNil   bool
		expectedError bool
	}{
		{nil, "", true, true},
		{new(hdkeychain.ExtendedKey), "", true, true},
		{validHdKey, "02f2a2942b9d1dba033d36ab0c193e680415f5c8c1ff5d854f805c8c42ed9dd1fd", false, false},
	}

	// Run tests
	var publicKey *bsvec.PublicKey
	for _, test := range tests {
		if publicKey, err = GetPublicKeyFromHDKey(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), test.input)
		} else if publicKey == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and was nil but not expected", t.Name(), test.input)
		} else if publicKey != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		} else if publicKey != nil && hex.EncodeToString(publicKey.SerializeCompressed()) != test.expectedKey {
			t.Errorf("%s Failed: [%v] inputted [%s] expected but got: %s", t.Name(), test.input, test.expectedKey, hex.EncodeToString(publicKey.SerializeCompressed()))
		}
	}
}

// ExampleGetPublicKeyFromHDKey example using GetPublicKeyFromHDKey()
func ExampleGetPublicKeyFromHDKey() {

	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var publicKey *bsvec.PublicKey
	if publicKey, err = GetPublicKeyFromHDKey(hdKey); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("public key: %s", hex.EncodeToString(publicKey.SerializeCompressed()))
	// Output:public key: 03a25f6c10eedcd41eebac22c6bbc5278690fa1aab3afc2bbe8f2277c85e5c5def
}

// BenchmarkGetPublicKeyFromHDKey benchmarks the method GetPublicKeyFromHDKey()
func BenchmarkGetPublicKeyFromHDKey(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for i := 0; i < b.N; i++ {
		_, _ = GetPublicKeyFromHDKey(hdKey)
	}
}

// TestGetAddressFromHDKey will test the method GetAddressFromHDKey()
func TestGetAddressFromHDKey(t *testing.T) {
	t.Parallel()

	validHdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Create the list of tests
	var tests = []struct {
		input           *hdkeychain.ExtendedKey
		expectedAddress string
		expectedNil     bool
		expectedError   bool
	}{
		{nil, "", true, true},
		{new(hdkeychain.ExtendedKey), "", true, true},
		{validHdKey, "13xHrMdZuqa2gpweHf37w8hu6tfv3JrnaW", false, false},
	}

	// Run tests
	var address *bsvutil.LegacyAddressPubKeyHash
	for _, test := range tests {
		if address, err = GetAddressFromHDKey(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), test.input)
		} else if address == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and was nil but not expected", t.Name(), test.input)
		} else if address != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		} else if address != nil && address.String() != test.expectedAddress {
			t.Errorf("%s Failed: [%v] inputted [%s] expected but got: %s", t.Name(), test.input, test.expectedAddress, address.String())
		}
	}
}

// ExampleGetAddressFromHDKey example using GetAddressFromHDKey()
func ExampleGetAddressFromHDKey() {

	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var address *bsvutil.LegacyAddressPubKeyHash
	if address, err = GetAddressFromHDKey(hdKey); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("address: %s", address.String())
	// Output:address: 18G2YRH3nRKRx8pnqVFUM5nAJhTZJ3YA4W
}

// BenchmarkGetAddressFromHDKey benchmarks the method GetAddressFromHDKey()
func BenchmarkGetAddressFromHDKey(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for i := 0; i < b.N; i++ {
		_, _ = GetAddressFromHDKey(hdKey)
	}
}
