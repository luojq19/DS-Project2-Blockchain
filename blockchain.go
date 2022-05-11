package main

// blockchain is a chain list of blocks
type Blockchain struct {
	blocks []*Block
}

// create a new bock
func NewBlockchain(difficultyList ...int64) *Blockchain {
	var difficulty int64
	difficulty = 24
	for _, num := range difficultyList {
		difficulty = num
	}
	return &Blockchain{[]*Block{NewGenesisBlock(difficulty)}}
}

// add a new block to the end of blockchain
func (bc *Blockchain) AddBlock(data string, difficultyList ...int64) {
	var difficulty int64

	prevBlock := bc.blocks[len(bc.blocks)-1]

	difficulty = prevBlock.Difficulty
	for _, num := range difficultyList {
		difficulty = num
	}

	newBlock := NewBlock(data, prevBlock.Hash, difficulty)
	bc.blocks = append(bc.blocks, newBlock)
}
