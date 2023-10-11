package blockchain

import (
	"fmt"
	"log"

	"github.com/BottleHub/Smart-Chain/internal/proof"
	"github.com/dgraph-io/badger"
)

const (
	path = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	DB       *badger.DB
}

type BlockChainIterator struct {
	Hash []byte
	DB   *badger.DB
}

func (b *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := b.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			log.Panic(err)
		}

		err = item.Value(func(val []byte) error {
			lastHash = val
			return err
		})

		return err
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := proof.CreateBlock(data, lastHash)

	err = b.DB.Update(func(txn *badger.Txn) error {
		err = txn.Set(newBlock.Hash, newBlock.Serialise())
		if err != nil {
			log.Panic(err)
		}

		err = txn.Set([]byte("lh"), newBlock.Hash)

		b.LastHash = newBlock.Hash

		return err
	})
}

func (b *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{b.LastHash, b.DB}

	return iter
}

func (iter *BlockChainIterator) Next() *proof.Block {
	var block *proof.Block

	err := iter.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.Hash)
		if err != nil {
			log.Panic(err)
		}

		err = item.Value(func(val []byte) error {
			block = proof.Deserialise(val)
			return err
		})

		return err
	})
	if err != nil {
		log.Panic(err)
	}

	iter.Hash = block.PrevHash

	return block
}

func Init() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(path)
	opts.Dir = path
	opts.ValueDir = path

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")

			genesis := proof.Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialise())
			if err != nil {
				log.Panic(err)
			}

			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			if err != nil {
				log.Panic(err)
			}
			err = item.Value(func(val []byte) error {
				lastHash = val
				return err
			})

			return err
		}
	})
	if err != nil {
		log.Panic(err)
	}

	blockChain := BlockChain{lastHash, db}
	return &blockChain
}
