package bitcoin

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/bitcoinsv/bsvd/bsvec"
)

// PrivateKeyFromString turns a private key (hex encoded string) into an bsvec.PrivateKey
func PrivateKeyFromString(privateKey string) (*bsvec.PrivateKey, error) {
	if len(privateKey) == 0 {
		return nil, errors.New("privateKey is missing")
	}
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

// PrivateAndPublicKeys will return both the private and public key in one method
// Expects a hex encoded privateKey
func PrivateAndPublicKeys(privateKey string) (*bsvec.PrivateKey, *bsvec.PublicKey, error) {

	// No key?
	if len(privateKey) == 0 {
		return nil, nil, errors.New("missing privateKey")
	}

	// Decode the private key into bytes
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Get the public and private key from the bytes
	rawKey, publicKey := bsvec.PrivKeyFromBytes(bsvec.S256(), privateKeyBytes)
	return rawKey, publicKey, nil
}
