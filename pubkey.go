package bitcoin

import (
	"encoding/hex"

	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
)

// PubKeyFromPrivateKeyString will derive a pubKey (hex encoded) from a given private key
func PubKeyFromPrivateKeyString(privateKey string, compressed bool) (string, error) {
	rawKey, err := PrivateKeyFromString(privateKey)
	if err != nil {
		return "", err
	}

	return PubKeyFromPrivateKey(rawKey, compressed), nil
}

// PubKeyFromPrivateKey will derive a pubKey (hex encoded) from a given private key
func PubKeyFromPrivateKey(privateKey *ec.PrivateKey, compressed bool) string {
	if compressed {
		return hex.EncodeToString(privateKey.PubKey().Compressed())
	}
	return hex.EncodeToString(privateKey.PubKey().Uncompressed())
}

// PubKeyFromString will convert a pubKey (string) into a pubkey (*ec.PublicKey)
func PubKeyFromString(pubKey string) (*ec.PublicKey, error) {
	// Invalid pubKey
	if len(pubKey) == 0 {
		return nil, ErrMissingPubKey
	}

	// Decode from hex string
	decoded, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}

	// Parse into a pubKey
	return ec.ParsePubKey(decoded)
}
