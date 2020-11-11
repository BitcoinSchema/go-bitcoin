package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvd/chaincfg"
	"github.com/bitcoinsv/bsvd/txscript"
	"github.com/bitcoinsv/bsvutil"
)

// A25 is a type for a 25 byte (not base58 encoded) bitcoin address.
type A25 [25]byte

// DoubleSHA256 computes a double sha256 hash of the first 21 bytes of the
// address.  This is the one function shared with the other bitcoin RC task.
// Returned is the full 32 byte sha256 hash.  (The bitcoin checksum will be
// the first four bytes of the slice.)
func (a *A25) doubleSHA256() []byte {
	h := sha256.New()
	_, _ = h.Write(a[:21])
	d := h.Sum([]byte{})
	h = sha256.New()
	_, _ = h.Write(d)
	return h.Sum(d[:0])
}

// Version returns the version byte of a A25 address
func (a *A25) Version() byte {
	return a[0]
}

// EmbeddedChecksum returns the 4 checksum bytes of a A25 address
func (a *A25) EmbeddedChecksum() (c [4]byte) {
	copy(c[:], a[21:])
	return
}

// Tmpl and Set58 are adapted from the C solution.
// Go has big integers but this technique seems better.
var tmpl = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Set58 takes a base58 encoded address and decodes it into the receiver.
// Errors are returned if the argument is not valid base58 or if the decoded
// value does not fit in the 25 byte address.  The address is not otherwise
// checked for validity.
func (a *A25) Set58(s []byte) error {
	for _, s1 := range s {
		c := bytes.IndexByte(tmpl, s1)
		if c < 0 {
			return errors.New("bad char")
		}
		for j := 24; j >= 0; j-- {
			c += 58 * int(a[j])
			a[j] = byte(c % 256)
			c /= 256
		}
		if c > 0 {
			return errors.New("too long")
		}
	}
	return nil
}

// ComputeChecksum returns a four byte checksum computed from the first 21
// bytes of the address.  The embedded checksum is not updated.
func (a *A25) ComputeChecksum() (c [4]byte) {
	copy(c[:], a.doubleSHA256())
	return
}

// ValidA58 validates a base58 encoded bitcoin address.  An address is valid
// if it can be decoded into a 25 byte address, the version number is 0,
// and the checksum validates.  Return value ok will be true for valid
// addresses.  If ok is false, the address is invalid and the error value
// may indicate why.
func ValidA58(a58 []byte) (bool, error) {
	var a A25
	if err := a.Set58(a58); err != nil {
		return false, err
	}
	if a.Version() != 0 {
		return false, errors.New("not version 0")
	}
	return a.EmbeddedChecksum() == a.ComputeChecksum(), nil
}

// GetAddressFromPrivateKey takes a bsvec private key and returns a Bitcoin address
func GetAddressFromPrivateKey(privateKey *bsvec.PrivateKey, compressed bool) (string, error) {
	address, err := GetAddressFromPubKey(privateKey.PubKey(), compressed)
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}

// GetAddressFromPrivateKeyString takes a private key string and returns a Bitcoin address
func GetAddressFromPrivateKeyString(privateKey string, compressed bool) (string, error) {
	rawKey, err := PrivateKeyFromString(privateKey)
	if err != nil {
		return "", err
	}
	var address *bsvutil.LegacyAddressPubKeyHash
	if address, err = GetAddressFromPubKey(rawKey.PubKey(), compressed); err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}

// GetAddressFromPubKey gets a bsvutil.LegacyAddressPubKeyHash from a bsvec.PublicKey
func GetAddressFromPubKey(publicKey *bsvec.PublicKey, compressed bool) (*bsvutil.LegacyAddressPubKeyHash, error) {
	if publicKey == nil {
		return nil, fmt.Errorf("publicKey cannot be nil")
	} else if publicKey.X == nil {
		return nil, fmt.Errorf("publicKey.X cannot be nil")
	}
	var serializedPublicKey []byte
	if compressed {
		serializedPublicKey = publicKey.SerializeCompressed()
	} else {
		serializedPublicKey = publicKey.SerializeUncompressed()
	}

	return bsvutil.NewLegacyAddressPubKeyHash(bsvutil.Hash160(serializedPublicKey), &chaincfg.MainNetParams)
}

// GetAddressFromPubKeyString is a convenience function to use a hex string pubKey
func GetAddressFromPubKeyString(pubKey string, compressed bool) (*bsvutil.LegacyAddressPubKeyHash, error) {
	rawPubKey, err := PubKeyFromString(pubKey)
	if err != nil {
		return nil, err
	}
	return GetAddressFromPubKey(rawPubKey, compressed)
}

// GetAddressFromScript will take an output script and extract a standard bitcoin address
func GetAddressFromScript(script string) (string, error) {

	// No script?
	if len(script) == 0 {
		return "", errors.New("missing script")
	}

	// Decode the hex string into bytes
	scriptBytes, err := hex.DecodeString(script)
	if err != nil {
		return "", err
	}

	// Extract the components from the script
	var addresses []bsvutil.Address
	_, addresses, _, err = txscript.ExtractPkScriptAddrs(scriptBytes, &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}

	// Missing an address?
	if len(addresses) == 0 {
		// This error case should not occur since the error above will occur when no address is found,
		// however we ensure that we have an address for the NewLegacyAddressPubKeyHash() below
		return "", fmt.Errorf("invalid output script, missing an address")
	}

	// Extract the address from the pubkey hash
	var address *bsvutil.LegacyAddressPubKeyHash
	if address, err = bsvutil.NewLegacyAddressPubKeyHash(
		addresses[0].ScriptAddress(),
		&chaincfg.MainNetParams,
	); err != nil {
		return "", err
	}

	// Use the encoded version of the address
	return address.EncodeAddress(), nil
}
