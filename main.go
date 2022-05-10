package main

import (
	"fmt"
	"blockchain/blockchain"
)

func main() {
	fmt.Println("Hello, playground")

	bc := blockchain.CreateBlockchain(5)
	bc.AddBlock("A", "B", 100)
	bc.AddBlock("B", "C", 200)

	for _, block := range bc.Chain {
		fmt.Println(block.Timestamp)
		fmt.Printf("Previous hash: %x\n", block.PreviousHash)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}
