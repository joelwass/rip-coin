package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	// sum := sha256.Sum256([]byte("hello world\n"))
	// fmt.Println(sum)
	// t := time.Now().Unix()
	// fmt.Println(t)

	kev := CreateWallet()
	pete := CreateWallet()
	nate := CreateWallet()
	db := CreateWallet()

	rip := Rip{
		Rip: "Why do you live so far away from home? Because I'm not afraid to go to work.",
		Votes: []Vote{
			Vote{
				Address:  kev,
				Approval: true,
			},
			Vote{
				Address:  pete,
				Approval: true,
			},
			Vote{
				Address:  db,
				Approval: false,
			},
			Vote{
				Address:  nate,
				Approval: true,
			},
		},
	}

	ripBytes, _ := json.Marshal(rip)
	tx := Tx{
		Rip:       rip,
		Hash:      sha256.Sum256(ripBytes),
		Timestamp: time.Now().Unix(),
	}

	fmt.Println(tx)
}
