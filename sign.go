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
//
// Spec: https://docs.moneybutton.com/docs/bsv-message.html
func SignMessage(privateKey string, message string) (string, error) {
	if len(privateKey) == 0 {
		return "", errors.New("privateKey is empty")
	}

	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, hBSV)
	wire.WriteVarString(&buf, 0, message)

	messageHash := chainhash.DoubleHashB(buf.Bytes())

	ecdsaPrivateKey, err := PrivateKeyFromString(privateKey)
	if err != nil {
		return "", err
	}
	var sigBytes []byte
	sigBytes, err = bsvec.SignCompact(bsvec.S256(), ecdsaPrivateKey, messageHash, true)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sigBytes), nil
}
