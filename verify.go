package bitcoin

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
	"github.com/piotrnar/gocoin/lib/secp256k1"
	"golang.org/x/crypto/ripemd160"
)

const (
	// H_BSV is the magic header string required fore Bitcoin Signed Messages
	hBSV string = "Bitcoin Signed Message:\n"
)

// VerifyMessage verifies a string and address against the provided signature and assumes Bitcoin Signed Message encoding
func VerifyMessage(address, signature, data string) error {
	if len(address) == 0 {
		return errors.New("address is missing")
	} else if len(signature) == 0 {
		return errors.New("signature is missing")
	}
	addresses, err := sigMessageToAddress(signature, data)
	if err != nil {
		return err
	}
	for _, addr2 := range addresses {
		if address == addr2 {
			return nil
		}
	}
	return fmt.Errorf("address: %s not found", address)
}

// sha256d is a double sha256
func sha256d(body []byte) []byte {
	msgHash1 := sha256.Sum256(body)
	msgHash2 := sha256.Sum256(msgHash1[:])
	return msgHash2[:]
}

// messageHash will compute a hash for the given message & header
func messageHash(message, header string) ([]byte, error) {
	headerLength := len(header)
	if headerLength >= 0xfd {
		return nil, fmt.Errorf("long header is not supported")
	}
	messageLength := len(message)
	if messageLength >= 0xfd {
		return nil, fmt.Errorf("long message is not supported")
	}
	bitcoinMsg := string([]byte{byte(headerLength)})
	bitcoinMsg += header
	bitcoinMsg += string([]byte{byte(messageLength)})
	bitcoinMsg += message
	return sha256d([]byte(bitcoinMsg)), nil
}

// parseSignature will parse the given signature
func parseSignature(signature string) (sig secp256k1.Signature, recID int, err error) {
	// todo: is this 64 or 65? is it always the same?
	// panic was occurring: sig.R.SetBytes(sigRaw[1 : 1+32])
	if len(signature) < 64 {
		err = fmt.Errorf("signature is less than %d characters", 64)
		return
	}
	var sigRaw []byte
	if sigRaw, err = base64.StdEncoding.DecodeString(signature); err != nil {
		return
	}
	r0 := sigRaw[0] - 27
	recID = int(r0 & 3)
	compressed := (r0 & 4) == 1
	if compressed {
		err = fmt.Errorf("compressed type is not supported")
		return
	}
	sig.R.SetBytes(sigRaw[1 : 1+32])
	sig.S.SetBytes(sigRaw[1+32 : 1+32+32])
	return
}

// pubKeyToAddress will convert a pubkey to an address
func pubKeyToAddress(pubkeyXy2 secp256k1.XY, compressed bool, magic []byte) (byteCopy []byte) {
	size := 65
	if compressed {
		size = 33
	}
	out := make([]byte, size)
	pubkeyXy2.GetPublicKey(out)
	sha256H := sha256.New()
	sha256H.Reset()
	sha256H.Write(out)
	pubHash1 := sha256H.Sum(nil)
	ripemd160H := ripemd160.New()
	ripemd160H.Reset()
	ripemd160H.Write(pubHash1)
	pubHash2 := ripemd160H.Sum(nil)
	byteCopy = append(magic, pubHash2...)
	hash2 := sha256d(byteCopy)
	byteCopy = append(byteCopy, hash2[0:4]...)
	return
}

// addressToString will convert a raw address to a string version
func addressToString(byteCopy []byte) (s string, err error) {
	z := new(big.Int)
	z.SetBytes(byteCopy)
	enc := base58.BitcoinEncoding
	var encodeResults []byte
	if encodeResults, err = enc.Encode([]byte(z.String())); err != nil {
		return
	}
	s = string(encodeResults)
	for _, v := range byteCopy {
		if v != 0 {
			break
		}
		s = "1" + s
	}
	return
}

// This function is copied from "piotrnar/gocoin/lib/secp256k1".
// And modified for local package.
// License is:
//   https://github.com/piotrnar/gocoin/blob/master/lib/secp256k1/COPYING
func getBin(num *secp256k1.Number, le int) ([]byte, error) {
	bts := num.Bytes()
	if len(bts) > le {
		return nil, fmt.Errorf("buffer too small")
	}
	if len(bts) == le {
		return bts, nil
	}
	return append(make([]byte, le-len(bts)), bts...), nil
}

// This function is copied from "piotrnar/gocoin/lib/secp256k1".
// And modified for local package.
// License is:
//   https://github.com/piotrnar/gocoin/blob/master/lib/secp256k1/COPYING
func recoverSig(sig *secp256k1.Signature, pubkey *secp256k1.XY, m *secp256k1.Number, recID int) (bool, error) {
	var theCurveP secp256k1.Number
	theCurveP.SetBytes([]byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFC, 0x2F})
	var rx, rn, u1, u2 secp256k1.Number
	var fx secp256k1.Field
	var X secp256k1.XY
	var xj, qj secp256k1.XYZ

	rx.Set(&sig.R.Int)
	if (recID & 2) != 0 {
		rx.Add(&rx.Int, &secp256k1.TheCurve.Order.Int)
		if rx.Cmp(&theCurveP.Int) >= 0 {
			return false, nil // todo: is this actually an error?
		}
	}

	bin, err := getBin(&rx, 32)
	if err != nil {
		return false, err
	}

	fx.SetB32(bin)

	X.SetXO(&fx, (recID&1) != 0)
	if !X.IsValid() {
		return false, nil // todo: is this actually an error?
	}

	xj.SetXY(&X)
	rn.ModInverse(&sig.R.Int, &secp256k1.TheCurve.Order.Int)

	u1.Mul(&rn.Int, &m.Int)
	u1.Mod(&u1.Int, &secp256k1.TheCurve.Order.Int)

	u1.Sub(&secp256k1.TheCurve.Order.Int, &u1.Int)

	u2.Mul(&rn.Int, &sig.S.Int)
	u2.Mod(&u2.Int, &secp256k1.TheCurve.Order.Int)

	xj.ECmult(&qj, &u2, &u1)
	pubkey.SetXYZ(&qj)

	return true, nil
}

// sigMessageToAddress will convert a signature & message to a list of addresses
func sigMessageToAddress(signature, message string) ([]string, error) {

	// Get message hash
	msgHash, err := messageHash(message, hBSV)
	if err != nil {
		return nil, err
	}

	// Parse the signature
	var sig secp256k1.Signature
	var recID int
	sig, recID, err = parseSignature(signature)
	if err != nil {
		return nil, err
	}
	var msg secp256k1.Number
	msg.SetBytes(msgHash)

	var pubkeyXy2 secp256k1.XY
	var ret bool
	ret, err = recoverSig(&sig, &pubkeyXy2, &msg, recID)
	if err != nil {
		return nil, err
	} else if !ret {
		return nil, fmt.Errorf("recover pubkey failed")
	}

	addresses := make([]string, 2)
	for i, compressed := range []bool{true, false} {
		byteCopy := pubKeyToAddress(pubkeyXy2, compressed, []byte{byte(0)})

		var addressString string
		addressString, err = addressToString(byteCopy)
		if err != nil {
			return nil, err
		}
		addresses[i] = addressString
	}
	return addresses, nil
}
