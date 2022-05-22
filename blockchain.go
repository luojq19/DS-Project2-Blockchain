package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const defultDifficulty = 20

//const knownDb = "blockchain.db"

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "This is the coinbase data for Genesis Block. "

type Blockchain struct {
	genesisHash []byte
	db          *bolt.DB
}

func getDifficulty(difficultyList ...int64) int64 {
	var difficulty int64
	difficulty = defultDifficulty
	for _, num := range difficultyList {
		difficulty = num
	}
	return difficulty
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// create a new blockchain and save it in the knownDb since we suppose everyone knows it
func CreateBlockchain(address, nodeID string, difficultyList ...int64) *Blockchain {

	difficulty := getDifficulty(difficultyList...)

	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) {
		//if dbExists(knownDb) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var genesisHash []byte

	cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := NewGenesisBlock(cbtx, difficulty)

	db, err := bolt.Open(dbFile, 0600, nil)
	printError(err)

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		printError(err)

		err = b.Put(genesis.Hash, genesis.Serialize())
		printError(err)

		err = b.Put([]byte("l"), genesis.Hash)
		printError(err)

		genesisHash = genesis.Hash

		return nil
	})
	printError(err)

	bc := Blockchain{genesisHash, db}
	return &bc
}

// create a new Blockchain with genesis Block
func NewBlockchain(nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if !dbExists(dbFile) {
		//if !dbExists(knownDb) {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var genesisHash []byte
	//db, err := bolt.Open(dbFile, 0600, nil)
	db, err := bolt.Open(dbFile, 0600, nil)
	printError(err)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		genesisHash = b.Get([]byte("l"))

		return nil
	})

	printError(err)
	bc := Blockchain{genesisHash, db}

	return &bc
}

// AddBlock saves the block into the blockchain
func (bc *Blockchain) AddBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockInDb := b.Get(block.Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := block.Serialize()
		err := b.Put(block.Hash, blockData)
		printError(err)

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height {
			err = b.Put([]byte("l"), block.Hash)
			printError(err)
			bc.genesisHash = block.Hash
		}

		return nil
	})
	printError(err)
}

// Find a transaction by its ID
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

// Find all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Get the next block starting from the genesisHash
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	printError(err)
	i.currentHash = block.PrevBlockHash

	return block
}

// create a blockchain iterator
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.genesisHash, bc.db}

	return bci
}

// get the height of the latest block
func (bc *Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

// find a block by its hash and returns it
func (bc *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		foundBlock := b.Get(blockHash)
		if foundBlock == nil {
			return errors.New("Block is not found in the blockchain. ")
		}

		block = *DeserializeBlock(foundBlock)
		return nil
	})

	if err != nil {
		return block, err
	}

	return block, nil
}

// get a list of hashes of all the blocks in the chain
func (bc *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()

	for {
		block := bci.Next()
		blocks = append(blocks, block.Hash)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

//mine a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction, diffcultyList ...int64) *Block {

	difficulty := getDifficulty(diffcultyList...)

	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		foundBlock := b.Get(lastHash)
		block := DeserializeBlock(foundBlock)

		lastHeight = block.Height

		return nil
	})
	printError(err)

	newBlock := NewBlock(transactions, lastHash, lastHeight+1, difficulty)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		printError(err)

		err = b.Put([]byte("l"), newBlock.Hash)
		printError(err)
		bc.genesisHash = newBlock.Hash

		return nil
	})
	printError(err)

	return newBlock
}

// sign the inputs of a Transaction
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		printError(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

// verify transaction input signatures
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		printError(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}
