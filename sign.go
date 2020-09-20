package bitcoin

import (
	"encoding/base64"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvd/chaincfg/chainhash"
)

// SignMessage signs a string with the provided private key using Bitcoin Signed Message encoding
func SignMessage(privKey string, message string) string {
	prefixBytes := []byte("Bitcoin Signed Message:\n")
	messageBytes := []byte(message)
	bytes := []byte{}
	bytes = append(bytes, byte(len(prefixBytes)))
	bytes = append(bytes, prefixBytes...)
	bytes = append(bytes, byte(len(messageBytes)))
	bytes = append(bytes, messageBytes...)
	ecdsaPrivKey := PrivateKey(privKey)

	sigbytes, err := bsvec.SignCompact(bsvec.S256(), ecdsaPrivKey, chainhash.DoubleHashB(bytes), true)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(sigbytes)
}
