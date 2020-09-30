package bitcoin

import (
	"fmt"
	"testing"
)

// TestSignMessage will test the method SignMessage()
func TestSignMessage(t *testing.T) {

	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputKey          string
		inputMessage      string
		expectedSignature string
		expectedError     bool
	}{
		{"ef0b8bad0be285099534277fde328f8f19b3be9cadcd4c08e6ac0b5f863745ac", "This is a test message", "H+zZagsyz7ioC/ZOa5EwsaKice0vs2BvZ0ljgkFHxD3vGsMlGeD4sXHEcfbI4h8lP29VitSBdf4A+nHXih7svf4=", false},
		{"93596babb564cbbdc84f2370c710b9bcc94333495b60af719b5fcf9ba00ba82c", "This is a test message", "IIuDw09ffPgEDuxEw5yHVp1+mi4QpuhAwLyQdpMTfsHCOkMqTKXuP7dSNWMEJqZsiQ8eKMDRvf2wZ4e5bxcu4O0=", false},
		{"50381cf8f52936faae4a05a073a03d688a9fa206d005e87a39da436c75476d78", "This is a test message", "ILBmbjCY2Z7eSXGXZoBI3x2ZRaYUYOGtEaDjXetaY+zNDtMOvagsOGEHnVT3f5kXlEbuvmPydHqLnyvZP3cDOWk=", false},
		{"c7726663147afd1add392d129086e57c0b05aa66a6ded564433c04bd55741434", "This is a test message", "IOI207QUnTLr2Ll+s4kUxNgLgorkc/Z5Pc+XNvUBYLy2TxaU6oHEJ2TTJ1mZVrtUyHm6e315v1tIjeosW3Odfqw=", false},
		{"c7726663147afd1add392d129086e57c0b05aa66a6ded564433c04bd55741434", "1", "IMcRFG1VNN9TDGXpCU+9CqKLNOuhwQiXI5hZpkTOuYHKBDOWayNuAABofYLqUHYTMiMf9mYFQ0sPgFJZz3F7ELQ=", false},
		{"", "This is a test message", "", true},
		{"0", "This is a test message", "", true},
		{"0000000", "This is a test message", "", true},
		{"c7726663147afd1add392d129086e57c0b", "This is a test message", "H6N+iPf23i2YkLsNzF/yyeBm9eSYBoY/HFV1Md1F0ElWBXW5E5mkdRtgjoRuq0yNb1CCFNWWlkn2gZknFJNUFJ8=", false},
	}

	// Run tests
	for _, test := range tests {
		if signature, err := SignMessage(test.inputKey, test.inputMessage); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error not expected but got: %s", t.Name(), test.inputKey, test.inputMessage, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error was expected", t.Name(), test.inputKey, test.inputMessage)
		} else if signature != test.expectedSignature {
			t.Errorf("%s Failed: [%s] [%s] inputted [%s] expected but got: %s", t.Name(), test.inputKey, test.inputMessage, test.expectedSignature, signature)
		}
	}
}

// ExampleSignMessage example using SignMessage()
func ExampleSignMessage() {
	signature, err := SignMessage("ef0b8bad0be285099534277fde328f8f19b3be9cadcd4c08e6ac0b5f863745ac", "This is a test message")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("signature created: %s", signature)
	// Output:signature created: H+zZagsyz7ioC/ZOa5EwsaKice0vs2BvZ0ljgkFHxD3vGsMlGeD4sXHEcfbI4h8lP29VitSBdf4A+nHXih7svf4=
}

// BenchmarkSignMessage benchmarks the method SignMessage()
func BenchmarkSignMessage(b *testing.B) {
	key, _ := CreatePrivateKeyString()
	for i := 0; i < b.N; i++ {
		_, _ = SignMessage(key, "This is a test message")
	}
}
