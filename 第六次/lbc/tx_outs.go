package lbc

import (
	"bytes"
	"encoding/gob"
	"log"
)

type XHQ_TXOutputs struct {
	UTXOS []*UTXO
}

// 将区块序列化成字节数组
func (txOutputs *XHQ_TXOutputs) XHQ_Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(txOutputs)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func DeserializeXHQ_TXOutputs(txOutputsBytes []byte) *XHQ_TXOutputs {

	var txOutputs XHQ_TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}

	return &txOutputs
}
