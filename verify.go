package bitcoin

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvd/chaincfg"
	"github.com/bitcoinsv/bsvd/chaincfg/chainhash"
	"github.com/bitcoinsv/bsvd/wire"
	"github.com/bitcoinsv/bsvutil"
)

const (
	// H_BSV is the magic header string required fore Bitcoin Signed Messages
	hBSV string = "Bitcoin Signed Message:\n"
)

// VerifyMessage verifies a string and address against the provided
// signature and assumes Bitcoin Signed Message encoding
//
// Error will occur if verify fails or verification is not successful (no bool)
// Spec: https://docs.moneybutton.com/docs/bsv-message.html
func VerifyMessage(address, sig, data string) error {

	decodedSig, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}

	// Validate the signature - this just shows that it was valid at all.
	// we will compare it with the key next.
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, hBSV)
	// TODO!!! The 0 here controls the variable length integer
	if len(data) > 0xFD {
		// https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
		log.Println("Long message! Change varint!")
	}
	wire.WriteVarString(&buf, 0, data)
	expectedMessageHash := chainhash.DoubleHashB(buf.Bytes())

	pk, wasCompressed, err := bsvec.RecoverCompact(bsvec.S256(), decodedSig, expectedMessageHash)
	if err != nil {
		return err
	}

	// Reconstruct the pubkey hash.
	var serializedPK []byte
	if wasCompressed {
		serializedPK = pk.SerializeCompressed()
	} else {
		serializedPK = pk.SerializeUncompressed()
	}
	bsvecAddress, err := bsvutil.NewAddressPubKey(serializedPK, &chaincfg.MainNetParams)
	if err != nil {
		return err
	}

	// Return nil if addresses match.
	if bsvecAddress.EncodeAddress() == address {
		return nil
	}
	return fmt.Errorf("address: %s not found vs %s", address, bsvecAddress.EncodeAddress())
}

// VerifyMessageDER will take a message string, a public key string and a signature string
// (in strict DER format) and verify that the message was signed by the public key.
//
// Copyright (c) 2019 Bitcoin Association
// License: https://github.com/bitcoin-sv/merchantapi-reference/blob/master/LICENSE
//
// Source: https://github.com/bitcoin-sv/merchantapi-reference/blob/master/handler/global.go
func VerifyMessageDER(hash [32]byte, pubKey string, signature string) (verified bool, err error) {

	// Decode the signature string
	var sigBytes []byte
	if sigBytes, err = hex.DecodeString(signature); err != nil {
		return
	}

	// Parse the signature
	var sig *bsvec.Signature
	if sig, err = bsvec.ParseDERSignature(sigBytes, bsvec.S256()); err != nil {
		return
	}

	// Decode the pubKey
	var pubKeyBytes []byte
	if pubKeyBytes, err = hex.DecodeString(pubKey); err != nil {
		return
	}

	// Parse the pubKey
	var rawPubKey *bsvec.PublicKey
	if rawPubKey, err = bsvec.ParsePubKey(pubKeyBytes, bsvec.S256()); err != nil {
		return
	}

	// Verify the signature against the pubKey
	verified = sig.Verify(hash[:], rawPubKey)
	return
}

// messageHash will compute a hash for the given message & header
// func messageHash(message, header string) ([]byte, error) {
// 	headerLength := len(header)
// 	if headerLength >= 0xfd {
// 		return nil, fmt.Errorf("long header is not supported")
// 	}
// 	messageLength := len(message)
// 	// if messageLength >= 0xfd {
// 	// 	return nil, fmt.Errorf("long message is not supported")
// 	// }
// 	bitcoinMsg := string([]byte{byte(headerLength)})
// 	bitcoinMsg += header
// 	bitcoinMsg += string([]byte{byte(messageLength)})
// 	bitcoinMsg += message
// 	return chainhash.DoubleHashB([]byte(bitcoinMsg)), nil
// }
