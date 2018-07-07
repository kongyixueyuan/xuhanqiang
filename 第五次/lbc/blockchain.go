package lbc

import (
	"fmt"
	"log"

	"encoding/hex"
	"math/big"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
	"time"
)

// 数据库名字
const dbName = "blockchain.db"

// 表的名字
const blockTableName = "blocks"

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

// 区块链的结构体，
type Blockchain struct {
	Tip []byte
	Db  *bolt.DB
}

// 区块链的迭代器结构，
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// 添加一个区块(带有交易信息的数据)到区块链内，
func (bc *Blockchain) AddBlock(txs []*Transaction) {
	var lastHash []byte


	/*
	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		// ⚠️，先获取最新区块
		blockBytes := b.Get(blc.Tip)
		// 反序列化
		block := DeserializeBlock(blockBytes)

		newBlock := NewBlock(txs, block.Height,lastHash)


		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	*/



	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		//
		lastHash = b.Get([]byte("l"))

		// ⚠️，先获取最新区块
		blockBytes := b.Get(bc.Tip)
		// 反序列化
		block := DeserializeBlock(blockBytes)

		newBlock := NewBlock(txs, block.Height+1,lastHash)

		//

		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.Tip = newBlock.Hash

		return nil
	})


	if err != nil {
		log.Panic(err)
	}


}

// 区块链的迭代器函数。
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.Tip, bc.Db}

	return bci
}

// 区块链的迭代器的下一个区块函数。
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

// 创建一个带有创世区块的区块链。
func NewBlockchain(address string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("没有现成的区块链，创建一个新的...")
			//genesis := NewGenesisBlock()
			txCoinbase := NewCoinbaseTransaction(address)
			genesis := CreateGenesisBlock([]*Transaction{txCoinbase})

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// 判断数据库是否存在
func DBExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}

///
//1. 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(address string) *Blockchain {

	// 判断数据库是否存在
	if DBExists() {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块.......")

	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var genesisHash []byte

	// 关闭数据库
	err = db.Update(func(tx *bolt.Tx) error {

		// 创建数据库表
		b, err := tx.CreateBucket([]byte(blockTableName))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块
			// 创建了一个coinbase Transaction
			txCoinbase := NewCoinbaseTransaction(address)

			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			// 将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// 存储最新的区块的hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			genesisHash = genesisBlock.Hash
		}

		return nil
	})

	return &Blockchain{genesisHash, db}

}

// 返回Blockchain对象
func BlockchainObject() *Blockchain {

	db, err := bolt.Open(dbName, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte("l"))

		}

		return nil
	})

	return &Blockchain{tip, db}
}

// 如果一个地址对应的TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (blockchain *Blockchain) UnUTXOs(address string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _, tx := range txs {

		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//是否能够解锁
				if in.UnLockWithAddress(address) {

					key := hex.EncodeToString(in.TxHash)

					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}

			}
		}
	}

	for _, tx := range txs {

	Work1:
		for index, out := range tx.Vouts {

			if out.UnLockScriptPubKeyWithAddress(address) {
				fmt.Println("查看...")
				fmt.Println(address)

				fmt.Println(spentTXOutputs)

				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentTXOutputs {

						txHashStr := hex.EncodeToString(tx.TxHash)

						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _, outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}

								if isUnSpentUTXO == false {
									utxo := &UTXO{tx.TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}

		}

	}

	blockIterator := blockchain.Iterator()

	for {

		block := blockIterator.Next()

		fmt.Println(block)
		fmt.Println()

		for i := len(block.Txs) - 1; i >= 0; i-- {

			tx := block.Txs[i]
			// txHash
			// Vins
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					//是否能够解锁
					if in.UnLockWithAddress(address) {

						key := hex.EncodeToString(in.TxHash)

						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}

				}
			}

			// Vouts

		work:
			for index, out := range tx.Vouts {

				if out.UnLockScriptPubKeyWithAddress(address) {

					fmt.Println(out)
					fmt.Println(spentTXOutputs)

					//&{2 zhangqiang}
					//map[]

					if spentTXOutputs != nil {

						//map[cea12d33b2e7083221bf3401764fb661fd6c34fab50f5460e77628c42ca0e92b:[0]]

						if len(spentTXOutputs) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range spentTXOutputs {

								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}

							if isSpentUTXO == false {

								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)

							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}
				}

			}

		}

		fmt.Println(spentTXOutputs)

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}

	}

	return unUTXOs
}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpendableUTXOS(from string, amount int, txs []*Transaction) (int64, map[string][]int) {

	//1. 现获取所有的UTXO

	utxos := blockchain.UnUTXOs(from, txs)

	spendableUTXO := make(map[string][]int)

	//2. 遍历utxos

	var value int64

	for _, utxo := range utxos {

		value = value + utxo.Output.Value

		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)

		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) {

		fmt.Printf("%s's fund is 不足\n", from)
		os.Exit(1)
	}

	return value, spendableUTXO
}

// 挖掘新的区块
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {

	//	$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	//1.建立一笔交易

	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(amount)

	var txs []*Transaction

	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], value, blockchain, txs)
		txs = append(txs, tx)
		//fmt.Println(tx)
	}

	//1. 通过相关算法建立Transaction数组
	var block *Block

	blockchain.Db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = DeserializeBlock(blockBytes)

		}

		return nil
	})

	//2. 建立新的区块
	//	block = NewBlock(txs, block.Height+1, block.Hash)
	block = NewBlock(txs, block.Height+1,block.Hash)

	//将新区块存储到数据库
	blockchain.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			b.Put(block.Hash, block.Serialize())

			b.Put([]byte("l"), block.Hash)

			blockchain.Tip = block.Hash

		}
		return nil
	})

}

// 查询余额
func (blockchain *Blockchain) GetBalance(address string) int64 {

	utxos := blockchain.UnUTXOs(address, []*Transaction{})

	var amount int64

	for _, utxo := range utxos {

		amount = amount + utxo.Output.Value
	}

	return amount
}




// 遍历输出所有区块的信息
func (blc *Blockchain) Printchain() {

	fmt.Println("遍历区块链，打印：")
	blockchainIterator := blc.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("Height：%d\n", block.Height)
		fmt.Printf("PrevBlockHash：%x\n", block.PrevBlockHash)
		fmt.Printf("Timestamp：%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.Hash)
		fmt.Printf("Nonce：%d\n", block.Nonce)
		fmt.Println("Txs:")
		for _, tx := range block.Txs {

			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%s\n", in.ScriptSig)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {

				fmt.Printf("%s:%d\n",out.ScriptPubKey,out.Value)
				//fmt.Println(out.Value)

			}
		}

		fmt.Println("------------------------------")

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

}