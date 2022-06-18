package main

import(
	"time"
	"fmt"
)

func mine_once(fast bool) int64 {
	start := time.Now().Unix()
	tx := NewCoinbaseTX("test", "test")
	block := NewGenesisBlock(tx)
	pow = NewProofOfWork(block)
	if fast {
		
	} else {
		
	}
	end := time.Now().Unix()
	total := end - start

	return total
}