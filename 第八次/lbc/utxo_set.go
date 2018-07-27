package lbc

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"os"
	"bytes"
)

const utxoTableName = "utxoTableName"

type XHQ_UTXOSet struct {
	Blockchain *Blockchain
}
/*

// 重置数据库表
func (utxoSet *XHQ_UTXOSet) ResetXHQ_UTXOSet() {

	err := utxoSet.Blockchain.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err != nil {
				log.Panic(err)
			}
		}

		b, _ = tx.CreateBucket([]byte(utxoTableName))
		if b != nil {

			//[string]*XHQ_TXOutputs
			txOutputsMap := utxoSet.Blockchain.XHQ_FindUTXOMap()

			for keyHash, outs := range txOutputsMap {

				txHash, _ := hex.DecodeString(keyHash)

				b.Put(txHash, outs.XHQ_Serialize())
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
*/



// 重置数据库表
func (utxoSet *XHQ_UTXOSet) ResetXHQ_UTXOSet()  {

	err := utxoSet.Blockchain.Db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {


			err := tx.DeleteBucket([]byte(utxoTableName))

			if err!= nil {
				log.Panic(err)
			}

		}

		b ,_ = tx.CreateBucket([]byte(utxoTableName))
		if b != nil {

			//[string]*TXOutputs
			txOutputsMap := utxoSet.Blockchain.XHQ_FindUTXOMap()


			for keyHash,outs := range txOutputsMap {

				txHash,_ := hex.DecodeString(keyHash)

				b.Put(txHash,outs.XHQ_Serialize())

			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("重置失败....")
		log.Panic(err)
	}

}



func (utxoSet *XHQ_UTXOSet) XHQ_GetBalance(address string) int64 {

	UTXOS := utxoSet.findUTXOForAddress(address)

	var amount int64

	for _, utxo := range UTXOS {
		amount += utxo.Output.Value
	}

	return amount
}

func (utxoSet *XHQ_UTXOSet) findUTXOForAddress(address string) []*UTXO {

	var utxos []*UTXO

	utxoSet.Blockchain.Db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))

		if b == nil{
			println("找不到utxosets库，请运行resetutxo方法!")
			os.Exit(1)
		}

		// 游标
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			txOutputs := DeserializeXHQ_TXOutputs(v)

			for _, utxo := range txOutputs.UTXOS {

				if utxo.Output.XHQ_UnLockScriptPubKeyWithAddress(address) {
					utxos = append(utxos, utxo)
				}
			}
		}

		return nil
	})

	return utxos
}

func (utxoSet *XHQ_UTXOSet) XHQ_FindSpendableUTXOS(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {

	unPackageUTXOS := utxoSet.XHQ_FindUnPackageSpendableUTXOS(from, txs)

	spentableUTXO := make(map[string][]int)

	var money int64 = 0

	for _, UTXO := range unPackageUTXOS {

		money += UTXO.Output.Value
		txHash := hex.EncodeToString(UTXO.TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash], UTXO.Index)
		if money >= amount {
			return money, spentableUTXO
		}
	}

	// 钱还不够
	utxoSet.Blockchain.Db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {

			c := b.Cursor()
		UTXOBREAK:
			for k, v := c.First(); k != nil; k, v = c.Next() {

				txOutputs := DeserializeXHQ_TXOutputs(v)

				for _, utxo := range txOutputs.UTXOS {

					money += utxo.Output.Value
					txHash := hex.EncodeToString(utxo.TxHash)
					spentableUTXO[txHash] = append(spentableUTXO[txHash], utxo.Index)

					if money >= amount {
						break UTXOBREAK
					}
				}
			}

		}

		return nil
	})

	if money < amount {
		log.Panic("余额不足......")
	}

	return money, spentableUTXO
}

// 返回要凑多少钱，对应XHQ_TXOutput的TX的Hash和index
func (utxoSet *XHQ_UTXOSet) XHQ_FindUnPackageSpendableUTXOS(from string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO

	spentXHQ_TXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _, tx := range txs {

		if tx.XHQ_IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//是否能够解锁
				publicKeyHash := XHQ_Base58Decode([]byte(from))

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

			if out.XHQ_UnLockScriptPubKeyWithAddress(from) {
				fmt.Println("看看是否是俊诚...")
				fmt.Println(from)

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

	return unUTXOs

}





/*

// 更新
func (utxoSet *XHQ_UTXOSet) Update()  {

	// blocks
	//


	// 最新的Block
	block := utxoSet.Blockchain.Iterator().Next()


	// utxoTable
	//

	ins := []*XHQ_TXInput{}

	outsMap := make(map[string]*XHQ_TXOutputs)

	// 找到所有我要删除的数据
	for _,tx := range block.Txs {

		for _,in := range tx.Vins {
			ins = append(ins,in)
		}
	}

	for _,tx := range block.Txs  {


		utxos := []*UTXO{}

		for index,out := range tx.Vouts  {

			isSpent := false

			for _,in := range ins  {

				if in.Vout == index && bytes.Compare(tx.TxHash ,in.TxHash) == 0 && bytes.Compare(out.Ripemd160Hash,Ripemd160Hash(in.PublicKey)) == 0 {

					isSpent = true
					continue
				}
			}

			if isSpent == false {
				utxo := &UTXO{tx.TxHash,index,out}
				utxos = append(utxos,utxo)
			}

		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.TxHash)
			outsMap[txHash] = &XHQ_TXOutputs{utxos}
		}

	}



	err := utxoSet.Blockchain.Db.Update(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {


			// 删除
			for _,in := range ins {

				txOutputsBytes := b.Get(in.TxHash)

				if len(txOutputsBytes) == 0 {
					continue
				}

				fmt.Println("DeserializeXHQ_TXOutputs")
				fmt.Println(txOutputsBytes)

				txOutputs := DeserializeXHQ_TXOutputs(txOutputsBytes)

				fmt.Println(txOutputs)

				UTXOS := []*UTXO{}

				// 判断是否需要
				isNeedDelete := false

				for _,utxo := range txOutputs.UTXOS  {

					if in.Vout == utxo.Index && bytes.Compare(utxo.Output.Ripemd160Hash,Ripemd160Hash(in.PublicKey)) == 0 {

						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS,utxo)
					}
				}



				if isNeedDelete {
					b.Delete(in.TxHash)
					if len(UTXOS) > 0 {

						preXHQ_TXOutputs := outsMap[hex.EncodeToString(in.TxHash)]

						preXHQ_TXOutputs.UTXOS = append(preXHQ_TXOutputs.UTXOS,UTXOS...)

						outsMap[hex.EncodeToString(in.TxHash)] = preXHQ_TXOutputs

					}
				}

			}

			// 新增

			for keyHash,outPuts := range outsMap  {
				keyHashBytes,_ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes,outPuts.XHQ_Serialize())
			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}
*/





// 更新
func (utxoSet *XHQ_UTXOSet) Update()  {

	// blocks
	//


	// 最新的Block
	block := utxoSet.Blockchain.Iterator().Next()


	// utxoTable
	//

	ins := []*XHQ_TXInput{}

	outsMap := make(map[string]*XHQ_TXOutputs)

	// 找到所有我要删除的数据
	for _,tx := range block.Txs {

		for _,in := range tx.Vins {
			ins = append(ins,in)
		}
	}

	for _,tx := range block.Txs  {


		utxos := []*UTXO{}

		for index,out := range tx.Vouts  {

			isSpent := false

			for _,in := range ins  {

				if in.Vout == index && bytes.Compare(tx.TxHash ,in.TxHash) == 0 && bytes.Compare(out.Ripemd160Hash,Ripemd160Hash(in.PublicKey)) == 0 {

					isSpent = true
					continue
				}
			}

			if isSpent == false {
				utxo := &UTXO{tx.TxHash,index,out}
				utxos = append(utxos,utxo)
			}

		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.TxHash)
			outsMap[txHash] = &XHQ_TXOutputs{utxos}
		}

	}



	err := utxoSet.Blockchain.Db.Update(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {


			// 删除
			for _,in := range ins {

				txOutputsBytes := b.Get(in.TxHash)

				if len(txOutputsBytes) == 0 {
					continue
				}

				fmt.Println("DeserializeTXOutputs")
				fmt.Println(txOutputsBytes)

				txOutputs := DeserializeXHQ_TXOutputs(txOutputsBytes)

				fmt.Println(txOutputs)

				UTXOS := []*UTXO{}

				// 判断是否需要
				isNeedDelete := false

				for _,utxo := range txOutputs.UTXOS  {

					if in.Vout == utxo.Index && bytes.Compare(utxo.Output.Ripemd160Hash,Ripemd160Hash(in.PublicKey)) == 0 {

						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS,utxo)
					}
				}



				if isNeedDelete {
					b.Delete(in.TxHash)
					if len(UTXOS) > 0 {

						preTXOutputs := outsMap[hex.EncodeToString(in.TxHash)]

						preTXOutputs.UTXOS = append(preTXOutputs.UTXOS,UTXOS...)

						outsMap[hex.EncodeToString(in.TxHash)] = preTXOutputs

					}
				}

			}

			// 新增

			for keyHash,outPuts := range outsMap  {
				keyHashBytes,_ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes,outPuts.XHQ_Serialize())
			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}



