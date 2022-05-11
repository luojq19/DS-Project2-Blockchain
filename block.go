package main

import (
	"time"
)

// Block keeps block headers
type Block struct {
	Timestamp     int64
	Difficulty    int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewBlock creates and returns Block
func NewBlock(data string, prevBlockHash []byte, difficulty int64) *Block {
	block := &Block{time.Now().Unix(), difficulty, []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(difficulty int64) *Block {
	return NewBlock("Genesis Block", []byte{}, difficulty)
}
