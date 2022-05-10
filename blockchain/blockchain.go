package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Data         map[string]interface{}
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	Nonce          int
}

type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int
}

func (b Block) CalculateHash() string {
	data, _ := json.Marshal(b.Data)
	blockData := b.PreviousHash + string(data) + b.Timestamp.String() + strconv.Itoa(b.Nonce)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) Mine(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
			b.Nonce++
			b.Hash = b.CalculateHash()
	}
}

func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
			Hash:      "0",
			Timestamp: time.Now(),
	}
	return Blockchain{
			genesisBlock,
			[]Block{genesisBlock},
			difficulty,
	}
}

func (b *Blockchain) AddBlock(from, to string, amount float64) {
	blockData := map[string]interface{}{
			"from":   from,
			"to":     to,
			"amount": amount,
	}
	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
			Data:         blockData,
			PreviousHash: lastBlock.Hash,
			Timestamp:    time.Now(),
	}
	newBlock.Mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)
}

func (b Blockchain) CheckValid() bool {
	for i := range b.Chain[1:] {
			previousBlock := b.Chain[i]
			currentBlock := b.Chain[i+1]
			if currentBlock.Hash != currentBlock.CalculateHash() || currentBlock.PreviousHash != previousBlock.Hash {
					return false
			}
	}
	return true
}