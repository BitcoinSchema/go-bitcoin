package bitcoin

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/bitcoinsv/bsvd/bsvec"
)

// PrivateKey turns a private key string into an bsvec.PrivateKey
func PrivateKey(privKey string) (ecdsaPrivKey *bsvec.PrivateKey) {
	privKeyBytes := HexDecode(privKey)
	x, y := bsvec.S256().ScalarBaseMult(privKeyBytes)
	ecdsaPubKey := ecdsa.PublicKey{
		Curve: bsvec.S256(),
		X:     x,
		Y:     y,
	}
	ecdsaPrivKey = &bsvec.PrivateKey{
		PublicKey: ecdsaPubKey,
		D:         new(big.Int).SetBytes(privKeyBytes),
	}
	return
}
