package bitcoin

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"testing"

	"github.com/piotrnar/gocoin/lib/secp256k1"
)

const (
	testDERSignature = "3045022100b976be863fffd361716b375a9a5c4e77073dfaa29d2b9af9addef94f029c2d0902205b1fffc58343f3d4bd8fc48a118e998072c655d318061e13e1ef0902fb42e15c"
	testDERPubKey    = "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"
)

// TestVerifyMessage will test the method VerifyMessage()
func TestVerifyMessage(t *testing.T) {

	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputAddress   string
		inputSignature string
		inputData      string
		expectedError  bool
	}{
		{
			"12SsqqYk43kggMBpSvWHwJwR31NsgMePKS",
			"HFxPx8JHsCiivB+DW/RgNpCLT6yG3j436cUNWKekV3ORBrHNChIjeVReyAco7PVmmDtVD3POs9FhDlm/nk5I6O8=",
			"test message",
			false,
		},
		{
			"1LN5p7Eg9Zju1b4g4eFPTBMPoMZGCxzrET",
			"IKmgOFrfWRffRNjrQcJQHSBD7WL2di+4doWdaz/a/p5RUiT7ErpUqbYeLi0yzmONFaV8uLWF2vydTjA8W8KnjZU=",
			"This time I'm writing a new message that is obnoxiously long af. This time I'm writing a " +
				"new message that is obnoxiously long af. This time I'm writing a new message that is obnoxiously " +
				"long af. This time I'm writing a new message that is obnoxiously long af. This time I'm writing a " +
				"new message that is obnoxiously long af. This time I'm writing a new message that is obnoxiously " +
				"long af. This time I'm writing a new message that is obnoxiously long af. This time I'm writing a " +
				"new message that is obnoxiously long af. This time I'm writing a new message that is obnoxiously " +
				"long af. This time I'm writing a new message that is obnoxiously long af. This time I'm writing a " +
				"new message that is obnoxiously long af. This time I'm writing a new message that is obnoxiously " +
				"long af. This time I'm writing a new message that is obnoxiously long af. This time I'm writing a " +
				"new message that is obnoxiously long af.",
			false,
		},
		{
			"1LN5p7Eg9Zju1b4g4eFPTBMPoMZGCxzrET",
			"IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=",
			"Testing!",
			true,
		},
		{
			"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
			"",
			"Testing!",
			true,
		},
		{
			"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
			"IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=",
			"",
			true,
		},
		{
			"0",
			"IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=",
			"Testing!",
			true,
		},
		{
			"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
			"GBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4naZ=",
			"Testing!",
			true,
		},
		{
			"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
			"GBD=",
			"Testing!",
			true,
		},
		{
			"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
			"GBse5w0f839t8wej8f2D=",
			"Testing!",
			true,
		},
	}

	// Run tests
	for _, test := range tests {
		if err := VerifyMessage(test.inputAddress, test.inputSignature, test.inputData); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%s] [%s] [%s] inputted and error not expected but got: %s", t.Name(), test.inputAddress, test.inputSignature, test.inputData, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%s] [%s] [%s] inputted and error was expected", t.Name(), test.inputAddress, test.inputSignature, test.inputData)
		}
	}

}

// ExampleVerifyMessage example using VerifyMessage()
func ExampleVerifyMessage() {
	if err := VerifyMessage(
		"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
		"IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=",
		"Testing!",
	); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("verification passed")
	// Output:verification passed
}

// BenchmarkVerifyMessage benchmarks the method VerifyMessage()
func BenchmarkVerifyMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = VerifyMessage(
			"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
			"IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=",
			"Testing!",
		)
	}
}

// TestVerifyMessageSigRecover will test the method recoverSig()
//
// From: https://github.com/piotrnar/gocoin/blob/master/lib/secp256k1/sig_test.go
func TestVerifyMessageSigRecover(t *testing.T) {
	var vs = [][6]string{
		{
			"6028b9e3a31c9e725fcbd7d5d16736aaaafcc9bf157dfb4be62bcbcf0969d488",
			"036d4a36fa235b8f9f815aa6f5457a607f956a71a035bf0970d8578bf218bb5a",
			"9cff3da1a4f86caf3683f865232c64992b5ed002af42b321b8d8a48420680487",
			"0",
			"56dc5df245955302893d8dda0677cc9865d8011bc678c7803a18b5f6faafec08",
			"54b5fbdcd8fac6468dac2de88fadce6414f5f3afbb103753e25161bef77705a6",
		},
		{
			"b470e02f834a3aaafa27bd2b49e07269e962a51410f364e9e195c31351a05e50",
			"560978aed76de9d5d781f87ed2068832ed545f2b21bf040654a2daff694c8b09",
			"9ce428d58e8e4caf619dc6fc7b2c2c28f0561654d1f80f322c038ad5e67ff8a6",
			"1",
			"15b7e7d00f024bffcd2e47524bb7b7d3a6b251e23a3a43191ed7f0a418d9a578",
			"bf29a25e2d1f32c5afb18b41ae60112723278a8af31275965a6ec1d95334e840",
		},
	}

	var sig secp256k1.Signature
	var pubkey, exp secp256k1.XY
	var msg secp256k1.Number

	for i := range vs {
		sig.R.SetHex(vs[i][0])
		sig.S.SetHex(vs[i][1])
		msg.SetHex(vs[i][2])
		rid, _ := strconv.ParseInt(vs[i][3], 10, 32)
		exp.X.SetHex(vs[i][4])
		exp.Y.SetHex(vs[i][5])

		success := secp256k1.RecoverPublicKey(sig.R.Bytes(), sig.S.Bytes(), msg.Bytes(), int(rid), &pubkey)

		if success {
			if !exp.X.Equals(&pubkey.X) {
				t.Error("X mismatch at vector", i)
			}
			if !exp.Y.Equals(&pubkey.Y) {
				t.Error("Y mismatch at vector", i)
			}
		} else {
			t.Error("sig.recover failed")
		}
	}
}

// TestVerifyMessageSigRecoverFailed will test the method recoverSig()
func TestVerifyMessageSigRecoverFailed(t *testing.T) {

	// Bad signature cases
	var vs = [][6]string{
		{
			"0",
			"0",
			"0",
			"0",
			"0",
			"0",
		},
		{
			"000",
			"560978aed76de9d5d781f87ed2068832ed545f2b21bf040654a2daff694c8b09",
			"9ce428d58e8e4caf619dc6fc7b2c2c28f0561654d1f80f322c038ad5e67ff8a6",
			"1",
			"000",
			"bf29a25e2d1f32c5afb18b41ae60112723278a8af31275965a6ec1d95334e840",
		},
		{
			"000",
			"1234567",
			"1234567",
			"1",
			"000",
			"1234567",
		},
		{
			"bf29a25e2d1f32c5afb18b41ae60112723278a8af31275965a6ec1d95334e840",
			"bf29a25e2d1f32c5afb18b41ae60112723278a8af31275965a6ec1d95334e840",
			"bf29a25e2d1f32c5afb18b41ae60112723278a8af31275965a6ec1d95334e840",
			"1",
			"bf29a25e2d1f32c5afb18b41ae60112723278a8af31275965a6ec1d95334e840",
			"bf29a25e2d1f32c5afb18b41ae60112723278a8af31275965a6ec1d95334e840",
		},
	}

	var sig secp256k1.Signature
	var pubkey, exp secp256k1.XY
	var msg secp256k1.Number

	for i := range vs {
		sig.R.SetHex(vs[i][0])
		sig.S.SetHex(vs[i][1])
		msg.SetHex(vs[i][2])
		rid, err := strconv.ParseInt(vs[i][3], 10, 32)
		if err != nil {
			t.Fatalf("failed in ParseInt: %s", err.Error())
		}
		exp.X.SetHex(vs[i][4])
		exp.Y.SetHex(vs[i][5])

		if secp256k1.RecoverPublicKey(sig.R.Bytes(), sig.S.Bytes(), msg.Bytes(), int(rid), &pubkey) {
			t.Fatalf("sigRecover should have failed")
		}
	}

}

// TestVerifyMessageDER will test the method VerifyMessageDER()
func TestVerifyMessageDER(t *testing.T) {

	// Example message (payload from Merchant API)
	message := []byte(`{"apiVersion":"0.1.0","timestamp":"2020-10-08T14:25:31.539Z","expiryTime":"2020-10-08T14:35:31.539Z","minerId":"` + testDERPubKey + `","currentHighestBlockHash":"0000000000000000021af4ee1f179a64e530bf818ef67acd09cae24a89124519","currentHighestBlockHeight":656007,"minerReputation":null,"fees":[{"id":1,"feeType":"standard","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}},{"id":2,"feeType":"data","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}}]}`)
	invalidMessage := []byte("invalid-message")
	validHash := sha256.Sum256(message)

	// Test a valid signature
	verified, err := VerifyMessageDER(validHash, testDERPubKey, testDERSignature)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if !verified {
		t.Fatalf("expected verified to be true")
	}

	// Test an invalid pubkey
	verified, err = VerifyMessageDER(validHash, testDERPubKey+"00", testDERSignature)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verified {
		t.Fatalf("expected verified to be false")
	}

	// Test an invalid pubkey
	verified, err = VerifyMessageDER(validHash, "0", testDERSignature)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verified {
		t.Fatalf("expected verified to be false")
	}

	// Test an invalid signature
	verified, err = VerifyMessageDER(validHash, testDERPubKey, "0"+testDERSignature)
	if verified {
		t.Fatalf("expected verified to be false but got: %v", verified)
	} else if err == nil {
		t.Fatalf("expected error not be nil")
	}

	// Test an invalid signature
	verified, err = VerifyMessageDER(validHash, testDERPubKey, testDERSignature+"-1")
	if verified {
		t.Fatalf("expected verified to be false but got: %v", verified)
	} else if err == nil {
		t.Fatalf("expected error not be nil")
	}

	// Test an invalid signature
	verified, err = VerifyMessageDER(validHash, testDERPubKey, "1234567")
	if verified {
		t.Fatalf("expected verified to be false but got: %v", verified)
	} else if err == nil {
		t.Fatalf("expected error not be nil")
	}

	// Test an invalid message
	verified, err = VerifyMessageDER(sha256.Sum256(invalidMessage), testDERPubKey, testDERSignature)
	if verified {
		t.Fatalf("expected verified to be false")
	} else if err != nil {
		t.Fatalf("expected error to be nil, but got: %s", err.Error())
	}
}

// ExampleVerifyMessageDER example using VerifyMessageDER()
func ExampleVerifyMessageDER() {
	message := []byte(`{"apiVersion":"0.1.0","timestamp":"2020-10-08T14:25:31.539Z","expiryTime":"2020-10-08T14:35:31.539Z","minerId":"` + testDERPubKey + `","currentHighestBlockHash":"0000000000000000021af4ee1f179a64e530bf818ef67acd09cae24a89124519","currentHighestBlockHeight":656007,"minerReputation":null,"fees":[{"id":1,"feeType":"standard","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}},{"id":2,"feeType":"data","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}}]}`)

	verified, err := VerifyMessageDER(sha256.Sum256(message), testDERPubKey, testDERSignature)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	} else if !verified {
		fmt.Printf("verification failed")
		return
	}
	fmt.Printf("verification passed")
	// Output:verification passed
}

// BenchmarkVerifyMessageDER benchmarks the method VerifyMessageDER()
func BenchmarkVerifyMessageDER(b *testing.B) {
	message := []byte(`{"apiVersion":"0.1.0","timestamp":"2020-10-08T14:25:31.539Z","expiryTime":"2020-10-08T14:35:31.539Z","minerId":"` + testDERPubKey + `","currentHighestBlockHash":"0000000000000000021af4ee1f179a64e530bf818ef67acd09cae24a89124519","currentHighestBlockHeight":656007,"minerReputation":null,"fees":[{"id":1,"feeType":"standard","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}},{"id":2,"feeType":"data","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}}]}`)

	for i := 0; i < b.N; i++ {
		_, _ = VerifyMessageDER(sha256.Sum256(message), testDERPubKey, testDERSignature)
	}
}
