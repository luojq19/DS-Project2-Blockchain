package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"time"
)

var (
	maxNonce = math.MaxInt64
)

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block  *Block
	target uint8
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	target := b.Difficulty

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(pow.block.Difficulty)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run performs a proof-of-work
func (pow *ProofOfWork) Run() (int, []byte) {
	var hash [32]byte
	nonce := 0

	begintime := time.Now().UnixNano()

	fmt.Printf("Mining a new block")
	for nonce < maxNonce {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)

		fmt.Printf("\r Current Try: %x", hash)

		if HasValidHash(hash, pow.block.Difficulty) {
			break
		} else {
			nonce++
		}
	}
	endtime := time.Now().UnixNano()

	total := endtime - begintime
	total = total / 1000000

	fmt.Printf("\rCurrent Try: %x, Total time: %d ms", hash, total)

	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate validates block's PoW
func (pow *ProofOfWork) Validate() bool {

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)

	return HasValidHash(hash, pow.block.Difficulty)
}

// Whether the hash is valdi undercurrent difficulty
func HasValidHash(hash [32]byte, diff uint8) bool {
	var temp1, temp2 uint8
	temp1 = (diff - diff%8) / 8
	temp2 = diff % 8
	var i uint8
	for i = 0; i < temp1; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	if temp2 != 0 {
		if hash[temp1]>>(8-temp2) != 0 {
			return false
		}
	}

	return true
}
