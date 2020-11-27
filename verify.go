package bitcoin

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvd/chaincfg/chainhash"
	"github.com/bitcoinsv/bsvd/wire"
	"github.com/bitcoinsv/bsvutil"
)

const (
	// hBSV is the magic header string required fore Bitcoin Signed Messages
	hBSV string = "Bitcoin Signed Message:\n"
)

// PubKeyFromSignature gets a publickey for a signature and tells you whether is was compressed
func PubKeyFromSignature(sig, data string) (pubKey *bsvec.PublicKey, wasCompressed bool, err error) {

	var decodedSig []byte
	if decodedSig, err = base64.StdEncoding.DecodeString(sig); err != nil {
		return nil, false, err
	}

	// Validate the signature - this just shows that it was valid at all
	// we will compare it with the key next
	var buf bytes.Buffer
	if err = wire.WriteVarString(&buf, 0, hBSV); err != nil {
		return nil, false, err
	}
	if err = wire.WriteVarString(&buf, 0, data); err != nil {
		return nil, false, err
	}

	// Create the hash
	expectedMessageHash := chainhash.DoubleHashB(buf.Bytes())
	return bsvec.RecoverCompact(bsvec.S256(), decodedSig, expectedMessageHash)
}

// VerifyMessage verifies a string and address against the provided
// signature and assumes Bitcoin Signed Message encoding.
// The key referenced by the signature must relate to the address provided.
// Do not provide an address from an uncompressed key along with
// a signature from a compressed key
//
// Error will occur if verify fails or verification is not successful (no bool)
// Spec: https://docs.moneybutton.com/docs/bsv-message.html
func VerifyMessage(address, sig, data string) error {

	// Reconstruct the pubkey
	publicKey, wasCompressed, err := PubKeyFromSignature(sig, data)
	if err != nil {
		return err
	}

	// Get the address
	var bsvecAddress *bsvutil.LegacyAddressPubKeyHash
	if bsvecAddress, err = GetAddressFromPubKey(publicKey, wasCompressed); err != nil {
		return err
	}

	// Return nil if addresses match.
	if bsvecAddress.String() == address {
		return nil
	}
	return fmt.Errorf(
		"address (%s) not found - compressed: %t\n%s was found instead",
		address,
		wasCompressed,
		bsvecAddress.EncodeAddress(),
	)
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
