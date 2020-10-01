package bitcoin

import (
	"encoding/hex"
	"fmt"
	"testing"
)

// TestCreatePrivateKey will test the method CreatePrivateKey()
func TestCreatePrivateKey(t *testing.T) {
	rawKey, err := CreatePrivateKey()
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	if rawKey == nil {
		t.Fatalf("private key was nil")
	}
	if len(rawKey.Serialize()) == 0 {
		t.Fatalf("key length was invalid")
	}
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
	for i := 0; i < b.N; i++ {
		_, _ = CreatePrivateKey()
	}
}

// TestCreatePrivateKeyString will test the method CreatePrivateKeyString()
func TestCreatePrivateKeyString(t *testing.T) {
	key, err := CreatePrivateKeyString()
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	if len(key) == 0 {
		t.Fatalf("private key is empty")
	}
	if len(key) != 64 {
		t.Fatalf("key length is not 64")
	}
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

	// Create the list of tests
	var tests = []struct {
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

	// Run tests
	for _, test := range tests {
		if rawKey, err := PrivateKeyFromString(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if rawKey == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was nil but not expected", t.Name(), test.input)
		} else if rawKey != nil && test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		} else if rawKey != nil && hex.EncodeToString(rawKey.Serialize()) != test.expectedKey {
			t.Errorf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedKey, hex.EncodeToString(rawKey.Serialize()))
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
	fmt.Printf("key converted: %s", hex.EncodeToString(key.Serialize()))
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

	// Create the list of tests
	var tests = []struct {
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

	// Run tests
	for _, test := range tests {
		if privateKey, publicKey, err := PrivateAndPublicKeys(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if (privateKey == nil || publicKey == nil) && !test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was nil but not expected", t.Name(), test.input)
		} else if (privateKey != nil || publicKey != nil) && test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		} else if privateKey != nil && hex.EncodeToString(privateKey.Serialize()) != test.expectedPrivateKey {
			t.Errorf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedPrivateKey, hex.EncodeToString(privateKey.Serialize()))
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
	fmt.Printf("private key: %s public key: %s", hex.EncodeToString(privateKey.Serialize()), hex.EncodeToString(publicKey.SerializeCompressed()))

	// Output:private key: 54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd public key: 031b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f
}

// BenchmarkPrivateAndPublicKeys benchmarks the method PrivateAndPublicKeys()
func BenchmarkPrivateAndPublicKeys(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for i := 0; i < b.N; i++ {
		_, _, _ = PrivateAndPublicKeys(key)
	}
}
