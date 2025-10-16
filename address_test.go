package bitcoin

import (
	"fmt"
	"testing"

	"github.com/libsv/go-bk/bec"
	"github.com/stretchr/testify/assert"
)

// TestValidA58 will test the method ValidA58()
func TestValidA58(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		input         string
		expectedValid bool
		expectedError bool
	}{
		{"1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi2", true, false},
		{"1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi", false, false},
		{"1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi", false, true},
		{"1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi", false, true},
		{"1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi", false, true},
		{"1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi", false, true},
		{"1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi", false, true},
		{"1KCEAmV", false, false},
		{"", false, false},
		{"0", false, true},
	}

	for _, test := range tests {
		if valid, err := ValidA58([]byte(test.input)); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if valid && !test.expectedValid {
			t.Fatalf("%s Failed: [%s] inputted and was valid but should NOT be valid", t.Name(), test.input)
		} else if !valid && test.expectedValid {
			t.Fatalf("%s Failed: [%s] inputted and was invalid but should be valid", t.Name(), test.input)
		}
	}
}

// ExampleValidA58 example using ValidA58()
func ExampleValidA58() {
	valid, err := ValidA58([]byte("1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi2"))
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	} else if !valid {
		fmt.Printf("address is not valid: %s", "1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi2")
		return
	} else {
		fmt.Printf("address is valid!")
	}
	// Output:address is valid!
}

// BenchmarkValidA58 benchmarks the method ValidA58()
func BenchmarkValidA58(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ValidA58([]byte("1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi2"))
	}
}

// TestGetAddressFromPrivateKey will test the method GetAddressFromPrivateKey()
func TestGetAddressFromPrivateKey(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		input           string
		expectedAddress string
		compressed      bool
		expectedError   bool
	}{
		{"0", "", true, true},
		{"00000", "", true, true},
		{"12345678", "1BHxe5Yw72oYoV8tFjySYrV9Y2JwMpAZEy", true, false},
		{"54035dd4c7dda99ac473905a3d82", "1L5GmmuGeS3HwoEDv7zkWcheayXrRsurUm", true, false},
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9", "13dnka5SaugRchayN84EED7a2E8dCNMLXQ", true, false},
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", "1DfGxKmgL3ETwUdNnXLBueEvNpjcDGcKgK", true, false},
	}

	for _, test := range tests {
		if address, err := GetAddressFromPrivateKeyString(test.input, test.compressed, true); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if address != test.expectedAddress {
			t.Fatalf("%s Failed: [%s] inputted and [%s] expected, but got: %s", t.Name(), test.input, test.expectedAddress, address)
		}
	}
}

// TestGetAddressFromPrivateKeyCompression will test the method GetAddressFromPrivateKey()
func TestGetAddressFromPrivateKeyCompression(t *testing.T) {

	privateKey, err := bec.NewPrivateKey(bec.S256())
	assert.NoError(t, err)

	var addressUncompressed string
	addressUncompressed, err = GetAddressFromPrivateKey(privateKey, false, true)
	assert.NoError(t, err)

	var addressCompressed string
	addressCompressed, err = GetAddressFromPrivateKey(privateKey, true, true)
	assert.NoError(t, err)

	assert.NotEqual(t, addressCompressed, addressUncompressed)

	addressCompressed, err = GetAddressFromPrivateKey(&bec.PrivateKey{}, true, true)
	assert.Error(t, err)
	assert.Equal(t, "", addressCompressed)
}

// ExampleGetAddressFromPrivateKey example using GetAddressFromPrivateKey()
func ExampleGetAddressFromPrivateKey() {
	address, err := GetAddressFromPrivateKeyString("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", true, true)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("address found: %s", address)
	// Output:address found: 1DfGxKmgL3ETwUdNnXLBueEvNpjcDGcKgK
}

// BenchmarkGetAddressFromPrivateKey benchmarks the method GetAddressFromPrivateKey()
func BenchmarkGetAddressFromPrivateKey(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for i := 0; i < b.N; i++ {
		_, _ = GetAddressFromPrivateKeyString(key, true, true)
	}
}

// testGetPublicKeyFromPrivateKey is a helper method for tests
func testGetPublicKeyFromPrivateKey(privateKey string) *bec.PublicKey {
	rawKey, err := PrivateKeyFromString(privateKey)
	if err != nil {
		return nil
	}
	return rawKey.PubKey()
}

// TestGetAddressFromPubKey will test the method GetAddressFromPubKey()
func TestGetAddressFromPubKey(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		input           *bec.PublicKey
		expectedAddress string
		expectedNil     bool
		expectedError   bool
	}{
		{&bec.PublicKey{}, "", true, true},
		{testGetPublicKeyFromPrivateKey("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"), "1DfGxKmgL3ETwUdNnXLBueEvNpjcDGcKgK", false, false},
		{testGetPublicKeyFromPrivateKey("000000"), "15wJjXvfQzo3SXqoWGbWZmNYND1Si4siqV", false, false},
		{testGetPublicKeyFromPrivateKey("0"), "15wJjXvfQzo3SXqoWGbWZmNYND1Si4siqV", true, true},
	}

	// todo: add more error cases of invalid *bec.PublicKey

	for _, test := range tests {
		if rawKey, err := GetAddressFromPubKey(test.input, true, true); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%v] inputted and error was expected", t.Name(), test.input)
		} else if rawKey == nil && !test.expectedNil {
			t.Fatalf("%s Failed: [%v] inputted and was nil but not expected", t.Name(), test.input)
		} else if rawKey != nil && test.expectedNil {
			t.Fatalf("%s Failed: [%v] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		} else if rawKey != nil && rawKey.AddressString != test.expectedAddress {
			t.Fatalf("%s Failed: [%v] inputted [%s] expected but failed comparison of addresses, got: %s", t.Name(), test.input, test.expectedAddress, rawKey.AddressString)
		}
	}
}

// ExampleGetAddressFromPubKey example using GetAddressFromPubKey()
func ExampleGetAddressFromPubKey() {
	rawAddress, err := GetAddressFromPubKey(testGetPublicKeyFromPrivateKey("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"), true, true)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("address found: %s", rawAddress.AddressString)
	// Output:address found: 1DfGxKmgL3ETwUdNnXLBueEvNpjcDGcKgK
}

// BenchmarkGetAddressFromPubKey benchmarks the method GetAddressFromPubKey()
func BenchmarkGetAddressFromPubKey(b *testing.B) {
	pubKey := testGetPublicKeyFromPrivateKey("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	for i := 0; i < b.N; i++ {
		_, _ = GetAddressFromPubKey(pubKey, true, true)
	}
}

// TestGetAddressFromScript will test the method GetAddressFromScript()
func TestGetAddressFromScript(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		inputScript     string
		expectedAddress string
		expectedError   bool
	}{
		{"", "", true},
		{"0", "", true},
		{"76a9141a9d62736746f85ca872dc555ff51b1fed2471e288ac", "13Rj7G3pn2GgG8KE6SFXLc7dCJdLNnNK7M", false},
		{"76a914b424110292f4ea2ac92beb9e83cf5e6f0fa2996388ac", "1HRVqUGDzpZSMVuNSZxJVaB9xjneEShfA7", false},
		{"76a914b424110292f4ea2ac92beb9e83cf5e6f0fa2", "", true},
		{"76a914b424110292f4ea2ac92beb9e83", "", true},
		{"76a914b424110292f", "", true},
		{"1HRVqUGDzpZSMVuNSZxJVaB9xjneEShfA7", "", true},
		{"514104cc71eb30d653c0c3163990c47b976f3fb3f37cccdcbedb169a1dfef58bbfbfaff7d8a473e7e2e6d317b87bafe8bde97e3cf8f065dec022b51d11fcdd0d348ac4410461cbdcc5409fb4b4d42b51d33381354d80e550078cb532a34bfa2fcfdeb7d76519aecc62770f5b0e4ef8551946d8a540911abe3e7854a26f39f58b25c15342af52ae", "", true},
		{"410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3", "", true},
		{"47304402204e45e16932b8af514961a1d3a1a25fdf3f4f7732e9d624c6c61548ab5fb8cd410220181522ec8eca07de4860a4acdd12909d831cc56cbbac4622082221a8768d1d0901", "", true},
	}

	for _, test := range tests {
		if address, err := GetAddressFromScript(test.inputScript); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.inputScript, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%v] inputted and error was expected", t.Name(), test.inputScript)
		} else if address != test.expectedAddress {
			t.Fatalf("%s Failed: [%v] inputted [%s] expected but failed comparison of addresses, got: %s", t.Name(), test.inputScript, test.expectedAddress, address)
		}
	}
}

// ExampleGetAddressFromScript example using GetAddressFromScript()
func ExampleGetAddressFromScript() {
	address, err := GetAddressFromScript("76a914b424110292f4ea2ac92beb9e83cf5e6f0fa2996388ac")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("address found: %s", address)
	// Output:address found: 1HRVqUGDzpZSMVuNSZxJVaB9xjneEShfA7
}

// BenchmarkAddressFromScript benchmarks the method GetAddressFromScript()
func BenchmarkGetAddressFromScript(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetAddressFromScript("76a914b424110292f4ea2ac92beb9e83cf5e6f0fa2996388ac")
	}
}

// TestGetAddressFromPubKeyString will test the method GetAddressFromPubKeyString()
func TestGetAddressFromPubKeyString(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		input           string
		expectedAddress string
		expectedNil     bool
		expectedError   bool
	}{
		{"", "", true, true},
		{"0", "", true, true},
		{"03ce8a73eb5e4d45966d719ac3ceb431cd0ee203e6395357a167b9abebc4baeacf", "17HeHWVDqDqexLJ31aG4qtVMoX8pKMGSuJ", false, false},
		{"0000", "", true, true},
	}

	for _, test := range tests {
		if rawKey, err := GetAddressFromPubKeyString(test.input, true, true); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%v] inputted and error was expected", t.Name(), test.input)
		} else if rawKey == nil && !test.expectedNil {
			t.Fatalf("%s Failed: [%v] inputted and was nil but not expected", t.Name(), test.input)
		} else if rawKey != nil && test.expectedNil {
			t.Fatalf("%s Failed: [%v] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		} else if rawKey != nil && rawKey.AddressString != test.expectedAddress {
			t.Fatalf("%s Failed: [%v] inputted [%s] expected but failed comparison of addresses, got: %s", t.Name(), test.input, test.expectedAddress, rawKey.AddressString)
		}
	}
}

// ExampleGetAddressFromPubKeyString example using GetAddressFromPubKeyString()
func ExampleGetAddressFromPubKeyString() {
	rawAddress, err := GetAddressFromPubKeyString("03ce8a73eb5e4d45966d719ac3ceb431cd0ee203e6395357a167b9abebc4baeacf", true, true)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("address found: %s", rawAddress.AddressString)
	// Output:address found: 17HeHWVDqDqexLJ31aG4qtVMoX8pKMGSuJ
}

// BenchmarkGetAddressFromPubKeyString benchmarks the method GetAddressFromPubKeyString()
func BenchmarkGetAddressFromPubKeyString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetAddressFromPubKeyString("03ce8a73eb5e4d45966d719ac3ceb431cd0ee203e6395357a167b9abebc4baeacf", true, true)
	}
}
