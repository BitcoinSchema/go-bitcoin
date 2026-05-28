package bitcoin

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// knownPyellipticPrivKey / knownPyellipticCipher are a vector produced by the
// original github.com/libsv/go-bk/bec.Encrypt (Pyelliptic / ANSI-X9.63 format).
// Decrypting it with the ported helper proves byte-for-byte format compatibility.
const (
	knownPyellipticPrivKey = "bb66a48a9f6dd7b8fb469a6f08a75c25770591dc509c72129b2aaeca77a5269e"
	knownPyellipticCipher  = "4fab3e3534e101cef1ec936628894a2c02ca00206b60babbfa6ccac9a8078d55a50a32b5e15b1a3c9a360749" +
		"9066bc4eb70721f4002005cdc7638c7e6c051be80cc83092ed5af5a9a015e24c3ed4af289d59c0c65cb608aa8f07c9318e6e52a18" +
		"e60dbcf7a5889304f4bfc01cad735a5f2279f06eb5da5ed7454320da87becbdcd889708fcab"
	knownPyellipticPlain = "test-data"
)

// TestEciesDecryptKnownVector proves the ported Pyelliptic ECIES decrypts data
// that was encrypted by the original go-bk/bec implementation (wire-compatible).
func TestEciesDecryptKnownVector(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(knownPyellipticPrivKey)
	require.NoError(t, err)

	cipherBytes, err := hex.DecodeString(knownPyellipticCipher)
	require.NoError(t, err)

	plain, err := eciesDecrypt(priv, cipherBytes)
	require.NoError(t, err)
	assert.Equal(t, knownPyellipticPlain, string(plain))
}

// TestEciesEncryptDecryptRoundTrip checks a variety of payloads round-trip.
func TestEciesEncryptDecryptRoundTrip(t *testing.T) {
	t.Parallel()

	priv, err := CreatePrivateKey()
	require.NoError(t, err)

	tests := []struct {
		name string
		data []byte
	}{
		{"empty", []byte("")},
		{"single space", []byte(" ")},
		{"newline", []byte("\n")},
		{"ascii", []byte("hello world")},
		{"json", []byte(`{"some":"data","n":1}`)},
		{"unicode", []byte("日本語テスト — éàü 🚀")}, //nolint:gosmopolitan // intentionally testing multi-byte UTF-8 payloads
		{"binary", []byte{0x00, 0x01, 0x02, 0xff, 0xfe, 0x10, 0x00}},
		{"block aligned 16", bytes.Repeat([]byte("A"), 16)},
		{"large 10k", bytes.Repeat([]byte("z"), 10_000)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encrypted, encErr := eciesEncrypt(priv.PubKey(), test.data)
			require.NoError(t, encErr)
			require.NotEmpty(t, encrypted)

			decrypted, decErr := eciesDecrypt(priv, encrypted)
			require.NoError(t, decErr)
			assert.Equal(t, test.data, decrypted)
		})
	}
}

// TestEciesEncryptIsNonDeterministic confirms the ephemeral key makes each
// ciphertext unique while both still decrypt to the same plaintext.
func TestEciesEncryptIsNonDeterministic(t *testing.T) {
	t.Parallel()

	priv, err := CreatePrivateKey()
	require.NoError(t, err)

	a, err := eciesEncrypt(priv.PubKey(), []byte("same message"))
	require.NoError(t, err)
	b, err := eciesEncrypt(priv.PubKey(), []byte("same message"))
	require.NoError(t, err)

	assert.NotEqual(t, a, b, "ephemeral key should make ciphertexts differ")

	da, err := eciesDecrypt(priv, a)
	require.NoError(t, err)
	db, err := eciesDecrypt(priv, b)
	require.NoError(t, err)
	assert.Equal(t, "same message", string(da))
	assert.Equal(t, "same message", string(db))
}

// TestEciesEncryptFormat verifies the Pyelliptic structure markers.
func TestEciesEncryptFormat(t *testing.T) {
	t.Parallel()

	priv, err := CreatePrivateKey()
	require.NoError(t, err)

	encrypted, err := eciesEncrypt(priv.PubKey(), []byte("x"))
	require.NoError(t, err)

	// Minimum length: 16 (IV) + 70 (pubkey block) + 16 (1 cipher block) + 32 (HMAC)
	require.GreaterOrEqual(t, len(encrypted), 134)
	// Curve marker 0x02CA at offset 16, X length 0x0020 at offset 18.
	assert.Equal(t, []byte{0x02, 0xCA}, encrypted[16:18])
	assert.Equal(t, []byte{0x00, 0x20}, encrypted[18:20])
	// Y length 0x0020 at offset 52.
	assert.Equal(t, []byte{0x00, 0x20}, encrypted[52:54])
}

// TestGenerateSharedSecretSymmetry verifies ECDH symmetry.
func TestGenerateSharedSecretSymmetry(t *testing.T) {
	t.Parallel()

	a, err := CreatePrivateKey()
	require.NoError(t, err)
	b, err := CreatePrivateKey()
	require.NoError(t, err)

	secretAB := generateSharedSecret(a, b.PubKey())
	secretBA := generateSharedSecret(b, a.PubKey())

	require.NotEmpty(t, secretAB)
	assert.Equal(t, secretAB, secretBA)
}

// TestEciesDecryptFailures covers the error paths of eciesDecrypt.
func TestEciesDecryptFailures(t *testing.T) {
	t.Parallel()

	priv, err := CreatePrivateKey()
	require.NoError(t, err)

	valid, err := eciesEncrypt(priv.PubKey(), []byte("the quick brown fox"))
	require.NoError(t, err)

	t.Run("too short", func(t *testing.T) {
		_, decErr := eciesDecrypt(priv, make([]byte, 50))
		require.ErrorIs(t, decErr, errInputTooShort)
	})

	t.Run("wrong private key", func(t *testing.T) {
		other, keyErr := CreatePrivateKey()
		require.NoError(t, keyErr)
		_, decErr := eciesDecrypt(other, valid)
		require.ErrorIs(t, decErr, ErrInvalidMAC)
	})

	t.Run("corrupted hmac", func(t *testing.T) {
		corrupted := append([]byte(nil), valid...)
		corrupted[len(corrupted)-1] ^= 0xff
		_, decErr := eciesDecrypt(priv, corrupted)
		require.ErrorIs(t, decErr, ErrInvalidMAC)
	})

	t.Run("bad curve marker", func(t *testing.T) {
		corrupted := append([]byte(nil), valid...)
		corrupted[16] ^= 0xff
		_, decErr := eciesDecrypt(priv, corrupted)
		require.ErrorIs(t, decErr, errUnsupportedCurve)
	})

	t.Run("bad X length", func(t *testing.T) {
		corrupted := append([]byte(nil), valid...)
		corrupted[18] ^= 0xff
		_, decErr := eciesDecrypt(priv, corrupted)
		require.ErrorIs(t, decErr, errInvalidXLength)
	})

	t.Run("bad Y length", func(t *testing.T) {
		corrupted := append([]byte(nil), valid...)
		corrupted[52] ^= 0xff
		_, decErr := eciesDecrypt(priv, corrupted)
		require.ErrorIs(t, decErr, errInvalidYLength)
	})
}

// TestEciesMatchesPublicWrappers confirms the unexported helpers and the public
// EncryptWithPrivateKey / DecryptWithPrivateKey wrappers interoperate.
func TestEciesMatchesPublicWrappers(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(testPrivateKeyHex)
	require.NoError(t, err)

	encrypted, err := eciesEncrypt(priv.PubKey(), []byte("interop"))
	require.NoError(t, err)

	decrypted, err := DecryptWithPrivateKey(priv, hex.EncodeToString(encrypted))
	require.NoError(t, err)
	assert.Equal(t, "interop", decrypted)
}

// TestPKCSPadding exercises the padding helpers directly.
func TestPKCSPadding(t *testing.T) {
	t.Parallel()

	for _, n := range []int{0, 1, 15, 16, 17, 31, 32} {
		in := bytes.Repeat([]byte{0xab}, n)
		padded := addPKCSPadding(in)
		require.Zero(t, len(padded)%16, "padded length must be a multiple of 16")

		unpadded, err := removePKCSPadding(padded)
		require.NoError(t, err)
		assert.Equal(t, in, unpadded)
	}

	t.Run("invalid padding", func(t *testing.T) {
		// A full block whose final byte claims a padding length larger than the block.
		bad := bytes.Repeat([]byte{0x20}, 16)
		_, err := removePKCSPadding(bad)
		require.ErrorIs(t, err, errInvalidPadding)
	})
}

// helper to assert a base64/hex-free quick sanity on shared secret length
func TestGenerateSharedSecretLength(t *testing.T) {
	t.Parallel()

	a, err := CreatePrivateKey()
	require.NoError(t, err)
	b, err := CreatePrivateKey()
	require.NoError(t, err)

	secret := generateSharedSecret(a, b.PubKey())
	// secp256k1 x-coordinate is at most 32 bytes.
	require.LessOrEqual(t, len(secret), 32)
	require.NotEmpty(t, secret)
}
