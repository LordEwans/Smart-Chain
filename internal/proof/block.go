package proof

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func handeHelper(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (b *Block) Serialise() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	handeHelper(err)

	return res.Bytes()
}

func Deserialise(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	handeHelper(err)

	return &block
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)

	nonce, hash := pow.Run()

	block.Nonce = nonce
	block.Hash = hash[:]

	return block
}
