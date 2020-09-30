package bitcoin

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	"github.com/bitcoinsv/bsvd/bsvec"
)

// PrivateKeyFromString turns a private key (hex encoded string) into an bsvec.PrivateKey
func PrivateKeyFromString(privateKey string) (*bsvec.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	x, y := bsvec.S256().ScalarBaseMult(privateKeyBytes)
	ecdsaPubKey := ecdsa.PublicKey{
		Curve: bsvec.S256(),
		X:     x,
		Y:     y,
	}
	return &bsvec.PrivateKey{PublicKey: ecdsaPubKey, D: new(big.Int).SetBytes(privateKeyBytes)}, nil
}

// CreatePrivateKey will create a new private key
func CreatePrivateKey() (*bsvec.PrivateKey, error) {
	return bsvec.NewPrivateKey(bsvec.S256())
}

// CreatePrivateKeyString will create a new private key (hex encoded)
func CreatePrivateKeyString() (string, error) {
	privateKey, err := bsvec.NewPrivateKey(bsvec.S256())
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(privateKey.Serialize()), nil
}
