package bitcoin

import (
	"testing"
)

// Identity Private Key
const privKey = "E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F2B75"

func TestSignMessage(t *testing.T) {

	// privKey string, message string, compress bool
	sig := SignMessage(privKey, "Testing!")

	t.Log("got sig", sig)

	if sig != "IBDscOd/Ov4yrd/YXantqajSAnW4fudpfr2KQy5GNo9pZybF12uNaal4KI822UpQLS/UJD+UK2SnNMn6Z3E4na8=" {
		t.Error("Failed to generate expected signature", sig)
	}
}
