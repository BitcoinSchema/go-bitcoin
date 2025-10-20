package bitcoin

import (
	"encoding/hex"

	"github.com/libsv/go-bk/bec"
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
func PubKeyFromPrivateKey(privateKey *bec.PrivateKey, compressed bool) string {
	if compressed {
		return hex.EncodeToString(privateKey.PubKey().SerialiseCompressed())
	}
	return hex.EncodeToString(privateKey.PubKey().SerialiseUncompressed())
}

// PubKeyFromString will convert a pubKey (string) into a pubkey (*bec.PublicKey)
func PubKeyFromString(pubKey string) (*bec.PublicKey, error) {
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
	return bec.ParsePubKey(decoded, bec.S256())
}
