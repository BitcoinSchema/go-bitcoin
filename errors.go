package bitcoin

import "errors"

// ErrPrivateKeyMissing is returned when a private key is missing
var ErrPrivateKeyMissing = errors.New("private key is missing")

// ErrWifMissing is returned when a wif is missing
var ErrWifMissing = errors.New("wif is missing")

// ErrBadCharacter is returned when a bad character is found
var ErrBadCharacter = errors.New("bad char")

// ErrTooLong is returned when a string is too long
var ErrTooLong = errors.New("too long")

// ErrNotVersion0 is returned when a string is not version 0
var ErrNotVersion0 = errors.New("not version 0")

// ErrMissingScript is returned when a script is missing
var ErrMissingScript = errors.New("missing script")

// ErrMissingPubKey is returned when a pubkey is missing
var ErrMissingPubKey = errors.New("missing pubkey")

// ErrMissingAddress is returned when an address is missing
var ErrMissingAddress = errors.New("missing address")

// ErrUtxosRequired is returned when utxos are required to create a tx
var ErrUtxosRequired = errors.New("utxo(s) are required to create a tx")

// ErrChangeAddressRequired is returned when a change address is required to create a tx
var ErrChangeAddressRequired = errors.New("change address is required")
