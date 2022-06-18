package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"
)

func mine_once(flag bool) int64 {
	start := time.Now().UnixNano()
	tx := NewCoinbaseTX("test", "test")
	block := NewGenesisBlock(tx, 16)
	pow := NewProofOfWork(block)

	if flag {
		var hash [32]byte
		nonce := 0

		fmt.Printf("Mining a new block")
		for nonce < maxNonce {
			data := pow.prepareData(nonce)

			hash = sha256.Sum256(data)

			fmt.Printf("\r Current Try: %x", hash)

			if HasValidHash(hash, 16) {
				break
			} else {
				nonce++
			}
		}

		fmt.Print("\n\n")
	} else {
		target := big.NewInt(1)
		target = target.Lsh(target, uint(256-16))

		var hashInt big.Int
		var hash [32]byte
		nonce := 0

		fmt.Printf("Mining a new block")

		for nonce < maxNonce {
			data := pow.prepareData(nonce)

			hash = sha256.Sum256(data)
			fmt.Printf("\r Current Try: %x", hash)
			hashInt.SetBytes(hash[:])

			if hashInt.Cmp(target) == -1 {
				break
			} else {
				nonce++
			}
		}
		fmt.Print("\n\n")
	}

	end := time.Now().UnixNano()

	total := end - start
	return total
}
