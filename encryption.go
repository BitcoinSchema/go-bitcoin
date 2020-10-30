package bitcoin

import (
	"encoding/hex"

	"github.com/bitcoinsv/bsvd/bsvec"
)

// EncryptWithPrivateKey will encrypt the data using a given private key
func EncryptWithPrivateKey(privateKey *bsvec.PrivateKey, data string) (string, error) {

	// Encrypt using bsvec
	encryptedData, err := bsvec.Encrypt(privateKey.PubKey(), []byte(data))
	if err != nil {
		return "", err
	}

	// Return the hex encoded value
	return hex.EncodeToString(encryptedData), nil
}

// DecryptWithPrivateKey is a wrapper to decrypt the previously encrypted
// information, given a corresponding private key
func DecryptWithPrivateKey(privateKey *bsvec.PrivateKey, data string) (string, error) {

	// Decode the hex encoded string
	rawData, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	// Decrypt the data
	var decrypted []byte
	if decrypted, err = bsvec.Decrypt(privateKey, rawData); err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// EncryptWithPrivateKeyString is a convenience wrapper for EncryptWithPrivateKey()
func EncryptWithPrivateKeyString(privateKey, data string) (string, error) {

	// Get the private key from string
	rawPrivateKey, err := PrivateKeyFromString(privateKey)
	if err != nil {
		return "", err
	}

	// Encrypt using bsvec
	return EncryptWithPrivateKey(rawPrivateKey, data)
}

// DecryptWithPrivateKeyString is a convenience wrapper for DecryptWithPrivateKey()
func DecryptWithPrivateKeyString(privateKey, data string) (string, error) {

	// Get private key
	rawPrivateKey, _, err := PrivateAndPublicKeys(privateKey)
	if err != nil {
		return "", err
	}

	// Decrypt
	return DecryptWithPrivateKey(rawPrivateKey, data)
}

// EncryptShared will encrypt data and provide shared keys for decryption
func EncryptShared(user1PrivateKey *bsvec.PrivateKey, user2PubKey *bsvec.PublicKey, data []byte) (
	*bsvec.PrivateKey, *bsvec.PublicKey, []byte, error) {

	// Generate shared keys that can be decrypted by either user
	sharedPrivKey, sharedPubKey := GenerateSharedKeyPair(user1PrivateKey, user2PubKey)

	// Encrypt data with shared key
	encryptedData, err := bsvec.Encrypt(sharedPubKey, data)
	return sharedPrivKey, sharedPubKey, encryptedData, err
}

// EncryptSharedString will encrypt a string to a hex encoded encrypted payload, and provide shared keys for decryption
func EncryptSharedString(user1PrivateKey *bsvec.PrivateKey, user2PubKey *bsvec.PublicKey, data string) (
	*bsvec.PrivateKey, *bsvec.PublicKey, string, error) {

	// Generate shared keys that can be decrypted by either user
	sharedPrivKey, sharedPubKey := GenerateSharedKeyPair(user1PrivateKey, user2PubKey)

	// Encrypt data with shared key
	encryptedData, err := bsvec.Encrypt(sharedPubKey, []byte(data))

	return sharedPrivKey, sharedPubKey, hex.EncodeToString(encryptedData), err
}
