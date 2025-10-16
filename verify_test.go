package bitcoin

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDERSignature = "3045022100b976be863fffd361716b375a9a5c4e77073dfaa29d2b9af9addef94f029c2d0902205b1fffc58343f3d4bd8fc48a118e998072c655d318061e13e1ef0902fb42e15c"
	testDERPubKey    = "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"
)

// TestVerifyMessage will test the method VerifyMessage()
func TestVerifyMessage(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		inputAddress   string
		inputSignature string
		inputData      string
		expectedError  bool
		mainnet        bool
	}{
		{
			"12SsqqYk43kggMBpSvWHwJwR31NsgMePKS",
			"HFxPx8JHsCiivB+DW/RgNpCLT6yG3j436cUNWKekV3ORBrHNChIjeVReyAco7PVmmDtVD3POs9FhDlm/nk5I6O8=",
			"test message",
			false,
			true,
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
			true,
		},
		{
			"mrH55sFASmaDJZ46RnKjMi11nA99b4d8GH",
			"IL3hIysOVTZu9NJp5YWkPh7PSrnX+kFuFArVB+ETNObqUeYtboWjfV2H7CVmOJGJkjo4REHJx26zCGrH71ySNRo=",
			"Testing!",
			false,
			false,
		},
		{
			"1LN5p7Eg9Zju1b4g4eFPTBMPoMZGCxzrET",
			"IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=",
			"Testing!",
			true,
			true,
		},
		{
			"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
			"",
			"Testing!",
			true,
			true,
		},
		{
			"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
			"IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=",
			"",
			true,
			true,
		},
		{
			"0",
			"IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=",
			"Testing!",
			true,
			true,
		},
		{
			"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
			"GBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4naZ=",
			"Testing!",
			true,
			true,
		},
		{
			"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
			"GBD=",
			"Testing!",
			true,
			true,
		},
		{
			"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ",
			"GBse5w0f839t8wej8f2D=",
			"Testing!",
			true,
			true,
		},
	}

	for _, test := range tests {
		if err := VerifyMessage(test.inputAddress, test.inputSignature, test.inputData, test.mainnet); err != nil && !test.expectedError {
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
		true,
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
			true,
		)
	}
}

// TestVerifyMessageDER will test the method VerifyMessageDER()
func TestVerifyMessageDER(t *testing.T) {

	// Example message (payload from Merchant API)
	message := []byte(`{"apiVersion":"0.1.0","timestamp":"2020-10-08T14:25:31.539Z","expiryTime":"2020-10-08T14:35:31.539Z","minerId":"` + testDERPubKey + `","currentHighestBlockHash":"0000000000000000021af4ee1f179a64e530bf818ef67acd09cae24a89124519","currentHighestBlockHeight":656007,"minerReputation":null,"fees":[{"id":1,"feeType":"standard","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}},{"id":2,"feeType":"data","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}}]}`)
	invalidMessage := []byte("invalid-message")
	validHash := sha256.Sum256(message)

	t.Run("valid signature", func(t *testing.T) {
		verified, err := VerifyMessageDER(validHash, testDERPubKey, testDERSignature)
		assert.NoError(t, err)
		assert.Equal(t, true, verified)
	})

	t.Run("invalid pubkey", func(t *testing.T) {
		verified, err := VerifyMessageDER(validHash, testDERPubKey+"00", testDERSignature)
		assert.Error(t, err)
		assert.Equal(t, false, verified)
	})

	t.Run("invalid pubkey 2", func(t *testing.T) {
		verified, err := VerifyMessageDER(validHash, "0", testDERSignature)
		assert.Error(t, err)
		assert.Equal(t, false, verified)
	})

	t.Run("invalid signature (prefix)", func(t *testing.T) {
		verified, err := VerifyMessageDER(validHash, testDERPubKey, "0"+testDERSignature)
		assert.Error(t, err)
		assert.Equal(t, false, verified)
	})

	t.Run("invalid signature (suffix)", func(t *testing.T) {
		verified, err := VerifyMessageDER(validHash, testDERPubKey, testDERSignature+"-1")
		assert.Error(t, err)
		assert.Equal(t, false, verified)
	})

	t.Run("invalid signature (length)", func(t *testing.T) {
		verified, err := VerifyMessageDER(validHash, testDERPubKey, "1234567")
		assert.Error(t, err)
		assert.Equal(t, false, verified)
	})

	t.Run("invalid message", func(t *testing.T) {
		verified, err := VerifyMessageDER(sha256.Sum256(invalidMessage), testDERPubKey, testDERSignature)
		assert.NoError(t, err)
		assert.Equal(t, false, verified)
	})
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
