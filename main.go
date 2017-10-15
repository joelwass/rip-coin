package main

import (
	"fmt"
	"net/http"

	"github.com/nathanjohnson320/rip-coin/ipfs"
)

func main() {
	// // Walk through the entire process
	// // Create three wallets
	// kev := Wallet{}
	// nate := Wallet{}
	// pete := Wallet{}
	// kev.New()
	// nate.New()
	// pete.New()

	// // Kev said he rips nate
	// kevPriv, kevPub := kev.Encode()
	// rip := "Some generic rip."
	// transaction := Tx{}
	// transaction.Initiate(rip, kevPriv, kevPub)

	// // Pete and Nate verify and vote on the rip, asyncronously
	// ticker := time.NewTicker(1 * time.Second)
	// vc := make(chan Vote)
	// current := 0
	// limit := 3 // Timeout at 10 seconds
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:

	// 			// These ifs would not be in prod
	// 			if current == 0 {
	// 				// Pete votes
	// 				vote, err := pete.Vote(&transaction, true)
	// 				if err == nil {
	// 					vc <- *vote
	// 				}
	// 			}
	// 			if current == 2 {
	// 				// Nate votes
	// 				vote, err := nate.Vote(&transaction, true)
	// 				if err == nil {
	// 					vc <- *vote
	// 				}
	// 			}

	// 			// Increment the timer
	// 			current++
	// 		}
	// 	}
	// }()

	// for vote := range vc {
	// 	transaction.Votes = append(transaction.Votes, vote)

	// 	// Close the channel
	// 	if current >= limit {
	// 		close(vc)
	// 	}
	// }

	// // Complete the transaction
	// transaction.Complete()
	// fmt.Println("Complete")
	c := make(chan string)
	go func() {
		for {
			m := <-c
			fmt.Println(m)
		}
	}()
	ipfs.Subscribe(c)

	app := http.FileServer(http.Dir("./rip-coin/dist"))
	http.ListenAndServe(":6969", app)
}
