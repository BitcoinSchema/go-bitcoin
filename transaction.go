package bitcoin

import (
	"errors"
	"math"
	"strings"

	"github.com/bitcoinsv/bsvutil"
	"github.com/libsv/libsv/transaction"
	"github.com/libsv/libsv/transaction/output"
	"github.com/libsv/libsv/transaction/signature"
)

const (

	// DefaultDataRate is the default rate for feeType: data (0.5 satoshis per byte)
	DefaultDataRate = 0.5

	// DefaultStandardRate is the default rate for feeType: standard (0.5 satoshis per byte)
	DefaultStandardRate = 0.5
)

// Utxo is an unspent transaction output
type Utxo struct {
	Satoshis  uint64 `json:"satoshis"`
	ScriptSig string `json:"string"`
	TxID      string `json:"tx_id"`
	Vout      uint32 `json:"vout"`
}

// PayToAddress is the pay-to-address
type PayToAddress struct {
	Address  string `json:"address"`
	Satoshis uint64 `json:"satoshis"`
}

// OpReturnData is the op return data to include in the tx
type OpReturnData [][]byte

// TxFromHex will return a libsv.tx from a raw hex string
func TxFromHex(rawHex string) (*transaction.Transaction, error) {
	return transaction.NewFromString(rawHex)
}

// CreateTx will create a basic transaction and return the raw transaction (*transaction.Transaction)
//
// Get the raw hex version: tx.ToString()
// Get the tx id: tx.GetTxID()
func CreateTx(utxos []*Utxo, addresses []*PayToAddress, opReturns []OpReturnData, wif string) (*transaction.Transaction, error) {

	// Missing utxos
	if len(utxos) == 0 {
		return nil, errors.New("utxos are required to create a tx")
	}

	// Start creating a new transaction
	tx := transaction.New()

	// Loop all utxos and add to the transaction
	var err error
	for _, utxo := range utxos {
		if err = tx.From(utxo.TxID, utxo.Vout, utxo.ScriptSig, utxo.Satoshis); err != nil {
			return nil, err
		}
	}

	// Loop any pay addresses
	for _, address := range addresses {
		if err = tx.PayTo(address.Address, address.Satoshis); err != nil {
			return nil, err
		}
	}

	// Loop any op returns
	var outPut *output.Output
	for _, op := range opReturns {
		if outPut, err = output.NewOpReturnParts(op); err != nil {
			return nil, err
		}
		tx.AddOutput(outPut)
	}

	// Decode the WIF
	var decodedWif *bsvutil.WIF
	if decodedWif, err = bsvutil.DecodeWIF(wif); err != nil {
		return nil, err
	}

	// Sign the transaction
	signer := signature.InternalSigner{PrivateKey: decodedWif.PrivKey, SigHashFlag: 0}
	if err = tx.SignAuto(&signer); err != nil {
		return nil, err
	}

	// Return the transaction as a raw string
	return tx, nil
}

// CalculateFeeForTx will estimate a fee for the given transaction
//
// If tx is nil this will panic
// Rate(s) can be derived from MinerAPI (default is DefaultDataRate and DefaultStandardRate)
// Reference: https://github.com/bitcoin-sv-specs/brfc-misc/tree/master/feespec#calculating-tx-fee-with-different-feetypes
func CalculateFeeForTx(tx *transaction.Transaction, standardRate, dataRate float64) uint64 {

	// Set the totals
	var totalFee float64
	var totalDataBytes int

	// Set the total bytes of the tx
	totalBytes := len(tx.ToBytes())

	// Loop all outputs and accumulate size (find data related outputs)
	for _, out := range tx.GetOutputs() {
		// todo: once libsv has outs.data.ToBytes() this can be removed/optimized
		outHexString := out.GetLockingScriptHexString()
		if strings.HasPrefix(outHexString, "006a") || strings.HasPrefix(outHexString, "6a") {
			totalDataBytes = totalDataBytes + len(out.ToBytes())
		}
	}

	// Got some data bytes?
	if totalDataBytes > 0 {
		totalBytes = totalBytes - totalDataBytes
		totalFee = totalFee + math.Ceil(float64(totalDataBytes)*dataRate)
	}

	// Still have regular standard bytes?
	if totalBytes > 0 {
		totalFee = totalFee + math.Ceil(float64(totalBytes)*standardRate)
	}

	// Return the total fee as a uint (easier to use with satoshi values)
	return uint64(totalFee)
}
