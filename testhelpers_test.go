package bitcoin

import (
	"testing"

	bip32 "github.com/bsv-blockchain/go-sdk/compat/bip32"
	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/stretchr/testify/require"
)

// newTestUtxo returns a *Utxo seeded with the shared test tx id / script and the
// given satoshi value. It centralizes the Utxo literal repeated across the
// transaction tests.
func newTestUtxo(satoshis uint64) *Utxo {
	return &Utxo{
		TxID:         testTxID,
		Vout:         0,
		ScriptPubKey: testScriptPubKey,
		Satoshis:     satoshis,
	}
}

// newTestOpReturns returns the standard two-element OpReturnData slice used
// throughout the transaction tests.
func newTestOpReturns() []OpReturnData {
	return []OpReturnData{
		{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}},
		{[]byte("prefix2"), []byte("more example data")},
	}
}

// mustTestPrivKey decodes the shared test WIF into an *ec.PrivateKey, failing the
// test on error.
func mustTestPrivKey(t *testing.T) *ec.PrivateKey {
	t.Helper()
	privateKey, err := WifToPrivateKey(testWIF)
	require.NoError(t, err)
	return privateKey
}

// mustHDKey generates an HD master key with the recommended seed length, failing
// the test on error.
func mustHDKey(t *testing.T) *bip32.ExtendedKey {
	t.Helper()
	key, err := GenerateHDKey(RecommendedSeedLength)
	require.NoError(t, err)
	return key
}
