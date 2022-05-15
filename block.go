package main

import (
	"bytes"
	"encoding/gob"
	"time"
)

// Block keeps block headers
type Block struct {
	Timestamp     int64
	Difficulty    int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Height        int
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int, difficulty int64) *Block {
	block := &Block{time.Now().Unix(), difficulty, transactions, prevBlockHash, []byte{}, 0, height}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// create and returns genesis Block
func NewGenesisBlock(Coinbase *Transaction, difficulty int64) *Block {
	return NewBlock([]*Transaction{Coinbase}, []byte{}, 0, difficulty)
}

// compute the hash of the transactions and return the root of the Merkle tree
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	temp := NewMerkleTree(transactions)

	return temp.RootNode.Data
}

// serialize the block for data transforming
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	printError(err)

	return result.Bytes()
}

// deserilize the block, as a reverse of serilization
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	printError(err)

	return &block
}
