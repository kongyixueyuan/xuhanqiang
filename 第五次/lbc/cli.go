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
	fmt.Println("  addresslists -- 打印所有钱包地址.")
	fmt.Println("  createwallet -- 创建 钱包.")
	fmt.Println("  addblock -address - 对区块链添加一个区块")
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

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	addresslistsCmd := flag.NewFlagSet("addresslists", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)

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

	switch os.Args[1] {

	case "addresslists":
		err := addresslistsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
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
		/*
			if *addBlockData == "" {
				addBlockCmd.Usage()
				os.Exit(1)
			}
		*/

		if IsValidForAdress([]byte(*addBlockData)) == false {
			fmt.Println("地址无效....")
			cli.printUsage()
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
		cli.send(from, to, amount)

	}

	if addresslistsCmd.Parsed() {
		cli.addressLists()
	}

	if createWalletCmd.Parsed() {
		// 创建钱包
		cli.createWallet()
	}

	//the end of run
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
func (cli *CLI) send(from []string, to []string, amount []string) {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.Db.Close()

	blockchain.MineNewBlock(from, to, amount)

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

// 打印所有的钱包地址
func (cli *CLI) addressLists() {

	fmt.Println("打印所有的钱包地址:")

	wallets, _ := NewWallets()

	for address, _ := range wallets.WalletsMap {

		fmt.Println(address)
	}
}

func (cli *CLI) createGenesisBlockchain(address string) {
	blockchain := CreateBlockchainWithGenesisBlock(address)
	defer blockchain.Db.Close()
}

func (cli *CLI) createWallet() {
	wallets, _ := NewWallets()

	wallets.CreateNewWallet()

	fmt.Println(len(wallets.WalletsMap))
}
