package lbc

//"version"

//12

type Version struct {
	Version    int64  // 版本
	BestHeight int64  // 当前节点区块的高度
	AddrFrom   string //当前节点的地址
}

type GetBlocks struct {
	AddrFrom string
}

type Inv struct {
	AddrFrom string   //自己的地址
	Type     string   //类型 block tx
	Items    [][]byte //hash二维数组
}

type GetData struct {
	AddrFrom string
	Type     string
	Hash     []byte
}

type BlockData struct {
	AddrFrom string
	Block    *Block
}
