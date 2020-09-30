package bitcoin

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	"github.com/bitcoinsv/bsvd/bsvec"
)

// PrivateKey turns a private key string into an bsvec.PrivateKey
func PrivateKey(privateKey string) (*bsvec.PrivateKey, error) {
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
