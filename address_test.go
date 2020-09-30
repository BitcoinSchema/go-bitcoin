package bitcoin

import (
	"fmt"
	"testing"
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
