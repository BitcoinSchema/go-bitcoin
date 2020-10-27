package bitcoin

import (
	"encoding/hex"
	"errors"

	"github.com/bitcoinsv/bsvd/bsvec"
)

// PubKeyFromPrivateKeyString will derive a pubKey (hex encoded) from a given private key
func PubKeyFromPrivateKeyString(privateKey string) (string, error) {
	rawKey, err := PrivateKeyFromString(privateKey)
	if err != nil {
		return "", err
	}

	return PubKeyFromPrivateKey(rawKey), nil
}

// PubKeyFromPrivateKey will derive a pubKey (hex encoded) from a given private key
func PubKeyFromPrivateKey(privateKey *bsvec.PrivateKey) string {
	return hex.EncodeToString(privateKey.PubKey().SerializeCompressed())
}

// PubKeyFromString will convert a pubKey (string) into a pubkey (*bsvec.PublicKey)
func PubKeyFromString(pubKey string) (*bsvec.PublicKey, error) {

	// Invalid pubKey
	if len(pubKey) == 0 {
		return nil, errors.New("missing pubkey")
	}

	// Decode from hex string
	decoded, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}

	// Parse into a pubKey
	return bsvec.ParsePubKey(decoded, bsvec.S256())
}
