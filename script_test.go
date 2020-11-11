package bitcoin

import (
	"fmt"
	"testing"
)

// TestScriptFromAddress will test the method ScriptFromAddress()
func TestScriptFromAddress(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputAddress   string
		expectedScript string
		expectedError  bool
	}{
		{"", "", true},
		{"0", "", true},
		{"1234567", "", true},
		{"1HRVqUGDzpZSMVuNSZxJVaB9xjneEShfA7", "76a914b424110292f4ea2ac92beb9e83cf5e6f0fa2996388ac", false},
		{"13Rj7G3pn2GgG8KE6SFXLc7dCJdLNnNK7M", "76a9141a9d62736746f85ca872dc555ff51b1fed2471e288ac", false},
	}

	// Run tests
	for _, test := range tests {
		if script, err := ScriptFromAddress(test.inputAddress); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.inputAddress, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%v] inputted and error was expected", t.Name(), test.inputAddress)
		} else if script != test.expectedScript {
			t.Fatalf("%s Failed: [%v] inputted [%s] expected but failed comparison of scripts, got: %s", t.Name(), test.inputAddress, test.expectedScript, script)
		}
	}
}

// ExampleScriptFromAddress example using ScriptFromAddress()
func ExampleScriptFromAddress() {
	script, err := ScriptFromAddress("1HRVqUGDzpZSMVuNSZxJVaB9xjneEShfA7")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("script generated: %s", script)
	// Output:script generated: 76a914b424110292f4ea2ac92beb9e83cf5e6f0fa2996388ac
}

// BenchmarkScriptFromAddress benchmarks the method ScriptFromAddress()
func BenchmarkScriptFromAddress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ScriptFromAddress("1HRVqUGDzpZSMVuNSZxJVaB9xjneEShfA7")
	}
}
