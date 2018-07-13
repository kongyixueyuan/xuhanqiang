package lbc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
)

// UTXO
type Transaction struct {

	//1. 交易hash
	TxHash []byte

	//2. 输入
	Vins []*XHQ_TXInput

	//3. 输出
	Vouts []*XHQ_TXOutput
}

//[]byte{}

// 判断当前的交易是否是Coinbase交易
func (tx *Transaction) XHQ_IsCoinbaseTransaction() bool {

	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

//1. Transaction 创建分两种情况
//1. 创世区块创建时的Transaction
func XHQ_NewCoinbaseTransaction(address string) *Transaction {

	//代表消费

	txInput := &XHQ_TXInput{[]byte{}, -1, nil, []byte{}}

	txOutput := NewXHQ_TXOutput(20, address)

	//	txOutput := &XHQ_TXOutput{1000,address}

	txCoinbase := &Transaction{[]byte{}, []*XHQ_TXInput{txInput}, []*XHQ_TXOutput{txOutput}}

	//设置hash值
	txCoinbase.XHQ_HashTransaction()

	return txCoinbase
}

func (tx *Transaction) XHQ_HashTransaction() {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())

	tx.TxHash = hash[:]
}



//2. 转账时产生的Transaction

func XHQ_NewSimpleTransaction(from string,to string,amount int64,utxoSet *XHQ_UTXOSet,txs []*Transaction) *Transaction {

	//$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	wallets,_ := XHQ_NewWallets()
	wallet := wallets.WalletsMap[from]


	// 通过一个函数，返回
	money,spendableUTXODic := utxoSet.XHQ_FindSpendableUTXOS(from,amount,txs)
	//
	//	{hash1:[0],hash2:[2,3]}

	var txIntputs []*XHQ_TXInput
	var txOutputs []*XHQ_TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {

			txInput := &XHQ_TXInput{txHashBytes,index,nil,wallet.PublicKey}
			txIntputs = append(txIntputs,txInput)
		}

	}

	// 转账
	txOutput := NewXHQ_TXOutput(int64(amount),to)
	txOutputs = append(txOutputs,txOutput)

	// 找零
	txOutput = NewXHQ_TXOutput(int64(money) - int64(amount),from)
	txOutputs = append(txOutputs,txOutput)

	tx := &Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.XHQ_HashTransaction()

	//进行签名
	utxoSet.Blockchain.XHQ_SignTransaction(tx, wallet.PrivateKey,txs)

	return tx

}



/*
func XHQ_NewSimpleTransaction(from string, to string, amount int, blockchain *Blockchain, txs []*Transaction) *Transaction {

	//$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	//验证接收地址是否有效

	if XHQ_IsValidForAdress([]byte(to)) == false{
		fmt.Println("接收地址无效....")
		os.Exit(1)
	}



	wallets, _ := XHQ_NewWallets()
	wallet := wallets.WalletsMap[from]

	// 通过一个函数，返回
	money, spendableUTXODic := blockchain.XHQ_FindSpendableUTXOS(from, amount, txs)
	//
	//	{hash1:[0],hash2:[2,3]}

	var txIntputs []*XHQ_TXInput
	var txOutputs []*XHQ_TXOutput

	for txHash, indexArray := range spendableUTXODic {

		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {

			txInput := &XHQ_TXInput{txHashBytes, index, nil, wallet.PublicKey}
			txIntputs = append(txIntputs, txInput)
		}

	}

	// 转账
	txOutput := NewXHQ_TXOutput(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)

	// 找零
	txOutput = NewXHQ_TXOutput(int64(money)-int64(amount), from)
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{[]byte{}, txIntputs, txOutputs}

	//设置hash值
	tx.XHQ_HashTransaction()

	//进行签名
	blockchain.XHQ_SignTransaction(tx, wallet.PrivateKey)

	return tx

}*/

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {

	if tx.XHQ_IsCoinbaseTransaction() {
		return
	}

	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vins[inID].Signature = signature
	}
}

// 拷贝一份新的Transaction用于签名                                    T
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*XHQ_TXInput
	var outputs []*XHQ_TXOutput

	for _, vin := range tx.Vins {
		inputs = append(inputs, &XHQ_TXInput{vin.TxHash, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vouts {
		outputs = append(outputs, &XHQ_TXOutput{vout.Value, vout.Ripemd160Hash})
	}

	txCopy := Transaction{tx.TxHash, inputs, outputs}

	return txCopy
}

// 数字签名验证

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.XHQ_IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		// 私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {
			return false
		}
	}

	return true
}

func (tx *Transaction) Hash() []byte {

	txCopy := tx

	txCopy.TxHash = []byte{}

	hash := sha256.Sum256(txCopy.XHQ_Serialize())
	return hash[:]
}

func (tx *Transaction) XHQ_Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}
