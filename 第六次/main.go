package main

import (
	"fmt"

	"./lbc"
)

func main() {

	//mybc := lbc.NewBlockchain("myaddress")
	//defer mybc.Db.Close()

	cli := lbc.CLI{}
	cli.Run()

	fmt.Println("执行结束!")
}
