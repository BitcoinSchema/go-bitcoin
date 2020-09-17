package bitcoin

import (
	"testing"
)

// Identity Private Key
const privKey = "E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F2B75"

func TestSignMessage(t *testing.T) {

	// privKey string, message string, compress bool

	tx, err := SignMessage(privKey, "Test message!", true)
	if err != nil {
		t.Error("Failed to sign message:", err)
	}

	if tx.GetTxID() != "e97ed4acb8d01a822dd5070e6addf762949f48a696311a954b85cd4a9c993a23" {
		t.Error("Failed to generate expected signature")
	}
}
