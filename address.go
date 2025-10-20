package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/crypto"
	"github.com/libsv/go-bt/v2/bscript"
)

var (
	// ErrPublicKeyNil is returned when the public key is nil
	ErrPublicKeyNil = errors.New("public key cannot be nil")
	// ErrPublicKeyXNil is returned when the public key X coordinate is nil
	ErrPublicKeyXNil = errors.New("public key X coordinate cannot be nil")
	// ErrInvalidOutputScript is returned when the output script is missing an address
	ErrInvalidOutputScript = errors.New("invalid output script, missing an address")
)

// A25 is a type for a 25 byte (not base58 encoded) bitcoin address.
type A25 [25]byte

// Version returns the version byte of an A25 address
func (a *A25) Version() byte {
	return a[0]
}

// EmbeddedChecksum returns the 4 checksum bytes of an A25 address
func (a *A25) EmbeddedChecksum() [4]byte {
	var c [4]byte
	copy(c[:], a[21:])
	return c
}

// Tmpl and Set58 are adapted from the C solution.
// Go has big integers but this technique seems better.
//
//nolint:gochecknoglobals // base58 alphabet constant
var tmpl = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Set58 takes a base58 encoded address and decodes it into the receiver.
// Errors are returned if the argument is not valid base58 or if the decoded
// value does not fit in the 25 byte address.  The address is not otherwise
// checked for validity.
func (a *A25) Set58(s []byte) error {
	for _, s1 := range s {
		c := bytes.IndexByte(tmpl, s1)
		if c < 0 {
			return ErrBadCharacter
		}
		for j := 24; j >= 0; j-- {
			c += 58 * int(a[j])
			a[j] = byte(c % 256)
			c /= 256
		}
		if c > 0 {
			return ErrTooLong
		}
	}
	return nil
}

// ComputeChecksum returns a four byte checksum computed from the first 21
// bytes of the address.  The embedded checksum is not updated.
func (a *A25) ComputeChecksum() [4]byte {
	var c [4]byte
	copy(c[:], a.doubleSHA256())
	return c
}

// doubleSHA256 computes a double sha256 hash of the first 21 bytes of the
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
		return false, ErrNotVersion0
	}
	return a.EmbeddedChecksum() == a.ComputeChecksum(), nil
}

// GetAddressFromPrivateKey takes a bec private key and returns a Bitcoin address
func GetAddressFromPrivateKey(privateKey *bec.PrivateKey, compressed, mainnet bool) (string, error) {
	address, err := GetAddressFromPubKey(privateKey.PubKey(), compressed, mainnet)
	if err != nil {
		return "", err
	}
	return address.AddressString, nil
}

// GetAddressFromPrivateKeyString takes a private key string and returns a Bitcoin address
func GetAddressFromPrivateKeyString(privateKey string, compressed, mainnet bool) (string, error) {
	rawKey, err := PrivateKeyFromString(privateKey)
	if err != nil {
		return "", err
	}
	var address *bscript.Address
	if address, err = GetAddressFromPubKey(rawKey.PubKey(), compressed, mainnet); err != nil {
		return "", err
	}
	return address.AddressString, nil
}

// GetAddressFromPubKey gets a bscript.Address from a bec.PublicKey
func GetAddressFromPubKey(publicKey *bec.PublicKey, compressed, mainnet bool) (*bscript.Address, error) {
	if publicKey == nil {
		return nil, ErrPublicKeyNil
	} else if publicKey.X == nil {
		return nil, ErrPublicKeyXNil
	}

	if !compressed {
		// go-bt/v2/bscript does not have a function that exports the uncompressed address
		// https://github.com/libsv/go-bt/blob/master/bscript/address.go#L98
		hash := crypto.Hash160(publicKey.SerialiseUncompressed())
		bb := make([]byte, 1)
		//nolint: makezero // we need to set up the array with 1
		bb = append(bb, hash...)
		return &bscript.Address{
			AddressString: bscript.Base58EncodeMissingChecksum(bb),
			PublicKeyHash: hex.EncodeToString(hash),
		}, nil
	}

	return bscript.NewAddressFromPublicKey(publicKey, mainnet)
}

// GetAddressFromPubKeyString is a convenience function to use a hex string pubKey
func GetAddressFromPubKeyString(pubKey string, compressed, mainnet bool) (*bscript.Address, error) {
	rawPubKey, err := PubKeyFromString(pubKey)
	if err != nil {
		return nil, err
	}
	return GetAddressFromPubKey(rawPubKey, compressed, mainnet)
}

// GetAddressFromScript will take an output script and extract a standard bitcoin address
func GetAddressFromScript(script string) (string, error) {
	// No script?
	if len(script) == 0 {
		return "", ErrMissingScript
	}

	// Decode the hex string into bytes
	scriptBytes, err := hex.DecodeString(script)
	if err != nil {
		return "", err
	}

	// Extract the addresses from the script
	bScript := bscript.NewFromBytes(scriptBytes)
	var addresses []string
	addresses, err = bScript.Addresses()
	if err != nil {
		return "", err
	}

	// Missing an address?
	if len(addresses) == 0 {
		// This error case should not occur since the error above will occur when no address is found,
		// however we ensure that we have an address for the NewLegacyAddressPubKeyHash() below
		return "", ErrInvalidOutputScript
	}

	// Use the encoded version of the address
	return addresses[0], nil
}
