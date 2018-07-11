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

func (cli *CLI) XHQ_printUsage() {

	fmt.Println("使用说明：")
	fmt.Println("  addresslists -- 打印所有钱包地址.")
	fmt.Println("  createwallet -- 创建 钱包.")
	fmt.Println("  addblock -address - 对区块链添加一个区块")
	fmt.Println("  printchain - 打印区块链的所有区块信息")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -- 发送币.")
	fmt.Println("  version - 查看软件版本")
	fmt.Println("  getbalance -address -- 查看帐号余额.")

	fmt.Println("  resetutxo -address -- 重设utxo测试.")

}

func (cli *CLI) XHQ_validateArgs() {
	if len(os.Args) < 2 {
		cli.XHQ_printUsage()
		os.Exit(1)
	}
}

/*func (cli *CLI) addBlock(data string) {
	cli.Bc.XHQ_AddBlock(data)
	fmt.Println("Success!")
}*/

func (cli *CLI) printChain() {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := XHQ_BlockchainObject()

	defer blockchain.Db.Close()

	blockchain.XHQ_Printchain()

}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.XHQ_validateArgs()

	addresslistsCmd := flag.NewFlagSet("addresslists", flag.ExitOnError)
	XHQ_createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)

	addBlockData := addBlockCmd.String("address", "", "初始化地址")

	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账人......")
	flagTo := sendBlockCmd.String("to", "", "转账对像......")
	flagAmount := sendBlockCmd.String("amount", "", "金额......")

	//	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address","","创建创世区块的地址")
	getbalanceWithAdress := getbalanceCmd.String("address", "", "要查询某一个账号的余额.......")

	resetutxoCmd := flag.NewFlagSet("resetutxo", flag.ExitOnError)
	//XHQ_resetUtxo := resetutxoCmd.String("resetutxo", "", "resetutxo.......")

	switch os.Args[1] {

	case "addresslists":
		err := addresslistsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "createwallet":
		err := XHQ_createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

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
		println("V1.0.5")
		break

	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		break

	case "resetutxo":
		err := resetutxoCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		break

	default:
		cli.XHQ_printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		/*
			if *addBlockData == "" {
				addBlockCmd.Usage()
				os.Exit(1)
			}
		*/

		if XHQ_IsValidForAdress([]byte(*addBlockData)) == false {
			fmt.Println("地址无效....")
			cli.XHQ_printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getbalanceCmd.Parsed() {

		if *getbalanceWithAdress == "" {
			fmt.Println("地址不能为空....")
			cli.XHQ_printUsage()
			os.Exit(1)
		}

		cli.XHQ_getBalance(*getbalanceWithAdress)
	}

	if sendBlockCmd.Parsed() {
		cli.checkArgs("from", flagFrom)
		cli.checkArgs("to", flagTo)
		cli.checkArgs("amount", flagAmount)
		//cli.send(flagFrom, flagTo, flagAmount)

		from := XHQ_JSONToArray(*flagFrom)
		to := XHQ_JSONToArray(*flagTo)
		amount := XHQ_JSONToArray(*flagAmount)
		cli.send(from, to, amount)

	}

	if addresslistsCmd.Parsed() {
		cli.XHQ_addressLists()
	}

	if XHQ_createWalletCmd.Parsed() {
		// 创建钱包
		cli.XHQ_createWallet()
	}

	if resetutxoCmd.Parsed() {
		// 创建钱包
		cli.XHQ_resetUtxo()
	}

	//the end of run
}

func (cli *CLI) checkArgs(flag string, arg *string) {

	if *arg == "" {
		fmt.Println("参数" + flag + "不能为空....")
		cli.XHQ_printUsage()
		os.Exit(1)
	}

}

// 先用它去查询余额
func (cli *CLI) XHQ_getBalance(address string) {

	fmt.Println("地址：" + address)

	blockchain := XHQ_BlockchainObject()

	defer blockchain.Db.Close()

//	amount := blockchain.XHQ_GetBalance(address)

	utxoSet := &XHQ_UTXOSet{blockchain}
	amount := utxoSet.XHQ_GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n", address, amount)

}
/*

// 转账中间函数。
func (cli *CLI) send(from []string, to []string, amount []string) {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := XHQ_BlockchainObject()
	defer blockchain.Db.Close()

	blockchain.XHQ_MineNewBlock(from, to, amount)

}
*/

func (cli *CLI) send(from []string,to []string,amount []string)  {


	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := XHQ_BlockchainObject()
	defer blockchain.Db.Close()

	blockchain.XHQ_MineNewBlock(from,to,amount)

	//utxoSet := &XHQ_UTXOSet{blockchain}

	//转账成功以后，需要更新一下
	//utxoSet.Update()

}

// json to array
func XHQ_JSONToArray(jsonString string) []string {

	//json 到 []string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}

// 打印所有的钱包地址
func (cli *CLI) XHQ_addressLists() {

	fmt.Println("打印所有的钱包地址:")

	wallets, _ := XHQ_NewWallets()

	for address, _ := range wallets.WalletsMap {

		fmt.Println(address)
	}
}

func (cli *CLI) createGenesisBlockchain(address string) {
	blockchain := XHQ_CreateBlockchainWithGenesisBlock(address)
	defer blockchain.Db.Close()
}

func (cli *CLI) XHQ_createWallet() {
	wallets, _ := XHQ_NewWallets()

	wallets.XHQ_CreateNewWallet()

	fmt.Println(len(wallets.WalletsMap))
}

func (cli *CLI) XHQ_resetUtxo() {
	fmt.Println("XHQ_resetUtxo:")
	blockchain := XHQ_BlockchainObject()
	defer blockchain.Db.Close()
	utxoSet := &XHQ_UTXOSet{blockchain}
	utxoSet.ResetXHQ_UTXOSet()
	//fmt.Println(blockchain.XHQ_FindUTXOMap())
}
