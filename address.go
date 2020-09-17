package bitcoin

import (
	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvd/chaincfg"
	"github.com/bitcoinsv/bsvutil"
)

func AddressFromPrivKey(privKey string) string {
	pubKey := PrivateKey(privKey).PubKey()
	return Address(pubKey).EncodeAddress()
}

// Address gets a bsvutil.LegacyAddressPubKeyHash
func Address(publicKey *bsvec.PublicKey) (address *bsvutil.LegacyAddressPubKeyHash) {
	publicKeyHash := bsvutil.Hash160(publicKey.SerializeCompressed())
	address, _ = bsvutil.NewLegacyAddressPubKeyHash(publicKeyHash, &chaincfg.MainNetParams)
	return
}
