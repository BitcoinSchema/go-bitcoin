package bitcoin

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"io"

	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
)

// This file preserves the Pyelliptic / ANSI-X9.63 ECIES wire format that was
// previously provided by github.com/libsv/go-bk/bec, rebuilt on top of
// github.com/bsv-blockchain/go-sdk's ec package so that data encrypted by older
// versions of this library still decrypts byte-for-byte.

var (
	// ErrInvalidMAC occurs when the Message Authentication Check (MAC) fails
	// during decryption, caused by an invalid private key or corrupt ciphertext.
	ErrInvalidMAC = errors.New("invalid mac hash")

	errInputTooShort    = errors.New("ciphertext too short")
	errUnsupportedCurve = errors.New("unsupported curve")
	errInvalidXLength   = errors.New("invalid X length, must be 32")
	errInvalidYLength   = errors.New("invalid Y length, must be 32")
	errInvalidPadding   = errors.New("invalid PKCS#7 padding")

	// ciphCurveBytes is 0x02CA (= 714 = secp256k1, from OpenSSL).
	ciphCurveBytes = [2]byte{0x02, 0xCA} //nolint:gochecknoglobals // fixed ECIES wire-format marker; byte arrays cannot be const
	// ciphCoordLength is 0x0020 (= 32).
	ciphCoordLength = [2]byte{0x00, 0x20} //nolint:gochecknoglobals // fixed ECIES wire-format marker; byte arrays cannot be const
)

// generateSharedSecret generates a shared secret based on a private key and a
// public key using Diffie-Hellman key exchange (ECDH) (RFC 4753). RFC 5903
// Section 9 states we should only return x.
//
// X is intentionally returned unpadded (big.Int.Bytes()) to match go-bk/Pyelliptic
// byte-for-byte; left-padding to the 32-byte curve size would change the SHA-512
// key material for shared points whose X has leading zeros and break decryption of
// older ciphertexts.
func generateSharedSecret(privKey *ec.PrivateKey, pubKey *ec.PublicKey) []byte {
	x, _ := ec.S256().ScalarMult(pubKey.X, pubKey.Y, privKey.D.Bytes())
	return x.Bytes()
}

// eciesEncrypt encrypts data for the target public key using AES-256-CBC. It
// also generates an ephemeral private key (the pubkey of which is in the
// output). The only supported curve is secp256k1. The structure it encodes
// everything into is:
//
//	struct {
//		IV [16]byte                 // AES-256-CBC initialization vector
//		PublicKey [70]byte          // curve(2) + lenX(2) + X(32) + lenY(2) + Y(32)
//		Data []byte                 // cipher text
//		HMAC [32]byte               // HMAC-SHA-256
//	}
//
// The format is byte-compatible with Pyelliptic (ANSI X9.63 section 5.8.1).
func eciesEncrypt(pubKey *ec.PublicKey, in []byte) ([]byte, error) {
	ephemeral, err := ec.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	ecdhKey := generateSharedSecret(ephemeral, pubKey)
	derivedKey := sha512.Sum512(ecdhKey)
	keyE := derivedKey[:32]
	keyM := derivedKey[32:]

	paddedIn := addPKCSPadding(in)
	// IV + Curve params/X/Y + padded plaintext/ciphertext + HMAC-256
	out := make([]byte, aes.BlockSize+70+len(paddedIn)+sha256.Size)
	iv := out[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	// start writing public key
	pb := ephemeral.PubKey().Uncompressed()
	offset := aes.BlockSize

	// curve and X length
	copy(out[offset:offset+4], append(ciphCurveBytes[:], ciphCoordLength[:]...))
	offset += 4
	// X
	copy(out[offset:offset+32], pb[1:33])
	offset += 32
	// Y length
	copy(out[offset:offset+2], ciphCoordLength[:])
	offset += 2
	// Y
	copy(out[offset:offset+32], pb[33:])
	offset += 32

	// start encryption
	block, err := aes.NewCipher(keyE)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(out[offset:len(out)-sha256.Size], paddedIn)

	// start HMAC-SHA-256
	hm := hmac.New(sha256.New, keyM)
	if _, err = hm.Write(out[:len(out)-sha256.Size]); err != nil { // everything is hashed
		return nil, err
	}
	copy(out[len(out)-sha256.Size:], hm.Sum(nil)) // write checksum

	return out, nil
}

// eciesDecrypt decrypts data that was encrypted using eciesEncrypt.
func eciesDecrypt(priv *ec.PrivateKey, in []byte) ([]byte, error) {
	// IV + Curve params/X/Y + 1 block + HMAC-256
	if len(in) < aes.BlockSize+70+aes.BlockSize+sha256.Size {
		return nil, errInputTooShort
	}

	// read iv
	iv := in[:aes.BlockSize]
	offset := aes.BlockSize

	// start reading pubkey
	if !bytes.Equal(in[offset:offset+2], ciphCurveBytes[:]) {
		return nil, errUnsupportedCurve
	}
	offset += 2

	if !bytes.Equal(in[offset:offset+2], ciphCoordLength[:]) {
		return nil, errInvalidXLength
	}
	offset += 2

	xBytes := in[offset : offset+32]
	offset += 32

	if !bytes.Equal(in[offset:offset+2], ciphCoordLength[:]) {
		return nil, errInvalidYLength
	}
	offset += 2

	yBytes := in[offset : offset+32]
	offset += 32

	pb := make([]byte, 65)
	pb[0] = 0x04 // uncompressed
	copy(pb[1:33], xBytes)
	copy(pb[33:], yBytes)
	// check if (X, Y) lies on the curve and create a Pubkey if it does
	pubKey, err := ec.ParsePubKey(pb)
	if err != nil {
		return nil, err
	}

	// check for cipher text length
	if (len(in)-aes.BlockSize-offset-sha256.Size)%aes.BlockSize != 0 {
		return nil, errInvalidPadding // not padded to 16 bytes
	}

	// read hmac
	messageMAC := in[len(in)-sha256.Size:]

	// generate shared secret
	ecdhKey := generateSharedSecret(priv, pubKey)
	derivedKey := sha512.Sum512(ecdhKey)
	keyE := derivedKey[:32]
	keyM := derivedKey[32:]

	// verify mac
	hm := hmac.New(sha256.New, keyM)
	if _, err = hm.Write(in[:len(in)-sha256.Size]); err != nil { // everything is hashed
		return nil, err
	}
	if !hmac.Equal(messageMAC, hm.Sum(nil)) {
		return nil, ErrInvalidMAC
	}

	// start decryption
	block, err := aes.NewCipher(keyE)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	// same length as ciphertext
	plaintext := make([]byte, len(in)-offset-sha256.Size)
	mode.CryptBlocks(plaintext, in[offset:len(in)-sha256.Size])

	return removePKCSPadding(plaintext)
}

// addPKCSPadding adds PKCS#7 padding to a block of data (AES block size).
func addPKCSPadding(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// removePKCSPadding removes padding from data added with addPKCSPadding.
func removePKCSPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errInvalidPadding
	}
	padLength := int(src[length-1])
	if padLength == 0 || padLength > aes.BlockSize || padLength > length {
		return nil, errInvalidPadding
	}
	// Every padding byte must equal the pad length (full PKCS#7 validation).
	for _, b := range src[length-padLength:] {
		if int(b) != padLength {
			return nil, errInvalidPadding
		}
	}

	return src[:length-padLength], nil
}
