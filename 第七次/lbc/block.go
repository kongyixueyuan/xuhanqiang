package lbc

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Height int64

	XHQ_Timestamp     int64
	Txs           []*Transaction
	XHQ_PrevBlockHash []byte
	Hash          []byte
	XHQ_Nonce         int
}

//新建一个区块(经过工作量证明验证的区块)，

func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block {

	block := &Block{height, time.Now().Unix(), txs, prevBlockHash, []byte{}, 0}

	pow := NewXHQ_ProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]

	block.XHQ_Nonce = nonce

	return block
}

//新建创世区块
func NewGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, 0, []byte{})
}

//对区块进行序列化。
func (b *Block) XHQ_Serialize() []byte {

	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//对区块进行反序列化。
func XHQ_DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

//
func XHQ_CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

/*
// 需要将Txs转换成[]byte
func (block *Block) XHQ_XHQ_HashTransactions() []byte  {

	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]

}
*/

//
func (block *Block) XHQ_XHQ_HashTransactions() []byte {

	var transactions [][]byte

	for _, tx := range block.Txs {
		transactions = append(transactions, tx.XHQ_Serialize())
	}
	mTree := NewXHQ_MerkleTree(transactions)

	return mTree.RootNode.Data

}
