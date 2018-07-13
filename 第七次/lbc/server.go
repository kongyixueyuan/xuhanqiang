package lbc

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

//localhost:3000 主节点的地址
//var knowNodes = []string{"localhost:3000"}
//var nodeAddress string //全局变量，节点地址

func XHQ_startServer(nodeID string, minerAdd string) {

	// 当前节点的IP地址
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)

	ln, err := net.Listen(PROTOCOL, nodeAddress)

	if err != nil {
		log.Panic(err)
	}

	defer ln.Close()

	bc := XHQ_BlockchainObject(nodeID)

	// 第一个终端：端口为3000,启动的就是主节点
	// 第二个终端：端口为3001，钱包节点
	// 第三个终端：端口号为3002，矿工节点
	if nodeAddress != knowNodes[0] {
		// 此节点是钱包节点或者矿工节点，需要向主节点发送请求同步数据

		XHQ_sendVersion(knowNodes[0], bc)
	}

	for {
		// 收到的数据的格式是固定的，12字节+结构体字节数组

		// 接收客户端发送过来的数据
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}

		go XHQ_handleConnection(conn, bc)

	}

}

func XHQ_handleConnection(conn net.Conn, bc *Blockchain) {
	// 读取客户端发送过来的所有的数据
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Receive a Message:%s\n", request[:COMMANDLENGTH])

	command := bytesToCommand(request[:COMMANDLENGTH])

	switch command {
	case COMMAND_VERSION:
		XHQ_handleVersion(request, bc)
	case COMMAND_ADDR:
		XHQ_handleAddr(request, bc)
	case COMMAND_BLOCK:
		XHQ_handleBlock(request, bc)
	case COMMAND_GETBLOCKS:
		XHQ_handleGetblocks(request, bc)
	case COMMAND_GETDATA:
		XHQ_handleGetData(request, bc)
	case COMMAND_INV:
		XHQ_handleInv(request, bc)
	case COMMAND_TX:
		XHQ_handleTx(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	defer conn.Close()
}
