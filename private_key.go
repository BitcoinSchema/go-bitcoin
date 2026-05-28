package bitcoin

import (
	"encoding/hex"

	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	chaincfg "github.com/bsv-blockchain/go-sdk/transaction/chaincfg"
)

// GenerateSharedKeyPair creates shared keys that can be used to encrypt/decrypt data
// that can be decrypted by yourself (privateKey) and also the owner of the given public key
func GenerateSharedKeyPair(privateKey *ec.PrivateKey,
	pubKey *ec.PublicKey,
) (*ec.PrivateKey, *ec.PublicKey) {
	return ec.PrivateKeyFromBytes(
		generateSharedSecret(privateKey, pubKey),
	)
}

// PrivateKeyFromString turns a private key (hex encoded string) into an ec.PrivateKey
func PrivateKeyFromString(privateKey string) (*ec.PrivateKey, error) {
	if len(privateKey) == 0 {
		return nil, ErrPrivateKeyMissing
	}
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	rawKey, _ := ec.PrivateKeyFromBytes(privateKeyBytes)
	return rawKey, nil
}

// CreatePrivateKey will create a new private key (*ec.PrivateKey)
func CreatePrivateKey() (*ec.PrivateKey, error) {
	return ec.NewPrivateKey()
}

// CreatePrivateKeyString will create a new private key (hex encoded)
func CreatePrivateKeyString() (string, error) {
	privateKey, err := CreatePrivateKey()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(privateKey.Serialize()), nil
}

// CreateWif will create a new WIF (*WIF)
func CreateWif() (*WIF, error) {
	privateKey, err := CreatePrivateKey()
	if err != nil {
		return nil, err
	}

	return NewWIF(privateKey, &chaincfg.MainNet, false)
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
func PrivateAndPublicKeys(privateKey string) (*ec.PrivateKey, *ec.PublicKey, error) {
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
	rawKey, publicKey := ec.PrivateKeyFromBytes(privateKeyBytes)
	return rawKey, publicKey, nil
}

// PrivateKeyToWif will convert a private key to a WIF (*WIF)
func PrivateKeyToWif(privateKey string) (*WIF, error) {
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
	rawKey, _ := ec.PrivateKeyFromBytes(decodedKey)

	// Create a new WIF (error never gets hit since (net) is set correctly)
	return NewWIF(rawKey, &chaincfg.MainNet, false)
}

// PrivateKeyToWifString will convert a private key to a WIF (string)
func PrivateKeyToWifString(privateKey string) (string, error) {
	privateWif, err := PrivateKeyToWif(privateKey)
	if err != nil {
		return "", err
	}

	return privateWif.String(), nil
}

// WifToPrivateKey will convert a WIF to a private key (*ec.PrivateKey)
func WifToPrivateKey(wifKey string) (*ec.PrivateKey, error) {
	// Missing wif?
	if len(wifKey) == 0 {
		return nil, ErrWifMissing
	}

	// Decode the wif
	decodedWif, err := DecodeWIF(wifKey)
	if err != nil {
		return nil, err
	}

	// Return the private key
	return decodedWif.PrivKey, nil
}

// WifToPrivateKeyString will convert a WIF to private key (string)
func WifToPrivateKeyString(wifKey string) (string, error) {
	// Convert the wif to private key
	privateKey, err := WifToPrivateKey(wifKey)
	if err != nil {
		return "", err
	}

	// Return the hex (string) version of the private key
	return hex.EncodeToString(privateKey.Serialize()), nil
}

// WifFromString will convert a WIF (string) to a WIF (*WIF)
func WifFromString(wifKey string) (*WIF, error) {
	// Missing wif?
	if len(wifKey) == 0 {
		return nil, ErrWifMissing
	}

	// Decode the WIF
	decodedWif, err := DecodeWIF(wifKey)
	if err != nil {
		return nil, err
	}

	return decodedWif, nil
}
