package bitcoin

import (
	"encoding/base64"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvd/chaincfg/chainhash"
)

// SignMessage signs a string with the provided private key using Bitcoin Signed Message encoding
func SignMessage(privateKey string, message string) (string, error) {
	prefixBytes := []byte(hBSV)
	messageBytes := []byte(message)
	var bytes []byte
	bytes = append(bytes, byte(len(prefixBytes)))
	bytes = append(bytes, prefixBytes...)
	bytes = append(bytes, byte(len(messageBytes)))
	bytes = append(bytes, messageBytes...)
	ecdsaPrivateKey, err := PrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	var sigBytes []byte
	sigBytes, err = bsvec.SignCompact(bsvec.S256(), ecdsaPrivateKey, chainhash.DoubleHashB(bytes), true)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sigBytes), nil
}
