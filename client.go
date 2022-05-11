package main

import (
	"log"
	"fmt"
	"strings"
	"strconv"
	"time"
	"io/ioutil"
)

func getAddress(name string) (string){
	byteadd, err := ioutil.ReadFile("address_"+name+".wal")
	if err != nil {
		log.Panic(err)
	}
	return string(byteadd)
}

func setAddress(name, address string){
	err := ioutil.WriteFile("address_"+name+".wal", []byte(address), 0644)
	if err != nil {
		panic(err)
	}
}

type CLI struct{
	nodeID string
	server bool
}

func NewCLI(ID string) (this *CLI){
	return &CLI{
		nodeID : ID,
		server : false,
	}
}

func (this *CLI) createBlockchain(name string) {
	address := getAddress(name)
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := CreateBlockchain(address, this.nodeID)
	defer bc.db.Close()

	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("Done!")
}

func (this *CLI) createWallet(name string) {
	wallets, _ := NewWallets(this.nodeID)
	address := wallets.CreateWallet()
	wallets.SaveToFile(this.nodeID)

	setAddress(name, address)

	fmt.Printf("Your new address: %s\n", address)
}

func (this *CLI) getBalance(name string) {
	address := getAddress(name)
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

	fmt.Printf("Balance of '%s': %d\n", name, balance)
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

func (this *CLI) send(fromName, toName string, amount int, mineNow bool) {
	from := getAddress(fromName)
	to := getAddress(toName)
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

func (this *CLI) startNode() {
	fmt.Printf("Starting node %s\n", this.nodeID)
	StartServer(this, "")
}

func (this *CLI) startMine(minerName string) {
	fmt.Printf("Starting node %s\n", this.nodeID)
	minerAddress := getAddress(minerName)
	if ValidateAddress(minerAddress) {
		fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
	} else {
		log.Panic("Wrong miner address!")
	}
	StartServer(this, minerAddress)
}

func (this *CLI) Atomic(text string) (){
	command := strings.SplitN(text, " ", 2)
	switch command[0]{
	case "**createblockchain":
		fmt.Println("Creating New Blockchain.")
		this.createBlockchain(command[1])
	case "**createwallet":
		fmt.Println("Creating A New wallet.")
		this.createWallet(command[1])
	case "**getbalance":
		fmt.Println("Getting balance.")
		this.getBalance(command[1])
	case "**listaddresses":
		fmt.Println("Listing existed addresses.")
		this.listAddresses()
	case "**printchain":
		fmt.Println("Printing the whole Blockchain.")
		this.printChain()
	case "**reindexUTXO":
		fmt.Println("Re-indexing UTXOs.")
		this.reindexUTXO()
	case "**send":
		fmt.Println("Performing a transaction.")
		ifmine := strings.HasSuffix(command[1], "mine")
		args := strings.Split(command[1], " ")
		amount, _ := strconv.Atoi(args[2])
		this.send(args[0], args[1], amount, ifmine)
	case "**startnode":
		args := strings.Split(command[1], " ")
		switch args[0]{
		case "normal":
			if this.server{
				for{}
			}else{
				this.startNode()
			}
		case "thread":
			if this.server == false{
				go this.startNode()
			}
		case "miner":
			this.startMine(args[1])
		}
	case "**sleep":
		period, _ := strconv.Atoi(command[1])
		fmt.Println("Sleeping...")
		time.Sleep(time.Millisecond * time.Duration(period))
		fmt.Println("Awake.")
	}
	return
}