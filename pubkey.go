package bitcoin

import (
	"encoding/hex"
)

// PubKeyFromPrivateKey will derive a pubKey (hex encoded) from a given private key
func PubKeyFromPrivateKey(privateKey string) (string, error) {
	rawKey, err := PrivateKeyFromString(privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rawKey.PubKey().SerializeCompressed()), nil
}
