package challenger

import (
	"errors"
	"regexp"
	"time"
)

const ChallengeExpirationDelta = 5 * time.Minute

type Challenge struct {
	Value string
	Exp   time.Time
}

var AddressRegexp = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
var SignatureRegexp = regexp.MustCompile("^0x[0-9a-fA-F]{130}$")
var SHA256HexRegexp = regexp.MustCompile("^0x[0-9a-fA-F]{64}$")

var (
	ErrChallengeWasNotRequested = errors.New("challenge was not requested")
	ErrChallengeExpired         = errors.New("challenge expired")
	ErrDecodeHex                = errors.New("failed to decode hex string")
	ErrBadLength                = errors.New("bad signature length")
	ErrBadRecoverByte           = errors.New("bad recovery byte")
	ErrMissMatched              = errors.New("recovered address didn't match any of the given ones")
)
