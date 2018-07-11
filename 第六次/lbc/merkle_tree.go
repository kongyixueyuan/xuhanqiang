package lbc

import (
	"crypto/sha256"
)

type XHQ_MerkleTree struct {
	RootNode *XHQ_MerkleNode
}

// Block  [tx1 tx2 tx3 tx3]

//XHQ_MerkleNode{nil,nil,tx1Bytes}
//XHQ_MerkleNode{nil,nil,tx2Bytes}
//XHQ_MerkleNode{nil,nil,tx3Bytes}
//XHQ_MerkleNode{nil,nil,tx3Bytes}
//
//

//
//XHQ_MerkleNode:
//	left: XHQ_MerkleNode{XHQ_MerkleNode{nil,nil,tx1Bytes},XHQ_MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
//
//	right: XHQ_MerkleNode{XHQ_MerkleNode{nil,nil,tx3Bytes},XHQ_MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
//
//	sha256(sha256(tx1Bytes,tx2Bytes)+sha256(tx3Bytes,tx3Bytes))

type XHQ_MerkleNode struct {
	Left  *XHQ_MerkleNode
	Right *XHQ_MerkleNode
	Data  []byte
}

func NewXHQ_MerkleTree(data [][]byte) *XHQ_MerkleTree {

	//[tx1,tx2,tx3]

	var nodes []XHQ_MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
		//[tx1,tx2,tx3,tx3]
	}

	// 创建叶子节点
	for _, datum := range data {
		node := NewXHQ_MerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	//XHQ_MerkleNode{nil,nil,tx1Bytes}
	//XHQ_MerkleNode{nil,nil,tx2Bytes}
	//XHQ_MerkleNode{nil,nil,tx3Bytes}
	//XHQ_MerkleNode{nil,nil,tx3Bytes}

	// 　循环两次
	for i := 0; i < len(data)/2; i++ {

		var newLevel []XHQ_MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewXHQ_MerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		//XHQ_MerkleNode{XHQ_MerkleNode{nil,nil,tx1Bytes},XHQ_MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
		//
		//XHQ_MerkleNode{XHQ_MerkleNode{nil,nil,tx3Bytes},XHQ_MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
		//

		nodes = newLevel
	}

	//XHQ_MerkleNode:
	//	left: XHQ_MerkleNode{XHQ_MerkleNode{nil,nil,tx1Bytes},XHQ_MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
	//
	//	right: XHQ_MerkleNode{XHQ_MerkleNode{nil,nil,tx3Bytes},XHQ_MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
	//
	//	sha256(sha256(tx1Bytes,tx2Bytes)+sha256(tx3Bytes,tx3Bytes))

	mTree := XHQ_MerkleTree{&nodes[0]}

	return &mTree
}

func NewXHQ_MerkleNode(left, right *XHQ_MerkleNode, data []byte) *XHQ_MerkleNode {
	mNode := XHQ_MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right

	return &mNode
}
