package bitcoin

import (
	"encoding/hex"
	"fmt"
	"testing"

	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreatePrivateKey will test the method CreatePrivateKey()
func TestCreatePrivateKey(t *testing.T) {
	rawKey, err := CreatePrivateKey()
	require.NoError(t, err)
	assert.NotNil(t, rawKey)
	assert.Len(t, rawKey.Serialize(), 32)
}

// ExampleCreatePrivateKey example using CreatePrivateKey()
func ExampleCreatePrivateKey() {
	rawKey, err := CreatePrivateKey()
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	} else if len(rawKey.Serialize()) > 0 {
		fmt.Printf("key created successfully!")
	}
	// Output:key created successfully!
}

// BenchmarkCreatePrivateKey benchmarks the method CreatePrivateKey()
func BenchmarkCreatePrivateKey(b *testing.B) {
	for b.Loop() {
		_, _ = CreatePrivateKey()
	}
}

// TestCreatePrivateKeyString will test the method CreatePrivateKeyString()
func TestCreatePrivateKeyString(t *testing.T) {
	key, err := CreatePrivateKeyString()
	require.NoError(t, err)
	assert.Len(t, key, 64)
}

// ExampleCreatePrivateKeyString example using CreatePrivateKeyString()
func ExampleCreatePrivateKeyString() {
	key, err := CreatePrivateKeyString()
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	} else if len(key) > 0 {
		fmt.Printf("key created successfully!")
	}
	// Output:key created successfully!
}

// BenchmarkCreatePrivateKeyString benchmarks the method CreatePrivateKeyString()
func BenchmarkCreatePrivateKeyString(b *testing.B) {
	for b.Loop() {
		_, _ = CreatePrivateKeyString()
	}
}

// TestPrivateKeyFromString will test the method PrivateKeyFromString()
func TestPrivateKeyFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectedKey   string
		expectedNil   bool
		expectedError bool
	}{
		{"valid key", testPrivateKeyHex, testPrivateKeyHex, false, false},
		{"uppercase hex", "E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F2B75", "e83385af76b2b1997326b567461fb73dd9c27eab9e1e86d26779f4650c5f2b75", false, false},
		{"short key left-padded", "E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F", "0000e83385af76b2b1997326b567461fb73dd9c27eab9e1e86d26779f4650c5f", false, false},
		{"odd length hex", "E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F", "", true, true},
		{"too short odd hex", "1234567", "", true, true},
		{caseSingleZero, "0", "", true, true},
		{caseEmpty, "", "", true, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rawKey, err := PrivateKeyFromString(test.input)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, rawKey)
			} else {
				require.NotNil(t, rawKey)
				assert.Equal(t, test.expectedKey, hex.EncodeToString(rawKey.Serialize()))
			}
		})
	}
}

// ExamplePrivateKeyFromString example using PrivateKeyFromString()
func ExamplePrivateKeyFromString() {
	key, err := PrivateKeyFromString(testPrivateKeyHex)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("key converted: %s", hex.EncodeToString(key.Serialize()))
	// Output:key converted: 54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd
}

// BenchmarkPrivateKeyFromString benchmarks the method PrivateKeyFromString()
func BenchmarkPrivateKeyFromString(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for b.Loop() {
		_, _ = PrivateKeyFromString(key)
	}
}

// TestPrivateAndPublicKeys will test the method PrivateAndPublicKeys()
func TestPrivateAndPublicKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		input              string
		expectedPrivateKey string
		expectedNil        bool
		expectedError      bool
	}{
		{caseEmpty, "", "", true, true},
		{caseSingleZero, "0", "", true, true},
		{"short zeros", "00000", "", true, true},
		{"dashes", "0-0-0-0-0", "", true, true},
		{"invalid hex chars", "z4035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abz", "", true, true},
		{"valid key", testPrivateKeyHex, testPrivateKeyHex, false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			privateKey, publicKey, err := PrivateAndPublicKeys(test.input)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, privateKey)
				assert.Nil(t, publicKey)
			} else {
				require.NotNil(t, privateKey)
				require.NotNil(t, publicKey)
				assert.Equal(t, test.expectedPrivateKey, hex.EncodeToString(privateKey.Serialize()))
			}
		})
	}
}

// ExamplePrivateAndPublicKeys example using PrivateAndPublicKeys()
func ExamplePrivateAndPublicKeys() {
	privateKey, publicKey, err := PrivateAndPublicKeys(testPrivateKeyHex)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("private key: %s public key: %s", hex.EncodeToString(privateKey.Serialize()), hex.EncodeToString(publicKey.Compressed()))

	// Output:private key: 54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd public key: 031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f
}

// BenchmarkPrivateAndPublicKeys benchmarks the method PrivateAndPublicKeys()
func BenchmarkPrivateAndPublicKeys(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for b.Loop() {
		_, _, _ = PrivateAndPublicKeys(key)
	}
}

// TestPrivateKeyToWif will test the method PrivateKeyToWif()
func TestPrivateKeyToWif(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectedWif   string
		expectedNil   bool
		expectedError bool
	}{
		{caseEmpty, "", "", true, true},
		{caseSingleZero, "0", "", true, true},
		{"zeros key", "000000", "5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAbuatmU", false, false},
		{"ascii key bytes", "6D792070726976617465206B6579", "5HpHagT65TZzG1PH3CSu63k8DbuTZnNJf6HgyQNymvXmALAsm9s", false, false},
		{"invalid hex chars", "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8azz", "", true, true},
		{"valid key", testPrivateKeyHex, testUncompressedWIF, false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			privateWif, err := PrivateKeyToWif(test.input)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.expectedNil {
				assert.Nil(t, privateWif)
			} else {
				require.NotNil(t, privateWif)
				assert.Equal(t, test.expectedWif, privateWif.String())
			}
		})
	}
}

// ExamplePrivateKeyToWif example using PrivateKeyToWif()
func ExamplePrivateKeyToWif() {
	privateWif, err := PrivateKeyToWif(testPrivateKeyHex)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("converted wif: %s", privateWif.String())

	// Output:converted wif: 5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei
}

// BenchmarkPrivateKeyToWif benchmarks the method PrivateKeyToWif()
func BenchmarkPrivateKeyToWif(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for b.Loop() {
		_, _ = PrivateKeyToWif(key)
	}
}

// TestPrivateKeyToWifString will test the method PrivateKeyToWifString()
func TestPrivateKeyToWifString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectedWif   string
		expectedError bool
	}{
		{caseEmpty, "", "", true},
		{caseSingleZero, "0", "", true},
		{"zeros key", "000000", "5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAbuatmU", false},
		{"ascii key bytes", "6D792070726976617465206B6579", "5HpHagT65TZzG1PH3CSu63k8DbuTZnNJf6HgyQNymvXmALAsm9s", false},
		{"invalid hex chars", "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8azz", "", true},
		{"valid key", testPrivateKeyHex, testUncompressedWIF, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			privateWif, err := PrivateKeyToWifString(test.input)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.expectedWif, privateWif)
		})
	}
}

// ExamplePrivateKeyToWifString example using PrivateKeyToWifString()
func ExamplePrivateKeyToWifString() {
	privateWif, err := PrivateKeyToWifString(testPrivateKeyHex)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("converted wif: %s", privateWif)

	// Output:converted wif: 5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei
}

// BenchmarkPrivateKeyToWifString benchmarks the method PrivateKeyToWifString()
func BenchmarkPrivateKeyToWifString(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for b.Loop() {
		_, _ = PrivateKeyToWifString(key)
	}
}

// TestWifToPrivateKey will test the method WifToPrivateKey()
func TestWifToPrivateKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectedKey   string
		expectedNil   bool
		expectedError bool
	}{
		{caseEmpty, "", "", true, true},
		{caseSingleZero, "0", "", true, true},
		{"zeros wif", "5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAbuatmU", "0000000000000000000000000000000000000000000000000000000000000000", false, false},
		{"ascii key wif", "5HpHagT65TZzG1PH3CSu63k8DbuTZnNJf6HgyQNymvXmALAsm9s", "0000000000000000000000000000000000006d792070726976617465206b6579", false, false},
		{"invalid wif", "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8azz", "", true, true},
		{"valid uncompressed wif", testUncompressedWIF, testPrivateKeyHex, false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			privateKey, err := WifToPrivateKey(test.input)
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

// ExampleWifToPrivateKey example using WifToPrivateKey()
func ExampleWifToPrivateKey() {
	privateKey, err := WifToPrivateKey(testUncompressedWIF)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("private key: %s", hex.EncodeToString(privateKey.Serialize()))

	// Output:private key: 54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd
}

// BenchmarkWifToPrivateKey benchmarks the method WifToPrivateKey()
func BenchmarkWifToPrivateKey(b *testing.B) {
	for b.Loop() {
		_, _ = WifToPrivateKey(testUncompressedWIF)
	}
}

// TestWifToPrivateKeyString will test the method WifToPrivateKeyString()
func TestWifToPrivateKeyString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectedKey   string
		expectedError bool
	}{
		{caseEmpty, "", "", true},
		{caseSingleZero, "0", "", true},
		{"zeros wif", "5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAbuatmU", "0000000000000000000000000000000000000000000000000000000000000000", false},
		{"ascii key wif", "5HpHagT65TZzG1PH3CSu63k8DbuTZnNJf6HgyQNymvXmALAsm9s", "0000000000000000000000000000000000006d792070726976617465206b6579", false},
		{"invalid wif", "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8azz", "", true},
		{"valid uncompressed wif", testUncompressedWIF, testPrivateKeyHex, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			privateKey, err := WifToPrivateKeyString(test.input)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.expectedKey, privateKey)
		})
	}
}

// ExampleWifToPrivateKeyString example using WifToPrivateKeyString()
func ExampleWifToPrivateKeyString() {
	privateKey, err := WifToPrivateKeyString(testUncompressedWIF)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("private key: %s", privateKey)

	// Output:private key: 54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd
}

// BenchmarkWifToPrivateKeyString benchmarks the method WifToPrivateKeyString()
func BenchmarkWifToPrivateKeyString(b *testing.B) {
	for b.Loop() {
		_, _ = WifToPrivateKeyString(testUncompressedWIF)
	}
}

// TestCreateWif will test the method CreateWif()
func TestCreateWif(t *testing.T) {
	t.Run("TestCreateWif", func(t *testing.T) {
		t.Parallel()

		// Create a WIF
		wifKey, err := CreateWif()
		require.NoError(t, err)
		require.NotNil(t, wifKey)
		// t.Log("WIF:", wifKey.String())
		require.Lenf(t, wifKey.String(), 51, "WIF should be 51 characters long, got: %d", len(wifKey.String()))
	})

	t.Run("TestWifToPrivateKey", func(t *testing.T) {
		t.Parallel()

		// Create a WIF
		wifKey, err := CreateWif()
		require.NoError(t, err)
		require.NotNil(t, wifKey)
		// t.Log("WIF:", wifKey.String())
		require.Lenf(t, wifKey.String(), 51, "WIF should be 51 characters long, got: %d", len(wifKey.String()))

		// Convert WIF to Private Key
		var privateKey *ec.PrivateKey
		privateKey, err = WifToPrivateKey(wifKey.String())
		require.NoError(t, err)
		require.NotNil(t, privateKey)
		privateKeyString := hex.EncodeToString(privateKey.Serialize())
		// t.Log("Private Key:", privateKeyString)
		require.Lenf(t, privateKeyString, 64, "Private Key should be 64 characters long, got: %d", len(privateKeyString))
	})
}

// ExampleCreateWif example using CreateWif()
func ExampleCreateWif() {
	wifKey, err := CreateWif()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("WIF Key Generated Length:", len(wifKey.String()))
	// Output: WIF Key Generated Length: 51
}

// BenchmarkCreateWif benchmarks the method CreateWif()
func BenchmarkCreateWif(b *testing.B) {
	for b.Loop() {
		_, _ = CreateWif()
	}
}

// TestCreateWifString will test the method CreateWifString()
func TestCreateWifString(t *testing.T) {
	t.Run("TestCreateWifString", func(t *testing.T) {
		t.Parallel()

		// Create a WIF
		wifKey, err := CreateWifString()
		require.NoError(t, err)
		require.NotNil(t, wifKey)
		// t.Log("WIF:", wifKey)
		require.Lenf(t, wifKey, 51, "WIF should be 51 characters long, got: %d", len(wifKey))
	})

	t.Run("TestWifToPrivateKeyString", func(t *testing.T) {
		t.Parallel()

		// Create a WIF
		wifKey, err := CreateWifString()
		require.NoError(t, err)
		require.NotNil(t, wifKey)
		// t.Log("WIF:", wifKey)
		require.Lenf(t, wifKey, 51, "WIF should be 51 characters long, got: %d", len(wifKey))

		// Convert WIF to Private Key
		var privateKeyString string
		privateKeyString, err = WifToPrivateKeyString(wifKey)
		require.NoError(t, err)
		require.NotNil(t, privateKeyString)
		// t.Log("Private Key:", privateKeyString)
		require.Lenf(t, privateKeyString, 64, "Private Key should be 64 characters long, got: %d", len(privateKeyString))
	})
}

// ExampleCreateWifString example using CreateWifString()
func ExampleCreateWifString() {
	wifKey, err := CreateWifString()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("WIF Key Generated Length:", len(wifKey))
	// Output: WIF Key Generated Length: 51
}

// BenchmarkCreateWifString benchmarks the method CreateWifString()
func BenchmarkCreateWifString(b *testing.B) {
	for b.Loop() {
		_, _ = CreateWifString()
	}
}

// TestWifFromString will test the method WifFromString()
func TestWifFromString(t *testing.T) {
	t.Run("TestCreateWifFromPrivateKey", func(t *testing.T) {
		t.Parallel()

		// Create a Private Key
		privateKey, err := CreatePrivateKeyString()
		require.NoError(t, err)
		require.NotNil(t, privateKey)

		// Create a WIF
		var wifKey *WIF
		wifKey, err = PrivateKeyToWif(privateKey)
		require.NoError(t, err)
		require.NotNil(t, wifKey)
		wifKeyString := wifKey.String()
		t.Log("WIF:", wifKeyString)
		require.Lenf(t, wifKeyString, 51, "WIF should be 51 characters long, got: %d", len(wifKeyString))

		// Convert WIF to Private Key
		var privateKeyString string
		privateKeyString, err = WifToPrivateKeyString(wifKeyString)
		require.NoError(t, err)
		require.NotNil(t, privateKeyString)
		t.Log("Private Key:", privateKeyString)
		require.Lenf(t, privateKeyString, 64, "Private Key should be 64 characters long, got: %d", len(privateKeyString))

		// Compare Private Keys
		require.Equalf(t, privateKey, privateKeyString, "Private Key should be equal, got: %s", privateKeyString)

		// Decode WIF
		var decodedWif *WIF
		decodedWif, err = WifFromString(wifKeyString)
		require.NoError(t, err)
		require.NotNil(t, decodedWif)
		require.Equalf(t, wifKeyString, decodedWif.String(), "WIF should be equal, got: %s", decodedWif.String())
	})

	t.Run("TestWifFromStringMissingWIF", func(t *testing.T) {
		t.Parallel()

		_, err := WifFromString("")
		require.Error(t, err)
		require.Equal(t, ErrWifMissing, err)
	})

	t.Run("TestWifFromStringInvalidWIF", func(t *testing.T) {
		t.Parallel()

		_, err := WifFromString("invalid")
		require.Error(t, err)
		require.Equal(t, "malformed private key", err.Error())
	})
}

// ExampleWifFromString example using WifFromString()
func ExampleWifFromString() {
	// Create a Private Key
	privateKey, err := CreatePrivateKeyString()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Private Key Generated Length:", len(privateKey))

	// Create a WIF
	var wifKey *WIF
	wifKey, err = PrivateKeyToWif(privateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("WIF Key Generated Length:", len(wifKey.String()))

	// Decode WIF
	var decodedWif *WIF
	decodedWif, err = WifFromString(wifKey.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("WIF Key Decoded Length:", len(decodedWif.String()))
	// Output: Private Key Generated Length: 64
	// WIF Key Generated Length: 51
	// WIF Key Decoded Length: 51
}

// BenchmarkWifFromString benchmarks the method WifFromString()
func BenchmarkWifFromString(b *testing.B) {
	wifKey, _ := CreateWif()
	wifString := wifKey.String()
	for b.Loop() {
		_, _ = WifFromString(wifString)
	}
}
