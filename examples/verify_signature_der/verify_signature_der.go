package main

import (
	"crypto/sha256"
	"log"

	"github.com/bitcoinschema/go-bitcoin"
)

func main() {

	// Example values (from Merchant API request)
	message := []byte(`{"apiVersion":"0.1.0","timestamp":"2020-10-08T14:25:31.539Z","expiryTime":"2020-10-08T14:35:31.539Z","minerId":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","currentHighestBlockHash":"0000000000000000021af4ee1f179a64e530bf818ef67acd09cae24a89124519","currentHighestBlockHeight":656007,"minerReputation":null,"fees":[{"id":1,"feeType":"standard","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}},{"id":2,"feeType":"data","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}}]}`)
	signature := "3045022100b976be863fffd361716b375a9a5c4e77073dfaa29d2b9af9addef94f029c2d0902205b1fffc58343f3d4bd8fc48a118e998072c655d318061e13e1ef0902fb42e15c"
	pubKey := "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"

	// Verify the signature
	if verified, err := bitcoin.VerifyMessageDER(sha256.Sum256(message), pubKey, signature); err != nil {
		log.Fatalf("verify failed: %s", err.Error())
	} else if !verified {
		log.Fatalf("verification failed")
	}
	log.Println("verification passed")
}
