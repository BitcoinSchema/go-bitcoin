package bitcoin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSignMessageExtended provides comprehensive tests for SignMessage
func TestSignMessageExtended(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		privateKey           string
		message              string
		sigRefCompressedKey  bool
		shouldError          bool
		expectedErrorContent string
	}{
		{
			name:                "sign message with compressed key reference",
			privateKey:          "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd",
			message:             "Test message",
			sigRefCompressedKey: true,
			shouldError:         false,
		},
		{
			name:                "sign message with uncompressed key reference",
			privateKey:          "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd",
			message:             "Test message",
			sigRefCompressedKey: false,
			shouldError:         false,
		},
		{
			name:                "sign empty message",
			privateKey:          "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd",
			message:             "",
			sigRefCompressedKey: true,
			shouldError:         false,
		},
		{
			name:                "sign long message",
			privateKey:          "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd",
			message:             string(make([]byte, 1000)),
			sigRefCompressedKey: true,
			shouldError:         false,
		},
		{
			name:                 "empty private key should error",
			privateKey:           "",
			message:              "Test message",
			sigRefCompressedKey:  true,
			shouldError:          true,
			expectedErrorContent: "missing",
		},
		{
			name:                "invalid private key should error",
			privateKey:          "invalid-hex",
			message:             "Test message",
			sigRefCompressedKey: true,
			shouldError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signature, err := SignMessage(tt.privateKey, tt.message, tt.sigRefCompressedKey)

			if tt.shouldError {
				require.Error(t, err, "Expected error but got none")
				if tt.expectedErrorContent != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorContent)
				}
				assert.Empty(t, signature, "Signature should be empty on error")
			} else {
				require.NoError(t, err, "Unexpected error: %v", err)
				assert.NotEmpty(t, signature, "Signature should not be empty")
			}
		})
	}
}

// TestPubKeyFromSignatureExtended tests PubKeyFromSignature with more cases
func TestPubKeyFromSignatureExtended(t *testing.T) {
	t.Parallel()

	// Create a valid signature first
	privateKey := "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"
	message := "Test message"
	signature, err := SignMessage(privateKey, message, true)
	require.NoError(t, err)

	t.Run("valid signature", func(t *testing.T) {
		pubKey, wasCompressed, err := PubKeyFromSignature(signature, message)
		require.NoError(t, err)
		assert.NotNil(t, pubKey)
		assert.True(t, wasCompressed)
	})

	t.Run("invalid base64 signature", func(t *testing.T) {
		pubKey, _, err := PubKeyFromSignature("not-base64!@#", message)
		require.Error(t, err)
		assert.Nil(t, pubKey)
	})

	t.Run("empty signature", func(t *testing.T) {
		pubKey, _, err := PubKeyFromSignature("", message)
		require.Error(t, err)
		assert.Nil(t, pubKey)
	})

	t.Run("wrong message", func(t *testing.T) {
		// Using a different message should still recover a pubkey, but verification would fail
		pubKey, wasCompressed, err := PubKeyFromSignature(signature, "Different message")
		// This might succeed or fail depending on implementation
		if err != nil {
			t.Logf("Error as expected: %v", err)
		} else {
			assert.NotNil(t, pubKey)
			t.Logf("Recovered pubkey with compressed: %v", wasCompressed)
		}
	})
}

// TestVerifyMessageExtended provides comprehensive tests for VerifyMessage
func TestVerifyMessageExtended(t *testing.T) {
	t.Parallel()

	privateKey := "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"
	message := "Test message"

	// Get the address for this private key
	address, err := GetAddressFromPrivateKeyString(privateKey, true, true)
	require.NoError(t, err)

	// Create a valid signature
	signature, err := SignMessage(privateKey, message, true)
	require.NoError(t, err)

	t.Run("valid signature and address - mainnet", func(t *testing.T) {
		err := VerifyMessage(address, signature, message, true)
		require.NoError(t, err)
	})

	t.Run("wrong address", func(t *testing.T) {
		err := VerifyMessage("1C8bzHM8XFBHZ2ZZVvFy2NSoAZbwCXAicL", signature, message, true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("wrong message", func(t *testing.T) {
		err := VerifyMessage(address, signature, "Wrong message", true)
		require.Error(t, err)
	})

	t.Run("invalid signature", func(t *testing.T) {
		err := VerifyMessage(address, "invalid-signature", message, true)
		require.Error(t, err)
	})

	t.Run("empty signature", func(t *testing.T) {
		err := VerifyMessage(address, "", message, true)
		require.Error(t, err)
	})
}

// TestVerifyMessageDERExtended provides additional tests for VerifyMessageDER
func TestVerifyMessageDERExtended(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		hash        [32]byte
		pubKey      string
		signature   string
		shouldError bool
		expectValid bool
	}{
		{
			name: "invalid hex in signature",
			hash: [32]byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
				0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			},
			pubKey:      "0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798",
			signature:   "not-hex",
			shouldError: true,
			expectValid: false,
		},
		{
			name: "invalid hex in pubkey",
			hash: [32]byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
				0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			},
			pubKey:      "not-hex",
			signature:   "3044022044dc17b0d7a0b3a6c7f497836da77c2c4cb7a7c4f7c1d8f6d0f5e8e8d8c6e5d40220123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0",
			shouldError: true,
			expectValid: false,
		},
		{
			name: "empty signature",
			hash: [32]byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
				0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			},
			pubKey:      "0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798",
			signature:   "",
			shouldError: true,
			expectValid: false,
		},
		{
			name: "empty pubkey",
			hash: [32]byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
				0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			},
			pubKey:      "",
			signature:   "3044022044dc17b0d7a0b3a6c7f497836da77c2c4cb7a7c4f7c1d8f6d0f5e8e8d8c6e5d40220123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0",
			shouldError: true,
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verified, err := VerifyMessageDER(tt.hash, tt.pubKey, tt.signature)

			if tt.shouldError {
				require.Error(t, err, "Expected error but got none")
				assert.False(t, verified, "Should not be verified on error")
			} else {
				require.NoError(t, err, "Unexpected error: %v", err)
				assert.Equal(t, tt.expectValid, verified, "Verification result mismatch")
			}
		})
	}
}
