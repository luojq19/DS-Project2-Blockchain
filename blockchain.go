package main

// blockchain is a chain list of blocks
type Blockchain struct {
	blocks []*Block
}

// create a new bock
func NewBlockchain(difficulty int64) *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock(difficulty)}}
}

// add a new block to the end of blockchain
func (bc *Blockchain) AddBlock(data string, difficulty int64) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash, difficulty)
	bc.blocks = append(bc.blocks, newBlock)
}
