package bitcoin

import (
	"crypto/ecdsa"
	"encoding/base64"
	"math/big"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvd/chaincfg/chainhash"
)

func SignMessage(privKey string, message string, compress bool) string {
	prefixBytes := []byte("Bitcoin Signed Message:\n")
	messageBytes := []byte(message)
	bytes := []byte{}
	bytes = append(bytes, byte(len(prefixBytes)))
	bytes = append(bytes, prefixBytes...)
	bytes = append(bytes, byte(len(messageBytes)))
	bytes = append(bytes, messageBytes...)
	privKeyBytes := HexDecode(privKey)
	x, y := bsvec.S256().ScalarBaseMult(privKeyBytes)
	ecdsaPubKey := ecdsa.PublicKey{
		Curve: bsvec.S256(),
		X:     x,
		Y:     y,
	}
	ecdsaPrivKey := &bsvec.PrivateKey{
		PublicKey: ecdsaPubKey,
		D:         new(big.Int).SetBytes(privKeyBytes),
	}

	sigbytes, err := bsvec.SignCompact(bsvec.S256(), ecdsaPrivKey, chainhash.DoubleHashB(bytes), compress)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(sigbytes)
}
