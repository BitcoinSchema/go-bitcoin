package bitcoin

import (
	"github.com/bitcoinsv/bsvd/chaincfg"
	"github.com/bitcoinsv/bsvutil/hdkeychain"
)

const (
	// RecommendedSeedLength is the recommended length in bytes for a seed to a master node.
	RecommendedSeedLength = 32 // 256 bits

	// SecureSeedLength is the max size of a seed length (most secure
	SecureSeedLength = 64 // 512 bits
)

// GenerateHDKey will create a new master node for use in creating a hierarchical deterministic key chain
func GenerateHDKey(seedLength uint8) (hdKey *hdkeychain.ExtendedKey, err error) {

	// Missing or invalid seed length
	if seedLength == 0 {
		seedLength = RecommendedSeedLength
	}

	// Generate a new seed (added extra security from 256 to 512 bits for seed length)
	var seed []byte
	if seed, err = hdkeychain.GenerateSeed(seedLength); err != nil {
		return
	}

	// Generate a new master key
	return hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
}

// GenerateHDKeyPair will generate a new xPub HD master node (xPrivateKey & xPublicKey)
func GenerateHDKeyPair(seedLength uint8) (xPrivateKey, xPublicKey string, err error) {

	// Generate an HD master key
	var masterKey *hdkeychain.ExtendedKey
	if masterKey, err = GenerateHDKey(seedLength); err != nil {
		return
	}

	// Set the xPriv
	xPrivateKey = masterKey.String()

	// Create the extended public key
	var pubKey *hdkeychain.ExtendedKey
	if pubKey, err = masterKey.Neuter(); err != nil {
		// Error should nearly never occur since it's using a safely derived masterKey
		return
	}

	// Set the actual xPub
	xPublicKey = pubKey.String()

	return
}
