package main

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid"
)

// Tx is the transaction
type Tx struct {
	Hash      [32]byte // The hash for validation
	Rip                // The context of the rip
	Timestamp int64    // Current time that we can use for expiration
}

// Rip is what we use to generate a rip coin
type Rip struct {
	Rip   string // Why do you live so far away from home? Because I'm not afraid to go to work.
	Votes []Vote
}

// Vote is when someone actually votes on the current rip
type Vote struct {
	Address  string // The person voting
	Approval bool   // 0 or 1
}

func CreateWallet() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
}
