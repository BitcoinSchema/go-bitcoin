package bitcoin

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/chaincfg"
	"github.com/libsv/go-bk/wif"
)

// GenerateSharedKeyPair creates shared keys that can be used to encrypt/decrypt data
// that can be decrypted by yourself (privateKey) and also the owner of the given public key
func GenerateSharedKeyPair(privateKey *bec.PrivateKey,
	pubKey *bec.PublicKey,
) (*bec.PrivateKey, *bec.PublicKey) {
	return bec.PrivKeyFromBytes(
		bec.S256(),
		bec.GenerateSharedSecret(privateKey, pubKey),
	)
}

// PrivateKeyFromString turns a private key (hex encoded string) into an bec.PrivateKey
func PrivateKeyFromString(privateKey string) (*bec.PrivateKey, error) {
	if len(privateKey) == 0 {
		return nil, ErrPrivateKeyMissing
	}
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	x, y := bec.S256().ScalarBaseMult(privateKeyBytes)
	ecdsaPubKey := ecdsa.PublicKey{
		Curve: bec.S256(),
		X:     x,
		Y:     y,
	}
	return &bec.PrivateKey{PublicKey: ecdsaPubKey, D: new(big.Int).SetBytes(privateKeyBytes)}, nil
}

// CreatePrivateKey will create a new private key (*bec.PrivateKey)
func CreatePrivateKey() (*bec.PrivateKey, error) {
	return bec.NewPrivateKey(bec.S256())
}

// CreatePrivateKeyString will create a new private key (hex encoded)
func CreatePrivateKeyString() (string, error) {
	privateKey, err := CreatePrivateKey()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(privateKey.Serialise()), nil
}

// CreateWif will create a new WIF (*wif.WIF)
func CreateWif() (*wif.WIF, error) {
	privateKey, err := CreatePrivateKey()
	if err != nil {
		return nil, err
	}

	return wif.NewWIF(privateKey, &chaincfg.MainNet, false)
}

// CreateWifString will create a new WIF (string)
func CreateWifString() (string, error) {
	wifKey, err := CreateWif()
	if err != nil {
		return "", err
	}

	return wifKey.String(), nil
}

// PrivateAndPublicKeys will return both the private and public key in one method
// Expects a hex encoded privateKey
func PrivateAndPublicKeys(privateKey string) (*bec.PrivateKey, *bec.PublicKey, error) {
	// No key?
	if len(privateKey) == 0 {
		return nil, nil, ErrPrivateKeyMissing
	}

	// Decode the private key into bytes
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Get the public and private key from the bytes
	rawKey, publicKey := bec.PrivKeyFromBytes(bec.S256(), privateKeyBytes)
	return rawKey, publicKey, nil
}

// PrivateKeyToWif will convert a private key to a WIF (*wif.WIF)
func PrivateKeyToWif(privateKey string) (*wif.WIF, error) {
	// Missing private key
	if len(privateKey) == 0 {
		return nil, ErrPrivateKeyMissing
	}

	// Decode the private key
	decodedKey, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	// Get the private key from bytes
	rawKey, _ := bec.PrivKeyFromBytes(bec.S256(), decodedKey)

	// Create a new WIF (error never gets hit since (net) is set correctly)
	return wif.NewWIF(rawKey, &chaincfg.MainNet, false)
}

// PrivateKeyToWifString will convert a private key to a WIF (string)
func PrivateKeyToWifString(privateKey string) (string, error) {
	privateWif, err := PrivateKeyToWif(privateKey)
	if err != nil {
		return "", err
	}

	return privateWif.String(), nil
}

// WifToPrivateKey will convert a WIF to a private key (*bec.PrivateKey)
func WifToPrivateKey(wifKey string) (*bec.PrivateKey, error) {
	// Missing wif?
	if len(wifKey) == 0 {
		return nil, ErrWifMissing
	}

	// Decode the wif
	decodedWif, err := wif.DecodeWIF(wifKey)
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
	return hex.EncodeToString(privateKey.Serialise()), nil
}

// WifFromString will convert a WIF (string) to a WIF (*wif.WIF)
func WifFromString(wifKey string) (*wif.WIF, error) {
	// Missing wif?
	if len(wifKey) == 0 {
		return nil, ErrWifMissing
	}

	// Decode the WIF
	decodedWif, err := wif.DecodeWIF(wifKey)
	if err != nil {
		return nil, err
	}

	return decodedWif, nil
}
