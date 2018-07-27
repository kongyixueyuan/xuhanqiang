package lbc

import (
	"fmt"
	"log"

	"encoding/hex"
	"math/big"
	"os"
	"strconv"

	"bytes"
	"crypto/ecdsa"
	"time"

	"github.com/boltdb/bolt"
)

// 数据库名字
//const dbName = "blockchain.db"
const dbName = "blockchain_%s.db"

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
type XHQ_BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// 添加一个区块(带有交易信息的数据)到区块链内，
func (bc *Blockchain) XHQ_AddBlock(txs []*Transaction) {
	var lastHash []byte

	/*
		err := bc.Db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			lastHash = b.Get([]byte("l"))

			// ⚠️，先获取最新区块
			blockBytes := b.Get(blc.Tip)
			// 反序列化
			block := XHQ_DeserializeBlock(blockBytes)

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
		block := XHQ_DeserializeBlock(blockBytes)

		newBlock := NewBlock(txs, block.Height+1, lastHash)

		//

		err := b.Put(newBlock.Hash, newBlock.XHQ_Serialize())
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
func (bc *Blockchain) Iterator() *XHQ_BlockchainIterator {
	bci := &XHQ_BlockchainIterator{bc.Tip, bc.Db}

	return bci
}

// 区块链的迭代器的下一个区块函数。
func (i *XHQ_BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = XHQ_DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.XHQ_PrevBlockHash

	return block
}

// 创建一个带有创世区块的区块链。
func NewBlockchain(address string) *Blockchain {
	var tip []byte

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("没有现成的区块链，创建一个新的...")
			//genesis := NewGenesisBlock()
			txCoinbase := XHQ_NewCoinbaseTransaction(address)
			genesis := XHQ_CreateGenesisBlock([]*Transaction{txCoinbase})

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.XHQ_Serialize())
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
func DBExists(nodeID string) bool {
	dbName := fmt.Sprintf(dbName,nodeID)
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}

///
//1. 创建带有创世区块的区块链
func XHQ_CreateBlockchainWithGenesisBlock(address string,nodeID string) *Blockchain {

	dbName := fmt.Sprintf(dbName,nodeID)

	// 判断数据库是否存在
	if DBExists(nodeID) {
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
			txCoinbase := XHQ_NewCoinbaseTransaction(address)

			genesisBlock := XHQ_CreateGenesisBlock([]*Transaction{txCoinbase})
			// 将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.XHQ_Serialize())
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
func XHQ_BlockchainObject(nodeID string) *Blockchain {

	dbName := fmt.Sprintf(dbName, nodeID)

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

/*
// 如果一个地址对应的XHQ_TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (blockchain *Blockchain) XHQ_UnUTXOs(address string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO

	spentXHQ_TXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _, tx := range txs {

		if tx.XHQ_IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//是否能够解锁
				if in.UnLockWithAddress(address) {

					key := hex.EncodeToString(in.TxHash)

					spentXHQ_TXOutputs[key] = append(spentXHQ_TXOutputs[key], in.Vout)
				}

			}
		}
	}

	for _, tx := range txs {

	Work1:
		for index, out := range tx.Vouts {

			if out.XHQ_UnLockScriptPubKeyWithAddress(address) {
				fmt.Println("查看...")
				fmt.Println(address)

				fmt.Println(spentXHQ_TXOutputs)

				if len(spentXHQ_TXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentXHQ_TXOutputs {

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
			if tx.XHQ_IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					//是否能够解锁
					if in.UnLockWithAddress(address) {

						key := hex.EncodeToString(in.TxHash)

						spentXHQ_TXOutputs[key] = append(spentXHQ_TXOutputs[key], in.Vout)
					}

				}
			}

			// Vouts

		work:
			for index, out := range tx.Vouts {

				if out.XHQ_UnLockScriptPubKeyWithAddress(address) {

					fmt.Println(out)
					fmt.Println(spentXHQ_TXOutputs)

					//&{2 zhangqiang}
					//map[]

					if spentXHQ_TXOutputs != nil {

						//map[cea12d33b2e7083221bf3401764fb661fd6c34fab50f5460e77628c42ca0e92b:[0]]

						if len(spentXHQ_TXOutputs) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range spentXHQ_TXOutputs {

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

		fmt.Println(spentXHQ_TXOutputs)

		var hashInt big.Int
		hashInt.SetBytes(block.XHQ_PrevBlockHash)

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
}*/

// 如果一个地址对应的XHQ_TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (blockchain *Blockchain) XHQ_UnUTXOs(address string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO

	spentXHQ_TXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _, tx := range txs {

		if tx.XHQ_IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//是否能够解锁
				publicKeyHash := XHQ_Base58Decode([]byte(address))

				ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4]
				if in.XHQ_UnLockRipemd160Hash(ripemd160Hash) {

					key := hex.EncodeToString(in.TxHash)

					spentXHQ_TXOutputs[key] = append(spentXHQ_TXOutputs[key], in.Vout)
				}

			}
		}
	}

	for _, tx := range txs {

	Work1:
		for index, out := range tx.Vouts {

			if out.XHQ_UnLockScriptPubKeyWithAddress(address) {
				fmt.Println("看看是否是俊诚...")
				//	fmt.Println(address)

				//	fmt.Println(spentXHQ_TXOutputs)

				if len(spentXHQ_TXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentXHQ_TXOutputs {

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

		//	fmt.Println(block)
		//	fmt.Println()

		for i := len(block.Txs) - 1; i >= 0; i-- {

			tx := block.Txs[i]
			// txHash
			// Vins
			if tx.XHQ_IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					//是否能够解锁
					publicKeyHash := XHQ_Base58Decode([]byte(address))

					ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4]

					if in.XHQ_UnLockRipemd160Hash(ripemd160Hash) {

						key := hex.EncodeToString(in.TxHash)

						spentXHQ_TXOutputs[key] = append(spentXHQ_TXOutputs[key], in.Vout)
					}

				}
			}

			// Vouts

		work:
			for index, out := range tx.Vouts {

				if out.XHQ_UnLockScriptPubKeyWithAddress(address) {

					//	fmt.Println(out)
					//	fmt.Println(spentXHQ_TXOutputs)

					//&{2 zhangqiang}
					//map[]

					if spentXHQ_TXOutputs != nil {

						//map[cea12d33b2e7083221bf3401764fb661fd6c34fab50f5460e77628c42ca0e92b:[0]]

						if len(spentXHQ_TXOutputs) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range spentXHQ_TXOutputs {

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

		//	fmt.Println(spentXHQ_TXOutputs)

		var hashInt big.Int
		hashInt.SetBytes(block.XHQ_PrevBlockHash)

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
func (blockchain *Blockchain) XHQ_FindSpendableUTXOS(from string, amount int, txs []*Transaction) (int64, map[string][]int) {

	//1. 现获取所有的UTXO

	utxos := blockchain.XHQ_UnUTXOs(from, txs)

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

/*

// 挖掘新的区块
func (blockchain *Blockchain) XHQ_MineNewBlock(from []string, to []string, amount []string) {

	//	$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	//1.建立一笔交易

	utxoSet := &XHQ_UTXOSet{blockchain}

	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(amount)

	var txs []*Transaction

	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
	//	tx := XHQ_NewSimpleTransaction(address, to[index], value, blockchain, txs)

		tx := XHQ_NewSimpleTransaction(address, to[index], int64(value), utxoSet,txs)
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

			block = XHQ_DeserializeBlock(blockBytes)

		}

		return nil
	})

	// 在建立新区块之前对txs进行签名验证

	_txs := []*Transaction{}

	for _, tx := range txs {

		if blockchain.VerifyTransaction(tx,_txs) != true {

		//	if blockchain.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	//奖励
	//coinbaseTx := XHQ_NewCoinbaseTransaction(from[0])
	//txs = append(txs, coinbaseTx)

	//2. 建立新的区块
	//	block = NewBlock(txs, block.Height+1, block.Hash)
	block = NewBlock(txs, block.Height+1, block.Hash)

	//将新区块存储到数据库
	blockchain.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			b.Put(block.Hash, block.XHQ_Serialize())

			b.Put([]byte("l"), block.Hash)

			blockchain.Tip = block.Hash

		}
		return nil
	})

}

*/

// 挖掘新的区块
func (blockchain *Blockchain) XHQ_MineNewBlock(from []string, to []string, amount []string,nodeID string) {

	//	$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	//1.建立一笔交易

	utxoSet := &XHQ_UTXOSet{blockchain}

	var txs []*Transaction

	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := XHQ_NewSimpleTransaction(address, to[index], int64(value), utxoSet, txs)
		txs = append(txs, tx)
		//fmt.Println(tx)
	}

	//奖励
	//	tx := XHQ_NewCoinbaseTransaction(from[0])
	//	txs = append(txs,tx)

	//1. 通过相关算法建立Transaction数组
	var block *Block

	blockchain.Db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = XHQ_DeserializeBlock(blockBytes)

		}

		return nil
	})

	// 在建立新区块之前对txs进行签名验证

	_txs := []*Transaction{}

	for _, tx := range txs {

		if blockchain.VerifyTransaction(tx, _txs) != true {
			log.Panic("ERROR: Invalid transaction")
		}

		_txs = append(_txs, tx)
	}

	//2. 建立新的区块
	block = NewBlock(txs, block.Height+1, block.Hash)

	//将新区块存储到数据库
	blockchain.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			b.Put(block.Hash, block.XHQ_Serialize())

			b.Put([]byte("l"), block.Hash)

			blockchain.Tip = block.Hash

		}
		return nil
	})

}

// 查询余额
func (blockchain *Blockchain) XHQ_GetBalance(address string) int64 {

	if XHQ_IsValidForAdress([]byte(address)) == false {
		fmt.Println("地址无效....")
		os.Exit(1)
	}

	utxos := blockchain.XHQ_UnUTXOs(address, []*Transaction{})

	var amount int64

	for _, utxo := range utxos {

		amount = amount + utxo.Output.Value
	}

	return amount
}

// 遍历输出所有区块的信息
func (blc *Blockchain) XHQ_Printchain(nodeId string) {

	fmt.Println("遍历区块链，打印：")
	blockchainIterator := blc.Iterator()

		for {
		block := blockchainIterator.Next()

		fmt.Printf("Height：%d\n", block.Height)
		fmt.Printf("XHQ_PrevBlockHash：%x\n", block.XHQ_PrevBlockHash)
		fmt.Printf("XHQ_Timestamp：%s\n", time.Unix(block.XHQ_Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.Hash)
		fmt.Printf("XHQ_Nonce：%d\n", block.XHQ_Nonce)
		fmt.Println("Txs:")
		for _, tx := range block.Txs {

			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%x\n", in.PublicKey)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {

				//	fmt.Printf("%s:%d\n", out.Ripemd160Hash, out.Value)
				//fmt.Println(out.Value)

				fmt.Printf("%d\n", out.Value)
				//fmt.Println(out.Ripemd160Hash)
				fmt.Printf("%x\n", out.Ripemd160Hash)

			}
		}

		fmt.Println("------------------------------")

		var hashInt big.Int
		hashInt.SetBytes(block.XHQ_PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

}

/*

func (bclockchain *Blockchain) XHQ_SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {

	if tx.XHQ_IsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTX, err := bclockchain.XHQ_FindTransaction(vin.TxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	tx.Sign(privKey, prevTXs)

}

*/

func (bclockchain *Blockchain) XHQ_SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey, txs []*Transaction) {

	if tx.XHQ_IsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTX, err := bclockchain.XHQ_FindTransaction(vin.TxHash, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	tx.Sign(privKey, prevTXs)

}

/*

func (bc *Blockchain) XHQ_FindTransaction(ID []byte) (Transaction, error) {

	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Txs {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.XHQ_PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

	return Transaction{}, nil
}
*/

func (bc *Blockchain) XHQ_FindTransaction(ID []byte, txs []*Transaction) (Transaction, error) {

	for _, tx := range txs {
		if bytes.Compare(tx.TxHash, ID) == 0 {
			return *tx, nil
		}
	}

	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Txs {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.XHQ_PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

	return Transaction{}, nil
}

/*

// 验证数字签名
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTX, err := bc.XHQ_FindTransaction(vin.TxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	return tx.Verify(prevTXs)
}
*/

// 验证数字签名
func (bc *Blockchain) VerifyTransaction(tx *Transaction, txs []*Transaction) bool {

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTX, err := bc.XHQ_FindTransaction(vin.TxHash, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	return tx.Verify(prevTXs)
}

//
/*

func (blc *Blockchain) XHQ_FindUTXOMap() map[string]*XHQ_TXOutputs {

	blcIterator := blc.Iterator()
	// 存储已花费的UTXO的信息
	spentUTXOsMap := make(map[string][]*XHQ_TXInput)

	utxoMaps := make(map[string]*XHQ_TXOutputs)

	var isSpent bool

	for {
		block := blcIterator.Next()

		//每笔交易
		for i := len(block.Txs) - 1; i >= 0; i-- {

			txOutputs := &XHQ_TXOutputs{[]*UTXO{}}

			tx := block.Txs[i]

			// coinbase
			if tx.XHQ_IsCoinbaseTransaction() == false {
				for _, txInput := range tx.Vins {

					txHash := hex.EncodeToString(txInput.TxHash)
					spentUTXOsMap[txHash] = append(spentUTXOsMap[txHash], txInput)

				}
			}

			txHash := hex.EncodeToString(tx.TxHash)

			//	WorkOutLoop:

			//每个输出
			for index, out := range tx.Vouts {

				if tx.XHQ_IsCoinbaseTransaction() {

					fmt.Println("XHQ_IsCoinbaseTransaction")
					fmt.Println(out)
					fmt.Println(txHash)
				}

				//	txInputs := spentUTXOsMap[txHash]

				//fmt.Printf("输入长度：%d\n",len(txInputs))
				fmt.Printf("txHash：%v", txHash)
				println("---->")

				//	if len(txInputs) > 0 {
				//	if block.Height > 1 {
				//	if tx.XHQ_IsCoinbaseTransaction() == false{
				if true {

					if tx.XHQ_IsCoinbaseTransaction() {
						isSpent = false
					} else {
						isSpent = false
					}
					//
					for _, txInputs := range spentUTXOsMap {
						//

						for _, in := range txInputs {

							outPublicKey := out.Ripemd160Hash
							inPublicKey := in.PublicKey

							if bytes.Compare(outPublicKey, Ripemd160Hash(inPublicKey)) == 0 {
								if index == in.Vout {
									isSpent = true
									//continue WorkOutLoop
								} else {
									isSpent = false
								}
							}

						}

						//
					}

					if isSpent == false {
						utxo := &UTXO{tx.TxHash, index, out}

						println("print_trade_out:")
						fmt.Printf("本次输出交易的->hash%s\n", hex.EncodeToString(tx.TxHash))

						fmt.Printf("%d-->%x\n", out.Value, out.Ripemd160Hash)
						println("print_end:")
						txOutputs.UTXOS = append(txOutputs.UTXOS, utxo)
					}

					//
				} else {
					utxo := &UTXO{tx.TxHash, index, out}

					println("print_coinbase:")

					fmt.Printf("本次输出交易的hash->%s\n", hex.EncodeToString(tx.TxHash))
					fmt.Printf("%d-->%x\n", out.Value, out.Ripemd160Hash)
					println("print_end:")

					txOutputs.UTXOS = append(txOutputs.UTXOS, utxo)
				}
			}
			// 设置键值对
			utxoMaps[txHash] = txOutputs

		}

		// 找到创世区块时退出
		var hashInt big.Int
		hashInt.SetBytes(block.XHQ_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return utxoMaps
}
*/



// [string]*TXOutputs
func (blc *Blockchain) XHQ_FindUTXOMap() map[string]*XHQ_TXOutputs  {

	blcIterator := blc.Iterator()

	// 存储已花费的UTXO的信息
	spentableUTXOsMap := make(map[string][]*XHQ_TXInput)


	utxoMaps := make(map[string]*XHQ_TXOutputs)


	for {
		block := blcIterator.Next()

		for i := len(block.Txs) - 1; i >= 0 ;i-- {

			txOutputs := &XHQ_TXOutputs{[]*UTXO{}}

			tx := block.Txs[i]

			// coinbase
			if tx.XHQ_IsCoinbaseTransaction() == false {
				for _,txInput := range tx.Vins {

					txHash := hex.EncodeToString(txInput.TxHash)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash],txInput)

				}
			}

			txHash := hex.EncodeToString(tx.TxHash)

			txInputs := spentableUTXOsMap[txHash]

			if len(txInputs) > 0 {


			WorkOutLoop:
				for index,out := range tx.Vouts  {

					for _,in := range  txInputs {

						outPublicKey := out.Ripemd160Hash
						inPublicKey := in.PublicKey


						if bytes.Compare(outPublicKey,Ripemd160Hash(inPublicKey)) == 0 {
							if index == in.Vout {

								continue WorkOutLoop
							} else {

								utxo := &UTXO{tx.TxHash,index,out}
								txOutputs.UTXOS = append(txOutputs.UTXOS,utxo)
							}
						}
					}


				}

			} else {

				for index,out := range tx.Vouts {
					utxo := &UTXO{tx.TxHash,index,out}
					txOutputs.UTXOS = append(txOutputs.UTXOS,utxo)
				}
			}


			// 设置键值对
			utxoMaps[txHash] = txOutputs

		}


		// 找到创世区块时退出
		var hashInt big.Int
		hashInt.SetBytes(block.XHQ_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}



	}

	return utxoMaps
}






// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Blockchain) XHQ_FindUTXOMap_rwq() map[string]*XHQ_TXOutputs {
	// 未花费的交易输出
	// key:交易hash   txID
	UTXOR := make(map[string]*XHQ_TXOutputs)
	// 已经花费的交易txID : XHQ_TXOutputs.index
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		// 循环区块中的交易
		for _, tx := range block.Txs {
			// 将区块中的交易hash，转为字符串
			txID := hex.EncodeToString(tx.TxHash)

		Outputs:
			for outIdx, out := range tx.Vouts { // 循环交易中的 XHQ_TXOutputs
				// Was the output spent?
				// 如果已经花费的交易输出中，有此输出，证明已经花费
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx { // 如果花费的正好是此笔输出
							continue Outputs // 继续下一次循环
						}
					}
				}

				outs := UTXOR[txID] // 获取UTXO指定txID对应的XHQ_TXOutputs

				utxo := &UTXO{tx.TxHash, outIdx, out}
				outs.UTXOS = append(outs.UTXOS, utxo)
				UTXOR[txID] = outs
			}

			if tx.XHQ_IsCoinbaseTransaction() == false { // 非创世区块
				for _, in := range tx.Vins {
					inTxID := hex.EncodeToString(in.TxHash)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}
		// 如果上一区块的hash为0，代表已经到创世区块，循环结束
		if len(block.XHQ_PrevBlockHash) == 0 {
			break
		}
	}

	return UTXOR
}

//----------

func (bc *Blockchain) GetBestHeight() int64 {

	block := bc.Iterator().Next()

	return block.Height
}

func (bc *Blockchain) GetBlockHashes() [][]byte {

	blockIterator := bc.Iterator()

	var blockHashs [][]byte

	for {
		block := blockIterator.Next()

		blockHashs = append(blockHashs, block.Hash)

		var hashInt big.Int
		hashInt.SetBytes(block.XHQ_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return blockHashs
}





func (bc *Blockchain) GetBlock(blockHash []byte) (*Block ,error) {

	var block *Block

	err := bc.Db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {

			blockBytes := b.Get(blockHash)

			block = XHQ_DeserializeBlock(blockBytes)

		}

		return nil
	})

	return block,err
}

func (bc *Blockchain) AddBlock(block *Block) error {

	err := bc.Db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {

			blockExist := b.Get(block.Hash)

			if blockExist != nil {
				// 如果存在，不需要做任何过多的处理
				return nil
			}

			err := b.Put(block.Hash,block.XHQ_Serialize())

			if err != nil {
				log.Panic(err)
			}

			// 最新的区块链的Hash
			blockHash := b.Get([]byte("l"))

			blockBytes := b.Get(blockHash)

			blockInDB := XHQ_DeserializeBlock(blockBytes)

			if blockInDB.Height < block.Height {

				b.Put([]byte("l"),block.Hash)
			}
		}

		return nil
	})

	return err
}