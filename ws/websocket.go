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
		ipfs.Publish(string(tB))
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
	c := make(chan string)
	go func() {
		for {
			m := <-c
			err = conn.WriteMessage(1, []byte(m))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()
	ipfs.Subscribe(c)

	// Handle incoming websocket messages
	go func() {
		for {
			payload := new(Payload)
			err := conn.ReadJSON(&payload)
			if err != nil {
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
