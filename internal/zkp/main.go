package zkp

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"sync"
	"time"

	zkptypes "github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/verifier"
	"github.com/rarimo/web3-auth-svc/pkg/circuit"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var verificationKey []byte

func init() {
	var err error
	if verificationKey, err = circuit.VerificationKey.ReadFile(circuit.VerificationKeyFileName); err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to parse: %s", circuit.VerificationKeyFileName)))
	}
}

type Verifier struct {
	mu sync.Mutex

	Enabled bool

	// Map for storing challenges to be verified in auth proofs. No need to store in db - very short-live data.
	challenges map[string]*Challenge
}

func (v *Verifier) Challenge(user string) (string, error) {
	challenge := make([]byte, 31)
	if _, err := rand.Read(challenge); err != nil {
		return "", err
	}

	challengeStr := base64.StdEncoding.EncodeToString(challenge)

	v.mu.Lock()
	defer v.mu.Unlock()

	v.challenges[user] = &Challenge{
		Value:    challengeStr,
		Exp:      time.Now().UTC().Add(ChallengeExpirationDelta),
		Verified: false,
	}

	return challengeStr, nil
}

// VerifyProof performs ZK Groth16 proof verification based on specified verification key and hardcoded/passed parameters.
func (v *Verifier) VerifyProof(user string, proof *zkptypes.ZKProof) (err error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	challenge, ok := v.challenges[user]
	if !ok {
		return ErrChallengeWasNotRequested
	}

	if challenge.Verified || challenge.Exp.Before(time.Now().UTC()) {
		return ErrChallengeIsInvalid
	}

	if !v.Enabled {
		return nil
	}

	// no error can appear
	chal, _ := base64.StdEncoding.DecodeString(challenge.Value)
	chalDec := new(big.Int).SetBytes(chal).String()

	switch {
	case proof.PubSignals[NullifierSignalsIndex] != user:
		return fmt.Errorf("expected user=%s, got %s", user, proof.PubSignals[NullifierSignalsIndex])
	case proof.PubSignals[EventIDSignalsIndex] != EventID:
		return fmt.Errorf("expected eventID=%s, got %s", EventID, proof.PubSignals[EventIDSignalsIndex])
	case proof.PubSignals[EventDataSignalsIndex] != chalDec:
		return fmt.Errorf("expected challenge=%s, got %s", chalDec, proof.PubSignals[EventDataSignalsIndex])
	}

	if err = verifier.VerifyGroth16(*proof, verificationKey); err != nil {
		return errors.Wrap(err, "failed to verify generated proof")
	}

	challenge.Verified = true

	return nil
}
