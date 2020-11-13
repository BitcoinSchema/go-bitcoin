package bitcoin

import (
	"errors"

	"github.com/libsv/go-bt/bscript"
)

// ScriptFromAddress will create an output P2PKH script from an address string
func ScriptFromAddress(address string) (string, error) {
	// Missing address?
	if len(address) == 0 {
		return "", errors.New("missing address")
	}

	// Generate a script from address
	rawScript, err := bscript.NewP2PKHFromAddress(address)
	if err != nil {
		return "", err
	}

	// Return the string version
	return rawScript.ToString(), nil
}
