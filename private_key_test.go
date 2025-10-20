package bitcoin

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/wif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreatePrivateKey will test the method CreatePrivateKey()
func TestCreatePrivateKey(t *testing.T) {
	rawKey, err := CreatePrivateKey()
	require.NoError(t, err)
	assert.NotNil(t, rawKey)
	assert.Len(t, rawKey.Serialise(), 32) //nolint:misspell // external library method name
}

// ExampleCreatePrivateKey example using CreatePrivateKey()
func ExampleCreatePrivateKey() {
	rawKey, err := CreatePrivateKey()
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	} else if len(rawKey.Serialise()) > 0 { //nolint:misspell // external library method name
		fmt.Printf("key created successfully!")
	}
	// Output:key created successfully!
}

// BenchmarkCreatePrivateKey benchmarks the method CreatePrivateKey()
func BenchmarkCreatePrivateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < b.N; i++ {
		_, _ = CreatePrivateKeyString()
	}
}

// TestPrivateKeyFromString will test the method PrivateKeyFromString()
func TestPrivateKeyFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input         string
		expectedKey   string
		expectedNil   bool
		expectedError bool
	}{
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", false, false},
		{"E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F2B75", "e83385af76b2b1997326b567461fb73dd9c27eab9e1e86d26779f4650c5f2b75", false, false},
		{"E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F", "0000e83385af76b2b1997326b567461fb73dd9c27eab9e1e86d26779f4650c5f", false, false},
		{"E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F", "", true, true},
		{"1234567", "", true, true},
		{"0", "", true, true},
		{"", "", true, true},
	}

	for _, test := range tests {
		rawKey, err := PrivateKeyFromString(test.input)
		if err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		}
		if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		}
		if rawKey == nil && !test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and was nil but not expected", t.Name(), test.input)
		}
		if rawKey != nil && test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		}
		if rawKey != nil && hex.EncodeToString(rawKey.Serialise()) != test.expectedKey { //nolint:misspell // external library method name
			t.Fatalf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedKey, hex.EncodeToString(rawKey.Serialise())) //nolint:misspell // external library method name
		}
	}
}

// ExamplePrivateKeyFromString example using PrivateKeyFromString()
func ExamplePrivateKeyFromString() {
	key, err := PrivateKeyFromString("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("key converted: %s", hex.EncodeToString(key.Serialise())) //nolint:misspell // external library method name
	// Output:key converted: 54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd
}

// BenchmarkPrivateKeyFromString benchmarks the method PrivateKeyFromString()
func BenchmarkPrivateKeyFromString(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for i := 0; i < b.N; i++ {
		_, _ = PrivateKeyFromString(key)
	}
}

// TestPrivateAndPublicKeys will test the method PrivateAndPublicKeys()
func TestPrivateAndPublicKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input              string
		expectedPrivateKey string
		expectedNil        bool
		expectedError      bool
	}{
		{"", "", true, true},
		{"0", "", true, true},
		{"00000", "", true, true},
		{"0-0-0-0-0", "", true, true},
		{"z4035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abz", "", true, true},
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", false, false},
	}

	for _, test := range tests {
		privateKey, publicKey, err := PrivateAndPublicKeys(test.input)
		if err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		}
		if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		}
		if (privateKey == nil || publicKey == nil) && !test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and was nil but not expected", t.Name(), test.input)
		}
		if (privateKey != nil || publicKey != nil) && test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		}
		if privateKey != nil && hex.EncodeToString(privateKey.Serialise()) != test.expectedPrivateKey { //nolint:misspell // external library method name
			t.Fatalf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedPrivateKey, hex.EncodeToString(privateKey.Serialise())) //nolint:misspell // external library method name
		}
	}
}

// ExamplePrivateAndPublicKeys example using PrivateAndPublicKeys()
func ExamplePrivateAndPublicKeys() {
	privateKey, publicKey, err := PrivateAndPublicKeys("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("private key: %s public key: %s", hex.EncodeToString(privateKey.Serialise()), hex.EncodeToString(publicKey.SerialiseCompressed())) //nolint:misspell // Serialise is the actual method name

	// Output:private key: 54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd public key: 031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f
}

// BenchmarkPrivateAndPublicKeys benchmarks the method PrivateAndPublicKeys()
func BenchmarkPrivateAndPublicKeys(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for i := 0; i < b.N; i++ {
		_, _, _ = PrivateAndPublicKeys(key)
	}
}

// TestPrivateKeyToWif will test the method PrivateKeyToWif()
func TestPrivateKeyToWif(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input         string
		expectedWif   string
		expectedNil   bool
		expectedError bool
	}{
		{"", "", true, true},
		{"0", "", true, true},
		{"000000", "5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAbuatmU", false, false},
		{"6D792070726976617465206B6579", "5HpHagT65TZzG1PH3CSu63k8DbuTZnNJf6HgyQNymvXmALAsm9s", false, false},
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8azz", "", true, true},
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", "5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei", false, false},
	}

	for _, test := range tests {
		privateWif, err := PrivateKeyToWif(test.input)
		if err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		}
		if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		}
		if privateWif == nil && !test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and was nil but not expected", t.Name(), test.input)
		}
		if privateWif != nil && test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		}
		if privateWif != nil && privateWif.String() != test.expectedWif {
			t.Fatalf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedWif, privateWif.String())
		}
	}
}

// ExamplePrivateKeyToWif example using PrivateKeyToWif()
func ExamplePrivateKeyToWif() {
	privateWif, err := PrivateKeyToWif("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
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
	for i := 0; i < b.N; i++ {
		_, _ = PrivateKeyToWif(key)
	}
}

// TestPrivateKeyToWifString will test the method PrivateKeyToWifString()
func TestPrivateKeyToWifString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input         string
		expectedWif   string
		expectedError bool
	}{
		{"", "", true},
		{"0", "", true},
		{"000000", "5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAbuatmU", false},
		{"6D792070726976617465206B6579", "5HpHagT65TZzG1PH3CSu63k8DbuTZnNJf6HgyQNymvXmALAsm9s", false},
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8azz", "", true},
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", "5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei", false},
	}

	for _, test := range tests {
		if privateWif, err := PrivateKeyToWifString(test.input); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if privateWif != test.expectedWif {
			t.Fatalf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedWif, privateWif)
		}
	}
}

// ExamplePrivateKeyToWifString example using PrivateKeyToWifString()
func ExamplePrivateKeyToWifString() {
	privateWif, err := PrivateKeyToWifString("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
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
	for i := 0; i < b.N; i++ {
		_, _ = PrivateKeyToWifString(key)
	}
}

// TestWifToPrivateKey will test the method WifToPrivateKey()
func TestWifToPrivateKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input         string
		expectedKey   string
		expectedNil   bool
		expectedError bool
	}{
		{"", "", true, true},
		{"0", "", true, true},
		{"5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAbuatmU", "0000000000000000000000000000000000000000000000000000000000000000", false, false},
		{"5HpHagT65TZzG1PH3CSu63k8DbuTZnNJf6HgyQNymvXmALAsm9s", "0000000000000000000000000000000000006d792070726976617465206b6579", false, false},
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8azz", "", true, true},
		{"5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei", "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", false, false},
	}

	for _, test := range tests {
		privateKey, err := WifToPrivateKey(test.input)
		if err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		}
		if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		}
		if privateKey == nil && !test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and was nil but not expected", t.Name(), test.input)
		}
		if privateKey != nil && test.expectedNil {
			t.Fatalf("%s Failed: [%s] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		}
		if privateKey != nil && hex.EncodeToString(privateKey.Serialise()) != test.expectedKey { //nolint:misspell // Serialise is the actual method name
			t.Fatalf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedKey, hex.EncodeToString(privateKey.Serialise())) //nolint:misspell // Serialise is the actual method name
		}
	}
}

// ExampleWifToPrivateKey example using WifToPrivateKey()
func ExampleWifToPrivateKey() {
	privateKey, err := WifToPrivateKey("5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("private key: %s", hex.EncodeToString(privateKey.Serialise())) //nolint:misspell // Serialise is the actual method name

	// Output:private key: 54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd
}

// BenchmarkWifToPrivateKey benchmarks the method WifToPrivateKey()
func BenchmarkWifToPrivateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = WifToPrivateKey("5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei")
	}
}

// TestWifToPrivateKeyString will test the method WifToPrivateKeyString()
func TestWifToPrivateKeyString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input         string
		expectedKey   string
		expectedError bool
	}{
		{"", "", true},
		{"0", "", true},
		{"5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAbuatmU", "0000000000000000000000000000000000000000000000000000000000000000", false},
		{"5HpHagT65TZzG1PH3CSu63k8DbuTZnNJf6HgyQNymvXmALAsm9s", "0000000000000000000000000000000000006d792070726976617465206b6579", false},
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8azz", "", true},
		{"5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei", "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", false},
	}

	for _, test := range tests {
		if privateKey, err := WifToPrivateKeyString(test.input); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if privateKey != test.expectedKey {
			t.Fatalf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedKey, privateKey)
		}
	}
}

// ExampleWifToPrivateKeyString example using WifToPrivateKeyString()
func ExampleWifToPrivateKeyString() {
	privateKey, err := WifToPrivateKeyString("5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("private key: %s", privateKey)

	// Output:private key: 54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd
}

// BenchmarkWifToPrivateKeyString benchmarks the method WifToPrivateKeyString()
func BenchmarkWifToPrivateKeyString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = WifToPrivateKeyString("5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei")
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
		var privateKey *bec.PrivateKey
		privateKey, err = WifToPrivateKey(wifKey.String())
		require.NoError(t, err)
		require.NotNil(t, privateKey)
		privateKeyString := hex.EncodeToString(privateKey.Serialise()) //nolint:misspell // external library method name
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
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < b.N; i++ {
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
		var wifKey *wif.WIF
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
		var decodedWif *wif.WIF
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
	var wifKey *wif.WIF
	wifKey, err = PrivateKeyToWif(privateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("WIF Key Generated Length:", len(wifKey.String()))

	// Decode WIF
	var decodedWif *wif.WIF
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
	for i := 0; i < b.N; i++ {
		_, _ = WifFromString(wifString)
	}
}
