package bitcoin

import (
	"encoding/hex"
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

	// Testing private methods
	var messageTests = []struct {
		inputMessage  string
		inputHeader   string
		expectedHash  string
		expectedError bool
	}{
		{"example message", hBSV, "002c483012a3c71d36349682d5ef6495926d4712b5cd2462f1e3c9f57bd4449f", false},
		{"", hBSV, "80e795d4a4caadd7047af389d9f7f220562feb6196032e2131e10563352c4bcc", false},
		{"example message", "", "f91e1e5a01b6aad5ec785946e4233b0613bf6183ffde8da9879949cbf7d7ca57", false},
		{"4qdD3HdK7SC4R9wTgfhr4QkNqRCKunbtRFlYPRYY6lGPiTbA9wZplnscnazyK0NMAx3KtvjDwWIX4J8djkSIYZaSNFEmztekNoe8NR0MLydp21U6Ayfm97oHelvTBcI5hQYccY45oI2KKEB1gyS0V6pbxoDtgjbCAGcnQvLB2iFykNcdU7A6Yntx812tKp90KilPADcEoKfkexMddqJ1pMz262MNhpTWmC4QOFMlB3xB5iTy2fxm6DgT3QLkiesk3kwM", "", "", true},
		{"", "4qdD3HdK7SC4R9wTgfhr4QkNqRCKunbtRFlYPRYY6lGPiTbA9wZplnscnazyK0NMAx3KtvjDwWIX4J8djkSIYZaSNFEmztekNoe8NR0MLydp21U6Ayfm97oHelvTBcI5hQYccY45oI2KKEB1gyS0V6pbxoDtgjbCAGcnQvLB2iFykNcdU7A6Yntx812tKp90KilPADcEoKfkexMddqJ1pMz262MNhpTWmC4QOFMlB3xB5iTy2fxm6DgT3QLkiesk3kwM", "", true},
	}

	// Run tests
	for _, test := range messageTests {
		if output, err := messageHash(test.inputMessage, test.inputHeader); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error not expected but got: %s", t.Name(), test.inputMessage, test.inputHeader, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] inputted and error was expected", t.Name(), test.inputMessage, test.inputHeader)
		} else if test.expectedHash != hex.EncodeToString(output) {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%s] expected, but got: %s", t.Name(), test.inputMessage, test.inputHeader, test.expectedHash, hex.EncodeToString(output))
		}
	}

	// Testing private methods
	var parseTests = []struct {
		inputSignature string
		expectedID     int
		expectedError  bool
	}{
		{"0", 0, true},
		{"1234567", 0, true},
		{"Op0u5nr4CPukjekOsOojjxuyEOG1HbIAf8qEteGXjA7tlqFinprrEcvdSlJOkZ8zMb", 0, true},
		{"000000000000000000000000000000000000000000000000000000000000000000", 0, true},
		{"-000-0-000-0---0--00-0-0-000-0--0-000-0-00-0-00---0-0-0--000-00-0-", 0, true},
		{"IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=", 1, false},
		{"IBDscOd-Ov4yrd-YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS-UJD+UK2SnNMn6Z3E4na8=", 0, true},
	}

	// Run tests
	for _, test := range parseTests {
		if _, output, err := parseSignature(test.inputSignature); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.inputSignature, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.inputSignature)
		} else if output != test.expectedID {
			t.Errorf("%s Failed: [%s] inputted and [%d] expected, but got: %d", t.Name(), test.inputSignature, test.expectedID, output)
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
