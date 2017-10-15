package main

import (
	"fmt"
	"net/http"
	"os"
	"os/user"

	"github.com/gorilla/websocket"
	"github.com/nathanjohnson320/rip-coin/rip"
	"github.com/nathanjohnson320/rip-coin/ws"
)

func init() {
	// Do checks for wallet
	fmt.Println("Loading wallet...")

	usr, _ := user.Current()
	dir := usr.HomeDir

	_, err := os.Open(dir + "/.rip/wallet.dat")
	if err != nil {
		fmt.Println("No wallet found, generating...")
		wally := rip.Wallet{}
		wally.New()

		wally.Save(dir + "/.rip/")
	} else {
		fmt.Println("Wallet loaded!")
	}
}

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
	// c := make(chan string)
	// go func() {
	// 	for {
	// 		m := <-c
	// 		fmt.Println(m)
	// 	}
	// }()
	// ipfs.Subscribe(c)

	// Websockets
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.Handle(w, r, upgrader)
	})

	// UI routes
	http.Handle("/", http.FileServer(http.Dir("./rip-coin/dist")))
	http.ListenAndServe(":6969", nil)
}
