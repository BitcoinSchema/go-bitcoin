package bitcoin

import (
	"testing"

	"github.com/bsv-blockchain/go-bt/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDefaultStandardFee verifies the default standard fee returned by DefaultStandardFee.
func TestDefaultStandardFee(t *testing.T) {
	t.Parallel()

	fee := DefaultStandardFee()
	require.NotNil(t, fee)
	assert.Equal(t, bt.FeeTypeStandard, fee.FeeType)
	assert.Equal(t, 5, fee.MiningFee.Satoshis)
	assert.Equal(t, 10, fee.MiningFee.Bytes)
	assert.Equal(t, 5, fee.RelayFee.Satoshis)
	assert.Equal(t, 10, fee.RelayFee.Bytes)
}

// TestGenerateSharedKeyPair verifies that the shared key pair is symmetric: the
// pair derived from (user1 private, user2 public) matches the pair derived from
// (user2 private, user1 public).
func TestGenerateSharedKeyPair(t *testing.T) {
	t.Parallel()

	user1, err := CreatePrivateKey()
	require.NoError(t, err)

	user2, err := CreatePrivateKey()
	require.NoError(t, err)

	priv1, pub1 := GenerateSharedKeyPair(user1, user2.PubKey())
	require.NotNil(t, priv1)
	require.NotNil(t, pub1)

	priv2, pub2 := GenerateSharedKeyPair(user2, user1.PubKey())
	require.NotNil(t, priv2)
	require.NotNil(t, pub2)

	// Both users derive the identical shared secret (ECDH symmetry)
	assert.Equal(t, priv1.Serialize(), priv2.Serialize())
	assert.Equal(t, pub1.Compressed(), pub2.Compressed())
}

// TestA25Accessors exercises the A25 address type accessors on a known-valid
// mainnet address (version 0, self-consistent checksum).
func TestA25Accessors(t *testing.T) {
	t.Parallel()

	var a A25
	require.NoError(t, a.Set58([]byte("1KCEAmVS6FFggtc7W9as7sEENvjt7DqMi2")))

	// Mainnet P2PKH addresses are version 0
	assert.Equal(t, byte(0), a.Version())

	// A valid address has a self-consistent checksum
	assert.Equal(t, a.EmbeddedChecksum(), a.ComputeChecksum())
}
