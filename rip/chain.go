package rip

// Block is a block in the chain
type Block struct {
	Transactions []Tx `json:"transactions"`
}
