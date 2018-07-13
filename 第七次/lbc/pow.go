package lbc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"encoding/hex"
)

var (
	maxXHQ_Nonce = math.MaxInt64
)

const targetBits = 12

// 工作量证明的结构体
type XHQ_ProofOfWork struct {
	block  *Block
	target *big.Int
}

//整形转二进制，
func XHQ_Int2Hex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

//将区块中的各数据连接成二进制串，
func (pows *XHQ_ProofOfWork) prepareData(nonce int) []byte {
	binData := bytes.Join(
		[][]byte{
			pows.block.XHQ_PrevBlockHash,
			pows.block.XHQ_XHQ_HashTransactions(),
			XHQ_Int2Hex(pows.block.XHQ_Timestamp),
			XHQ_Int2Hex(int64(targetBits)),
			XHQ_Int2Hex(int64(nonce)),
		},
		[]byte{},
	)

	return binData
}

// 新建一个工作量证明的函数。
func NewXHQ_ProofOfWork(b *Block) *XHQ_ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &XHQ_ProofOfWork{b, target}

	return pow
}

//执行工作量证明的函数。
func (pow *XHQ_ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	//fmt.Printf("对含有 \"%s\" 信息的区块进行挖矿。\n", pow.block.XHQ_XHQ_HashTransactions())
	fmt.Printf("对含有 \"%s\" 信息的区块进行挖矿。\n", hex.EncodeToString(pow.block.XHQ_XHQ_HashTransactions()))

	for nonce < maxXHQ_Nonce {

		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// 检查是否小于指定的难度值。
func (pow *XHQ_ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.XHQ_Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
