package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"
)

func mine_once(flag bool) int64 {
	start := time.Now().UnixNano()
	from := "aliceasdfasdfasdf"
	to := "boaasdfasfasdfasdfb"
	tx := NewCoinbaseTX(from, to)
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

// func main() {
// 	total1 := 0
// 	for i := 0; i < 10; i++ {
// 		total1 += int(mine_once(true))
// 	}
// 	total2 := 0
// 	for i := 0; i < 10; i++ {
// 		total2 += int(mine_once(false))
// 	}
// 	fmt.Printf("Slow mine 10 blocks: average time: %d ms\n", total1/1000000/10)
// 	fmt.Printf("Fast mine 10 blocks: average time: %d ms\n", total2/1000000/10)
// }