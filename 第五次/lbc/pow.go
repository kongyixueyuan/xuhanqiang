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
	maxNonce = math.MaxInt64
)

const targetBits = 16

// 工作量证明的结构体
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

//整形转二进制，
func Int2Hex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

//将区块中的各数据连接成二进制串，
func (pows *ProofOfWork) prepareData(nonce int) []byte {
	binData := bytes.Join(
		[][]byte{
			pows.block.PrevBlockHash,
			pows.block.HashTransactions(),
			Int2Hex(pows.block.Timestamp),
			Int2Hex(int64(targetBits)),
			Int2Hex(int64(nonce)),
		},
		[]byte{},
	)

	return binData
}

// 新建一个工作量证明的函数。
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

//执行工作量证明的函数。
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	//fmt.Printf("对含有 \"%s\" 信息的区块进行挖矿。\n", pow.block.HashTransactions())
	fmt.Printf("对含有 \"%s\" 信息的区块进行挖矿。\n", hex.EncodeToString(pow.block.HashTransactions()))

	for nonce < maxNonce {

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
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
