package bitcoin

import (
	"errors"
	"fmt"
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

// CreateTxWithChange will automatically create the change output and calculate fees
//
// Use this if you don't want to figure out fees/change for a tx
func CreateTxWithChange(utxos []*Utxo, payToAddresses []*PayToAddress, opReturns []OpReturnData,
	changeAddress string, standardRate, dataRate float64, wif string) (*transaction.Transaction, error) {

	// Missing utxo(s) or change address
	if len(utxos) == 0 {
		return nil, errors.New("utxo(s) are required to create a tx")
	} else if len(changeAddress) == 0 {
		return nil, errors.New("change address is required")
	}

	// Accumulate the total satoshis from all utxo(s)
	var totalSatoshis uint64
	var totalPayToSatoshis uint64

	// Loop utxos and get total usable satoshis
	for _, utxo := range utxos {
		totalSatoshis += utxo.Satoshis
	}

	// Loop all payout address amounts
	for _, address := range payToAddresses {
		totalPayToSatoshis += address.Satoshis
	}

	// Sanity check - already not enough satoshis?
	if totalPayToSatoshis >= totalSatoshis {
		return nil, fmt.Errorf("not enough in utxo(s) to cover: %d + (fee) found: %d", totalPayToSatoshis, totalSatoshis)
	}

	// Add the change address as the difference (all change except 1 sat for Draft tx)
	payToAddresses = append(payToAddresses, &PayToAddress{
		Address:  changeAddress,
		Satoshis: totalSatoshis - (totalPayToSatoshis + 1),
	})

	// Create the "Draft tx"
	tx, err := CreateTx(utxos, payToAddresses, opReturns, wif)
	if err != nil {
		return nil, err
	}

	// Calculate the fees for the "Draft tx"
	fee := CalculateFeeForTx(tx, standardRate, dataRate)

	// Check that we have enough to cover the fee
	if (totalPayToSatoshis + fee) > totalSatoshis {
		return nil, fmt.Errorf("not enough in utxo(s) to cover: %d found: %d", totalPayToSatoshis+fee, totalSatoshis)
	}

	// Remove the change address (old version with original satoshis)
	payToAddresses = payToAddresses[:len(payToAddresses)-1]

	// Add the change address as the difference (now with adjusted fee)
	payToAddresses = append(payToAddresses, &PayToAddress{
		Address:  changeAddress,
		Satoshis: totalSatoshis - (totalPayToSatoshis + fee),
	})

	// Create the "Final tx" (or error)
	return CreateTx(utxos, payToAddresses, opReturns, wif)
}

// CreateTx will create a basic transaction and return the raw transaction (*transaction.Transaction)
//
// Note: this will NOT create a "change" address (it's assumed you have already specified an address)
// Note: this will NOT handle "fee" calculation (it's assumed you have already calculated the fee)
//
// Get the raw hex version: tx.ToString()
// Get the tx id: tx.GetTxID()
func CreateTx(utxos []*Utxo, addresses []*PayToAddress,
	opReturns []OpReturnData, wif string) (*transaction.Transaction, error) {

	// Missing utxo(s)
	if len(utxos) == 0 {
		return nil, errors.New("utxo(s) are required to create a tx")
	}

	// Start creating a new transaction
	tx := transaction.New()

	// Accumulate the total satoshis from all utxo(s)
	var totalSatoshis uint64

	// Loop all utxos and add to the transaction
	var err error
	for _, utxo := range utxos {
		if err = tx.From(utxo.TxID, utxo.Vout, utxo.ScriptSig, utxo.Satoshis); err != nil {
			return nil, err
		}
		totalSatoshis += utxo.Satoshis
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

	// Sanity check - not enough satoshis in utxo(s) to cover all paid amount(s)
	// They should never be equal, since the fee is the spread between the two amounts
	totalOutputSatoshis := tx.GetTotalOutputSatoshis() // Does not work properly
	if totalOutputSatoshis >= totalSatoshis {
		return nil, fmt.Errorf("not enough in utxo(s) to cover: %d + (fee) found: %d", totalOutputSatoshis, totalSatoshis)
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
// Reference: https://tncpw.co/c215a75c
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
			totalDataBytes += len(out.ToBytes())
		}
	}

	// Got some data bytes?
	if totalDataBytes > 0 {
		totalBytes = totalBytes - totalDataBytes
		totalFee += math.Ceil(float64(totalDataBytes) * dataRate)
	}

	// Still have regular standard bytes?
	if totalBytes > 0 {
		totalFee += math.Ceil(float64(totalBytes) * standardRate)
	}

	// Return the total fee as a uint (easier to use with satoshi values)
	return uint64(totalFee)
}
