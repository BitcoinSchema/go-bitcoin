// Package main demonstrates how to calculate the fee for a Bitcoin transaction.
package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {
	// Get the tx from hex string
	rawTx := "0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100e07b7661af4e4b521c012a146b25da2c7b9d606e9ceaae28fa73eb347ef6da6f0220527f0638a89ff11cbe53d5f8c4c2962484a370dcd9463a6330f45d31247c2512412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff0364030000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000"
	tx, err := bitcoin.TxFromHex(rawTx)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Calculate the fee using default rates (you can replace with MinerAPI rates)
	estimatedFee := bitcoin.CalculateFeeForTx(tx, nil, nil)

	// Success!
	log.Printf("tx id: %s estimated fee: %d satoshis", tx.TxID(), estimatedFee)
}
