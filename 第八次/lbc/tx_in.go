package lbc

import "bytes"

type XHQ_TXInput struct {
	// 1. 交易的Hash
	TxHash      []byte
	// 2. 存储XHQ_TXOutput在Vout里面的索引
	Vout      int
	// 3. 用户名
	//ScriptSig string

	Signature []byte // 数字签名

	PublicKey    []byte // 公钥，钱包里面

}


/*

// 判断当前的消费是谁的钱
func (txInput *XHQ_TXInput) UnLockWithAddress(address string) bool {

	return txInput.ScriptSig == address
}*/


// 判断当前的消费是谁的钱
func (txInput *XHQ_TXInput) XHQ_UnLockRipemd160Hash(ripemd160Hash []byte) bool {

	publicKey := Ripemd160Hash(txInput.PublicKey)

	return bytes.Compare(publicKey,ripemd160Hash) == 0
}