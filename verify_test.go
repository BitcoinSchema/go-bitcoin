package bitcoin

import (
	"fmt"
	"testing"
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
		{"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ", "IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=", "Testing!", false},
		{"", "IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=", "Testing!", true},
		{"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ", "", "Testing!", true},
		{"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ", "IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=", "", true},
		{"0", "IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=", "Testing!", true},
		{"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ", "GBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4naZ=", "Testing!", true},
		{"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ", "GBD=", "Testing!", true},
		{"1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ", "GBse5w0f839t8wej8f2D=", "Testing!", true},
	}

	// Run tests
	for _, test := range tests {
		if err := VerifyMessage(test.inputAddress, test.inputSignature, test.inputData); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] [%s] inputted and error not expected but got: %s", t.Name(), test.inputAddress, test.inputSignature, test.inputData, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] [%s] inputted and error was expected", t.Name(), test.inputAddress, test.inputSignature, test.inputData)
		}
	}
}

// ExampleVerifyMessage example using VerifyMessage()
func ExampleVerifyMessage() {
	if err := VerifyMessage("1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ", "IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=", "Testing!"); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("verification passed")
	// Output:verification passed
}

// BenchmarkVerifyMessage benchmarks the method VerifyMessage()
func BenchmarkVerifyMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = VerifyMessage("1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ", "IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=", "Testing!")
	}
}
