package bitcoin

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/unlocker"
)

var (
	// ErrInsufficientFunds is returned when UTXOs don't cover the output amount plus fees
	ErrInsufficientFunds = errors.New("insufficient funds in UTXOs to cover outputs and fees")
	// ErrAutoFeeNotApplicable is returned when auto-fee calculation cannot be applied
	ErrAutoFeeNotApplicable = errors.New("auto-fee could not be applied without removing an output")
	// ErrEmptyTxHex is returned when an empty hex string is provided
	ErrEmptyTxHex = errors.New("transaction hex string is empty")
	// ErrTxHexTooLarge is returned when the hex string exceeds maximum size
	ErrTxHexTooLarge = errors.New("transaction hex string too large")
	// ErrTxHexOddLength is returned when the hex string has an odd number of characters
	ErrTxHexOddLength = errors.New("transaction hex string has odd length")
	// ErrTxParsePanic is returned when the underlying library panics during transaction parsing
	ErrTxParsePanic = errors.New("transaction parsing panic")
	// ErrVarintOffsetOutOfBounds is returned when varint offset is out of bounds
	ErrVarintOffsetOutOfBounds = errors.New("varint offset out of bounds")
	// ErrVarintTruncated is returned when varint is truncated
	ErrVarintTruncated = errors.New("varint truncated")
	// ErrVarintExceedsMax is returned when varint value exceeds maximum allowed value
	ErrVarintExceedsMax = errors.New("varint value exceeds maximum")
	// ErrVarintUnexpectedPrefix is returned when varint has an unexpected prefix
	ErrVarintUnexpectedPrefix = errors.New("unexpected varint prefix")
)

const (

	// DustLimit is the minimum value for a tx that can be spent
	// This is being deprecated in the new node software (TBD)
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

// account is a struct/interface for implementing unlocker
type account struct {
	PrivateKey *bec.PrivateKey
}

// Unlocker get the correct un-locker for a given locking script
func (a *account) Unlocker(context.Context, *bscript.Script) (bt.Unlocker, error) {
	return &unlocker.Simple{
		PrivateKey: a.PrivateKey,
	}, nil
}

// OpReturnData is the op return data to include in the tx
type OpReturnData [][]byte

// validateTransactionVarints checks that varints in the transaction are within reasonable bounds.
// This prevents OOM attacks from malformed transactions that claim billions of inputs/outputs or script lengths.
// We check specific known varint positions (input count, output count) rather than scanning all bytes.
func validateTransactionVarints(data []byte) error {
	// Maximum reasonable count for inputs/outputs
	// Bitcoin transactions realistically have at most thousands of inputs/outputs
	const maxReasonableCount = 100_000

	// Need at least version (4 bytes) + input count varint (1 byte minimum)
	if len(data) < 5 {
		return nil // Too short, will fail in library anyway
	}

	// Check input count varint at position 4 (after 4-byte version)
	inputCount, inputCountBytes, err := readAndValidateVarint(data, 4, maxReasonableCount)
	if err != nil {
		return fmt.Errorf("invalid input count: %w", err)
	}

	// If input count is 0, we can also validate the output count
	// (which would be immediately after the input count varint)
	if inputCount == 0 {
		outputCountPos := 4 + inputCountBytes
		if outputCountPos < len(data) {
			_, _, err := readAndValidateVarint(data, outputCountPos, maxReasonableCount)
			if err != nil {
				return fmt.Errorf("invalid output count: %w", err)
			}
		}
	}

	// Validation coverage: We only validate input and output count varints in simple cases.
	// The underlying library may still encounter OOM errors with malformed varints in
	// other positions (script lengths, witness data, etc.). This is a limitation of
	// the underlying libsv/go-bt library that we cannot fully work around without
	// reimplementing the entire transaction parser.

	return nil
}

// readAndValidateVarint reads a Bitcoin varint and validates it's within bounds.
func readAndValidateVarint(data []byte, offset int, maxValue uint64) (uint64, int, error) {
	if offset >= len(data) {
		return 0, 0, ErrVarintOffsetOutOfBounds
	}

	first := data[offset]

	// Single byte varint (0x00-0xFC)
	if first < 0xFD {
		val := uint64(first)
		if val > maxValue {
			return 0, 0, fmt.Errorf("%w: %d exceeds maximum %d", ErrVarintExceedsMax, val, maxValue)
		}
		return val, 1, nil
	}

	// Multi-byte varint
	switch first {
	case 0xFD: // Next 2 bytes (little-endian uint16)
		if offset+3 > len(data) {
			return 0, 0, ErrVarintTruncated
		}
		val := uint64(data[offset+1]) | uint64(data[offset+2])<<8
		if val > maxValue {
			return 0, 0, fmt.Errorf("%w: %d exceeds maximum %d", ErrVarintExceedsMax, val, maxValue)
		}
		return val, 3, nil

	case 0xFE: // Next 4 bytes (little-endian uint32)
		if offset+5 > len(data) {
			return 0, 0, ErrVarintTruncated
		}
		val := uint64(data[offset+1]) |
			uint64(data[offset+2])<<8 |
			uint64(data[offset+3])<<16 |
			uint64(data[offset+4])<<24
		if val > maxValue {
			return 0, 0, fmt.Errorf("%w: %d exceeds maximum %d", ErrVarintExceedsMax, val, maxValue)
		}
		return val, 5, nil

	case 0xFF: // Next 8 bytes (little-endian uint64)
		if offset+9 > len(data) {
			return 0, 0, ErrVarintTruncated
		}
		val := uint64(data[offset+1]) |
			uint64(data[offset+2])<<8 |
			uint64(data[offset+3])<<16 |
			uint64(data[offset+4])<<24 |
			uint64(data[offset+5])<<32 |
			uint64(data[offset+6])<<40 |
			uint64(data[offset+7])<<48 |
			uint64(data[offset+8])<<56
		if val > maxValue {
			return 0, 0, fmt.Errorf("%w: %d exceeds maximum %d", ErrVarintExceedsMax, val, maxValue)
		}
		return val, 9, nil

	default:
		return 0, 0, fmt.Errorf("%w: 0x%02x", ErrVarintUnexpectedPrefix, first)
	}
}

// TxFromHex will return a libsv.tx from a raw hex string
func TxFromHex(rawHex string) (tx *bt.Tx, err error) {
	// Validate input is not empty
	if len(rawHex) == 0 {
		return nil, ErrEmptyTxHex
	}

	// Validate input size to prevent out-of-memory errors
	// Bitcoin transactions are typically a few KB; 200KB hex (100KB decoded) is a conservative upper bound
	// This also prevents malformed varint exploits in the underlying library
	const maxTxHexSize = 200_000
	if len(rawHex) > maxTxHexSize {
		return nil, fmt.Errorf("%w: %d characters (max %d)", ErrTxHexTooLarge, len(rawHex), maxTxHexSize)
	}

	// Validate hex contains only valid characters and has even length
	if len(rawHex)%2 != 0 {
		return nil, ErrTxHexOddLength
	}

	// Validate hex string can be decoded (early validation before passing to underlying library)
	// This catches invalid hex characters and prevents certain malformed inputs
	decoded, err := hex.DecodeString(rawHex)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}

	// Validate varints are within reasonable bounds to prevent OOM attacks
	// Check for varint prefixes that indicate extremely large values
	// 0xff prefix indicates an 8-byte varint follows
	if err = validateTransactionVarints(decoded); err != nil {
		return nil, err
	}

	// Recover from panics in the underlying library when parsing malformed transactions
	// This catches some "makeslice: len out of range" panics from malformed varint values
	// Limitation: Fatal OOM errors that call runtime.throw() cannot be recovered and may still
	// occur with certain malformed inputs. The validations above (size limits, hex validation,
	// input/output count checks) catch most problematic cases.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: %v", ErrTxParsePanic, r)
			tx = nil
		}
	}()

	return bt.NewTxFromString(rawHex)
}

// CreateTxWithChange will automatically create the change output and calculate fees
//
// Use this if you don't want to figure out fees/change for a tx
// USE AT YOUR OWN RISK - this will modify a "pay-to" output to accomplish auto-fees
func CreateTxWithChange(utxos []*Utxo, payToAddresses []*PayToAddress, opReturns []OpReturnData,
	changeAddress string, standardRate, dataRate *bt.Fee,
	privateKey *bec.PrivateKey,
) (*bt.Tx, error) {
	// Missing utxo(s) or change address
	if len(utxos) == 0 {
		return nil, ErrUtxosRequired
	} else if len(changeAddress) == 0 {
		return nil, ErrChangeAddressRequired
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
			"%w: need %d + (fee), found %d",
			ErrInsufficientFunds,
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
	if (totalPayToSatoshis + fee) <= totalSatoshis {
		// Remove the change address (old version with original satoshis)
		// Add the change address as the difference (now with adjusted fee)
		if hasChange {
			payToAddresses = payToAddresses[:len(payToAddresses)-1]

			payToAddresses = append(payToAddresses, &PayToAddress{
				Address:  changeAddress,
				Satoshis: totalSatoshis - (totalPayToSatoshis + fee),
			})
		}
		// Create the "Final tx" (or error)
		return CreateTx(utxos, payToAddresses, opReturns, privateKey)
	}

	// Not enough to cover the fee - need to adjust
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
			"%w: (payTo %d) (amount %d) (remainder %d) (fee %d) (total %d)",
			ErrAutoFeeNotApplicable,
			len(payToAddresses), totalPayToSatoshis, remainder, fee, totalSatoshis,
		)
	}

	// Create the "Final tx" (or error)
	return CreateTx(utxos, payToAddresses, opReturns, privateKey)
}

// draftTx is a helper method to create a draft tx and associated fees
func draftTx(utxos []*Utxo, payToAddresses []*PayToAddress, opReturns []OpReturnData,
	privateKey *bec.PrivateKey, standardRate, dataRate *bt.Fee,
) (uint64, error) {
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
	changeAddress string, standardRate, dataRate *bt.Fee, wif string,
) (*bt.Tx, error) {
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
// This will NOT create a change output (funds are sent to "addresses")
// This will NOT handle fee calculation (it's assumed you have already calculated the fee)
//
// Get the raw hex version: tx.ToString()
// Get the tx id: tx.GetTxID()
func CreateTx(utxos []*Utxo, addresses []*PayToAddress,
	opReturns []OpReturnData, privateKey *bec.PrivateKey,
) (*bt.Tx, error) {
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
		var a *bscript.Script
		a, err = bscript.NewP2PKHFromAddress(address.Address)
		if err != nil {
			return nil, err
		}

		if err = tx.PayTo(a, address.Satoshis); err != nil {
			return nil, err
		}
	}

	// Loop any op returns
	for _, op := range opReturns {
		if err = tx.AddOpReturnPartsOutput(op); err != nil {
			return nil, err
		}
	}

	// If inputs are supplied, make sure they are sufficient for this transaction
	if len(tx.Inputs) > 0 {
		// Sanity check - not enough satoshis in utxo(s) to cover all paid amount(s)
		// They should never be equal, since the fee is the spread between the two amounts
		totalOutputSatoshis := tx.TotalOutputSatoshis() // Does not work properly
		if totalOutputSatoshis > totalSatoshis {
			return nil, fmt.Errorf("%w: need %d + (fee), found %d", ErrInsufficientFunds, totalOutputSatoshis, totalSatoshis)
		}
	}

	// Sign the transaction
	if privateKey != nil {
		myAccount := &account{PrivateKey: privateKey}
		// todo: support context (ctx)
		if err = tx.FillAllInputs(context.Background(), myAccount); err != nil {
			return nil, err
		}
	}

	// Return the transaction as a raw string
	return tx, nil
}

// CreateTxUsingWif will create a basic transaction and return the raw transaction (*transaction.Transaction)
//
// This will NOT create a "change" address (it's assumed you have already specified an address)
// This will NOT handle "fee" calculation (it's assumed you have already calculated the fee)
//
// Get the raw hex version: tx.ToString()
// Get the tx id: tx.GetTxID()
func CreateTxUsingWif(utxos []*Utxo, addresses []*PayToAddress,
	opReturns []OpReturnData, wif string,
) (*bt.Tx, error) {
	// Decode the WIF
	privateKey, err := WifToPrivateKey(wif)
	if err != nil {
		return nil, err
	}

	// Create the Tx
	return CreateTx(utxos, addresses, opReturns, privateKey)
}

// DefaultStandardFee returns the default standard fees offered by most miners.
// this function is not public anymore in go-bt
func DefaultStandardFee() *bt.Fee {
	return &bt.Fee{
		FeeType: bt.FeeTypeStandard,
		MiningFee: bt.FeeUnit{
			Satoshis: 5,
			Bytes:    10,
		},
		RelayFee: bt.FeeUnit{
			Satoshis: 5,
			Bytes:    10,
		},
	}
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
		standardRate = DefaultStandardFee()
	}
	if dataRate == nil {
		dataRate = DefaultStandardFee()
		// todo: adjusted to 5/10 for now, since all miners accept that rate
		dataRate.FeeType = bt.FeeTypeData
	}

	// Set the total bytes of the tx
	totalBytes := len(tx.Bytes())

	// Loop all outputs and accumulate size (find data related outputs)
	for _, out := range tx.Outputs {
		outHexString := out.LockingScriptHexString()
		if strings.HasPrefix(outHexString, "006a") || strings.HasPrefix(outHexString, "6a") {
			totalDataBytes += len(out.Bytes())
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

	// Safety check for negative fee (should never happen in practice)
	if totalFee < 0 {
		totalFee = 1
	}

	// Return the total fee as an uint (easier to use with satoshi values)
	return uint64(totalFee) //nolint:gosec // totalFee is checked to be non-negative above
}
