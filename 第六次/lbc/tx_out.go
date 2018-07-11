package lbc

import "bytes"

//XHQ_TXOutput{100,"zhangbozhi"}
//XHQ_TXOutput{30,"xietingfeng"}
//XHQ_TXOutput{40,"zhangbozhi"}


type XHQ_TXOutput struct {
	Value int64
	//ScriptPubKey string  //用户名

	Ripemd160Hash []byte  //用户名
}

/*
// 解锁
func (txOutput *XHQ_TXOutput) XHQ_UnLockScriptPubKeyWithAddress(address string) bool {

	return txOutput.ScriptPubKey == address
}
*/



func (txOutput *XHQ_TXOutput)  Lock(address string)  {

	publicKeyHash := XHQ_Base58Decode([]byte(address))

	txOutput.Ripemd160Hash = publicKeyHash[1:len(publicKeyHash) - 4]
}


func NewXHQ_TXOutput(value int64,address string) *XHQ_TXOutput {

	txOutput := &XHQ_TXOutput{value,nil}

	// 设置Ripemd160Hash
	txOutput.Lock(address)

	return txOutput
}


// 解锁
func (txOutput *XHQ_TXOutput) XHQ_UnLockScriptPubKeyWithAddress(address string) bool {

	publicKeyHash := XHQ_Base58Decode([]byte(address))
	hash160 := publicKeyHash[1:len(publicKeyHash) - 4]

	return bytes.Compare(txOutput.Ripemd160Hash,hash160) == 0
}

