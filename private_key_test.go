package bitcoin

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bitcoinsv/bsvd/bsvec"
)

// TestGenerateSharedKeyPair will test creating a shared key that can be used to encrypt data that can be decrypted by yourself (privkey) and also the owner of the given public key
func TestGenerateSharedKeyPair(t *testing.T) {

	// The data that will be encrypted / shared
	testString := "testing 1, 2, 3..."

	// User 1
	privKey1, _ := CreatePrivateKey()

	// User 2
	privKey2, _ := CreatePrivateKey()

	_, user1SharedPubKey := GenerateSharedKeyPair(privKey1, privKey2.PubKey())

	// encrypt something with the shared public key
	eciesTest, err := bsvec.Encrypt(user1SharedPubKey, []byte(testString))
	if err != nil {
		t.Errorf("Failed to encrypt test data %s", err)
	}

	// user 2 decrypts it
	user2SharedPrivKey, _ := GenerateSharedKeyPair(privKey2, privKey1.PubKey())

	decryptedTestData, err := bsvec.Decrypt(user2SharedPrivKey, eciesTest)
	if err != nil {
		t.Errorf("Failed to decrypt test data %s", err)
	}

	if string(decryptedTestData) != testString {
		t.Errorf("Decrypted string doesnt match %s", decryptedTestData)
	}
}

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

// TestPrivateKeyToWif will test the method PrivateKeyToWif()
func TestPrivateKeyToWif(t *testing.T) {

	t.Parallel()

	// Create the list of tests
	var tests = []struct {
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

	// Run tests
	for _, test := range tests {
		if wif, err := PrivateKeyToWif(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if wif == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was nil but not expected", t.Name(), test.input)
		} else if wif != nil && test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		} else if wif != nil && wif.String() != test.expectedWif {
			t.Errorf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedWif, wif.String())
		}
	}

}

// ExamplePrivateKeyToWif example using PrivateKeyToWif()
func ExamplePrivateKeyToWif() {
	wif, err := PrivateKeyToWif("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("converted wif: %s", wif.String())

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

	// Create the list of tests
	var tests = []struct {
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

	// Run tests
	for _, test := range tests {
		if wif, err := PrivateKeyToWifString(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if wif != test.expectedWif {
			t.Errorf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedWif, wif)
		}
	}

}

// ExamplePrivateKeyToWifString example using PrivateKeyToWifString()
func ExamplePrivateKeyToWifString() {
	wif, err := PrivateKeyToWifString("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("converted wif: %s", wif)

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

	// Create the list of tests
	var tests = []struct {
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

	// Run tests
	for _, test := range tests {
		if privateKey, err := WifToPrivateKey(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if privateKey == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was nil but not expected", t.Name(), test.input)
		} else if privateKey != nil && test.expectedNil {
			t.Errorf("%s Failed: [%s] inputted and was NOT nil but expected to be nil", t.Name(), test.input)
		} else if privateKey != nil && hex.EncodeToString(privateKey.Serialize()) != test.expectedKey {
			t.Errorf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedKey, hex.EncodeToString(privateKey.Serialize()))
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
	fmt.Printf("private key: %s", hex.EncodeToString(privateKey.Serialize()))

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

	// Create the list of tests
	var tests = []struct {
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

	// Run tests
	for _, test := range tests {
		if privateKey, err := WifToPrivateKeyString(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		} else if privateKey != test.expectedKey {
			t.Errorf("%s Failed: [%s] inputted [%s] expected but failed comparison of keys, got: %s", t.Name(), test.input, test.expectedKey, privateKey)
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
