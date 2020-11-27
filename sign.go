package bitcoin

import (
	"bytes"
	"encoding/base64"
	"errors"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvd/chaincfg/chainhash"
	"github.com/bitcoinsv/bsvd/wire"
)

// SignMessage signs a string with the provided private key using Bitcoin Signed Message encoding
// sigRefCompressedKey bool determines whether the signature will reference a compressed or uncompresed key
// Spec: https://docs.moneybutton.com/docs/bsv-message.html
func SignMessage(privateKey string, message string, sigRefCompressedKey bool) (string, error) {
	if len(privateKey) == 0 {
		return "", errors.New("privateKey is empty")
	}

	var buf bytes.Buffer
	var err error
	if err = wire.WriteVarString(&buf, 0, hBSV); err != nil {
		return "", err
	}
	if err = wire.WriteVarString(&buf, 0, message); err != nil {
		return "", err
	}

	// Create the hash
	messageHash := chainhash.DoubleHashB(buf.Bytes())

	// Get the private key
	var ecdsaPrivateKey *bsvec.PrivateKey
	if ecdsaPrivateKey, err = PrivateKeyFromString(privateKey); err != nil {
		return "", err
	}

	// Sign
	var sigBytes []byte
	if sigBytes, err = bsvec.SignCompact(bsvec.S256(), ecdsaPrivateKey, messageHash, sigRefCompressedKey); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sigBytes), nil
}
