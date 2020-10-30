package bitcoin

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvd/chaincfg/chainhash"
	"github.com/piotrnar/gocoin/lib/secp256k1"
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
func VerifyMessage(address, signature, data string) error {
	if len(address) == 0 {
		return errors.New("address is missing")
	} else if len(signature) == 0 {
		return errors.New("signature is missing")
	}
	addresses, err := sigMessageToAddress(signature, data)
	if err != nil {
		return err
	}
	for _, testAddress := range addresses {
		if address == testAddress {
			return nil
		}
	}
	return fmt.Errorf("address: %s not found in %s", address, addresses)
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
func messageHash(message, header string) ([]byte, error) {
	headerLength := len(header)
	if headerLength >= 0xfd {
		return nil, fmt.Errorf("long header is not supported")
	}
	messageLength := len(message)
	if messageLength >= 0xfd {
		return nil, fmt.Errorf("long message is not supported")
	}
	bitcoinMsg := string([]byte{byte(headerLength)})
	bitcoinMsg += header
	bitcoinMsg += string([]byte{byte(messageLength)})
	bitcoinMsg += message
	return chainhash.DoubleHashB([]byte(bitcoinMsg)), nil
}

// parseSignature will parse the given signature
func parseSignature(signature string) (sig secp256k1.Signature, recID int, err error) {
	// todo: is this 64 or 65? is it always the same?
	// panic was occurring: sig.R.SetBytes(sigRaw[1 : 1+32])
	if len(signature) < 64 {
		err = fmt.Errorf("signature is less than %d characters", 64)
		return
	}
	var sigRaw []byte
	if sigRaw, err = base64.StdEncoding.DecodeString(signature); err != nil {
		return
	}
	r0 := sigRaw[0] - 27
	recID = int(r0 & 3)
	if (r0 & 4) == 1 {
		err = fmt.Errorf("compressed type is not supported")
		return
	}
	sig.R.SetBytes(sigRaw[1 : 1+32])
	sig.S.SetBytes(sigRaw[1+32 : 1+32+32])
	return
}

// pubKeyToAddress will convert a pubkey to an address
func pubKeyToAddress(pubkeyXy2 secp256k1.XY, compressed bool, magic []byte) (address string) {
	pubkey, _ := bsvec.ParsePubKey(pubkeyXy2.Bytes(compressed), bsvec.S256())
	bsvecAddress, _ := GetAddressFromPubKey(pubkey)
	return bsvecAddress.String()
}

// sigMessageToAddress will convert a signature & message to a list of addresses
func sigMessageToAddress(signature, message string) ([]string, error) {

	// Get message hash
	msgHash, err := messageHash(message, hBSV)
	if err != nil {
		return nil, err
	}

	// Parse the signature
	var sig secp256k1.Signature
	var recID int
	sig, recID, err = parseSignature(signature)
	if err != nil {
		return nil, err
	}

	var msg secp256k1.Number
	msg.SetBytes(msgHash)

	var pubkeyXy2 secp256k1.XY
	var ret bool
	ret = secp256k1.RecoverPublicKey(sig.R.Bytes(), sig.S.Bytes(), msgHash, recID, &pubkeyXy2)
	if !ret {
		return nil, fmt.Errorf("recover pubkey failed")
	}

	addresses := make([]string, 2)
	for i, compressed := range []bool{true, false} {
		addressString := pubKeyToAddress(pubkeyXy2, compressed, []byte{byte(0)})
		addresses[i] = addressString
	}
	return addresses, nil
}
