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

// privateKeyFromHex decodes a hex-encoded private key string into its
// ec.PrivateKey and matching ec.PublicKey. It returns ErrPrivateKeyMissing for
// an empty string and the hex decode error for malformed input.
func privateKeyFromHex(privateKey string) (*ec.PrivateKey, *ec.PublicKey, error) {
	if privateKey == "" {
		return nil, nil, ErrPrivateKeyMissing
	}
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, nil, err
	}
	rawKey, publicKey := ec.PrivateKeyFromBytes(privateKeyBytes)
	return rawKey, publicKey, nil
}

// PrivateKeyFromString turns a private key (hex encoded string) into an ec.PrivateKey
func PrivateKeyFromString(privateKey string) (*ec.PrivateKey, error) {
	rawKey, _, err := privateKeyFromHex(privateKey)
	if err != nil {
		return nil, err
	}
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

// CreateWif will create a new uncompressed mainnet WIF (*WIF).
//
// Use CreateWifWithCompression to choose compressed (52-char, K/L...) vs
// uncompressed (51-char, 5...) public-key encoding.
func CreateWif() (*WIF, error) {
	return CreateWifWithCompression(false)
}

// CreateWifWithCompression will create a new random mainnet WIF (*WIF) using the
// chosen public-key compression. Compressed WIFs are 52 characters and start
// with K or L; uncompressed WIFs are 51 characters and start with 5.
func CreateWifWithCompression(compress bool) (*WIF, error) {
	privateKey, err := CreatePrivateKey()
	if err != nil {
		return nil, err
	}

	return NewWIF(privateKey, &chaincfg.MainNet, compress)
}

// CreateWifString will create a new uncompressed mainnet WIF (string).
//
// Use CreateWifStringWithCompression to choose compression.
func CreateWifString() (string, error) {
	return CreateWifStringWithCompression(false)
}

// CreateWifStringWithCompression will create a new random mainnet WIF (string)
// using the chosen public-key compression.
func CreateWifStringWithCompression(compress bool) (string, error) {
	wifKey, err := CreateWifWithCompression(compress)
	if err != nil {
		return "", err
	}

	return wifKey.String(), nil
}

// PrivateAndPublicKeys will return both the private and public key in one method
// Expects a hex encoded privateKey
func PrivateAndPublicKeys(privateKey string) (*ec.PrivateKey, *ec.PublicKey, error) {
	return privateKeyFromHex(privateKey)
}

// PrivateKeyToWif will convert a private key to an uncompressed mainnet WIF (*WIF).
//
// Use PrivateKeyToWifWithCompression to choose compression.
func PrivateKeyToWif(privateKey string) (*WIF, error) {
	return PrivateKeyToWifWithCompression(privateKey, false)
}

// PrivateKeyToWifWithCompression will convert a hex private key to a mainnet WIF
// (*WIF) using the chosen public-key compression.
func PrivateKeyToWifWithCompression(privateKey string, compress bool) (*WIF, error) {
	// Decode the private key (returns ErrPrivateKeyMissing for an empty string)
	rawKey, _, err := privateKeyFromHex(privateKey)
	if err != nil {
		return nil, err
	}

	// Create a new WIF (error never gets hit since (net) is set correctly)
	return NewWIF(rawKey, &chaincfg.MainNet, compress)
}

// PrivateKeyToWifString will convert a private key to an uncompressed mainnet WIF (string).
//
// Use PrivateKeyToWifStringWithCompression to choose compression.
func PrivateKeyToWifString(privateKey string) (string, error) {
	return PrivateKeyToWifStringWithCompression(privateKey, false)
}

// PrivateKeyToWifStringWithCompression will convert a hex private key to a
// mainnet WIF (string) using the chosen public-key compression.
func PrivateKeyToWifStringWithCompression(privateKey string, compress bool) (string, error) {
	privateWif, err := PrivateKeyToWifWithCompression(privateKey, compress)
	if err != nil {
		return "", err
	}

	return privateWif.String(), nil
}

// WifToPrivateKey will convert a WIF to a private key (*ec.PrivateKey)
func WifToPrivateKey(wifKey string) (*ec.PrivateKey, error) {
	// Missing wif?
	if wifKey == "" {
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
	if wifKey == "" {
		return nil, ErrWifMissing
	}

	// Decode the WIF
	decodedWif, err := DecodeWIF(wifKey)
	if err != nil {
		return nil, err
	}

	return decodedWif, nil
}
