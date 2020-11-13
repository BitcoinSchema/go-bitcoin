package bitcoin

import (
	"errors"

	"github.com/libsv/go-bt"
)

// ScriptFromAddress will create an output P2PKH script from an address string
func ScriptFromAddress(address string) (string, error) {
	// Missing address?
	if len(address) == 0 {
		return "", errors.New("missing address")
	}

	// Generate a script from address
	rawScript, err := bt.NewP2PKHOutputFromAddress(address, 0)
	if err != nil {
		return "", err
	}

	// Return the string version
	return rawScript.GetLockingScriptHexString(), nil
}
