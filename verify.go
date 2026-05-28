package bitcoin

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/bsv-blockchain/go-bt/v2/bscript"
	bsm "github.com/bsv-blockchain/go-sdk/compat/bsm"
	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
)

// ErrAddressNotFound is returned when the signature address does not match the expected address
var ErrAddressNotFound = errors.New("address not found")

// PubKeyFromSignature gets a publickey for a signature and tells you whether is was compressed
func PubKeyFromSignature(sig, data string) (pubKey *ec.PublicKey, wasCompressed bool, err error) {
	var decodedSig []byte
	if decodedSig, err = base64.StdEncoding.DecodeString(sig); err != nil {
		return nil, false, err
	}

	// Validate the signature and recover the public key (Bitcoin Signed Message encoding)
	return bsm.PubKeyFromSignature(decodedSig, []byte(data))
}

// VerifyMessage verifies a string and address against the provided
// signature and assumes Bitcoin Signed Message encoding.
// The key referenced by the signature must relate to the address provided.
// Do not provide an address from an uncompressed key along with
// a signature from a compressed key
//
// Error will occur if verify fails or verification is not successful (no bool)
// Spec: https://docs.moneybutton.com/docs/bsv-message.html
func VerifyMessage(address, sig, data string, mainnet bool) error {
	// Reconstruct the pubkey
	publicKey, wasCompressed, err := PubKeyFromSignature(sig, data)
	if err != nil {
		return err
	}

	// Get the address
	var bscriptAddress *bscript.Address
	if bscriptAddress, err = GetAddressFromPubKey(publicKey, wasCompressed, mainnet); err != nil {
		return err
	}

	// Return nil if addresses match.
	if bscriptAddress.AddressString == address {
		return nil
	}
	return fmt.Errorf(
		"%w: expected %s (compressed: %t), found %s",
		ErrAddressNotFound,
		address,
		wasCompressed,
		bscriptAddress.AddressString,
	)
}

// VerifyMessageDER will take a message string, a public key string and a signature string
// (in strict DER format) and verify that the message was signed by the public key.
//
// Copyright (c) 2019 Bitcoin Association
// License: https://github.com/bitcoin-sv/merchantapi-reference/blob/master/LICENSE
//
// Source: https://github.com/bitcoin-sv/merchantapi-reference/blob/master/handler/global.go
func VerifyMessageDER(hash [32]byte, pubKey, signature string) (verified bool, err error) {
	// Decode the signature string
	var sigBytes []byte
	if sigBytes, err = hex.DecodeString(signature); err != nil {
		return false, err
	}

	// Parse the signature
	var sig *ec.Signature
	if sig, err = ec.ParseDERSignature(sigBytes); err != nil {
		return false, err
	}

	// Decode the pubKey
	var pubKeyBytes []byte
	if pubKeyBytes, err = hex.DecodeString(pubKey); err != nil {
		return false, err
	}

	// Parse the pubKey
	var rawPubKey *ec.PublicKey
	if rawPubKey, err = ec.ParsePubKey(pubKeyBytes); err != nil {
		return false, err
	}

	// Verify the signature against the pubKey
	verified = sig.Verify(hash[:], rawPubKey)
	return verified, nil
}
