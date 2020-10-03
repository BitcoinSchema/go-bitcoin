package bitcoin

import (
	"errors"

	"github.com/bitcoinsv/bsvutil"
	"github.com/libsv/libsv/transaction"
	"github.com/libsv/libsv/transaction/output"
	"github.com/libsv/libsv/transaction/signature"
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
