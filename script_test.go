package bitcoin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestScriptFromAddress will test the method ScriptFromAddress()
func TestScriptFromAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputAddress   string
		expectedScript string
		expectedError  bool
	}{
		{"empty address", "", "", true},
		{caseSingleZero, "0", "", true},
		{"too short", "1234567", "", true},
		{"valid address 1", "1HRVqUGDzpZSMVuNSZxJVaB9xjneEShfA7", "76a914b424110292f4ea2ac92beb9e83cf5e6f0fa2996388ac", false},
		{"valid address 2", "13Rj7G3pn2GgG8KE6SFXLc7dCJdLNnNK7M", "76a9141a9d62736746f85ca872dc555ff51b1fed2471e288ac", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			script, err := ScriptFromAddress(test.inputAddress)
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.expectedScript, script)
		})
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
	for b.Loop() {
		_, _ = ScriptFromAddress("1HRVqUGDzpZSMVuNSZxJVaB9xjneEShfA7")
	}
}
