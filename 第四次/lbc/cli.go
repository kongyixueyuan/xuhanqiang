package lbc

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

// CLI responsible for processing command line arguments
type CLI struct {
	Bc *Blockchain
}

func (cli *CLI) printUsage() {

	fmt.Println("使用说明：")
	fmt.Println("  addblock -data BLOCK_DATA - 对区块链添加一个区块")
	fmt.Println("  printchain - 打印区块链的所有区块信息")

	fmt.Println("  send -from FROM -to TO -amount AMOUNT -- 发送币.")

	fmt.Println("  version - 查看软件版本")

	fmt.Println("  getbalance -address -- 查看帐号余额.")

}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

/*func (cli *CLI) addBlock(data string) {
	cli.Bc.AddBlock(data)
	fmt.Println("Success!")
}*/

func (cli *CLI) printChain() {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	defer blockchain.Db.Close()

	blockchain.Printchain()

}

//引用数据 库之前的打印方法，
/*
func (cli *CLI) printChain() {
	bci := cli.Bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		//fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

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
				fmt.Println(out.Value)
				fmt.Println(out.ScriptPubKey)
			}
		}

		fmt.Println("------------------------------")

		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

*/

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账人......")
	flagTo := sendBlockCmd.String("to", "", "转账对像......")
	flagAmount := sendBlockCmd.String("amount", "", "金额......")

	//	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address","","创建创世区块的地址")
	getbalanceWithAdress := getbalanceCmd.String("address", "", "要查询某一个账号的余额.......")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		break

	case "version":
		println("V1.0.1")
		break

	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		break

	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}

		println("todo:addBlock implement!")
		//	cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getbalanceCmd.Parsed() {

		if *getbalanceWithAdress == "" {
			fmt.Println("地址不能为空....")
			cli.printUsage()
			os.Exit(1)
		}

		cli.getBalance(*getbalanceWithAdress)
	}

	if sendBlockCmd.Parsed() {
		cli.checkArgs("from", flagFrom)
		cli.checkArgs("to", flagTo)
		cli.checkArgs("amount", flagAmount)
		//cli.send(flagFrom, flagTo, flagAmount)


		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
		cli.send(from,to,amount)

	}

}

func (cli *CLI) checkArgs(flag string, arg *string) {

	if *arg == "" {
		fmt.Println("参数" + flag + "不能为空....")
		cli.printUsage()
		os.Exit(1)
	}

}

// 先用它去查询余额
func (cli *CLI) getBalance(address string) {

	fmt.Println("地址：" + address)

	blockchain := BlockchainObject()

	defer blockchain.Db.Close()

	amount := blockchain.GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n", address, amount)

}





// 转账中间函数。
func (cli *CLI) send(from []string,to []string,amount []string)  {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.Db.Close()

	blockchain.MineNewBlock(from,to,amount)

}




// json to array
func JSONToArray(jsonString string) []string {

	//json 到 []string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}


