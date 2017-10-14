package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"math/big"
	"time"
)

// Tx is the transaction
type Tx struct {
	PreviousHash    [32]byte // The last transaction hash
	Hash            [32]byte // The hash for block validation
	TotalAmount     int64    // All the RC generator for the transaction
	NumVotes        int64    // How many people voted
	RipperPublicKey []byte   // Who initiated the transaction
	Rip                      // The context of the rip
	Timestamp       int64    // Current time that we can use for expiration
	Signature                // The signature of the transaction used to reject bogus transactions
}

// Initiate starts a transaction, verifies it
func (t *Tx) Initiate(rip string, priv, pub []byte) {
	t.TotalAmount = 1
	t.Rip.Rip = rip
	t.RipperPublicKey = pub
	t.Sign(DecodePrivate(priv))
}

// Verify verifies that the incoming transaction is in fact valid
func (t *Tx) Verify() bool {
	// Get the public key
	key := DecodePublic(t.RipperPublicKey)

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
func (t *Tx) Complete() {
	yay := 0
	for _, vote := range t.Votes {
		if vote.Approval {
			yay++
		}
	}

	// Majority votes?
	if yay > (len(t.Votes) / 2) {
		AddToBlockchain(t)
	}
}

// AddToBlockchain sends the transaction off to IPFS
func AddToBlockchain(t *Tx) {
	t.Timestamp = time.Now().Unix()

	tBytes, _ := json.Marshal(t)
	t.Hash = sha256.Sum256(tBytes)

	// Todo: blockchain lol
}

// Rip is what we use to generate a rip coin
type Rip struct {
	Rip   string // Why do you live so far away from home? Because I'm not afraid to go to work.
	Votes []Vote
}

// Vote is when someone actually votes on the current rip
type Vote struct {
	Address  []byte // The person voting
	Approval bool   // 0 or 1
}

// Signature is used for validating transactions
type Signature struct {
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}
