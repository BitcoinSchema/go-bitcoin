package bitcoin

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvd/chaincfg"
	"github.com/bitcoinsv/bsvutil"
)

// GenerateSharedKeyPair creates shared keys that can be used to encrypt/decrypt data
// that can be decrypted by yourself (privateKey) and also the owner of the given public key
func GenerateSharedKeyPair(privateKey *bsvec.PrivateKey,
	pubKey *bsvec.PublicKey) (*bsvec.PrivateKey, *bsvec.PublicKey) {
	return bsvec.PrivKeyFromBytes(
		bsvec.S256(),
		bsvec.GenerateSharedSecret(privateKey, pubKey),
	)
}

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

// CreatePrivateKey will create a new private key (*bsvec.PrivateKey)
func CreatePrivateKey() (*bsvec.PrivateKey, error) {
	return bsvec.NewPrivateKey(bsvec.S256())
}

// CreatePrivateKeyString will create a new private key (hex encoded)
func CreatePrivateKeyString() (string, error) {
	privateKey, err := CreatePrivateKey()
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

// PrivateKeyToWif will convert a private key to a WIF (*bsvutil.WIF)
func PrivateKeyToWif(privateKey string) (*bsvutil.WIF, error) {

	// Missing private key
	if len(privateKey) == 0 {
		return nil, errors.New("missing privateKey")
	}

	// Decode the private key
	decodedKey, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	// Get the private key from bytes
	rawKey, _ := bsvec.PrivKeyFromBytes(bsvec.S256(), decodedKey)

	// Create a new WIF (error never gets hit since (net) is set correctly)
	return bsvutil.NewWIF(rawKey, &chaincfg.MainNetParams, false)
}

// PrivateKeyToWifString will convert a private key to a WIF (string)
func PrivateKeyToWifString(privateKey string) (string, error) {
	wif, err := PrivateKeyToWif(privateKey)
	if err != nil {
		return "", err
	}

	return wif.String(), nil
}

// WifToPrivateKey will convert a WIF to a private key (*bsvec.PrivateKey)
func WifToPrivateKey(wif string) (*bsvec.PrivateKey, error) {

	// Missing wif?
	if len(wif) == 0 {
		return nil, errors.New("missing wif")
	}

	// Decode the wif
	decodedWif, err := bsvutil.DecodeWIF(wif)
	if err != nil {
		return nil, err
	}

	// Return the private key
	return decodedWif.PrivKey, nil
}

// WifToPrivateKeyString will convert a WIF to private key (string)
func WifToPrivateKeyString(wif string) (string, error) {

	// Convert the wif to private key
	privateKey, err := WifToPrivateKey(wif)
	if err != nil {
		return "", err
	}

	// Return the hex (string) version of the private key
	return hex.EncodeToString(privateKey.Serialize()), nil
}
