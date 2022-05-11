package main

import (
	"fmt"
	"strconv"
)

const difficulty = 16

func main() {
	bc := NewBlockchain(difficulty)
	bc.AddBlock("Send 1 BTC to Ivan", difficulty)
	bc.AddBlock("Send 2 more BTC to Ivan")
	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
