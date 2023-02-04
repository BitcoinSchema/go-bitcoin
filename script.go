package bitcoin

import (
	"github.com/libsv/go-bt/v2/bscript"
)

// ScriptFromAddress will create an output P2PKH script from an address string
func ScriptFromAddress(address string) (string, error) {
	// Missing address?
	if len(address) == 0 {
		return "", ErrMissingAddress
	}

	// Generate a script from address
	rawScript, err := bscript.NewP2PKHFromAddress(address)
	if err != nil {
		return "", err
	}

	// Return the string version
	return rawScript.String(), nil
}
