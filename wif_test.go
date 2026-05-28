package bitcoin

import (
	"encoding/hex"
	"strings"
	"testing"

	base58 "github.com/bsv-blockchain/go-sdk/compat/base58"
	chaincfg "github.com/bsv-blockchain/go-sdk/transaction/chaincfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testPrivateKeyHex (defined in test_constants_test.go) encodes to these
// mainnet WIFs (uncompressed leads with '5', compressed with 'K'/'L').
const (
	testUncompressedWIF = "5JTHas7yTFMBLqgFogxZFf8Vc5uKEbkE7yQAQ2g3xPHo2sNG1Ei"
	testCompressedWIF   = "Kz32CUDnArL4eZrGM5NDMhJ5FrduV2MnumwEUePN3TP8AwSRRFvQ"
)

// TestNewWIFUncompressed verifies uncompressed WIF output is preserved (51 chars,
// leading '5') and round-trips back to the original private key.
func TestNewWIFUncompressed(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(testPrivateKeyHex)
	require.NoError(t, err)

	w, err := NewWIF(priv, &chaincfg.MainNet, false)
	require.NoError(t, err)
	require.NotNil(t, w)
	assert.False(t, w.CompressPubKey)

	s := w.String()
	require.Len(t, s, 51)
	assert.True(t, strings.HasPrefix(s, "5"), "uncompressed mainnet WIF starts with 5")
	assert.Equal(t, testUncompressedWIF, s)

	decoded, err := DecodeWIF(s)
	require.NoError(t, err)
	assert.False(t, decoded.CompressPubKey)
	assert.Equal(t, testPrivateKeyHex, hex.EncodeToString(decoded.PrivKey.Serialize()))
}

// TestNewWIFCompressed verifies compressed WIF output (52 chars) and round-trip.
func TestNewWIFCompressed(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(testPrivateKeyHex)
	require.NoError(t, err)

	w, err := NewWIF(priv, &chaincfg.MainNet, true)
	require.NoError(t, err)
	require.True(t, w.CompressPubKey)

	s := w.String()
	require.Len(t, s, 52)
	assert.Equal(t, testCompressedWIF, s, "compressed mainnet WIF of the test key")

	decoded, err := DecodeWIF(testCompressedWIF)
	require.NoError(t, err)
	assert.True(t, decoded.CompressPubKey)
	assert.Equal(t, testPrivateKeyHex, hex.EncodeToString(decoded.PrivKey.Serialize()))
}

// TestDecodeWIFKnownVectors checks decoding well-known WIF strings.
func TestDecodeWIFKnownVectors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		wif            string
		expectedKeyHex string
		expectCompress bool
	}{
		{
			name:           "uncompressed all zeros",
			wif:            "5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAbuatmU",
			expectedKeyHex: "0000000000000000000000000000000000000000000000000000000000000000",
			expectCompress: false,
		},
		{
			name:           "uncompressed test key",
			wif:            testUncompressedWIF,
			expectedKeyHex: testPrivateKeyHex,
			expectCompress: false,
		},
		{
			name:           "compressed test key",
			wif:            testCompressedWIF,
			expectedKeyHex: testPrivateKeyHex,
			expectCompress: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decoded, err := DecodeWIF(test.wif)
			require.NoError(t, err)
			assert.Equal(t, test.expectCompress, decoded.CompressPubKey)
			assert.Equal(t, test.expectedKeyHex, hex.EncodeToString(decoded.PrivKey.Serialize()))
			// Re-encoding must reproduce the original string.
			assert.Equal(t, test.wif, decoded.String())
		})
	}
}

// TestWIFIsForNet checks the network identifier comparison.
func TestWIFIsForNet(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(testPrivateKeyHex)
	require.NoError(t, err)

	w, err := NewWIF(priv, &chaincfg.MainNet, false)
	require.NoError(t, err)

	assert.True(t, w.IsForNet(&chaincfg.MainNet))
	assert.False(t, w.IsForNet(&chaincfg.TestNet))
}

// TestNewWIFNilNet ensures a nil network returns an error.
func TestNewWIFNilNet(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(testPrivateKeyHex)
	require.NoError(t, err)

	_, err = NewWIF(priv, nil, false)
	require.Error(t, err)
}

// TestDecodeWIFErrors covers malformed and bad-checksum inputs.
func TestDecodeWIFErrors(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		_, err := DecodeWIF("")
		require.ErrorIs(t, err, ErrMalformedPrivateKey)
	})

	t.Run("invalid base58 characters", func(t *testing.T) {
		// Contains 'l' which is not in the base58 alphabet.
		_, err := DecodeWIF("invalid")
		require.ErrorIs(t, err, ErrMalformedPrivateKey)
	})

	t.Run("wrong length", func(t *testing.T) {
		_, err := DecodeWIF(base58.Encode([]byte{0x80, 0x01, 0x02, 0x03}))
		require.ErrorIs(t, err, ErrMalformedPrivateKey)
	})

	t.Run("checksum mismatch", func(t *testing.T) {
		pkBytes, decErr := hex.DecodeString(testPrivateKeyHex)
		require.NoError(t, decErr)

		// netID + 32-byte key + intentionally wrong 4-byte checksum.
		raw := make([]byte, 0, 1+privKeyBytesLen+4)
		raw = append(raw, chaincfg.MainNet.PrivateKeyID)
		raw = append(raw, pkBytes...)
		raw = append(raw, 0xde, 0xad, 0xbe, 0xef)

		_, err := DecodeWIF(base58.Encode(raw))
		require.ErrorIs(t, err, ErrChecksumMismatch)
	})

	t.Run("bad compress magic", func(t *testing.T) {
		pkBytes, decErr := hex.DecodeString(testPrivateKeyHex)
		require.NoError(t, decErr)

		// netID + 32-byte key + bad compress flag (not 0x01) + checksum-length filler.
		raw := make([]byte, 0, 1+privKeyBytesLen+1+4)
		raw = append(raw, chaincfg.MainNet.PrivateKeyID)
		raw = append(raw, pkBytes...)
		raw = append(raw, 0x02) // invalid compress magic
		raw = append(raw, 0x00, 0x00, 0x00, 0x00)

		_, err := DecodeWIF(base58.Encode(raw))
		require.ErrorIs(t, err, ErrMalformedPrivateKey)
	})
}

// TestWIFSerializePubKey checks compressed vs uncompressed public key sizes.
func TestWIFSerializePubKey(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(testPrivateKeyHex)
	require.NoError(t, err)

	compressed, err := NewWIF(priv, &chaincfg.MainNet, true)
	require.NoError(t, err)
	assert.Len(t, compressed.SerializePubKey(), 33)

	uncompressed, err := NewWIF(priv, &chaincfg.MainNet, false)
	require.NoError(t, err)
	assert.Len(t, uncompressed.SerializePubKey(), 65)
}
