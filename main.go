package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/BottleHub/Smart-Chain/internal/blockchain"
	"github.com/BottleHub/Smart-Chain/internal/proof"
)

type CLI struct {
	blockchain *blockchain.BlockChain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage")
	fmt.Println(" add -block BLOCK_DATA - Adds a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CLI) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block")
}

func (cli *CLI) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Block Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := proof.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CLI) run() {
	cli.validateArgs()

	addBlockCMD := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCMD := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCMD.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCMD.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "print":
		err := printChainCMD.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCMD.Parsed() {
		if *addBlockData == "" {
			addBlockCMD.Usage()
			runtime.Goexit()
		}

		cli.addBlock(*addBlockData)
	}

	if printChainCMD.Parsed() {
		cli.printChain()
	}
}

func main() {
	//router := gin.Default()
	chain := blockchain.Init()
	defer chain.DB.Close()

	cli := CLI{chain}
	defer cli.run()

	//router.GET("/")
	//router.POST("/")
	//router.POST("/create")
}
