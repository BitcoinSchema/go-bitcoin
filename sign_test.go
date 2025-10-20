package bitcoin

import (
	"fmt"
	"testing"
)

func TestSigningCompression(t *testing.T) {
	testKey := "0499f8239bfe10eb0f5e53d543635a423c96529dd85fa4bad42049a0b435ebdd"
	testData := "test message"

	// Test sign uncompressed
	address, err := GetAddressFromPrivateKeyString(testKey, false, true)
	if err != nil {
		t.Errorf("Get address err %s", err)
	}
	sig, err := SignMessage(testKey, testData, false)
	if err != nil {
		t.Errorf("Failed to sign uncompressed %s", err)
	}

	err = VerifyMessage(address, sig, testData, true)
	if err != nil {
		t.Errorf("Failed to validate uncompressed %s", err)
	}

	// Test sign compressed
	address, err = GetAddressFromPrivateKeyString(testKey, true, true)
	if err != nil {
		t.Errorf("Get address err %s", err)
	}
	sig, err = SignMessage(testKey, testData, true)
	if err != nil {
		t.Errorf("Failed to sign compressed %s", err)
	}

	err = VerifyMessage(address, sig, testData, true)
	if err != nil {
		t.Errorf("Failed to validate compressed %s", err)
	}
}

// TestSignMessage will test the method SignMessage()
func TestSignMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inputKey          string
		inputMessage      string
		expectedSignature string
		expectedError     bool
	}{
		{
			"0499f8239bfe10eb0f5e53d543635a423c96529dd85fa4bad42049a0b435ebdd",
			"test message",
			"HFxPx8JHsCiivB+DW/RgNpCLT6yG3j436cUNWKekV3ORBrHNChIjeVReyAco7PVmmDtVD3POs9FhDlm/nk5I6O8=",
			false,
		},
		{
			"ef0b8bad0be285099534277fde328f8f19b3be9cadcd4c08e6ac0b5f863745ac",
			"This is a test message",
			"G+zZagsyz7ioC/ZOa5EwsaKice0vs2BvZ0ljgkFHxD3vGsMlGeD4sXHEcfbI4h8lP29VitSBdf4A+nHXih7svf4=",
			false,
		},
		{
			"0499f8239bfe10eb0f5e53d543635a423c96529dd85fa4bad42049a0b435ebdd",
			"This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af. This time I'm writing a new message that is obnixiously long af.",
			"GxRcFXQc7LHxFNpK5lzhR+LF5ixIvhB089bxYzTAV02yGHm/3ALxltz/W4lGp77Q5UTxdj+TU+96mdAcJ5b/fGs=",
			false,
		},
		{
			"93596babb564cbbdc84f2370c710b9bcc94333495b60af719b5fcf9ba00ba82c",
			"This is a test message",
			"HIuDw09ffPgEDuxEw5yHVp1+mi4QpuhAwLyQdpMTfsHCOkMqTKXuP7dSNWMEJqZsiQ8eKMDRvf2wZ4e5bxcu4O0=",
			false,
		},
		{
			"50381cf8f52936faae4a05a073a03d688a9fa206d005e87a39da436c75476d78",
			"This is a test message",
			"HLBmbjCY2Z7eSXGXZoBI3x2ZRaYUYOGtEaDjXetaY+zNDtMOvagsOGEHnVT3f5kXlEbuvmPydHqLnyvZP3cDOWk=",
			false,
		},
		{
			"c7726663147afd1add392d129086e57c0b05aa66a6ded564433c04bd55741434",
			"This is a test message",
			"HOI207QUnTLr2Ll+s4kUxNgLgorkc/Z5Pc+XNvUBYLy2TxaU6oHEJ2TTJ1mZVrtUyHm6e315v1tIjeosW3Odfqw=",
			false,
		},
		{
			"c7726663147afd1add392d129086e57c0b05aa66a6ded564433c04bd55741434",
			"1",
			"HMcRFG1VNN9TDGXpCU+9CqKLNOuhwQiXI5hZpkTOuYHKBDOWayNuAABofYLqUHYTMiMf9mYFQ0sPgFJZz3F7ELQ=",
			false,
		},
		{
			"",
			"This is a test message",
			"",
			true,
		},
		{
			"0",
			"This is a test message",
			"",
			true,
		},
		{
			"0000000",
			"This is a test message",
			"",
			true,
		},
		{
			"c7726663147afd1add392d129086e57c0b",
			"This is a test message",
			"G6N+iPf23i2YkLsNzF/yyeBm9eSYBoY/HFV1Md1F0ElWBXW5E5mkdRtgjoRuq0yNb1CCFNWWlkn2gZknFJNUFJ8=",
			false,
		},
	}

	for idx, test := range tests {
		if signature, err := SignMessage(test.inputKey, test.inputMessage, false); err != nil && !test.expectedError {
			t.Fatalf("%d %s Failed: [%s] [%s] inputted and error not expected but got: %s", idx, t.Name(), test.inputKey, test.inputMessage, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%d %s Failed: [%s] [%s] inputted and error was expected", idx, t.Name(), test.inputKey, test.inputMessage)
		} else if signature != test.expectedSignature {
			t.Fatalf("%d %s Failed: [%s] [%s] inputted [%s] expected but got: %s", idx, t.Name(), test.inputKey, test.inputMessage, test.expectedSignature, signature)
		}
	}
}

// ExampleSignMessage example using SignMessage()
func ExampleSignMessage() {
	signature, err := SignMessage("ef0b8bad0be285099534277fde328f8f19b3be9cadcd4c08e6ac0b5f863745ac", "This is a test message", false)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("signature created: %s", signature)
	// Output:signature created: G+zZagsyz7ioC/ZOa5EwsaKice0vs2BvZ0ljgkFHxD3vGsMlGeD4sXHEcfbI4h8lP29VitSBdf4A+nHXih7svf4=
}

// BenchmarkSignMessage benchmarks the method SignMessage()
func BenchmarkSignMessage(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for i := 0; i < b.N; i++ {
		_, _ = SignMessage(key, "This is a test message", false)
	}
}
