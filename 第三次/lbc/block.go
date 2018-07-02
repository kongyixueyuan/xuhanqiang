package lbc

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

//新建一个区块(经过工作量证明验证的区块)，

func NewBlock(data string, prevBlockHash []byte) *Block {

	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]

	block.Nonce = nonce

	return block
}



	//新建创世区块
func NewGenesisBlock() *Block  {
	return NewBlock("创世区块!",[]byte{})
}



//对区块进行序列化。
func (b *Block) Serialize()[]byte  {

	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err!= nil{
		log.Panic(err)
	}
	return result.Bytes()
}

//对区块进行反序列化。
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}