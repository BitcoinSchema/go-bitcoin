package bitcoin

import (
	"fmt"
	"testing"

	"github.com/bitcoinsv/bsvd/bsvec"
)

// TestValidA58 will test the method ValidA58()
func TestValidA58(t *testing.T) {

	t.Parallel()

	// Create the list of tests
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

	// Run tests
	for _, test := range tests {
		if valid, err := ValidA58([]byte(test.input)); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if valid && !test.expectedValid {
			t.Errorf("%s Failed: [%s] inputted and was valid but should NOT be valid", t.Name(), test.input)
		} else if !valid && test.expectedValid {
			t.Errorf("%s Failed: [%s] inputted and was invalid but should be valid", t.Name(), test.input)
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

// TestAddressFromPrivateKey will test the method AddressFromPrivateKey()
func TestAddressFromPrivateKey(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		input           string
		expectedAddress string
		expectedError   bool
	}{
		{"0", "", true},
		{"00000", "", true},
		{"12345678", "1BHxe5Yw72oYoV8tFjySYrV9Y2JwMpAZEy", false},
		{"54035dd4c7dda99ac473905a3d82", "1L5GmmuGeS3HwoEDv7zkWcheayXrRsurUm", false},
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9", "13dnka5SaugRchayN84EED7a2E8dCNMLXQ", false},
		{"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd", "1DfGxKmgL3ETwUdNnXLBueEvNpjcDGcKgK", false},
	}

	// Run tests
	for _, test := range tests {
		if address, err := AddressFromPrivateKey(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if address != test.expectedAddress {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, but got: %s", t.Name(), test.input, test.expectedAddress, address)
		}
	}
}

// ExampleAddressFromPrivateKey example using AddressFromPrivateKey()
func ExampleAddressFromPrivateKey() {
	address, err := AddressFromPrivateKey("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("address found: %s", address)
	// Output:address found: 1DfGxKmgL3ETwUdNnXLBueEvNpjcDGcKgK
}

// BenchmarkAddressFromPrivateKey benchmarks the method AddressFromPrivateKey()
func BenchmarkAddressFromPrivateKey(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for i := 0; i < b.N; i++ {
		_, _ = AddressFromPrivateKey(key)
	}
}

// testGetPublicKeyFromPrivateKey is a helper method for tests
func testGetPublicKeyFromPrivateKey(privateKey string) *bsvec.PublicKey {
	rawKey, err := PrivateKeyFromString(privateKey)
	if err != nil {
		return nil
	}
	return rawKey.PubKey()
}

// TestAddressFromPubKey will test the method AddressFromPubKey()
func TestAddressFromPubKey(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		input           *bsvec.PublicKey
		expectedAddress string
		expectedNil     bool
		expectedError   bool
	}{
		{&bsvec.PublicKey{}, "", true, true},
		{testGetPublicKeyFromPrivateKey("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"), "1DfGxKmgL3ETwUdNnXLBueEvNpjcDGcKgK", false, false},
		{testGetPublicKeyFromPrivateKey("000000"), "15wJjXvfQzo3SXqoWGbWZmNYND1Si4siqV", false, false},
		{testGetPublicKeyFromPrivateKey("0"), "15wJjXvfQzo3SXqoWGbWZmNYND1Si4siqV", true, true},
	}

	// todo: add more error cases of invalid *bsvec.PublicKey

	// Run tests
	for _, test := range tests {
		if rawKey, err := AddressFromPubKey(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), test.input)
		} else if rawKey == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and was nil but not expected", t.Name(), test.input)
		} else if rawKey != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		} else if rawKey != nil && rawKey.EncodeAddress() != test.expectedAddress {
			t.Errorf("%s Failed: [%v] inputted [%s] expected but failed comparison of addresses, got: %s", t.Name(), test.input, test.expectedAddress, rawKey.EncodeAddress())
		}
	}
}

// ExampleAddressFromPubKey example using AddressFromPubKey()
func ExampleAddressFromPubKey() {
	rawAddress, err := AddressFromPubKey(testGetPublicKeyFromPrivateKey("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"))
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("address found: %s", rawAddress.EncodeAddress())
	// Output:address found: 1DfGxKmgL3ETwUdNnXLBueEvNpjcDGcKgK
}

// BenchmarkAddressFromPubKey benchmarks the method AddressFromPubKey()
func BenchmarkAddressFromPubKey(b *testing.B) {
	pubKey := testGetPublicKeyFromPrivateKey("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	for i := 0; i < b.N; i++ {
		_, _ = AddressFromPubKey(pubKey)
	}
}
