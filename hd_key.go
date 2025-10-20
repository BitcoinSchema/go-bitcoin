package bitcoin

import (
	"encoding/hex"

	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/chaincfg"
	"github.com/libsv/go-bt/v2/bscript"
)

const (
	// RecommendedSeedLength is the recommended length in bytes for a seed to a master node
	RecommendedSeedLength = 32 // 256 bits

	// SecureSeedLength is the max size of a seed length (most secure)
	SecureSeedLength = 64 // 512 bits

	// DefaultExternalChain is the default external chain (for public use to accept incoming txs)
	// Reference: https://en.bitcoin.it/wiki/BIP_0032#The_default_wallet_layout
	DefaultExternalChain = 0

	// DefaultInternalChain is the default internal chain (for change, generating, other purposes...)
	// Reference: https://en.bitcoin.it/wiki/BIP_0032#The_default_wallet_layout
	DefaultInternalChain = 1
)

// GenerateHDKey will create a new master node for use in creating a hierarchical deterministic keychain
func GenerateHDKey(seedLength uint8) (hdKey *bip32.ExtendedKey, err error) {

	// Missing or invalid seed length
	if seedLength == 0 {
		seedLength = RecommendedSeedLength
	}

	// Generate a new seed (added extra security from 256 to 512 bits for seed length)
	var seed []byte
	if seed, err = bip32.GenerateSeed(seedLength); err != nil {
		return
	}

	// Generate a new master key
	return bip32.NewMaster(seed, &chaincfg.MainNet)
}

// GenerateHDKeyFromString will create a new master node for use in creating a
// hierarchical deterministic keychain from an xPrivKey string
func GenerateHDKeyFromString(xPriv string) (hdKey *bip32.ExtendedKey, err error) {
	return bip32.NewKeyFromString(xPriv)
}

// GenerateHDKeyPair will generate a new xPub HD master node (xPrivateKey & xPublicKey)
func GenerateHDKeyPair(seedLength uint8) (xPrivateKey, xPublicKey string, err error) {

	// Generate an HD master key
	var masterKey *bip32.ExtendedKey
	if masterKey, err = GenerateHDKey(seedLength); err != nil {
		return
	}

	// Set the xPriv (string)
	xPrivateKey = masterKey.String()

	// Set the xPub (string)
	xPublicKey, err = GetExtendedPublicKey(masterKey)

	return
}

// GetHDKeyByPath gets the corresponding HD key from a chain/num path
// Reference: https://en.bitcoin.it/wiki/BIP_0032#The_default_wallet_layout
func GetHDKeyByPath(hdKey *bip32.ExtendedKey, chain, num uint32) (*bip32.ExtendedKey, error) {

	// Derive the child key from the chain path
	childKeyChain, err := GetHDKeyChild(hdKey, chain)
	if err != nil {
		return nil, err
	}

	// Get the child key from the num path
	return GetHDKeyChild(childKeyChain, num)
}

// GetHDKeyChild gets the child hd key for a given num
// Note: For a hardened child, start at 0x80000000. (For reference, 0x8000000 = 0')
//
// Expects hdKey to not be nil (otherwise will panic)
func GetHDKeyChild(hdKey *bip32.ExtendedKey, num uint32) (*bip32.ExtendedKey, error) {
	return hdKey.Child(num)
}

// GetPrivateKeyByPath gets the key for a given derivation path (chain/num)
//
// Expects hdKey to not be nil (otherwise will panic)
func GetPrivateKeyByPath(hdKey *bip32.ExtendedKey, chain, num uint32) (*bec.PrivateKey, error) {

	// Get the child key from the num & chain
	childKeyNum, err := GetHDKeyByPath(hdKey, chain, num)
	if err != nil {
		return nil, err
	}

	// Get the private key
	return childKeyNum.ECPrivKey()
}

// GetPrivateKeyFromHDKey is a helper function to get the Private Key associated
// with a given hdKey
//
// Expects hdKey to not be nil (otherwise will panic)
func GetPrivateKeyFromHDKey(hdKey *bip32.ExtendedKey) (*bec.PrivateKey, error) {
	return hdKey.ECPrivKey()
}

// GetPrivateKeyStringFromHDKey is a helper function to get the Private Key (string)
// associated with a given hdKey
//
// Expects hdKey to not be nil (otherwise will panic)
func GetPrivateKeyStringFromHDKey(hdKey *bip32.ExtendedKey) (string, error) {
	key, err := GetPrivateKeyFromHDKey(hdKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key.Serialise()), nil
}

// GetPublicKeyFromHDKey is a helper function to get the Public Key associated with a given hdKey
//
// Expects hdKey to not be nil (otherwise will panic)
func GetPublicKeyFromHDKey(hdKey *bip32.ExtendedKey) (*bec.PublicKey, error) {
	return hdKey.ECPubKey()
}

// GetAddressFromHDKey is a helper function to get the Address associated with a given hdKey
//
// Expects hdKey to not be nil (otherwise will panic)
func GetAddressFromHDKey(hdKey *bip32.ExtendedKey, mainnet bool) (*bscript.Address, error) {
	pubKey, err := GetPublicKeyFromHDKey(hdKey)
	if err != nil {
		return nil, err
	}
	return GetAddressFromPubKey(pubKey, true, mainnet)
}

// GetAddressStringFromHDKey is a helper function to get the Address (string) associated with a given hdKey
//
// Expects hdKey to not be nil (otherwise will panic)
func GetAddressStringFromHDKey(hdKey *bip32.ExtendedKey, mainnet bool) (string, error) {
	address, err := GetAddressFromHDKey(hdKey, mainnet)
	if err != nil {
		return "", err
	}
	return address.AddressString, nil
}

// GetPublicKeysForPath gets the PublicKeys for a given derivation path
// Uses the standard m/0/0 (external) and m/0/1 (internal) paths
// Reference: https://en.bitcoin.it/wiki/BIP_0032#The_default_wallet_layout
func GetPublicKeysForPath(hdKey *bip32.ExtendedKey, num uint32) (pubKeys []*bec.PublicKey, err error) {

	//  m/0/x
	var childM0x *bip32.ExtendedKey
	if childM0x, err = GetHDKeyByPath(hdKey, DefaultExternalChain, num); err != nil {
		return
	}

	// Get the external pubKey from m/0/x
	var pubKey *bec.PublicKey
	if pubKey, err = childM0x.ECPubKey(); err != nil {
		// Should never error since the previous method ensures a valid hdKey
		return
	}
	pubKeys = append(pubKeys, pubKey)

	//  m/1/x
	var childM1x *bip32.ExtendedKey
	if childM1x, err = GetHDKeyByPath(hdKey, DefaultInternalChain, num); err != nil {
		// Should never error since the previous method ensures a valid hdKey
		return
	}

	// Get the internal pubKey from m/1/x
	if pubKey, err = childM1x.ECPubKey(); err != nil {
		// Should never error since the previous method ensures a valid hdKey
		return
	}
	pubKeys = append(pubKeys, pubKey)

	return
}

// GetAddressesForPath will get the corresponding addresses for the PublicKeys at the given path m/0/x
// Returns 2 keys, first is internal and second is external
func GetAddressesForPath(hdKey *bip32.ExtendedKey, num uint32, mainnet bool) (addresses []string, err error) {

	// Get the public keys for the corresponding chain/num (using default chain)
	var pubKeys []*bec.PublicKey
	if pubKeys, err = GetPublicKeysForPath(hdKey, num); err != nil {
		return
	}

	// Loop, get address and append to results
	var address *bscript.Address
	for _, key := range pubKeys {
		if address, err = GetAddressFromPubKey(key, true, mainnet); err != nil {
			// Should never error if the pubKeys are valid keys
			return
		}
		addresses = append(addresses, address.AddressString)
	}

	return
}

// GetExtendedPublicKey will get the extended public key (xPub)
func GetExtendedPublicKey(hdKey *bip32.ExtendedKey) (string, error) {

	// Neuter the extended public key from hd key
	pub, err := hdKey.Neuter()
	if err != nil {
		// Error should never occur if using a valid hd key
		return "", err
	}

	// Return the string version
	return pub.String(), nil
}

// GetHDKeyFromExtendedPublicKey will get the hd key from an existing extended public key (xPub)
func GetHDKeyFromExtendedPublicKey(xPublicKey string) (*bip32.ExtendedKey, error) {
	return bip32.NewKeyFromString(xPublicKey)
}
