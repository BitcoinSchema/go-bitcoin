package bitcoin

import (
	"testing"
)

// Test address
const address = "1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi2"

func TestValidA58(t *testing.T) {

	valid, err := ValidA58([]byte(address))

	if !valid {
		t.Error("Failed to validate address", err)
	}
}
