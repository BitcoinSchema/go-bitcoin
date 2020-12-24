package bitcoin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/libsv/go-bt"
)

const (

	// DustLimit is the minimum value for a tx that can be spent
	// Note: this is being deprecated in the new node software (TBD)
	DustLimit uint64 = 546
)

// Utxo is an unspent transaction output
type Utxo struct {
	Satoshis     uint64 `json:"satoshis"`
	ScriptPubKey string `json:"string"`
	TxID         string `json:"tx_id"`
	Vout         uint32 `json:"vout"`
}

// PayToAddress is the pay-to-address
type PayToAddress struct {
	Address  string `json:"address"`
	Satoshis uint64 `json:"satoshis"`
}

// OpReturnData is the op return data to include in the tx
type OpReturnData [][]byte

// TxFromHex will return a libsv.tx from a raw hex string
func TxFromHex(rawHex string) (*bt.Tx, error) {
	return bt.NewTxFromString(rawHex)
}

// CreateTxWithChange will automatically create the change output and calculate fees
//
// Use this if you don't want to figure out fees/change for a tx
// USE AT YOUR OWN RISK - this will modify a "pay-to" output to accomplish auto-fees
func CreateTxWithChange(utxos []*Utxo, payToAddresses []*PayToAddress, opReturns []OpReturnData,
	changeAddress string, standardRate, dataRate *bt.Fee,
	privateKey *bsvec.PrivateKey) (*bt.Tx, error) {

	// Missing utxo(s) or change address
	if len(utxos) == 0 {
		return nil, errors.New("utxo(s) are required to create a tx")
	} else if len(changeAddress) == 0 {
		return nil, errors.New("change address is required")
	}

	// Accumulate the total satoshis from all utxo(s)
	var totalSatoshis uint64
	var totalPayToSatoshis uint64
	var remainder uint64
	var hasChange bool

	// Loop utxos and get total usable satoshis
	for _, utxo := range utxos {
		totalSatoshis += utxo.Satoshis
	}

	// Loop all payout address amounts
	for _, address := range payToAddresses {
		totalPayToSatoshis += address.Satoshis
	}

	// Sanity check - already not enough satoshis?
	if totalPayToSatoshis > totalSatoshis {
		return nil, fmt.Errorf(
			"not enough in utxo(s) to cover: %d + (fee), total found: %d",
			totalPayToSatoshis,
			totalSatoshis,
		)
	}

	// Add the change address as the difference (all change except 1 sat for Draft tx)
	// Only if the tx is NOT for the full amount
	if totalPayToSatoshis != totalSatoshis {
		hasChange = true
		payToAddresses = append(payToAddresses, &PayToAddress{
			Address:  changeAddress,
			Satoshis: totalSatoshis - (totalPayToSatoshis + 1),
		})
	}

	// Create the "Draft tx"
	fee, err := draftTx(utxos, payToAddresses, opReturns, privateKey, standardRate, dataRate)
	if err != nil {
		return nil, err
	}

	// Check that we have enough to cover the fee
	if (totalPayToSatoshis + fee) > totalSatoshis {

		// Remove temporary change address first
		if hasChange {
			payToAddresses = payToAddresses[:len(payToAddresses)-1]
		}

		// Re-run draft tx with no change address
		if fee, err = draftTx(
			utxos, payToAddresses, opReturns, privateKey, standardRate, dataRate,
		); err != nil {
			return nil, err
		}

		// Get the remainder missing (handle negative overflow safer)
		totalToPay := totalPayToSatoshis + fee
		if totalToPay >= totalSatoshis {
			remainder = totalToPay - totalSatoshis
		} else {
			remainder = totalSatoshis - totalToPay
		}

		// Remove remainder from last used payToAddress (or continue until found)
		feeAdjusted := false
		for i := len(payToAddresses) - 1; i >= 0; i-- { // Working backwards
			if payToAddresses[i].Satoshis > remainder {
				payToAddresses[i].Satoshis = payToAddresses[i].Satoshis - remainder
				feeAdjusted = true
				break
			}
		}

		// Fee was not adjusted (all inputs do not cover the fee)
		if !feeAdjusted {
			return nil, fmt.Errorf(
				"auto-fee could not be applied without removing an output (payTo %d) "+
					"(amount %d) (remainder %d) (fee %d) (total %d)",
				len(payToAddresses), totalPayToSatoshis, remainder, fee, totalSatoshis,
			)
		}

	} else {

		// Remove the change address (old version with original satoshis)
		// Add the change address as the difference (now with adjusted fee)
		if hasChange {
			payToAddresses = payToAddresses[:len(payToAddresses)-1]

			payToAddresses = append(payToAddresses, &PayToAddress{
				Address:  changeAddress,
				Satoshis: totalSatoshis - (totalPayToSatoshis + fee),
			})
		}
	}

	// Create the "Final tx" (or error)
	return CreateTx(utxos, payToAddresses, opReturns, privateKey)
}

// draftTx is a helper method to create a draft tx and associated fees
func draftTx(utxos []*Utxo, payToAddresses []*PayToAddress, opReturns []OpReturnData,
	privateKey *bsvec.PrivateKey, standardRate, dataRate *bt.Fee) (uint64, error) {

	// Create the "Draft tx"
	tx, err := CreateTx(utxos, payToAddresses, opReturns, privateKey)
	if err != nil {
		return 0, err
	}

	// Calculate the fees for the "Draft tx"
	// todo: hack to add 1 extra sat - ensuring that fee is over the minimum with rounding issues in WOC and other systems
	fee := CalculateFeeForTx(tx, standardRate, dataRate) + 1
	return fee, nil
}

// CreateTxWithChangeUsingWif will automatically create the change output and calculate fees
//
// Use this if you don't want to figure out fees/change for a tx
// USE AT YOUR OWN RISK - this will modify a "pay-to" output to accomplish auto-fees
func CreateTxWithChangeUsingWif(utxos []*Utxo, payToAddresses []*PayToAddress, opReturns []OpReturnData,
	changeAddress string, standardRate, dataRate *bt.Fee, wif string) (*bt.Tx, error) {

	// Decode the WIF
	privateKey, err := WifToPrivateKey(wif)
	if err != nil {
		return nil, err
	}

	// Create the "Final tx" (or error)
	return CreateTxWithChange(utxos, payToAddresses, opReturns, changeAddress, standardRate, dataRate, privateKey)
}

// CreateTx will create a basic transaction and return the raw transaction (*transaction.Transaction)
//
// Note: this will NOT create a change output (funds are sent to "addresses")
// Note: this will NOT handle fee calculation (it's assumed you have already calculated the fee)
//
// Get the raw hex version: tx.ToString()
// Get the tx id: tx.GetTxID()
func CreateTx(utxos []*Utxo, addresses []*PayToAddress,
	opReturns []OpReturnData, privateKey *bsvec.PrivateKey) (*bt.Tx, error) {

	// Start creating a new transaction
	tx := bt.NewTx()

	// Accumulate the total satoshis from all utxo(s)
	var totalSatoshis uint64

	// Loop all utxos and add to the transaction
	var err error
	for _, utxo := range utxos {
		if err = tx.From(utxo.TxID, utxo.Vout, utxo.ScriptPubKey, utxo.Satoshis); err != nil {
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
	var outPut *bt.Output
	for _, op := range opReturns {
		if outPut, err = bt.NewOpReturnPartsOutput(op); err != nil {
			return nil, err
		}
		tx.AddOutput(outPut)
	}

	// If inputs are supplied, make sure they are sufficient for this transaction
	if len(tx.GetInputs()) > 0 {
		// Sanity check - not enough satoshis in utxo(s) to cover all paid amount(s)
		// They should never be equal, since the fee is the spread between the two amounts
		totalOutputSatoshis := tx.GetTotalOutputSatoshis() // Does not work properly
		if totalOutputSatoshis > totalSatoshis {
			return nil, fmt.Errorf("not enough in utxo(s) to cover: %d + (fee) found: %d", totalOutputSatoshis, totalSatoshis)
		}
	}

	// Sign the transaction
	if privateKey != nil {

		signer := bt.InternalSigner{PrivateKey: privateKey, SigHashFlag: 0}
		if _, err = tx.SignAuto(&signer); err != nil {
			return nil, err
		}
	}

	// Return the transaction as a raw string
	return tx, nil
}

// CreateTxUsingWif will create a basic transaction and return the raw transaction (*transaction.Transaction)
//
// Note: this will NOT create a "change" address (it's assumed you have already specified an address)
// Note: this will NOT handle "fee" calculation (it's assumed you have already calculated the fee)
//
// Get the raw hex version: tx.ToString()
// Get the tx id: tx.GetTxID()
func CreateTxUsingWif(utxos []*Utxo, addresses []*PayToAddress,
	opReturns []OpReturnData, wif string) (*bt.Tx, error) {

	// Decode the WIF
	privateKey, err := WifToPrivateKey(wif)
	if err != nil {
		return nil, err
	}

	// Create the Tx
	return CreateTx(utxos, addresses, opReturns, privateKey)
}

// CalculateFeeForTx will estimate a fee for the given transaction
//
// If tx is nil this will panic
// Rate(s) can be derived from MinerAPI (default is DefaultDataRate and DefaultStandardRate)
// If rate is nil it will use default rates (0.5 sat per byte)
// Reference: https://tncpw.co/c215a75c
func CalculateFeeForTx(tx *bt.Tx, standardRate, dataRate *bt.Fee) uint64 {

	// Set the totals
	var totalFee int
	var totalDataBytes int

	// Set defaults if not found
	if standardRate == nil {
		standardRate = bt.DefaultStandardFee()
	}
	if dataRate == nil {
		dataRate = bt.DefaultStandardFee()
		// todo: adjusted to 5/10 for now, since all miners accept that rate
		dataRate.FeeType = bt.FeeTypeData
	}

	// Set the total bytes of the tx
	totalBytes := len(tx.ToBytes())

	// Loop all outputs and accumulate size (find data related outputs)
	for _, out := range tx.GetOutputs() {
		outHexString := out.GetLockingScriptHexString()
		if strings.HasPrefix(outHexString, "006a") || strings.HasPrefix(outHexString, "6a") {
			totalDataBytes += len(out.ToBytes())
		}
	}

	// Got some data bytes?
	if totalDataBytes > 0 {
		totalBytes = totalBytes - totalDataBytes
		totalFee += (dataRate.MiningFee.Satoshis * totalDataBytes) / dataRate.MiningFee.Bytes
	}

	// Still have regular standard bytes?
	if totalBytes > 0 {
		totalFee += (standardRate.MiningFee.Satoshis * totalBytes) / standardRate.MiningFee.Bytes
	}

	// Safety check (possible division by zero?)
	if totalFee == 0 {
		totalFee = 1
	}

	// Return the total fee as a uint (easier to use with satoshi values)
	return uint64(totalFee)
}
