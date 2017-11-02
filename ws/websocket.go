// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nathanjohnson320/rip-coin/ipfs"
	"github.com/nathanjohnson320/rip-coin/rip"
)

// Message to function map
var wsHandles = map[string]func(rip.Tx) []byte{
	"new_rip": func(transaction rip.Tx) []byte {
		w, err := rip.LoadWallet()
		if err != nil {
			fmt.Println("Could not load wallet.")
			return nil
		}

		transaction.Initiate(string(transaction.Rip.Rip), w.Priv, w.Pub)
		tB, err := json.Marshal(transaction)
		if err != nil {
			fmt.Println("Could not marshal transaction.")
			return nil
		}

		// Publish the transaction
		ipfs.Publish("rip-coin-tx", string(tB))
		return nil
	},
	"new_block": func(transaction rip.Tx) []byte {
		tB, err := json.Marshal(transaction)
		if err != nil {
			fmt.Println("Could not marshal transaction.")
			return nil
		}
		ipfs.Publish("rip-coin-block", string(tB))
		return nil
	},
	"vote": func(transaction rip.Tx) []byte {
		tB, err := json.Marshal(transaction)
		if err != nil {
			fmt.Println("Could not marshal transaction.")
			return nil
		}
		ipfs.Publish("rip-coin-vote", string(tB))
		return nil
	},
}

// Payload is the inbound websocket payload format
type Payload struct {
	Type string `json:"type"`
	Data rip.Tx `json:"data"`
}

// Handle responds to websocket requests
func Handle(w http.ResponseWriter, r *http.Request, upgrader websocket.Upgrader) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Subscribe to IPFS transactions
	go func() {
		c := make(chan string)
		ipfs.Subscribe(c)
		for {
			m := <-c
			fmt.Println(m)

			tB, err := json.Marshal(Message{Label: "new_tx", Payload: m})
			if err != nil {
				fmt.Println(err)
				return
			}

			err = conn.WriteMessage(1, tB)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}()

	// Send up the user's public key
	sendKeys(conn)

	// Handle incoming websocket messages
	go func() {
		for {
			payload := new(Payload)
			err := conn.ReadJSON(&payload)
			if err != nil {
				fmt.Println("Error decoding inbound message")
				fmt.Println(err)
				return
			}

			// Get the response to send back
			response := wsHandles[payload.Type](payload.Data)

			if response != nil {
				conn.WriteMessage(1, response)
			}
		}
	}()
}

func sendKeys(conn *websocket.Conn) {
	// Grab the wallet
	w, err := rip.LoadWallet()
	if err != nil {
		fmt.Println(err)
		return
	}

	j, err := json.Marshal(w)
	if err != nil {
		fmt.Println(err)
		return
	}

	m := Message{
		Label:   "pub_key",
		Payload: string(j),
	}
	j, err = json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = conn.WriteMessage(1, j)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Message is the generic struct we send up to the UI
type Message struct {
	Label   string `json:"label"`
	Payload string `json:"payload"`
}
