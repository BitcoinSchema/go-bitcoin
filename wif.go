package bitcoin

import (
	"bytes"
	"errors"

	base58 "github.com/bsv-blockchain/go-sdk/compat/base58"
	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	hash "github.com/bsv-blockchain/go-sdk/primitives/hash"
	chaincfg "github.com/bsv-blockchain/go-sdk/transaction/chaincfg"
)

// This file preserves the Wallet Import Format (WIF) type and its uncompressed
// output that was previously provided by github.com/libsv/go-bk/wif, rebuilt on
// top of github.com/bsv-blockchain/go-sdk's ec package. go-sdk's own
// PrivateKey.Wif() only emits compressed WIFs, so the port is required to keep
// the existing public API and wire format intact.

var (
	// ErrChecksumMismatch describes an error where decoding failed due to a bad checksum.
	ErrChecksumMismatch = errors.New("checksum mismatch")

	// ErrMalformedPrivateKey describes an error where a WIF-encoded private key
	// cannot be decoded due to being improperly formatted: an incorrect byte
	// length or an unexpected magic number.
	ErrMalformedPrivateKey = errors.New("malformed private key")

	// ErrNoNetwork is returned when a nil network is passed to NewWIF.
	ErrNoNetwork = errors.New("no network")
)

// privKeyBytesLen is the length in bytes of a serialized private key.
const privKeyBytesLen = 32

// compressMagic is the magic byte used to identify a WIF encoding for an
// address created from a compressed serialized public key.
const compressMagic byte = 0x01

// WIF contains the individual components described by the Wallet Import Format
// (WIF). A WIF string is typically used to represent a private key and its
// associated address in a way that may be easily copied and imported into or
// exported from wallet software. WIF strings may be decoded into this structure
// by calling DecodeWIF or created with a user-provided private key by calling NewWIF.
type WIF struct {
	// PrivKey is the private key being imported or exported.
	PrivKey *ec.PrivateKey

	// CompressPubKey specifies whether the address controlled by the imported or
	// exported private key was created by hashing a compressed (33-byte)
	// serialized public key, rather than an uncompressed (65-byte) one.
	CompressPubKey bool

	// netID is the bitcoin network identifier byte used when WIF encoding the private key.
	netID byte
}

// NewWIF creates a new WIF structure to export an address and its private key as
// a string encoded in the Wallet Import Format. The compress argument specifies
// whether the address intended to be imported or exported was created by
// serializing the public key compressed rather than uncompressed.
func NewWIF(privKey *ec.PrivateKey, net *chaincfg.Params, compress bool) (*WIF, error) {
	if net == nil {
		return nil, ErrNoNetwork
	}
	return &WIF{privKey, compress, net.PrivateKeyID}, nil
}

// IsForNet returns whether the decoded WIF structure is associated with the
// passed bitcoin network.
func (w *WIF) IsForNet(net *chaincfg.Params) bool {
	return w.netID == net.PrivateKeyID
}

// DecodeWIF creates a new WIF structure by decoding the string encoding of the
// import format.
//
// The WIF string must be a base58-encoded string of the following byte sequence:
//   - 1 byte to identify the network, must be 0x80 for mainnet or 0xef for
//     testnet3 or the regression test network
//   - 32 bytes of a binary-encoded, big-endian, zero-padded private key
//   - Optional 1 byte (equal to 0x01) if the address was created from a
//     compressed (33-byte) public key
//   - 4 bytes of checksum, equal to the first four bytes of the double SHA256 of
//     every byte before the checksum in this sequence
//
// ErrMalformedPrivateKey is returned for an impossible length or a bad compress
// magic; ErrChecksumMismatch is returned if the checksum does not match.
func DecodeWIF(wif string) (*WIF, error) {
	decoded, err := base58.Decode(wif)
	if err != nil {
		return nil, ErrMalformedPrivateKey
	}
	decodedLen := len(decoded)
	var compress bool

	// Length of base58 decoded WIF must be 32 bytes + an optional 1 byte (0x01)
	// if compressed, plus 1 byte for netID + 4 bytes of checksum.
	switch decodedLen {
	case 1 + privKeyBytesLen + 1 + 4:
		if decoded[33] != compressMagic {
			return nil, ErrMalformedPrivateKey
		}
		compress = true
	case 1 + privKeyBytesLen + 4:
		compress = false
	default:
		return nil, ErrMalformedPrivateKey
	}

	// Checksum is first four bytes of double SHA256 of the identifier byte and
	// privKey. Verify this matches the final 4 bytes of the decoded private key.
	var tosum []byte
	if compress {
		tosum = decoded[:1+privKeyBytesLen+1]
	} else {
		tosum = decoded[:1+privKeyBytesLen]
	}
	cksum := hash.Sha256d(tosum)[:4]
	if !bytes.Equal(cksum, decoded[decodedLen-4:]) {
		return nil, ErrChecksumMismatch
	}

	netID := decoded[0]
	privKeyBytes := decoded[1 : 1+privKeyBytesLen]
	privKey, _ := ec.PrivateKeyFromBytes(privKeyBytes)
	return &WIF{privKey, compress, netID}, nil
}

// String creates the Wallet Import Format string encoding of a WIF structure.
// See DecodeWIF for a detailed breakdown of the format.
func (w *WIF) String() string {
	// Precalculate size: one byte for the network, 32 bytes of private key,
	// possibly one extra byte if the pubkey is to be compressed, and four bytes
	// of checksum.
	encodeLen := 1 + privKeyBytesLen + 4
	if w.CompressPubKey {
		encodeLen++
	}

	a := make([]byte, 0, encodeLen)
	a = append(a, w.netID)
	// Pad and append bytes manually to avoid another call to make.
	a = paddedAppend(privKeyBytesLen, a, w.PrivKey.D.Bytes())
	if w.CompressPubKey {
		a = append(a, compressMagic)
	}
	cksum := hash.Sha256d(a)[:4]
	a = append(a, cksum...)
	return base58.Encode(a)
}

// SerializePubKey serializes the associated public key of the imported or
// exported private key in either a compressed or uncompressed format depending
// on the value of w.CompressPubKey.
func (w *WIF) SerializePubKey() []byte {
	pk := w.PrivKey.PubKey()
	if w.CompressPubKey {
		return pk.Compressed()
	}
	return pk.Uncompressed()
}

// paddedAppend appends the src byte slice to dst, returning the new slice. If
// the length of the source is smaller than the passed size, leading zero bytes
// are appended to dst before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}
