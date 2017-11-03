package rip

// Genesis is the IPFS hash of the genesis block
const Genesis = "QmZangbZcrTWgU4UYn1ACxBUorfbivFmSetBKjS7yMpnX5"

// Block is a block in the chain
type Block struct {
	Hash         string `json:"hash"`
	Transactions []Tx   `json:"transactions"`
}
