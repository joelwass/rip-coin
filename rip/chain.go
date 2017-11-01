package rip

// Genesis is the IPFS hash of the genesis block
const Genesis = "QmbR1iGdnLWg3HaJtejJxw9TqPQkhquWYHUaStxHtezWih"

// Block is a block in the chain
type Block struct {
	Transactions []Tx `json:"transactions"`
}
