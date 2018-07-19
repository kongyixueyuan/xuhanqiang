package lbc

import (
	"bytes"
	"log"
	"encoding/gob"
)

func XHQ_handleVersion(request []byte,bc *Blockchain)  {

	var buff bytes.Buffer
	var payload Version

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	//Version
	//1. Version
	//2. BestHeight
	//3. 节点地址

	bestHeight := bc.GetBestHeight() //3
	foreignerBestHeight := payload.BestHeight // 1

	if bestHeight > foreignerBestHeight {
		XHQ_sendVersion(payload.AddrFrom,bc)
	} else if bestHeight < foreignerBestHeight {
		// 去向主节点要信息
		sendGetBlocks(payload.AddrFrom)
	}


}

func XHQ_handleAddr(request []byte,bc *Blockchain)  {




}

func XHQ_handleGetblocks(request []byte,bc *Blockchain)  {


	var buff bytes.Buffer
	var payload GetBlocks

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()

	//
	sendInv(payload.AddrFrom, BLOCK_TYPE, blocks)


}

func XHQ_handleGetData(request []byte,bc *Blockchain)  {

	var buff bytes.Buffer
	var payload GetData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == BLOCK_TYPE {

		block, err := bc.GetBlock([]byte(payload.Hash))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, block)
	}

	if payload.Type == "tx" {

	}
}

func XHQ_handleBlock(request []byte,bc *Blockchain)  {
	var buff bytes.Buffer
	var payload BlockData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	block := payload.Block

	bc.AddBlock(block)

	if len(transactionArray) == 0{

		utxoSet := &XHQ_UTXOSet{bc}
		utxoSet.ResetXHQ_UTXOSet()
	}


	if len(transactionArray) > 0 {
		XHQ_sendGetData(payload.AddrFrom,BLOCK_TYPE,transactionArray[0])

		transactionArray = transactionArray[1:]
	}

}

func XHQ_handleTx(request []byte,bc *Blockchain)  {

}


func XHQ_handleInv(request []byte,bc *Blockchain)  {

	var buff bytes.Buffer
	var payload Inv

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	// Ivn 3000 block hashes [][]

	if payload.Type == BLOCK_TYPE {

		//tansactionArray = payload.Items

		//payload.Items

		blockHash := payload.Items[0]
		XHQ_sendGetData(payload.AddrFrom, BLOCK_TYPE , blockHash)

		if len(payload.Items) >= 1 {
			transactionArray = payload.Items[1:]
		}
	}

	if payload.Type == TX_TYPE {

	}

}