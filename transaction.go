package bitcoin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/libsv/go-bt"
)

const (

	// Spec: https://github.com/bitcoin-sv-specs/brfc-misc/tree/master/feespec

	// DefaultDataRate is the default rate for feeType: data (500 satoshis per X bytes)
	DefaultDataRate int = 500

	// DefaultStandardRate is the default rate for feeType: standard (500 satoshis per X bytes)
	DefaultStandardRate int = 500

	// DefaultRateBytes is the default amount of bytes to calculate fees (X Satoshis per X bytes)
	DefaultRateBytes int = 1000

	// DustLimit is the minimum value for a tx that can be spent
	// Note: this is being deprecated in the new node software (TBD)
	DustLimit uint64 = 546
)

// FeeAmount is the actual fee for the given feeType (data or standard)
//
// Reference: https://github.com/tonicpow/go-minercraft/blob/b14d26a5d60436ecd3481f94d9cb468513dcf86b/fee_quote.go#L164
// Spec: https://github.com/bitcoin-sv-specs/brfc-misc/tree/master/feespec
type FeeAmount struct {
	Bytes    uint64 `json:"bytes"`
	Satoshis uint64 `json:"satoshis"`
}

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
	tx, err := CreateTx(utxos, payToAddresses, opReturns, privateKey)
	if err != nil {
		return nil, err
	}

	// Calculate the fees for the "Draft tx"
	fee := CalculateFeeForTx(tx, standardRate, dataRate)

	// todo: replace with go-bt way to create change tx (when released)
	// for now (hacking the fee) (ensure we are over the min fee for the miner)
	fee++

	// Check that we have enough to cover the fee
	if (totalPayToSatoshis + fee) > totalSatoshis {
		return nil, fmt.Errorf("not enough in utxo(s) to cover: %d found: %d", totalPayToSatoshis+fee, totalSatoshis)
	}

	// Remove the change address (old version with original satoshis)
	payToAddresses = payToAddresses[:len(payToAddresses)-1]

	// If change is less than dust...
	// if (totalSatoshis - (totalPayToSatoshis + fee)) < DustLimit {
	// todo: Warn about change amount being < dust ?
	// }

	// Add the change address as the difference (now with adjusted fee)
	payToAddresses = append(payToAddresses, &PayToAddress{
		Address:  changeAddress,
		Satoshis: totalSatoshis - (totalPayToSatoshis + fee),
	})

	// Create the "Final tx" (or error)
	return CreateTx(utxos, payToAddresses, opReturns, privateKey)
}

// CreateTxWithChangeUsingWif will automatically create the change output and calculate fees
//
// Use this if you don't want to figure out fees/change for a tx
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
