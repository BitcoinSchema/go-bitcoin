package bitcoin

import (
	"testing"
)

// Identity Private Key
const privateKey = "E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F2B75"

// TestSignMessage will test the method SignMessage()
func TestSignMessage(t *testing.T) {

	// privateKey string, message string, compress bool
	sig, err := SignMessage(privateKey, "Testing!")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if sig != "IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=" {
		t.Error("Failed to generate expected signature", sig)
	}
}
