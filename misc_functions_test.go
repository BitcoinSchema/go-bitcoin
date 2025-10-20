package bitcoin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEncryptDecryptExtended provides comprehensive encryption/decryption tests
func TestEncryptDecryptExtended(t *testing.T) {
	t.Parallel()

	privateKey, err := CreatePrivateKey()
	require.NoError(t, err)

	privateKeyString, err := CreatePrivateKeyString()
	require.NoError(t, err)

	testData := "Test data to encrypt"

	t.Run("encrypt and decrypt with private key object", func(t *testing.T) {
		encrypted, err := EncryptWithPrivateKey(privateKey, testData)
		require.NoError(t, err)
		assert.NotEmpty(t, encrypted)

		decrypted, err := DecryptWithPrivateKey(privateKey, encrypted)
		require.NoError(t, err)
		assert.Equal(t, testData, decrypted)
	})

	t.Run("encrypt and decrypt with private key string", func(t *testing.T) {
		encrypted, err := EncryptWithPrivateKeyString(privateKeyString, testData)
		require.NoError(t, err)
		assert.NotEmpty(t, encrypted)

		decrypted, err := DecryptWithPrivateKeyString(privateKeyString, encrypted)
		require.NoError(t, err)
		assert.Equal(t, testData, decrypted)
	})

	t.Run("encrypt with empty data", func(t *testing.T) {
		encrypted, err := EncryptWithPrivateKey(privateKey, "")
		require.NoError(t, err)
		assert.NotEmpty(t, encrypted)
	})

	t.Run("decrypt with invalid hex", func(t *testing.T) {
		decrypted, err := DecryptWithPrivateKey(privateKey, "not-hex")
		require.Error(t, err)
		assert.Empty(t, decrypted)
	})

	t.Run("encrypt with invalid private key string", func(t *testing.T) {
		encrypted, err := EncryptWithPrivateKeyString("invalid", testData)
		require.Error(t, err)
		assert.Empty(t, encrypted)
	})

	t.Run("decrypt with invalid private key string", func(t *testing.T) {
		decrypted, err := DecryptWithPrivateKeyString("invalid", "abcd1234")
		require.Error(t, err)
		assert.Empty(t, decrypted)
	})
}

// TestEncryptSharedExtended tests shared encryption functionality
func TestEncryptSharedExtended(t *testing.T) {
	t.Parallel()

	user1Key, err := CreatePrivateKey()
	require.NoError(t, err)

	user2Key, err := CreatePrivateKey()
	require.NoError(t, err)

	testData := []byte("Shared secret data")

	t.Run("encrypt and decrypt shared data", func(t *testing.T) {
		sharedPrivKey, sharedPubKey, encrypted, err := EncryptShared(user1Key, user2Key.PubKey(), testData)
		require.NoError(t, err)
		assert.NotNil(t, sharedPrivKey)
		assert.NotNil(t, sharedPubKey)
		assert.NotEmpty(t, encrypted)
	})

	t.Run("encrypt shared string", func(t *testing.T) {
		sharedPrivKey, sharedPubKey, encrypted, err := EncryptSharedString(user1Key, user2Key.PubKey(), string(testData))
		require.NoError(t, err)
		assert.NotNil(t, sharedPrivKey)
		assert.NotNil(t, sharedPubKey)
		assert.NotEmpty(t, encrypted)
		assert.IsType(t, "", encrypted, "Should return string")
	})
}

// TestCreateKeysExtended tests key creation functions
func TestCreateKeysExtended(t *testing.T) {
	t.Parallel()

	t.Run("create private key string", func(t *testing.T) {
		keyString, err := CreatePrivateKeyString()
		require.NoError(t, err)
		assert.NotEmpty(t, keyString)
		assert.Len(t, keyString, 64, "Private key should be 64 hex characters")
	})

	t.Run("create wif", func(t *testing.T) {
		wifObj, err := CreateWif()
		require.NoError(t, err)
		assert.NotNil(t, wifObj)
		assert.NotEmpty(t, wifObj.String())
	})

	t.Run("create wif string", func(t *testing.T) {
		wifString, err := CreateWifString()
		require.NoError(t, err)
		assert.NotEmpty(t, wifString)
	})

	t.Run("private key to wif and back", func(t *testing.T) {
		keyString, err := CreatePrivateKeyString()
		require.NoError(t, err)

		wifString, err := PrivateKeyToWifString(keyString)
		require.NoError(t, err)
		assert.NotEmpty(t, wifString)

		recoveredKey, err := WifToPrivateKeyString(wifString)
		require.NoError(t, err)
		assert.Equal(t, keyString, recoveredKey)
	})

	t.Run("private key to wif - empty key error", func(t *testing.T) {
		wif, err := PrivateKeyToWif("")
		require.Error(t, err)
		assert.Nil(t, wif)
		assert.Contains(t, err.Error(), "missing")
	})

	t.Run("private key to wif string - invalid hex", func(t *testing.T) {
		wifString, err := PrivateKeyToWifString("not-hex")
		require.Error(t, err)
		assert.Empty(t, wifString)
	})
}

// TestGetAddressesExtended tests address-related functions
func TestGetAddressesExtended(t *testing.T) {
	t.Parallel()

	privateKey := "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"

	t.Run("get address from private key string - mainnet compressed", func(t *testing.T) {
		address, err := GetAddressFromPrivateKeyString(privateKey, true, true)
		require.NoError(t, err)
		assert.NotEmpty(t, address)
	})

	t.Run("get address from private key string - mainnet uncompressed", func(t *testing.T) {
		address, err := GetAddressFromPrivateKeyString(privateKey, false, true)
		require.NoError(t, err)
		assert.NotEmpty(t, address)
	})

	t.Run("get address from private key string - testnet", func(t *testing.T) {
		address, err := GetAddressFromPrivateKeyString(privateKey, true, false)
		require.NoError(t, err)
		assert.NotEmpty(t, address)
	})

	t.Run("get address from private key string - invalid hex", func(t *testing.T) {
		address, err := GetAddressFromPrivateKeyString("not-hex", true, true)
		require.Error(t, err)
		assert.Empty(t, address)
	})

	t.Run("get address from script - valid", func(t *testing.T) {
		script := "76a914b424110292f4ea2ac92beb9e83cf5e6f0fa2996388ac"
		address, err := GetAddressFromScript(script)
		require.NoError(t, err)
		assert.NotEmpty(t, address)
	})

	t.Run("get address from script - empty script", func(t *testing.T) {
		address, err := GetAddressFromScript("")
		require.Error(t, err)
		assert.Empty(t, address)
		assert.Contains(t, err.Error(), "missing")
	})

	t.Run("get address from script - invalid hex", func(t *testing.T) {
		address, err := GetAddressFromScript("not-hex")
		require.Error(t, err)
		assert.Empty(t, address)
	})

	t.Run("get address from pubkey string - compressed", func(t *testing.T) {
		pk, err := PrivateKeyFromString(privateKey)
		require.NoError(t, err)

		pubKeyString := PubKeyFromPrivateKey(pk, true)
		address, err := GetAddressFromPubKeyString(pubKeyString, true, true)
		require.NoError(t, err)
		assert.NotNil(t, address)
		assert.NotEmpty(t, address.AddressString)
	})

	t.Run("get address from pubkey string - invalid", func(t *testing.T) {
		address, err := GetAddressFromPubKeyString("invalid", true, true)
		require.Error(t, err)
		assert.Nil(t, address)
	})
}

// TestHDKeyExtended provides comprehensive HD key tests
func TestHDKeyExtended(t *testing.T) {
	t.Parallel()

	t.Run("generate HD key with custom seed length", func(t *testing.T) {
		key, err := GenerateHDKey(RecommendedSeedLength)
		require.NoError(t, err)
		assert.NotNil(t, key)
	})

	t.Run("generate HD key with secure seed length", func(t *testing.T) {
		key, err := GenerateHDKey(SecureSeedLength)
		require.NoError(t, err)
		assert.NotNil(t, key)
	})

	t.Run("get extended public key from HD key", func(t *testing.T) {
		key, err := GenerateHDKey(RecommendedSeedLength)
		require.NoError(t, err)

		xPub, err := GetExtendedPublicKey(key)
		require.NoError(t, err)
		assert.NotEmpty(t, xPub)
		assert.Contains(t, xPub, "xpub")
	})

	t.Run("get private key by path", func(t *testing.T) {
		key, err := GenerateHDKey(RecommendedSeedLength)
		require.NoError(t, err)

		privKey, err := GetPrivateKeyByPath(key, 0, 0)
		require.NoError(t, err)
		assert.NotNil(t, privKey)
	})

	t.Run("get public keys for path", func(t *testing.T) {
		key, err := GenerateHDKey(RecommendedSeedLength)
		require.NoError(t, err)

		pubKeys, err := GetPublicKeysForPath(key, 0)
		require.NoError(t, err)
		assert.Len(t, pubKeys, 2, "Should return external and internal keys")
	})

	t.Run("get addresses for path - mainnet", func(t *testing.T) {
		key, err := GenerateHDKey(RecommendedSeedLength)
		require.NoError(t, err)

		addresses, err := GetAddressesForPath(key, 0, true)
		require.NoError(t, err)
		assert.Len(t, addresses, 2, "Should return external and internal addresses")
	})

	t.Run("get addresses for path - testnet", func(t *testing.T) {
		key, err := GenerateHDKey(RecommendedSeedLength)
		require.NoError(t, err)

		addresses, err := GetAddressesForPath(key, 0, false)
		require.NoError(t, err)
		assert.Len(t, addresses, 2, "Should return external and internal addresses")
	})

	t.Run("get address string from HD key", func(t *testing.T) {
		key, err := GenerateHDKey(RecommendedSeedLength)
		require.NoError(t, err)

		address, err := GetAddressStringFromHDKey(key, true)
		require.NoError(t, err)
		assert.NotEmpty(t, address)
	})

	t.Run("get private key string from HD key", func(t *testing.T) {
		key, err := GenerateHDKey(RecommendedSeedLength)
		require.NoError(t, err)

		privKeyString, err := GetPrivateKeyStringFromHDKey(key)
		require.NoError(t, err)
		assert.NotEmpty(t, privKeyString)
		assert.Len(t, privKeyString, 64)
	})
}

// TestPubKeyFunctions tests public key utility functions
func TestPubKeyFunctions(t *testing.T) {
	t.Parallel()

	privateKey := "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"

	t.Run("pubkey from private key string - compressed", func(t *testing.T) {
		pubKey, err := PubKeyFromPrivateKeyString(privateKey, true)
		require.NoError(t, err)
		assert.NotEmpty(t, pubKey)
		assert.Len(t, pubKey, 66, "Compressed pubkey should be 33 bytes (66 hex chars)")
	})

	t.Run("pubkey from private key string - uncompressed", func(t *testing.T) {
		pubKey, err := PubKeyFromPrivateKeyString(privateKey, false)
		require.NoError(t, err)
		assert.NotEmpty(t, pubKey)
		assert.Len(t, pubKey, 130, "Uncompressed pubkey should be 65 bytes (130 hex chars)")
	})

	t.Run("pubkey from string - invalid", func(t *testing.T) {
		pk, err := PubKeyFromString("invalid")
		require.Error(t, err)
		assert.Nil(t, pk)
	})

	t.Run("pubkey from string - empty", func(t *testing.T) {
		pk, err := PubKeyFromString("")
		require.Error(t, err)
		assert.Nil(t, pk)
		assert.Contains(t, err.Error(), "missing")
	})
}
