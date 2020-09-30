package bitcoin

import (
	"fmt"
	"testing"
)

// TestVerifyMessage will test the method VerifyMessage()
func TestVerifyMessage(t *testing.T) {

	var sig = "IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8="
	var address = "1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ"
	var data = "Testing!"

	if err := VerifyMessage(address, sig, data); err != nil {
		t.Fatalf("failed to verify message: %s", err.Error())
	}
}

// ExampleVerifyMessage example using VerifyMessage()
func ExampleVerifyMessage() {
	var sig = "IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8="
	var address = "1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ"
	var data = "Testing!"

	if err := VerifyMessage(address, sig, data); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("verification passed")
	// Output:verification passed
}

// BenchmarkVerifyMessage benchmarks the method VerifyMessage()
func BenchmarkVerifyMessage(b *testing.B) {
	var sig = "IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8="
	var address = "1FiyJnrgwBc3Ff83V1yRWAkmXBdGrDQnXQ"
	var data = "Testing!"

	for i := 0; i < b.N; i++ {
		_ = VerifyMessage(address, sig, data)
	}
}
