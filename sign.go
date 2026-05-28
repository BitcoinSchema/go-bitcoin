package bitcoin

import (
	"encoding/base64"

	bsm "github.com/bsv-blockchain/go-sdk/compat/bsm"
)

// SignMessage signs a string with the provided private key using Bitcoin Signed Message encoding
// sigRefCompressedKey bool determines whether the signature will reference a compressed or uncompresed key
// Spec: https://docs.moneybutton.com/docs/bsv-message.html
func SignMessage(privateKey, message string, sigRefCompressedKey bool) (string, error) {
	if len(privateKey) == 0 {
		return "", ErrPrivateKeyMissing
	}

	// Get the private key
	ecdsaPrivateKey, err := PrivateKeyFromString(privateKey)
	if err != nil {
		return "", err
	}

	// Sign using Bitcoin Signed Message encoding
	var sigBytes []byte
	if sigBytes, err = bsm.SignMessageWithCompression(ecdsaPrivateKey, []byte(message), sigRefCompressedKey); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sigBytes), nil
}
