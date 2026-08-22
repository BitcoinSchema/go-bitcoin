package bitcoin

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bsv-blockchain/go-bt/v2/bscript"
	bip32 "github.com/bsv-blockchain/go-sdk/compat/bip32"
	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateHDKey will test the method GenerateHDKey()
func TestGenerateHDKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		inputSeed     uint8
		expectedNil   bool
		expectedError bool
	}{
		{"zero seed defaults to recommended", 0, false, false},
		{"seed length 1 is invalid", 1, true, true},
		{"seed length 15 is invalid", 15, true, true},
		{"seed length 65 is invalid", 65, true, true},
		{"recommended seed length", RecommendedSeedLength, false, false},
		{"secure seed length", SecureSeedLength, false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			hdKey, err := GenerateHDKey(test.inputSeed)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, hdKey)
			} else {
				assert.NotNil(t, hdKey)
			}
		})
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
	for b.Loop() {
		_, _ = GenerateHDKey(RecommendedSeedLength)
	}
}

// BenchmarkGenerateHDKeySecure benchmarks the method GenerateHDKey()
func BenchmarkGenerateHDKeySecure(b *testing.B) {
	for b.Loop() {
		_, _ = GenerateHDKey(SecureSeedLength)
	}
}

// TestGenerateHDKeyPair will test the method GenerateHDKeyPair()
func TestGenerateHDKeyPair(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		inputSeed     uint8
		expectedError bool
	}{
		{"zero seed defaults to recommended", 0, false},
		{"seed length 1 is invalid", 1, true},
		{"seed length 15 is invalid", 15, true},
		{"seed length 65 is invalid", 65, true},
		{"recommended seed length", RecommendedSeedLength, false},
		{"secure seed length", SecureSeedLength, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			privateKey, publicKey, err := GenerateHDKeyPair(test.inputSeed)
			if test.expectedError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, privateKey)
			assert.NotEmpty(t, publicKey)
		})
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
	for b.Loop() {
		_, _, _ = GenerateHDKeyPair(RecommendedSeedLength)
	}
}

// BenchmarkGenerateHDKeyPairSecure benchmarks the method GenerateHDKeyPair()
func BenchmarkGenerateHDKeyPairSecure(b *testing.B) {
	for b.Loop() {
		_, _, _ = GenerateHDKeyPair(SecureSeedLength)
	}
}

// TestGetPrivateKeyByPath will test the method GetPrivateKeyByPath()
func TestGetPrivateKeyByPath(t *testing.T) {
	t.Parallel()

	// Generate a valid key
	validKey := mustHDKey(t)

	tests := []struct {
		name          string
		inputHDKey    *bip32.ExtendedKey
		inputChain    uint32
		inputNum      uint32
		expectedNil   bool
		expectedError bool
	}{
		// {"nil key", nil, 0, 0, true, true},
		{"chain 0 num 0", validKey, 0, 0, false, false},
		{"chain 10 num 10", validKey, 10, 10, false, false},
		{"chain 100 num 100", validKey, 100, 100, false, false},
		{"chain 1<<31+1 num 1<<32-1", validKey, 1<<31 + 1, 1<<32 - 1, false, false},
		{"chain 1<<8 num 1<<8", validKey, 1 << 8, 1 << 8, false, false},
		{"chain 1<<9 num 1<<9", validKey, 1 << 9, 1 << 9, false, false},
		{"chain 1<<10 num 1<<10", validKey, 1 << 10, 1 << 10, false, false},
		{"chain 1<<11 num 1<<11", validKey, 1 << 11, 1 << 11, false, false},
		{"chain 1<<12 num 1<<12", validKey, 1 << 12, 1 << 12, false, false},
		{"chain 1<<16 num 1<<16", validKey, 1 << 16, 1 << 16, false, false},
		{"chain 1<<32-1 num 1<<32-1", validKey, 1<<32 - 1, 1<<32 - 1, false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			privateKey, err := GetPrivateKeyByPath(test.inputHDKey, test.inputChain, test.inputNum)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, privateKey)
			} else {
				require.NotNil(t, privateKey)
				assert.NotEmpty(t, hex.EncodeToString(privateKey.Serialize()))
			}
		})
	}
}

// TestGetPrivateKeyByPathPanic tests for nil case in GetPrivateKeyByPath()
func TestGetPrivateKeyByPathPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, err := GetPrivateKeyByPath(nil, 0, 1)
		assert.Error(t, err)
	})
}

// TestGetPrivateKeyByPathHardenedFromPublic ensures GetPrivateKeyByPath returns
// the derivation error when a hardened child is requested from a public-only
// (neutered) extended key.
func TestGetPrivateKeyByPathHardenedFromPublic(t *testing.T) {
	t.Parallel()

	masterKey, err := GenerateHDKey(RecommendedSeedLength)
	require.NoError(t, err)

	// Neuter the key so only the public portion remains
	var xPub string
	xPub, err = GetExtendedPublicKey(masterKey)
	require.NoError(t, err)

	var publicKey *bip32.ExtendedKey
	publicKey, err = GetHDKeyFromExtendedPublicKey(xPub)
	require.NoError(t, err)

	// A hardened chain index cannot be derived from a public key
	privateKey, pathErr := GetPrivateKeyByPath(publicKey, bip32.HardenedKeyStart, 0)
	require.Error(t, pathErr)
	assert.Nil(t, privateKey)
}

// ExampleGetPrivateKeyByPath example using GetPrivateKeyByPath()
func ExampleGetPrivateKeyByPath() {
	hdKey, err := GenerateHDKey(SecureSeedLength)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get a private key at the path
	var privateKey *ec.PrivateKey
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
	for b.Loop() {
		_, _ = GetPrivateKeyByPath(hdKey, 0, 1)
	}
}

// TestGetHDKeyByPath will test the method GetHDKeyByPath()
func TestGetHDKeyByPath(t *testing.T) {
	t.Parallel()

	// Generate a valid key
	validKey, err := GenerateHDKey(RecommendedSeedLength)
	require.NoError(t, err)
	assert.NotNil(t, validKey)

	// Max depth key
	/*
		var maxKey *bip32.ExtendedKey
		maxKey, err = GetHDKeyByPath(validKey, 1<<9, 1<<9)
		if err != nil {
			t.Fatalf("error occurred: %s", err.Error())
		}
	*/

	// Test depth limit
	// todo: make a better test (after 126 maxKey is now nil)
	/*
		for i := 0; i < 1<<8-1; i++ {
			maxKey, err = GetHDKeyByPath(maxKey, uint32(i), uint32(i))
			if i >= 126 && err == nil {
				t.Fatalf("expected to hit depth limit on HD key index: %d", i)
			}
		}
	*/

	tests := []struct {
		name          string
		inputHDKey    *bip32.ExtendedKey
		inputChain    uint32
		inputNum      uint32
		expectedNil   bool
		expectedError bool
	}{
		{"chain 0 num 0", validKey, 0, 0, false, false},
		{"chain 10 num 10", validKey, 10, 10, false, false},
		{"chain 100 num 100", validKey, 100, 100, false, false},
		{"chain 1<<31+1 num 1<<32-1", validKey, 1<<31 + 1, 1<<32 - 1, false, false},
		{"chain 1<<8 num 1<<8", validKey, 1 << 8, 1 << 8, false, false},
		{"chain 1<<9 num 1<<9", validKey, 1 << 9, 1 << 9, false, false},
		{"chain 1<<10 num 1<<10", validKey, 1 << 10, 1 << 10, false, false},
		{"chain 1<<11 num 1<<11", validKey, 1 << 11, 1 << 11, false, false},
		{"chain 1<<12 num 1<<12", validKey, 1 << 12, 1 << 12, false, false},
		{"chain 1<<16 num 1<<16", validKey, 1 << 16, 1 << 16, false, false},
		{"chain 1<<32-1 num 1<<32-1", validKey, 1<<32 - 1, 1<<32 - 1, false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			hdKey, err := GetHDKeyByPath(test.inputHDKey, test.inputChain, test.inputNum)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, hdKey)
			} else {
				require.NotNil(t, hdKey)
				assert.NotEmpty(t, hdKey.String())
			}
		})
	}
}

// TestGetHDKeyByPathPanic tests for nil case in GetHDKeyByPath()
func TestGetHDKeyByPathPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, err := GetHDKeyByPath(nil, 0, 1)
		assert.Error(t, err)
	})
}

// ExampleGetHDKeyByPath example using GetHDKeyByPath()
func ExampleGetHDKeyByPath() {
	hdKey, err := GenerateHDKey(SecureSeedLength)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get a child key
	var childKey *bip32.ExtendedKey
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
	for b.Loop() {
		_, _ = GetHDKeyByPath(hdKey, 0, 1)
	}
}

// TestGetHDKeyChild will test the method GetHDKeyChild()
func TestGetHDKeyChild(t *testing.T) {
	t.Parallel()

	// Generate a valid key
	validKey, err := GenerateHDKey(RecommendedSeedLength)
	require.NoError(t, err)
	assert.NotNil(t, validKey)

	// Max depth key
	/*
		var maxKey *bip32.ExtendedKey
		maxKey, err = GetHDKeyByPath(validKey, 1<<9, 1<<9)
		if err != nil {
			t.Fatalf("error occurred: %s", err.Error())
		}
	*/

	// Test depth limit
	// todo: make a better test (after 126 maxKey is now nil)
	/*
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
	*/

	tests := []struct {
		name          string
		inputHDKey    *bip32.ExtendedKey
		inputNum      uint32
		expectedNil   bool
		expectedError bool
	}{
		// {"nil key", nil, 0, true, true},
		{"num 0", validKey, 0, false, false},
		{"num 10", validKey, 10, false, false},
		{"num 100", validKey, 100, false, false},
		{"num 1<<31+1", validKey, 1<<31 + 1, false, false},
		{"num 1<<8", validKey, 1 << 8, false, false},
		{"num 1<<9", validKey, 1 << 9, false, false},
		{"num 1<<10", validKey, 1 << 10, false, false},
		{"num 1<<11", validKey, 1 << 11, false, false},
		{"num 1<<12", validKey, 1 << 12, false, false},
		{"num 1<<16", validKey, 1 << 16, false, false},
		{"num 1<<32-1", validKey, 1<<32 - 1, false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			hdKey, err := GetHDKeyChild(test.inputHDKey, test.inputNum)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, hdKey)
			} else {
				require.NotNil(t, hdKey)
				assert.NotEmpty(t, hdKey.String())
			}
		})
	}
}

// TestGetHDKeyChildPanic tests for nil case in GetHDKeyChild()
func TestGetHDKeyChildPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, err := GetHDKeyChild(nil, 1)
		assert.Error(t, err)
	})
}

// ExampleGetHDKeyChild example using GetHDKeyChild()
func ExampleGetHDKeyChild() {
	hdKey, err := GenerateHDKey(SecureSeedLength)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get a child key
	var childKey *bip32.ExtendedKey
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
	for b.Loop() {
		_, _ = GetHDKeyChild(hdKey, 0)
	}
}

// TestGenerateHDKeyFromString will test the method GenerateHDKeyFromString()
func TestGenerateHDKeyFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectedNil   bool
		expectedError bool
	}{
		{"empty string", "", true, true},
		{"zero", "0", true, true},
		{"too short numeric", "1234567", true, true},
		{"valid xprv", "xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE", false, false},
		{"truncated xprv", "xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUv", true, true},
		{"valid xprv 2", "xprv9s21ZrQH143K3XJueaaswvbJ38UX3FhnXkcA7xF8kqeN62qEu116M1XnqaDpSE7SoKp8NxejVJG9dfpuvBC314VZNdB7W1kQN3Viwgkjr8L", false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			hdKey, err := GenerateHDKeyFromString(test.input)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, hdKey)
			} else {
				require.NotNil(t, hdKey)
				assert.Equal(t, test.input, hdKey.String())
			}
		})
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
	for b.Loop() {
		_, _ = GenerateHDKeyFromString(xPriv)
	}
}

// TestGetPrivateKeyFromHDKey will test the method GetPrivateKeyFromHDKey()
func TestGetPrivateKeyFromHDKey(t *testing.T) {
	t.Parallel()

	validHdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	require.NoError(t, err)
	assert.NotNil(t, validHdKey)

	tests := []struct {
		name          string
		input         *bip32.ExtendedKey
		expectedKey   string
		expectedNil   bool
		expectedError bool
	}{
		{"zeroed extended key", new(bip32.ExtendedKey), "", true, true},
		{"valid hd key", validHdKey, "8511f5e1e35ab748e7639aa68666df71857866af13fda1d081d5917948a6cd34", false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			privateKey, err := GetPrivateKeyFromHDKey(test.input)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, privateKey)
			} else {
				require.NotNil(t, privateKey)
				assert.Equal(t, test.expectedKey, hex.EncodeToString(privateKey.Serialize()))
			}
		})
	}
}

// TestGetPrivateKeyFromHDKeyPanic tests for nil case in GetPrivateKeyFromHDKey()
func TestGetPrivateKeyFromHDKeyPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, err := GetPrivateKeyFromHDKey(nil)
		assert.Error(t, err)
	})
}

// ExampleGetPrivateKeyFromHDKey example using GetPrivateKeyFromHDKey()
func ExampleGetPrivateKeyFromHDKey() {
	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var privateKey *ec.PrivateKey
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
	for b.Loop() {
		_, _ = GetPrivateKeyFromHDKey(hdKey)
	}
}

// TestGetPrivateKeyStringFromHDKey will test the method GetPrivateKeyStringFromHDKey()
func TestGetPrivateKeyStringFromHDKey(t *testing.T) {
	t.Parallel()

	validHdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	require.NoError(t, err)
	assert.NotNil(t, validHdKey)

	tests := []struct {
		name          string
		input         *bip32.ExtendedKey
		expectedKey   string
		expectedError bool
	}{
		{"zeroed extended key", new(bip32.ExtendedKey), "", true},
		{"valid hd key", validHdKey, "8511f5e1e35ab748e7639aa68666df71857866af13fda1d081d5917948a6cd34", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			privateKey, err := GetPrivateKeyStringFromHDKey(test.input)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.expectedKey, privateKey)
		})
	}
}

// TestGetPrivateKeyStringFromHDKeyPanic tests for nil case in GetPrivateKeyStringFromHDKey()
func TestGetPrivateKeyStringFromHDKeyPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, err := GetPrivateKeyStringFromHDKey(nil)
		assert.Error(t, err)
	})
}

// ExampleGetPrivateKeyStringFromHDKey example using GetPrivateKeyStringFromHDKey()
func ExampleGetPrivateKeyStringFromHDKey() {
	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var privateKey string
	if privateKey, err = GetPrivateKeyStringFromHDKey(hdKey); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("private key: %s", privateKey)
	// Output:private key: 0ccf07f2cbe10dbe6f6034b7efbf62fc83cac3d44f49d67aa22ac8893d294e7a
}

// BenchmarkGetPrivateKeyStringFromHDKey benchmarks the method GetPrivateKeyStringFromHDKey()
func BenchmarkGetPrivateKeyStringFromHDKey(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for b.Loop() {
		_, _ = GetPrivateKeyStringFromHDKey(hdKey)
	}
}

// TestGetPublicKeyFromHDKey will test the method GetPublicKeyFromHDKey()
func TestGetPublicKeyFromHDKey(t *testing.T) {
	t.Parallel()

	validHdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	require.NoError(t, err)
	assert.NotNil(t, validHdKey)

	tests := []struct {
		name          string
		input         *bip32.ExtendedKey
		expectedKey   string
		expectedNil   bool
		expectedError bool
	}{
		{"zeroed extended key", new(bip32.ExtendedKey), "", true, true},
		{"valid hd key", validHdKey, "02f2a2942b9d1dba033d36ab0c193e680415f5c8c1ff5d854f805c8c42ed9dd1fd", false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			publicKey, err := GetPublicKeyFromHDKey(test.input)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, publicKey)
			} else {
				require.NotNil(t, publicKey)
				assert.Equal(t, test.expectedKey, hex.EncodeToString(publicKey.Compressed()))
			}
		})
	}
}

// TestGetPublicKeyFromHDKeyPanic tests for nil case in GetPublicKeyFromHDKey()
func TestGetPublicKeyFromHDKeyPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, err := GetPublicKeyFromHDKey(nil)
		assert.Error(t, err)
	})
}

// ExampleGetPublicKeyFromHDKey example using GetPublicKeyFromHDKey()
func ExampleGetPublicKeyFromHDKey() {
	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var publicKey *ec.PublicKey
	if publicKey, err = GetPublicKeyFromHDKey(hdKey); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("public key: %s", hex.EncodeToString(publicKey.Compressed()))
	// Output:public key: 03a25f6c10eedcd41eebac22c6bbc5278690fa1aab3afc2bbe8f2277c85e5c5def
}

// BenchmarkGetPublicKeyFromHDKey benchmarks the method GetPublicKeyFromHDKey()
func BenchmarkGetPublicKeyFromHDKey(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for b.Loop() {
		_, _ = GetPublicKeyFromHDKey(hdKey)
	}
}

// TestGetAddressFromHDKey will test the method GetAddressFromHDKey()
func TestGetAddressFromHDKey(t *testing.T) {
	t.Parallel()

	validHdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	require.NoError(t, err)
	assert.NotNil(t, validHdKey)

	tests := []struct {
		name            string
		input           *bip32.ExtendedKey
		expectedAddress string
		expectedNil     bool
		expectedError   bool
		mainnet         bool
	}{
		{"zeroed extended key mainnet", new(bip32.ExtendedKey), "", true, true, true},
		{"valid hd key mainnet", validHdKey, "13xHrMdZuqa2gpweHf37w8hu6tfv3JrnaW", false, false, true},
		{"valid hd key testnet", validHdKey, "miUF9QiYis1HTwRG1E1Vm3vDxtGczs2oph", false, false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			address, err := GetAddressFromHDKey(test.input, test.mainnet)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, address)
			} else {
				require.NotNil(t, address)
				assert.Equal(t, test.expectedAddress, address.AddressString)
			}
		})
	}
}

// TestGetAddressFromHDKeyPanic tests for nil case in GetAddressFromHDKey()
func TestGetAddressFromHDKeyPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, err := GetAddressFromHDKey(nil, true)
		assert.Error(t, err)
	})
}

// ExampleGetAddressFromHDKey example using GetAddressFromHDKey()
func ExampleGetAddressFromHDKey() {
	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var address *bscript.Address
	if address, err = GetAddressFromHDKey(hdKey, true); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("address: %s", address.AddressString)
	// Output:address: 18G2YRH3nRKRx8pnqVFUM5nAJhTZJ3YA4W
}

// BenchmarkGetAddressFromHDKey benchmarks the method GetAddressFromHDKey()
func BenchmarkGetAddressFromHDKey(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for b.Loop() {
		_, _ = GetAddressFromHDKey(hdKey, true)
	}
}

// TestGetAddressStringFromHDKey will test the method GetAddressStringFromHDKey()
func TestGetAddressStringFromHDKey(t *testing.T) {
	t.Parallel()

	validHdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	require.NoError(t, err)
	assert.NotNil(t, validHdKey)

	tests := []struct {
		name            string
		input           *bip32.ExtendedKey
		expectedAddress string
		expectedError   bool
		mainnet         bool
	}{
		{"zeroed extended key mainnet", new(bip32.ExtendedKey), "", true, true},
		{"valid hd key mainnet", validHdKey, "13xHrMdZuqa2gpweHf37w8hu6tfv3JrnaW", false, true},
		{"valid hd key testnet", validHdKey, "miUF9QiYis1HTwRG1E1Vm3vDxtGczs2oph", false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			address, err := GetAddressStringFromHDKey(test.input, test.mainnet)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.expectedAddress, address)
		})
	}
}

// TestGetAddressStringFromHDKeyPanic tests for nil case in GetAddressStringFromHDKey()
func TestGetAddressStringFromHDKeyPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, err := GetAddressStringFromHDKey(nil, true)
		assert.Error(t, err)
	})
}

// ExampleGetAddressStringFromHDKey example using GetAddressStringFromHDKey()
func ExampleGetAddressStringFromHDKey() {
	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var address string
	if address, err = GetAddressStringFromHDKey(hdKey, true); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("address: %s", address)
	// Output:address: 18G2YRH3nRKRx8pnqVFUM5nAJhTZJ3YA4W
}

// BenchmarkGetAddressStringFromHDKey benchmarks the method GetAddressStringFromHDKey()
func BenchmarkGetAddressStringFromHDKey(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for b.Loop() {
		_, _ = GetAddressStringFromHDKey(hdKey, true)
	}
}

// TestGetPublicKeysForPath will test the method GetPublicKeysForPath()
func TestGetPublicKeysForPath(t *testing.T) {
	t.Parallel()

	validHdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	require.NoError(t, err)
	assert.NotNil(t, validHdKey)

	tests := []struct {
		name            string
		input           *bip32.ExtendedKey
		inputNum        uint32
		expectedPubKey1 string
		expectedPubKey2 string
		expectedNil     bool
		expectedError   bool
	}{
		{"zeroed extended key", new(bip32.ExtendedKey), 1, "", "", true, true},
		{"valid hd key num 1", validHdKey, 1, "03cc3334f0a6f0fae0420d1442ca0ce64fad0da76d652f2cc3b333e7ed95b97259", "02ceb23902f8dcf6fbff656597ee0343e05c907c6dfcdd8aaf6d033e14e85fd955", false, false},
		{"valid hd key num 2", validHdKey, 2, "020cb908e3b9f3de7c9b40e7bcce63708c5617536d85cf4ab5635e3d3819c02c37", "030007ae60fc6eef98ea17b4f80f9b791e61ea94936e8a9e6ec343eeaa50a875e0", false, false},
		{"valid hd key num 3", validHdKey, 3, "0342593453c476ac6c78eb1b1e586df00b20352e61c42536fe1b33c9fdf3bfbb6f", "03786a41dbf0b099256da26cb0019e10063628f6ce31b96801703f1bb2e1b17724", false, false},
		{"valid hd key num 4", validHdKey, 4, "0366dcdebfc8abfd34bffc181ccb54f1706839a80ad4f0842ae5a43f39fdd35c1e", "03a095db29ae9ee0b22c775118b4444b59db40acdea137fd9ecd9c68dacf50a644", false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			pubKeys, err := GetPublicKeysForPath(test.input, test.inputNum)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, pubKeys)
			} else {
				require.Len(t, pubKeys, 2)
				assert.Equal(t, test.expectedPubKey1, hex.EncodeToString(pubKeys[0].Compressed()))
				assert.Equal(t, test.expectedPubKey2, hex.EncodeToString(pubKeys[1].Compressed()))
			}
		})
	}
}

// TestGetPublicKeysForPathPanic tests for nil case in GetPublicKeysForPath()
func TestGetPublicKeysForPathPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, err := GetPublicKeysForPath(nil, 1)
		assert.Error(t, err)
	})
}

// ExampleGetPublicKeysForPath example using GetPublicKeysForPath()
func ExampleGetPublicKeysForPath() {
	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var publicKeys []*ec.PublicKey
	publicKeys, err = GetPublicKeysForPath(hdKey, 5)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("found [%d] keys! Key 1: %s Key 2: %s", len(publicKeys), hex.EncodeToString(publicKeys[0].Compressed()), hex.EncodeToString(publicKeys[1].Compressed()))
	// Output:found [2] keys! Key 1: 03f87ac38fb0cfca12988b51a2f1cd3e85bb4aeb1b05f549682190ac8205a67d30 Key 2: 02e78303aeef1acce1347c6493fadc1914e6d85ef3189a8856afb3accd53fbd9c5
}

// BenchmarkGetPublicKeysForPath benchmarks the method GetPublicKeysForPath()
func BenchmarkGetPublicKeysForPath(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for b.Loop() {
		_, _ = GetPublicKeysForPath(hdKey, 5)
	}
}

// TestGetAddressesForPath will test the method GetAddressesForPath()
func TestGetAddressesForPath(t *testing.T) {
	t.Parallel()

	validHdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	require.NoError(t, err)
	assert.NotNil(t, validHdKey)

	tests := []struct {
		name             string
		input            *bip32.ExtendedKey
		inputNum         uint32
		expectedAddress1 string
		expectedAddress2 string
		expectedNil      bool
		expectedError    bool
		mainnet          bool
	}{
		{"zeroed extended key mainnet", new(bip32.ExtendedKey), 1, "", "", true, true, true},
		{"valid hd key num 1 mainnet", validHdKey, 1, "1KMxfSfRCkC1jrBAuYaLde4XBzdsWApbdH", "174DL9ZbBWx568ssAg8w2YwW6FTTBwXGEu", false, false, true},
		{"valid hd key num 2 mainnet", validHdKey, 2, "18s3peTU7fMSwgui54avpnqm1126pRVccw", "1KgZZ3NsJDw3v1GPHBj8ASnxutA1kFxo2i", false, false, true},
		{"valid hd key num 1 testnet", validHdKey, 1, "mysuxVkQ1mdGWxend7YiTZGr3zEaTcMjrz", "mmaAdCeZzYPKsFMUtF7JrU9pxF4AAgMHK5", false, false, false},
		{"zeroed extended key testnet", new(bip32.ExtendedKey), 1, "", "", true, true, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			addresses, err := GetAddressesForPath(test.input, test.inputNum, test.mainnet)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, addresses)
			} else {
				require.Len(t, addresses, 2)
				assert.Equal(t, test.expectedAddress1, addresses[0])
				assert.Equal(t, test.expectedAddress2, addresses[1])
			}
		})
	}
}

// TestGetAddressesForPathPanic tests for nil case in GetAddressesForPath()
func TestGetAddressesForPathPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, err := GetAddressesForPath(nil, 1, true)
		assert.Error(t, err)
	})
}

// ExampleGetAddressesForPath example using GetAddressesForPath()
func ExampleGetAddressesForPath() {
	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var addresses []string
	addresses, err = GetAddressesForPath(hdKey, 5, true)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("found [%d] addresses! Address 1: %s Address 2: %s", len(addresses), addresses[0], addresses[1])
	// Output:found [2] addresses! Address 1: 1JHGJTqsiFHo4yQYJ1WbTvbxYMZC7nZKYb Address 2: 1DTHBcGeJFRmS26S11tt2EddhSkFM8tmze
}

// BenchmarkGetAddressesForPath benchmarks the method GetAddressesForPath()
func BenchmarkGetAddressesForPath(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for b.Loop() {
		_, _ = GetAddressesForPath(hdKey, 5, true)
	}
}

// TestGetExtendedPublicKey will test the method GetExtendedPublicKey()
func TestGetExtendedPublicKey(t *testing.T) {
	t.Parallel()

	validHdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K4FdJCmPQe1CFUvK3PKVrcp3b5xVr5Bs3cP5ab6ytszeHggTmHoqTXpaa8CgYPxZZzigSGCDjtyWdUDJqPogb1JGWAPkBLdF")
	require.NoError(t, err)
	assert.NotNil(t, validHdKey)

	tests := []struct {
		name          string
		input         *bip32.ExtendedKey
		expectedKey   string
		expectedError bool
	}{
		{"valid hd key", validHdKey, "xpub661MyMwAqRbcGjhmJnvR198z2x9XnnDhz2yBtLuTdXQ2VBQj8eJ9RnxmXxKnRPhYy6nLsmabmUfVkbajvP7aZASrrnoZkzmwgyjiNskiefG", false},
		{"zeroed extended key", new(bip32.ExtendedKey), "zeroed extended key", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			xPub, err := GetExtendedPublicKey(test.input)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.expectedKey, xPub)
		})
	}
}

// TestGetExtendedPublicKeyPanic tests for nil case in GetExtendedPublicKey()
func TestGetExtendedPublicKeyPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, err := GetExtendedPublicKey(nil)
		assert.Error(t, err)
	})
}

// ExampleGetExtendedPublicKey example using GetExtendedPublicKey()
func ExampleGetExtendedPublicKey() {
	hdKey, err := GenerateHDKeyFromString("xprv9s21ZrQH143K3PZSwbEeXEYq74EbnfMngzAiMCZcfjzyRpUvt2vQJnaHRTZjeuEmLXeN6BzYRoFsEckfobxE9XaRzeLGfQoxzPzTRyRb6oE")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var xPub string
	xPub, err = GetExtendedPublicKey(hdKey)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("xPub: %s", xPub)
	// Output:xPub: xpub661MyMwAqRbcFsdv3cmetNVZf656C85e4D6K9ayEE5XxJcp5RaEeratmGmg7ggt3ZShibYcsusYPom69yDG9hf3UE1i4LrXJbuA9d7hPujt
}

// BenchmarkGetExtendedPublicKey benchmarks the method GetExtendedPublicKey()
func BenchmarkGetExtendedPublicKey(b *testing.B) {
	hdKey, _ := GenerateHDKey(SecureSeedLength)
	for b.Loop() {
		_, _ = GetExtendedPublicKey(hdKey)
	}
}

// TestGetHDKeyFromExtendedPublicKey will test the method GetHDKeyFromExtendedPublicKey()
func TestGetHDKeyFromExtendedPublicKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectedKey   string
		expectedError bool
		expectedNil   bool
	}{
		{
			"valid xpub",
			"xpub661MyMwAqRbcGjhmJnvR198z2x9XnnDhz2yBtLuTdXQ2VBQj8eJ9RnxmXxKnRPhYy6nLsmabmUfVkbajvP7aZASrrnoZkzmwgyjiNskiefG",
			"xpub661MyMwAqRbcGjhmJnvR198z2x9XnnDhz2yBtLuTdXQ2VBQj8eJ9RnxmXxKnRPhYy6nLsmabmUfVkbajvP7aZASrrnoZkzmwgyjiNskiefG",
			false,
			false,
		},
		{
			"valid xpub 2",
			"xpub661MyMwAqRbcH3WGvLjupmr43L1GVH3MP2WQWvdreDraBeFJy64Xxv4LLX9ZVWWz3ZjZkMuZtSsc9qH9JZR74bR4PWkmtEvP423r6DJR8kA",
			"xpub661MyMwAqRbcH3WGvLjupmr43L1GVH3MP2WQWvdreDraBeFJy64Xxv4LLX9ZVWWz3ZjZkMuZtSsc9qH9JZR74bR4PWkmtEvP423r6DJR8kA",
			false,
			false,
		},
		{
			"empty string",
			"",
			"",
			true,
			true,
		},
		{
			"zero",
			"0",
			"",
			true,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			xPub, err := GetHDKeyFromExtendedPublicKey(test.input)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, xPub)
			} else {
				require.NotNil(t, xPub)
				assert.Equal(t, test.expectedKey, xPub.String())
			}
		})
	}
}

// ExampleGetHDKeyFromExtendedPublicKey example using GetHDKeyFromExtendedPublicKey()
func ExampleGetHDKeyFromExtendedPublicKey() {
	// Start with an existing xPub
	xPub := "xpub661MyMwAqRbcH3WGvLjupmr43L1GVH3MP2WQWvdreDraBeFJy64Xxv4LLX9ZVWWz3ZjZkMuZtSsc9qH9JZR74bR4PWkmtEvP423r6DJR8kA"

	// Convert to a HD key
	key, err := GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("key: %s", key.String())
	// Output:key: xpub661MyMwAqRbcH3WGvLjupmr43L1GVH3MP2WQWvdreDraBeFJy64Xxv4LLX9ZVWWz3ZjZkMuZtSsc9qH9JZR74bR4PWkmtEvP423r6DJR8kA
}

// BenchmarkGetHDKeyFromExtendedPublicKey benchmarks the method GetHDKeyFromExtendedPublicKey()
func BenchmarkGetHDKeyFromExtendedPublicKey(b *testing.B) {
	xPub := "xpub661MyMwAqRbcH3WGvLjupmr43L1GVH3MP2WQWvdreDraBeFJy64Xxv4LLX9ZVWWz3ZjZkMuZtSsc9qH9JZR74bR4PWkmtEvP423r6DJR8kA"
	for b.Loop() {
		_, _ = GetHDKeyFromExtendedPublicKey(xPub)
	}
}
