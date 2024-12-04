package challenger

import (
	"crypto/rand"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type AuthVerifier struct {
	// Map for storing challenges to be signed by user. No need to store in db - very short-live data.
	challenges map[string]*Challenge
	Disabled   bool

	mu sync.Mutex
}

func (v *AuthVerifier) Challenge(user string) (string, error) {
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		return "", err
	}

	challengeStr := hexutil.Encode(challenge)

	v.mu.Lock()
	defer v.mu.Unlock()

	v.challenges[user] = &Challenge{
		Value: challengeStr,
		Exp:   time.Now().UTC().Add(ChallengeExpirationDelta),
	}

	return challengeStr, nil
}

func ChallengeToHash(challenge string) []byte {
	message := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(challenge), challenge)
	hash := crypto.Keccak256Hash([]byte(message)).Bytes()
	return hash
}

func DecodeSignature(signature string) ([]byte, error) {
	signatureBytes, err := hexutil.Decode(signature)
	if err != nil {
		return nil, ErrDecodeHex
	}
	if len(signatureBytes) != 65 {
		return nil, ErrBadLength
	}
	if signatureBytes[64] == 0 || signatureBytes[64] == 1 {
		signatureBytes[64] = signatureBytes[64] + 27
	}
	if signatureBytes[64] != 27 && signatureBytes[64] != 28 {
		return nil, ErrBadRecoverByte
	}
	signatureBytes[64] -= 27

	return signatureBytes, nil
}

func (v *AuthVerifier) VerifySignature(signature string, address string) error {
	if v.Disabled {
		return nil
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	defer func() {
		delete(v.challenges, address)
	}()

	challenge, ok := v.challenges[address]
	if !ok {
		return ErrChallengeWasNotRequested
	}

	if challenge.Exp.Before(time.Now().UTC()) {
		return ErrChallengeExpired
	}

	signatureBytes, err := DecodeSignature(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	recoveredPubkey, err := crypto.SigToPub(ChallengeToHash(challenge.Value), signatureBytes)
	if err != nil {
		return fmt.Errorf("failed to recover pubkey from signed message: %w", err)
	}

	matched := false
	if strings.EqualFold(address, crypto.PubkeyToAddress(*recoveredPubkey).Hex()) {
		matched = true
	}

	if !matched {
		return ErrMissMatched
	}

	return nil
}
