package bitcoin

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
	"github.com/piotrnar/gocoin/lib/secp256k1"
	"golang.org/x/crypto/ripemd160"
)

func sha256d(body []byte) []byte {
	msghash1 := sha256.Sum256([]byte(body))
	msghash2 := sha256.Sum256(msghash1[:])
	return msghash2[:]
}

const (
	H_BSV string = "Bitcoin Signed Message:\n"
)

func VerifyMessage(address, signature, data string) (ok boolean) {
	addrs, err := sigmestoaddr(signature, data)
	if err != nil {
		return
	}
	for _, addr2 := range addrs {
		if address == addr2 {
			ok = true
			return
		}
	}
	return
}

func messagehash(message, header string) (msghash2 []byte, err error) {
	hlen := len(header)
	if hlen >= 0xfd {
		err = fmt.Errorf("long header is not supported")
		return
	}
	mlen := len(message)
	if mlen >= 0xfd {
		err = fmt.Errorf("long message is not supported")
		return
	}
	bitcoinmsg := string([]byte{byte(hlen)})
	bitcoinmsg += header
	bitcoinmsg += string([]byte{byte(mlen)})
	bitcoinmsg += message
	msghash2 = sha256d([]byte(bitcoinmsg))
	return
}

func parseSignature(signature string) (sig secp256k1.Signature, recid int,
	err error) {
	sigraw, err2 := base64.StdEncoding.DecodeString(signature)
	if err2 != nil {
		err = err2
		return
	}
	r0 := sigraw[0] - 27
	recid = int(r0 & 3)
	compressed := (r0 & 4) == 1
	if compressed {
		err = fmt.Errorf("compressed type is not supported")
		return
	}
	sig.R.SetBytes(sigraw[1 : 1+32])
	sig.S.SetBytes(sigraw[1+32 : 1+32+32])
	return
}

func pubtoaddr(pubkey_xy2 secp256k1.XY, compressed bool,
	magic []byte) (bcpy []byte) {
	size := 65
	if compressed {
		size = 33
	}
	out := make([]byte, size)
	pubkey_xy2.GetPublicKey(out)
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(out)
	pub_hash_1 := sha256_h.Sum(nil)
	ripemd160_h := ripemd160.New()
	ripemd160_h.Reset()
	ripemd160_h.Write(pub_hash_1)
	pub_hash_2 := ripemd160_h.Sum(nil)
	bcpy = append(magic, pub_hash_2...)
	hash2 := sha256d(bcpy)
	bcpy = append(bcpy, hash2[0:4]...)
	return
}

func addrtostr(bcpy []byte) (s string, err error) {
	z := new(big.Int)
	z.SetBytes(bcpy)
	enc := base58.BitcoinEncoding
	var encdd []byte
	encdd, err = enc.Encode([]byte(z.String()))
	if err != nil {
		return
	}
	s = string(encdd)
	for _, v := range bcpy {
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
func get_bin(num *secp256k1.Number, le int) []byte {
	bts := num.Bytes()
	if len(bts) > le {
		panic("buffer too small")
	}
	if len(bts) == le {
		return bts
	}
	return append(make([]byte, le-len(bts)), bts...)
}

// This function is copied from "piotrnar/gocoin/lib/secp256k1".
// And modified for local package.
// License is:
//   https://github.com/piotrnar/gocoin/blob/master/lib/secp256k1/COPYING
func recover(sig *secp256k1.Signature, pubkey *secp256k1.XY,
	m *secp256k1.Number, recid int) (ret bool) {
	var thecurve_p secp256k1.Number
	thecurve_p.SetBytes([]byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFC, 0x2F})
	var rx, rn, u1, u2 secp256k1.Number
	var fx secp256k1.Field
	var X secp256k1.XY
	var xj, qj secp256k1.XYZ

	rx.Set(&sig.R.Int)
	if (recid & 2) != 0 {
		rx.Add(&rx.Int, &secp256k1.TheCurve.Order.Int)
		if rx.Cmp(&thecurve_p.Int) >= 0 {
			return false
		}
	}

	fx.SetB32(get_bin(&rx, 32))

	X.SetXO(&fx, (recid&1) != 0)
	if !X.IsValid() {
		return false
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

	return true
}

func sigmestoaddr(signature, message string) (addrs []string, err error) {
	msghash2, err2 := messagehash(message, H_BSV)
	if err2 != nil {
		err = err2
		return
	}
	sig, recid, err2 := parseSignature(signature)
	if err2 != nil {
		err = err2
		return
	}
	var msg secp256k1.Number
	msg.SetBytes(msghash2)

	var pubkey_xy2 secp256k1.XY
	ret2 := recover(&sig, &pubkey_xy2, &msg, recid)
	if !ret2 {
		err = fmt.Errorf("recover pubkey failed")
		return
	}

	addrs = make([]string, 2)
	for i, compressed := range []bool{true, false} {
		bcpy := pubtoaddr(pubkey_xy2, compressed, []byte{byte(0)})
		s, err2 := addrtostr(bcpy)
		if err2 != nil {
			err = err2
			return
		}
		addrs[i] = s
	}
	return
}
