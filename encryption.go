package bitcoin

import (
	"encoding/hex"

	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
)

// EncryptWithPrivateKey will encrypt the data using a given private key
func EncryptWithPrivateKey(privateKey *ec.PrivateKey, data string) (string, error) {
	// Encrypt using ec
	encryptedData, err := eciesEncrypt(privateKey.PubKey(), []byte(data))
	if err != nil {
		return "", err
	}

	// Return the hex encoded value
	return hex.EncodeToString(encryptedData), nil
}

// DecryptWithPrivateKey is a wrapper to decrypt the previously encrypted
// information, given a corresponding private key
func DecryptWithPrivateKey(privateKey *ec.PrivateKey, data string) (string, error) {
	// Decode the hex encoded string
	rawData, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	// Decrypt the data
	var decrypted []byte
	if decrypted, err = eciesDecrypt(privateKey, rawData); err != nil {
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

	// Encrypt using bec
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
func EncryptShared(user1PrivateKey *ec.PrivateKey, user2PubKey *ec.PublicKey, data []byte) (
	*ec.PrivateKey, *ec.PublicKey, []byte, error,
) {
	// Generate shared keys that can be decrypted by either user
	sharedPrivKey, sharedPubKey := GenerateSharedKeyPair(user1PrivateKey, user2PubKey)

	// Encrypt data with shared key
	encryptedData, err := eciesEncrypt(sharedPubKey, data)
	return sharedPrivKey, sharedPubKey, encryptedData, err
}

// EncryptSharedString will encrypt a string to a hex encoded encrypted payload, and provide shared keys for decryption
func EncryptSharedString(user1PrivateKey *ec.PrivateKey, user2PubKey *ec.PublicKey, data string) (
	*ec.PrivateKey, *ec.PublicKey, string, error,
) {
	// Generate shared keys that can be decrypted by either user
	sharedPrivKey, sharedPubKey := GenerateSharedKeyPair(user1PrivateKey, user2PubKey)

	// Encrypt data with shared key
	encryptedData, err := eciesEncrypt(sharedPubKey, []byte(data))

	return sharedPrivKey, sharedPubKey, hex.EncodeToString(encryptedData), err
}
