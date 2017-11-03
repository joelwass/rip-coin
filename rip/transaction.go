package rip

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/oklog/ulid"
)

// Tx is the transaction
type Tx struct {
	ID              string             `json:"id"`
	PreviousHash    string             `json:"previousHash"`    // The last transaction hash
	Hash            string             `json:"hash"`            // The hash for block validation
	TotalAmount     int64              `json:"totalAmount"`     // All the RC generator for the transaction
	NumVotes        int64              `json:"numVotes"`        // How many people voted
	RipperPublicKey string             `json:"ripperPublicKey"` // Who initiated the transaction
	Rip             `json:"rip"`       // The context of the rip
	Timestamp       int64              `json:"timestamp"` // Current time that we can use for expiration
	Signature       `json:"signature"` // The signature of the transaction used to reject bogus transactions
}

// Initiate starts a transaction, verifies it
func (t *Tx) Initiate(rip, pub string, priv []byte) {
	t.TotalAmount = 1
	t.Rip.Rip = rip
	t.RipperPublicKey = pub
	t.ID = getULID().String()
	t.Sign(DecodePrivate(priv))
}

// Verify verifies that the incoming transaction is in fact valid
func (t *Tx) Verify() bool {
	// Get the public key
	// Base64 the trash because of golang's idiocy
	k, err := base64.StdEncoding.DecodeString(t.RipperPublicKey)
	if err != nil {
		fmt.Println("Error decoding public key")
		return false
	}
	key := DecodePublic(k)

	// Return the verification of said rip
	return ecdsa.Verify(key, []byte(t.Rip.Rip), t.Signature.R, t.Signature.S)
}

// Sign creates the signature
func (t *Tx) Sign(key *ecdsa.PrivateKey) {
	r, s, _ := ecdsa.Sign(rand.Reader, key, []byte(t.Rip.Rip))
	t.Signature = Signature{
		R: r,
		S: s,
	}
}

// Complete checks the votes and if they're all good seal
// the deal with a hash and send it off to IPFS
func (t *Tx) Complete(previousBlock Block) {
	yay := 0
	for _, vote := range t.Votes {
		if vote.Approval {
			yay++
		}
	}

	// Majority votes?
	if yay > (len(t.Votes) / 2) {
		AddToBlockchain(t, previousBlock)
	}
}

// AddToBlockchain sends the transaction off to IPFS
func AddToBlockchain(t *Tx, previousBlock Block) {
	t.Timestamp = time.Now().Unix()
	t.PreviousHash = previousBlock.Hash

	tBytes, _ := json.Marshal(t)
	t.Hash = fmt.Sprintf("%x", sha256.Sum256(tBytes))

	// Todo: blockchain lol
}

// Rip is what we use to generate a rip coin
type Rip struct {
	Rip   string `json:"rip"` // Why do you live so far away from home? Because I'm not afraid to go to work.
	Votes []Vote `json:"votes"`
}

// Vote is when someone actually votes on the current rip
type Vote struct {
	Address  []byte `json:"address"`  // The person voting
	Approval bool   `json:"approval"` // 0 or 1
}

// Signature is used for validating transactions
type Signature struct {
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

func getULID() ulid.ULID {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader)
}
