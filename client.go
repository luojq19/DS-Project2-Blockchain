package main

import (
	"log"
	"fmt"
	"strings"
	"strconv"
	"time"
)

type CLI struct{
	nodeID string
}

func NewCLI(ID string) (this *CLI){
	return &CLI{
		nodeID : ID,
	}
}

func (this *CLI) createBlockchain(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := CreateBlockchain(address, this.nodeID)
	defer bc.db.Close()

	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("Done!")
}

func (this *CLI) createWallet() {
	wallets, _ := NewWallets(this.nodeID)
	address := wallets.CreateWallet()
	wallets.SaveToFile(this.nodeID)

	fmt.Printf("Your new address: %s\n", address)
}

func (this *CLI) getBalance(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := NewBlockchain(this.nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (this *CLI) listAddresses() {
	wallets, err := NewWallets(this.nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (this *CLI) printChain() {
	bc := NewBlockchain(this.nodeID)
	defer bc.db.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (this *CLI) reindexUTXO() {
	bc := NewBlockchain(this.nodeID)
	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}

func (this *CLI) send(from, to string, amount int, mineNow bool) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockchain(this.nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	wallets, err := NewWallets(this.nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := NewUTXOTransaction(&wallet, to, amount, &UTXOSet)

	if mineNow {
		cbTx := NewCoinbaseTX(from, "")
		txs := []*Transaction{cbTx, tx}

		newBlock := bc.MineBlock(txs)
		UTXOSet.Update(newBlock)
	} else {
		sendTx(knownNodes[0], tx)
	}

	fmt.Println("Success!")
}

func (this *CLI) startNode(minerAddress string) {
	fmt.Printf("Starting node %s\n", this.nodeID)
	if len(minerAddress) > 0 {
		if ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	StartServer(this.nodeID, minerAddress)
}

func (this *CLI) Atomic(text string) (){
	command := strings.SplitN(text, " ", 2)
	switch command[0]{
	case "**createblockchain":
		this.createBlockchain(command[1])
	case "**createwallet":
		this.createWallet()
	case "**getbalance":
		this.getBalance(command[1])
	case "**listaddresses":
		this.listAddresses()
	case "**printchain":
		this.printChain()
	case "**reindexUTXO":
		this.reindexUTXO()
	case "**send":
		ifmine := strings.HasSuffix(command[1], "mine")
		args := strings.Split(command[1], " ")
		amount, _ := strconv.Atoi(args[2])
		this.send(args[0], args[1], amount, ifmine)
	case "**startnode":
		if strings.HasPrefix(command[1], "miner"){
			args := strings.Split(command[1], " ")
			this.startNode(args[1])
		}else{
			this.startNode("")
		}
	case "**sleep":
		period, _ := strconv.Atoi(command[1])
		fmt.Println("Sleeping...")
		time.Sleep(time.Millisecond * time.Duration(period))
		fmt.Println("Awake.")
	}
	return
}