package bitcoin

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bitcoinsv/bsvd/bsvec"
)

// TestEncryptWithPrivateKey will test the method EncryptWithPrivateKey()
func TestEncryptWithPrivateKey(t *testing.T) {
	t.Parallel()

	// Create a valid private key
	privateKey, err := CreatePrivateKey()
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Create the list of tests
	var tests = []struct {
		inputKey      *bsvec.PrivateKey
		inputData     string
		expectedError bool
	}{
		{privateKey, "", false},
		{privateKey, " ", false},
		{privateKey, "\n", false},
		{privateKey, "0", false},
		{privateKey, "-1", false},
		{privateKey, "test-data", false},
		{privateKey, `{"json":"data"}`, false},
		{privateKey, `PLXfZ8ZNN9emRVLapZlvsDpjti3TazobtrY1LWibG5ogf07R5u6xwTGSsBGH10E5naTJLOhj38dFFDeG7giLIvK9QbiF8xSqbXoxvILl09AeXSFfOx7154KdOQdoXMxsghUFi91Efh4uWNl53d7feec7NXh7WrRk68sYUQrfGNdvgxyJlYXBpp7pq5f9Yhhm2TxbNyTyh4dkgvXVJtVfFLwiOWYfO9PAwmfDAgoFyP1vpNcL7G2lCsDW2LuaO1PO0L8a0Z55jpqYrYwzUG3IK4HKnx61xWvawFSzapzGLMrrdxVK56WuToZPOVlLpXjGKYNXlfBXh6T4d1uYUQg9aKfrdIZWCXjVnCXrqDa8vrF0djypIPq24zE0KM0oQZfQJT2Vptu1OPJkes2nIpAUTKxtiLIiunbk8EXaRFBAgGYFFHYm3tHLQUaFQ0K6USoTN3Nr92SjgKUcXVJtSeH854fFpSAoi8je74XPLFHU54GJ0LB3ixHFmD1YQFbBJux7738gI05pxZTwl69KXoat8OfFamLhwNxg8BcA8GXpkN6i5UZW8VfvOQ11nfpd5RykehgYeFyZaewUizPTZkVfI8BT5fdBGQAGXfNT8Xo0OiRxb64rB3Q3YcKiwXjpD1gc14vZXb5EzP6y274yh9RODHPi3rZuVYeGDBd1woFHng2cFeu1ZZ1kBpB30jqqoxDvYzHWvUEpeP3mRyyL51pqYECgOoTEB9nZlb8J31sbAAyY5EXGp6mPuxnhYfR03wIPdkVaDgY9IUdBfKxiIfghTw9Zk9NYZWz1THnqe1GfnBqtHTWiHHuBt0vEJ5QNiMyvcsZbLb6djbbtVFV5sEMzRWf8cJInSEQFgcXPRFaZYow2bJiugFTGvt0ZZKscHqn2SJFNKGRgr3zlcTgF3Y33PVJVTg1uOaZw6y6mbVituSAR5LuDN0zlFuvkOvr63Hys3fE4dONIOxphiDyNG0Nvci5I3bB3E1H1gT8vUDCq`, false},
	}

	t.Log(hex.EncodeToString(privateKey.Serialize()))

	// Run tests
	var encrypted string
	for _, test := range tests {
		if encrypted, err = EncryptWithPrivateKey(test.inputKey, test.inputData); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error not expected but got: %s", t.Name(), test.inputKey, test.inputData, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error was expected", t.Name(), test.inputKey, test.inputData)
		} else if len(encrypted) == 0 {
			t.Errorf("%s Failed: [%s] [%s] inputted and expected length > 0, but got: 0", t.Name(), test.inputKey, test.inputData)
		}
		t.Log(test.inputData)
		t.Log(encrypted)
	}
}

// TestEncryptWithPrivateKeyPanic tests for nil case in EncryptWithPrivateKey()
func TestEncryptWithPrivateKeyPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()

	_, err := EncryptWithPrivateKey(nil, "")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
}

// ExampleEncryptWithPrivateKey example using EncryptWithPrivateKey()
func ExampleEncryptWithPrivateKey() {
	// Create a valid private key
	privateKey, err := CreatePrivateKey()
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Encrypt the text
	var encryptedText string
	encryptedText, err = EncryptWithPrivateKey(privateKey, "encrypt my message")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	// Can't show the real example, as the value changes each time
	// fmt.Printf("encrypted text length: %d %s", len(encryptedText), encryptedText)
	fmt.Printf("encrypted text length: %d %s", len(encryptedText), "40dd6dbad36cdfb7b1bf85a87887ac3202ca0020b6"+
		"03e6c1fbc30b9a586f06d72eb2fc5117d0f2e4c3193788ebd79c6eb095d88e00209c4ba9a883f631566dfe17d0b2dc83ee456b71927b5b"+
		"5dc6e5a8074d929f3d78c22a4bb7c5f8ffdb0895adf035b2e36fdb5504d3c3b99f3d03a9767eff1b93af14cc586efbf5e81e6fec839c55"+
		"13ab2d644bf3e56d6ddc92e2985c18dab3c6a9")
	// Output:encrypted text length: 300 40dd6dbad36cdfb7b1bf85a87887ac3202ca0020b603e6c1fbc30b9a586f06d72eb2fc5117d0f2e4c3193788ebd79c6eb095d88e00209c4ba9a883f631566dfe17d0b2dc83ee456b71927b5b5dc6e5a8074d929f3d78c22a4bb7c5f8ffdb0895adf035b2e36fdb5504d3c3b99f3d03a9767eff1b93af14cc586efbf5e81e6fec839c5513ab2d644bf3e56d6ddc92e2985c18dab3c6a9
}

// BenchmarkEncryptWithPrivateKey benchmarks the method EncryptWithPrivateKey()
func BenchmarkEncryptWithPrivateKey(b *testing.B) {
	key, _ := CreatePrivateKey()
	for i := 0; i < b.N; i++ {
		_, _ = EncryptWithPrivateKey(key, "some-data")
	}
}

// TestEncryptWithPrivateKeyString will test the method EncryptWithPrivateKeyString()
func TestEncryptWithPrivateKeyString(t *testing.T) {
	t.Parallel()

	// Create a valid private key
	privateKey, err := CreatePrivateKeyString()
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Create the list of tests
	var tests = []struct {
		inputKey      string
		inputData     string
		expectedError bool
	}{
		{"", "", true},
		{privateKey, "", false},
		{privateKey, " ", false},
		{privateKey, "\n", false},
		{privateKey, "0", false},
		{privateKey, "-1", false},
		{privateKey, "test-data", false},
		{privateKey, `{"json":"data"}`, false},
		{privateKey, `PLXfZ8ZNN9emRVLapZlvsDpjti3TazobtrY1LWibG5ogf07R5u6xwTGSsBGH10E5naTJLOhj38dFFDeG7giLIvK9QbiF8xSqbXoxvILl09AeXSFfOx7154KdOQdoXMxsghUFi91Efh4uWNl53d7feec7NXh7WrRk68sYUQrfGNdvgxyJlYXBpp7pq5f9Yhhm2TxbNyTyh4dkgvXVJtVfFLwiOWYfO9PAwmfDAgoFyP1vpNcL7G2lCsDW2LuaO1PO0L8a0Z55jpqYrYwzUG3IK4HKnx61xWvawFSzapzGLMrrdxVK56WuToZPOVlLpXjGKYNXlfBXh6T4d1uYUQg9aKfrdIZWCXjVnCXrqDa8vrF0djypIPq24zE0KM0oQZfQJT2Vptu1OPJkes2nIpAUTKxtiLIiunbk8EXaRFBAgGYFFHYm3tHLQUaFQ0K6USoTN3Nr92SjgKUcXVJtSeH854fFpSAoi8je74XPLFHU54GJ0LB3ixHFmD1YQFbBJux7738gI05pxZTwl69KXoat8OfFamLhwNxg8BcA8GXpkN6i5UZW8VfvOQ11nfpd5RykehgYeFyZaewUizPTZkVfI8BT5fdBGQAGXfNT8Xo0OiRxb64rB3Q3YcKiwXjpD1gc14vZXb5EzP6y274yh9RODHPi3rZuVYeGDBd1woFHng2cFeu1ZZ1kBpB30jqqoxDvYzHWvUEpeP3mRyyL51pqYECgOoTEB9nZlb8J31sbAAyY5EXGp6mPuxnhYfR03wIPdkVaDgY9IUdBfKxiIfghTw9Zk9NYZWz1THnqe1GfnBqtHTWiHHuBt0vEJ5QNiMyvcsZbLb6djbbtVFV5sEMzRWf8cJInSEQFgcXPRFaZYow2bJiugFTGvt0ZZKscHqn2SJFNKGRgr3zlcTgF3Y33PVJVTg1uOaZw6y6mbVituSAR5LuDN0zlFuvkOvr63Hys3fE4dONIOxphiDyNG0Nvci5I3bB3E1H1gT8vUDCq`, false},
	}

	// Run tests
	var encrypted string
	for _, test := range tests {
		if encrypted, err = EncryptWithPrivateKeyString(test.inputKey, test.inputData); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error not expected but got: %s", t.Name(), test.inputKey, test.inputData, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error was expected", t.Name(), test.inputKey, test.inputData)
		} else if err == nil && len(encrypted) == 0 {
			t.Errorf("%s Failed: [%s] [%s] inputted and expected length > 0, but got: 0", t.Name(), test.inputKey, test.inputData)
		}
	}
}

// ExampleEncryptWithPrivateKeyString example using EncryptWithPrivateKeyString()
func ExampleEncryptWithPrivateKeyString() {
	// Create a valid private key
	privateKey, err := CreatePrivateKeyString()
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Encrypt the text
	var encryptedText string
	encryptedText, err = EncryptWithPrivateKeyString(privateKey, "encrypt my message")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	// Can't show the real example, as the value changes each time
	// fmt.Printf("encrypted text length: %d %s", len(encryptedText), encryptedText)
	fmt.Printf("encrypted text length: %d %s", len(encryptedText), "40dd6dbad36cdfb7b1bf85a87887ac3202ca0020b6"+
		"03e6c1fbc30b9a586f06d72eb2fc5117d0f2e4c3193788ebd79c6eb095d88e00209c4ba9a883f631566dfe17d0b2dc83ee456b71927b5b"+
		"5dc6e5a8074d929f3d78c22a4bb7c5f8ffdb0895adf035b2e36fdb5504d3c3b99f3d03a9767eff1b93af14cc586efbf5e81e6fec839c55"+
		"13ab2d644bf3e56d6ddc92e2985c18dab3c6a9")
	// Output:encrypted text length: 300 40dd6dbad36cdfb7b1bf85a87887ac3202ca0020b603e6c1fbc30b9a586f06d72eb2fc5117d0f2e4c3193788ebd79c6eb095d88e00209c4ba9a883f631566dfe17d0b2dc83ee456b71927b5b5dc6e5a8074d929f3d78c22a4bb7c5f8ffdb0895adf035b2e36fdb5504d3c3b99f3d03a9767eff1b93af14cc586efbf5e81e6fec839c5513ab2d644bf3e56d6ddc92e2985c18dab3c6a9
}

// BenchmarkEncryptWithPrivateKeyString benchmarks the method EncryptWithPrivateKeyString()
func BenchmarkEncryptWithPrivateKeyString(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for i := 0; i < b.N; i++ {
		_, _ = EncryptWithPrivateKeyString(key, "some-data")
	}
}

// TestDecryptWithPrivateKey will test the method DecryptWithPrivateKey()
func TestDecryptWithPrivateKey(t *testing.T) {
	t.Parallel()

	// Create a valid private key
	privateKey, err := PrivateKeyFromString("bb66a48a9f6dd7b8fb469a6f08a75c25770591dc509c72129b2aaeca77a5269e")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Create the list of tests
	var tests = []struct {
		inputKey       *bsvec.PrivateKey
		inputEncrypted string
		expectedData   string
		expectedError  bool
	}{
		{privateKey,
			"",
			"",
			true,
		},
		{privateKey,
			"0",
			"",
			true,
		},
		{privateKey,
			"176984bb2f54bfe6275e2d52c8de3d42",
			"",
			true,
		},
		{privateKey,
			"7bea44d59dca4a95f179f9852d511e0f02ca0020961bd9266174d5275c11abd7f0efc19905c04199893ab784e70a7e01f4856b4e0020259296d1bceb4a8ad5267f5633366de03e45d5e7d29754cb35289409c96ec2256a5976bdd8e866f2640b5091b319385f02d59fb89f9d4c0790c70098d142e2a47139671bc4c775e22e5bee59e2310f9b",
			"",
			false,
		},
		{privateKey,
			"8372c24e143fad68496b970f52c971a102ca0020184e2516b9c6ec25873332795f6ae286ec498f72e1879a1b4a6178adbb167d2700204ffc66b0afc2c5e7ca860322f2684ac329901452ec879f5d36ff874285d3a1d77f860da05d454fcb3dc0bc2a880419d0a261236cdd2bf6946e4d2b426c5327bb2e06020b48866730147fc6cac2befdb0",
			" ",
			false,
		},
		{privateKey,
			"4bf9002dacbe9bd7cfeef112f90f1f3202ca0020b3ae7e409189c4fefc92d30dada67dbe6c25eb9bfcc7a0a1078fcafeb355704c00201533cfe68aeaee68caa0d502d7884acfa42101fd2ef5a933a7bcb08812e9795f3ae8913a72a92cfc5c2f4ab2f09903e2d13fde409df0f8395a474bae9a2d6f5f570b5a885dfbc1ea4a31861ccd6cf559",
			"\n",
			false,
		},
		{privateKey,
			"2a7da26bfd1dd55b448c1ce080365fb702ca002074a0063668f26ef3a9371d3043548e44059c4653ea0d7971739d3c122066e1990020084e47d2a2d7a7fca332eb9ea05e981a363ac30f4826fddfef897f5be95f81edb5344b66d884956d23b643ae70ca2f0c1c81f9b81b674920ae4ded7669713ad6f6c3ea0416f020e36b651c68d1c5940f",
			"0",
			false,
		},
		{privateKey,
			"07b4c9a606c2fbe56247847ad2eaf38d02ca0020907673758a824ad6040d699a20e08fb99147908453453a5bc566f52581dd7f390020dc2827b4c12c765ee1c32eb8600feb325ca12bbe81fa41ecd853b3f54b3294e351e0f7b60c353a146bc43e9c0f682af1ee441b3b1bd3b06088de36f933ed7bfeab2c58af1ee5ab7ca9c02d688dee887b",
			"-1",
			false,
		},
		{privateKey,
			"4fab3e3534e101cef1ec936628894a2c02ca00206b60babbfa6ccac9a8078d55a50a32b5e15b1a3c9a3607499066bc4eb70721f4002005cdc7638c7e6c051be80cc83092ed5af5a9a015e24c3ed4af289d59c0c65cb608aa8f07c9318e6e52a18e60dbcf7a5889304f4bfc01cad735a5f2279f06eb5da5ed7454320da87becbdcd889708fcab",
			"test-data",
			false,
		},
		{privateKey,
			"b644c914ce076f88a3ae57ae2726632a02ca002019c1f3296eb7b2e6a909c5efb419b031712e75ca32299e98cc7ea84f7be0ff6e00204b811188d218a0a1839ab11a07d27b4755e92dd730ec044629b9b34dbc0193ab5de4cd3cef9e3c9fdc446664d97644137b1df3cccbca47bc190a3ba41f0805fc4eb2e74d458ec67437cb07ac7d61af75",
			`{"json":"data"}`,
			false,
		},
		{privateKey,
			"b157ee308057c09700f91d6007ec52b802ca0020ca9a8f5267e9e14b141aa9269132a13f119919887ca5d78e1231466bde29f00400208143ed3063b295d90cc76e8e270328e59b390a824fb0c3034f7a8aa3c71716595db36debc32faafb709712413ce2771c384cd2de265cd10a05d6eeb3a806825f20a153bfd298a8e91c158f7b3a8abd2bfe7595f0ba5c2690170f7b725edd28aa60aa4c9e0e7b84c15d6b06222cb0ca3a19fc5a82361002cfd528e9ede4b68a2a3eb0b500536ed2890373e4a5aa55c2ffee3e0da4b2c32657cb03b9c7f51a4baead7d5915076a3db6f3c8b64688b44dfe8d70d4c96ac4c1d06e63dd66a27cbd493b5408eed21e68dc691ba9c23d87740f2352af53f4a0dc080d39c483a09a3299c846aea6fff88602f04a960c8b1fc889cd4bf26365159bdab537004ab74bc81cae4d8fae0616beb0ee9e7dd10255132afd6d326a9e961aba3d2d0efe49b87c8831d13c6ff9503a428e3e646cd5a3239d2e138446a1ca2592b8fb6b8cb19cf91db0e46bcb5e9abe9ac5afd2d1e654e85c7da098ed094b7291ec229767d5636b06796072222ad6f87d3ad3e3b6529fd6a1b8e873f90f5a2aa992e4fc0ec930cb22be95895bb955f8853ab40117b697e70830763cfb9b755278d1a3f9581bf54a81228eb6dfc83d99508c05264cbaab88be2f78fa78aa3e2e86bc893250c7370d8b80d88b5d18a7c43e0f123ea00a44b42c86857e317250fbdc990261ea2f42f5380782cf3056caf66a61eb9ebb3bf43867d230b8415a34832f5072c9e12db7caa0ef7c032b079cbdbdfcda2eb9397e3f7febc11d172b7659732faa0aa95c84a4c4162d043f5fa3fce3788ac2ad5cf4cb5d2a8f82e36273e091c8480a3e5d3c35e01e2e61f4b33d66f3b49f62dc005fa3d55e25157dceecc7b58bf43b36ebf0bea610aa016ec57f9dfc9fd97c9b006fcf6bd2530a231f983bfb430c867f05dfa599eea8720a9ae9829b8b63cdec57766940e1c663966c2d557d2fdb80d67a2d928a929dce84af6b596aecd5cebb0cc21a2dd3bc08c6c48b455e71cac3d863f9c4e1a847b31c2d86ab6b8a691cfdc2d3c0f367d253e9ea0c6668125ae1640434eb6d14d095a128537019f387e9e7af1d85a35fe0139d0949bb67e2d1fd3a3b747f7faaa42de02fcc2ceb0df79ccf2e743f8d5099b222f537dd6c146ca7a54bea4e7fc2cdf99f61055817c296c8fab8ff5f242ff31a5853e7d204007e77265e773001de335d520009861246cbbfbdbb4ef1bbf6aaad1e80dadac0415717b633fc04d3fd0d7217a6d47d38b234e832856e2dd1b2c8ab5ba8c23620ccb82236f1ca227cbef4d59cdd58dcdb4ae2e1efae725eff755c7b929190bda4fa9a80a1d752cb5a23174977badceb2ab0a30a64337f46c903e6ae2906f096d6fe475cbd6742d67871fb58bf3c71d598d6f0f6433c1818a3424758acda5d505d815686c0ab43f9c0e783dfa9f918fecc2b831210812c34606b4f3c1ae55853cabe2cc07ea0e97897446d37606ee097c72e3173213774df1539372fadff75e1bef84d9c2d3de2a22d67e416fb23dbddd74a96ec51b4394225",
			`PLXfZ8ZNN9emRVLapZlvsDpjti3TazobtrY1LWibG5ogf07R5u6xwTGSsBGH10E5naTJLOhj38dFFDeG7giLIvK9QbiF8xSqbXoxvILl09AeXSFfOx7154KdOQdoXMxsghUFi91Efh4uWNl53d7feec7NXh7WrRk68sYUQrfGNdvgxyJlYXBpp7pq5f9Yhhm2TxbNyTyh4dkgvXVJtVfFLwiOWYfO9PAwmfDAgoFyP1vpNcL7G2lCsDW2LuaO1PO0L8a0Z55jpqYrYwzUG3IK4HKnx61xWvawFSzapzGLMrrdxVK56WuToZPOVlLpXjGKYNXlfBXh6T4d1uYUQg9aKfrdIZWCXjVnCXrqDa8vrF0djypIPq24zE0KM0oQZfQJT2Vptu1OPJkes2nIpAUTKxtiLIiunbk8EXaRFBAgGYFFHYm3tHLQUaFQ0K6USoTN3Nr92SjgKUcXVJtSeH854fFpSAoi8je74XPLFHU54GJ0LB3ixHFmD1YQFbBJux7738gI05pxZTwl69KXoat8OfFamLhwNxg8BcA8GXpkN6i5UZW8VfvOQ11nfpd5RykehgYeFyZaewUizPTZkVfI8BT5fdBGQAGXfNT8Xo0OiRxb64rB3Q3YcKiwXjpD1gc14vZXb5EzP6y274yh9RODHPi3rZuVYeGDBd1woFHng2cFeu1ZZ1kBpB30jqqoxDvYzHWvUEpeP3mRyyL51pqYECgOoTEB9nZlb8J31sbAAyY5EXGp6mPuxnhYfR03wIPdkVaDgY9IUdBfKxiIfghTw9Zk9NYZWz1THnqe1GfnBqtHTWiHHuBt0vEJ5QNiMyvcsZbLb6djbbtVFV5sEMzRWf8cJInSEQFgcXPRFaZYow2bJiugFTGvt0ZZKscHqn2SJFNKGRgr3zlcTgF3Y33PVJVTg1uOaZw6y6mbVituSAR5LuDN0zlFuvkOvr63Hys3fE4dONIOxphiDyNG0Nvci5I3bB3E1H1gT8vUDCq`,
			false,
		},
	}

	// Run tests
	var decrypted string
	for _, test := range tests {
		if decrypted, err = DecryptWithPrivateKey(test.inputKey, test.inputEncrypted); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error not expected but got: %s", t.Name(), test.inputKey, test.inputEncrypted, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error was expected", t.Name(), test.inputKey, test.inputEncrypted)
		} else if decrypted != test.expectedData {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%s] expected, but got: %s", t.Name(), test.inputKey, test.inputEncrypted, test.expectedData, decrypted)
		}
	}
}

// ExampleDecryptWithPrivateKey example using DecryptWithPrivateKey()
func ExampleDecryptWithPrivateKey() {
	// Start with a private key
	privateKey, err := PrivateKeyFromString("bb66a48a9f6dd7b8fb469a6f08a75c25770591dc509c72129b2aaeca77a5269e")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Decrypt the encrypted text
	var decryptedText string
	decryptedText, err = DecryptWithPrivateKey(privateKey, "4fab3e3534e101cef1ec936628894a2c02ca00206b60babbfa6cc"+
		"ac9a8078d55a50a32b5e15b1a3c9a3607499066bc4eb70721f4002005cdc7638c7e6c051be80cc83092ed5af5a9a015e24c3ed4af289d59"+
		"c0c65cb608aa8f07c9318e6e52a18e60dbcf7a5889304f4bfc01cad735a5f2279f06eb5da5ed7454320da87becbdcd889708fcab")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("decrypted text: %s", decryptedText)
	// Output:decrypted text: test-data
}

// BenchmarkDecryptWithPrivateKey benchmarks the method DecryptWithPrivateKey()
func BenchmarkDecryptWithPrivateKey(b *testing.B) {
	key, _ := PrivateKeyFromString("bb66a48a9f6dd7b8fb469a6f08a75c25770591dc509c72129b2aaeca77a5269e")
	for i := 0; i < b.N; i++ {
		_, _ = DecryptWithPrivateKey(key, "4fab3e3534e101cef1ec936628894a2c02ca00206b60babbfa6cc"+
			"ac9a8078d55a50a32b5e15b1a3c9a3607499066bc4eb70721f4002005cdc7638c7e6c051be80cc83092ed5af5a9a015e24c3ed4af289d59"+
			"c0c65cb608aa8f07c9318e6e52a18e60dbcf7a5889304f4bfc01cad735a5f2279f06eb5da5ed7454320da87becbdcd889708fcab")
	}
}

// TestDecryptWithPrivateKeyString will test the method DecryptWithPrivateKeyString()
func TestDecryptWithPrivateKeyString(t *testing.T) {
	t.Parallel()

	// Create a valid private key
	privateKey := "bb66a48a9f6dd7b8fb469a6f08a75c25770591dc509c72129b2aaeca77a5269e"

	// Create the list of tests
	var tests = []struct {
		inputKey       string
		inputEncrypted string
		expectedData   string
		expectedError  bool
	}{
		{"",
			"",
			"",
			true,
		},
		{"0",
			"",
			"",
			true,
		},
		{privateKey,
			"",
			"",
			true,
		},
		{privateKey,
			"0",
			"",
			true,
		},
		{privateKey,
			"176984bb2f54bfe6275e2d52c8de3d42",
			"",
			true,
		},
		{privateKey,
			"7bea44d59dca4a95f179f9852d511e0f02ca0020961bd9266174d5275c11abd7f0efc19905c04199893ab784e70a7e01f4856b4e0020259296d1bceb4a8ad5267f5633366de03e45d5e7d29754cb35289409c96ec2256a5976bdd8e866f2640b5091b319385f02d59fb89f9d4c0790c70098d142e2a47139671bc4c775e22e5bee59e2310f9b",
			"",
			false,
		},
		{privateKey,
			"8372c24e143fad68496b970f52c971a102ca0020184e2516b9c6ec25873332795f6ae286ec498f72e1879a1b4a6178adbb167d2700204ffc66b0afc2c5e7ca860322f2684ac329901452ec879f5d36ff874285d3a1d77f860da05d454fcb3dc0bc2a880419d0a261236cdd2bf6946e4d2b426c5327bb2e06020b48866730147fc6cac2befdb0",
			" ",
			false,
		},
		{privateKey,
			"4bf9002dacbe9bd7cfeef112f90f1f3202ca0020b3ae7e409189c4fefc92d30dada67dbe6c25eb9bfcc7a0a1078fcafeb355704c00201533cfe68aeaee68caa0d502d7884acfa42101fd2ef5a933a7bcb08812e9795f3ae8913a72a92cfc5c2f4ab2f09903e2d13fde409df0f8395a474bae9a2d6f5f570b5a885dfbc1ea4a31861ccd6cf559",
			"\n",
			false,
		},
		{privateKey,
			"2a7da26bfd1dd55b448c1ce080365fb702ca002074a0063668f26ef3a9371d3043548e44059c4653ea0d7971739d3c122066e1990020084e47d2a2d7a7fca332eb9ea05e981a363ac30f4826fddfef897f5be95f81edb5344b66d884956d23b643ae70ca2f0c1c81f9b81b674920ae4ded7669713ad6f6c3ea0416f020e36b651c68d1c5940f",
			"0",
			false,
		},
		{privateKey,
			"07b4c9a606c2fbe56247847ad2eaf38d02ca0020907673758a824ad6040d699a20e08fb99147908453453a5bc566f52581dd7f390020dc2827b4c12c765ee1c32eb8600feb325ca12bbe81fa41ecd853b3f54b3294e351e0f7b60c353a146bc43e9c0f682af1ee441b3b1bd3b06088de36f933ed7bfeab2c58af1ee5ab7ca9c02d688dee887b",
			"-1",
			false,
		},
		{privateKey,
			"4fab3e3534e101cef1ec936628894a2c02ca00206b60babbfa6ccac9a8078d55a50a32b5e15b1a3c9a3607499066bc4eb70721f4002005cdc7638c7e6c051be80cc83092ed5af5a9a015e24c3ed4af289d59c0c65cb608aa8f07c9318e6e52a18e60dbcf7a5889304f4bfc01cad735a5f2279f06eb5da5ed7454320da87becbdcd889708fcab",
			"test-data",
			false,
		},
		{privateKey,
			"b644c914ce076f88a3ae57ae2726632a02ca002019c1f3296eb7b2e6a909c5efb419b031712e75ca32299e98cc7ea84f7be0ff6e00204b811188d218a0a1839ab11a07d27b4755e92dd730ec044629b9b34dbc0193ab5de4cd3cef9e3c9fdc446664d97644137b1df3cccbca47bc190a3ba41f0805fc4eb2e74d458ec67437cb07ac7d61af75",
			`{"json":"data"}`,
			false,
		},
		{privateKey,
			"b157ee308057c09700f91d6007ec52b802ca0020ca9a8f5267e9e14b141aa9269132a13f119919887ca5d78e1231466bde29f00400208143ed3063b295d90cc76e8e270328e59b390a824fb0c3034f7a8aa3c71716595db36debc32faafb709712413ce2771c384cd2de265cd10a05d6eeb3a806825f20a153bfd298a8e91c158f7b3a8abd2bfe7595f0ba5c2690170f7b725edd28aa60aa4c9e0e7b84c15d6b06222cb0ca3a19fc5a82361002cfd528e9ede4b68a2a3eb0b500536ed2890373e4a5aa55c2ffee3e0da4b2c32657cb03b9c7f51a4baead7d5915076a3db6f3c8b64688b44dfe8d70d4c96ac4c1d06e63dd66a27cbd493b5408eed21e68dc691ba9c23d87740f2352af53f4a0dc080d39c483a09a3299c846aea6fff88602f04a960c8b1fc889cd4bf26365159bdab537004ab74bc81cae4d8fae0616beb0ee9e7dd10255132afd6d326a9e961aba3d2d0efe49b87c8831d13c6ff9503a428e3e646cd5a3239d2e138446a1ca2592b8fb6b8cb19cf91db0e46bcb5e9abe9ac5afd2d1e654e85c7da098ed094b7291ec229767d5636b06796072222ad6f87d3ad3e3b6529fd6a1b8e873f90f5a2aa992e4fc0ec930cb22be95895bb955f8853ab40117b697e70830763cfb9b755278d1a3f9581bf54a81228eb6dfc83d99508c05264cbaab88be2f78fa78aa3e2e86bc893250c7370d8b80d88b5d18a7c43e0f123ea00a44b42c86857e317250fbdc990261ea2f42f5380782cf3056caf66a61eb9ebb3bf43867d230b8415a34832f5072c9e12db7caa0ef7c032b079cbdbdfcda2eb9397e3f7febc11d172b7659732faa0aa95c84a4c4162d043f5fa3fce3788ac2ad5cf4cb5d2a8f82e36273e091c8480a3e5d3c35e01e2e61f4b33d66f3b49f62dc005fa3d55e25157dceecc7b58bf43b36ebf0bea610aa016ec57f9dfc9fd97c9b006fcf6bd2530a231f983bfb430c867f05dfa599eea8720a9ae9829b8b63cdec57766940e1c663966c2d557d2fdb80d67a2d928a929dce84af6b596aecd5cebb0cc21a2dd3bc08c6c48b455e71cac3d863f9c4e1a847b31c2d86ab6b8a691cfdc2d3c0f367d253e9ea0c6668125ae1640434eb6d14d095a128537019f387e9e7af1d85a35fe0139d0949bb67e2d1fd3a3b747f7faaa42de02fcc2ceb0df79ccf2e743f8d5099b222f537dd6c146ca7a54bea4e7fc2cdf99f61055817c296c8fab8ff5f242ff31a5853e7d204007e77265e773001de335d520009861246cbbfbdbb4ef1bbf6aaad1e80dadac0415717b633fc04d3fd0d7217a6d47d38b234e832856e2dd1b2c8ab5ba8c23620ccb82236f1ca227cbef4d59cdd58dcdb4ae2e1efae725eff755c7b929190bda4fa9a80a1d752cb5a23174977badceb2ab0a30a64337f46c903e6ae2906f096d6fe475cbd6742d67871fb58bf3c71d598d6f0f6433c1818a3424758acda5d505d815686c0ab43f9c0e783dfa9f918fecc2b831210812c34606b4f3c1ae55853cabe2cc07ea0e97897446d37606ee097c72e3173213774df1539372fadff75e1bef84d9c2d3de2a22d67e416fb23dbddd74a96ec51b4394225",
			`PLXfZ8ZNN9emRVLapZlvsDpjti3TazobtrY1LWibG5ogf07R5u6xwTGSsBGH10E5naTJLOhj38dFFDeG7giLIvK9QbiF8xSqbXoxvILl09AeXSFfOx7154KdOQdoXMxsghUFi91Efh4uWNl53d7feec7NXh7WrRk68sYUQrfGNdvgxyJlYXBpp7pq5f9Yhhm2TxbNyTyh4dkgvXVJtVfFLwiOWYfO9PAwmfDAgoFyP1vpNcL7G2lCsDW2LuaO1PO0L8a0Z55jpqYrYwzUG3IK4HKnx61xWvawFSzapzGLMrrdxVK56WuToZPOVlLpXjGKYNXlfBXh6T4d1uYUQg9aKfrdIZWCXjVnCXrqDa8vrF0djypIPq24zE0KM0oQZfQJT2Vptu1OPJkes2nIpAUTKxtiLIiunbk8EXaRFBAgGYFFHYm3tHLQUaFQ0K6USoTN3Nr92SjgKUcXVJtSeH854fFpSAoi8je74XPLFHU54GJ0LB3ixHFmD1YQFbBJux7738gI05pxZTwl69KXoat8OfFamLhwNxg8BcA8GXpkN6i5UZW8VfvOQ11nfpd5RykehgYeFyZaewUizPTZkVfI8BT5fdBGQAGXfNT8Xo0OiRxb64rB3Q3YcKiwXjpD1gc14vZXb5EzP6y274yh9RODHPi3rZuVYeGDBd1woFHng2cFeu1ZZ1kBpB30jqqoxDvYzHWvUEpeP3mRyyL51pqYECgOoTEB9nZlb8J31sbAAyY5EXGp6mPuxnhYfR03wIPdkVaDgY9IUdBfKxiIfghTw9Zk9NYZWz1THnqe1GfnBqtHTWiHHuBt0vEJ5QNiMyvcsZbLb6djbbtVFV5sEMzRWf8cJInSEQFgcXPRFaZYow2bJiugFTGvt0ZZKscHqn2SJFNKGRgr3zlcTgF3Y33PVJVTg1uOaZw6y6mbVituSAR5LuDN0zlFuvkOvr63Hys3fE4dONIOxphiDyNG0Nvci5I3bB3E1H1gT8vUDCq`,
			false,
		},
	}

	// Run tests
	for _, test := range tests {
		if decrypted, err := DecryptWithPrivateKeyString(test.inputKey, test.inputEncrypted); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error not expected but got: %s", t.Name(), test.inputKey, test.inputEncrypted, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error was expected", t.Name(), test.inputKey, test.inputEncrypted)
		} else if decrypted != test.expectedData {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%s] expected, but got: %s", t.Name(), test.inputKey, test.inputEncrypted, test.expectedData, decrypted)
		}
	}
}

// ExampleDecryptWithPrivateKeyString example using DecryptWithPrivateKeyString()
func ExampleDecryptWithPrivateKeyString() {
	// Decrypt the encrypted text
	decryptedText, err := DecryptWithPrivateKeyString(
		"bb66a48a9f6dd7b8fb469a6f08a75c25770591dc509c72129b2aaeca77a5269e",
		"4fab3e3534e101cef1ec936628894a2c02ca00206b60babbfa6cc"+
			"ac9a8078d55a50a32b5e15b1a3c9a3607499066bc4eb70721f4002005cdc7638c7e6c051be80cc83092ed5af5a9a015e24c3ed4af289d59"+
			"c0c65cb608aa8f07c9318e6e52a18e60dbcf7a5889304f4bfc01cad735a5f2279f06eb5da5ed7454320da87becbdcd889708fcab",
	)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("decrypted text: %s", decryptedText)
	// Output:decrypted text: test-data
}

// BenchmarkDecryptWithPrivateKeyString benchmarks the method DecryptWithPrivateKeyString()
func BenchmarkDecryptWithPrivateKeyString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DecryptWithPrivateKeyString(
			"bb66a48a9f6dd7b8fb469a6f08a75c25770591dc509c72129b2aaeca77a5269e",
			"4fab3e3534e101cef1ec936628894a2c02ca00206b60babbfa6cc"+
				"ac9a8078d55a50a32b5e15b1a3c9a3607499066bc4eb70721f4002005cdc7638c7e6c051be80cc83092ed5af5a9a015e24c3ed4af289d59"+
				"c0c65cb608aa8f07c9318e6e52a18e60dbcf7a5889304f4bfc01cad735a5f2279f06eb5da5ed7454320da87becbdcd889708fcab",
		)
	}
}
