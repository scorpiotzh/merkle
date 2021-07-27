package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

func main() {
	hashList := []string{
		"bb792343a07430ee5eda4562b79ca1bfd321f897327b1c0b98bcd63438de5133",
		"0c8bad290fb41cdf0a87b9c835db2cbb6b286b8ab0f08248392aa5e5bc283ae9",
		"7cbcf4a6bb4b896fa3d058152e87ff0a0046ada1e058379468a09781a87e72b4",
		"f4bcbb7227712decf5ae611bf0c4dc53fef679a06233d2fa5618c4f1cb3de17f",
		"af57028f1eb27e75b91d30e90edd6c90280e14c27bae76880972edd88a83f3e6",
		"b01245151b1918f5c0b39edfd4b392a29e37e5a5409d9f993162c20dd4f08f64",
		"ef853016fe848084ffc4785f8c3217648ec4ad17abfc18301bbdc19b13cdcad3",
		"6db31b7fa35a741a27439205d3b9990dfa3d805d2ace3154a152a3ef59babb7b",
		"a2a489288e11f6da6a3b376be9dca75a0c8c3684fcea8003f3f17f90a52d5d54",
		"d1f834846766e619c6cb9a7ae4abaa98dd8a1e919c23c529d8a59d7b3421a77d",
		"6565cb0e22d5ba9e565ceb9fcc1952e8e46d60de937b22db0fd92acd9667aa32",
	}
	sort.Strings(hashList)
	var dataList [][]byte
	for _, v := range hashList {
		data, _ := hex.DecodeString(v)
		dataList = append(dataList, data)
	}
	tree := NewMerkleTree(dataList)
	fmt.Println("root:", hex.EncodeToString(tree.Node.Data), tree.Node.Flag)
	//
	var res = make(map[string]*MerkleNode)
	tree.Node.Traverse(res)
	fmt.Println(res)
	//
	for _, v := range tree.LeafNodes {
		fmt.Println(v.Flag)
	}
	//
	tree.MerklePath()
}

// 默克尔树
type MerkleTree struct {
	Node      *MerkleNode
	LeafNodes []MerkleNode
}

// 默克尔节点
type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNode *MerkleNode
	Data      []byte //保存当前节点哈希
	Flag      string
}

func NewMerkleTree(dataList [][]byte) *MerkleTree {
	//包含整个树上的节点的容器
	var nodes []MerkleNode
	//生成所有的叶子节点
	for i, data := range dataList {
		node := NewMerkleNode(nil, nil, data)
		node.Flag = fmt.Sprintf(",%d,", i)
		nodes = append(nodes, *node)
	}
	j := 0
	//生成分支节点
	for nSize := len(dataList); nSize > 1; nSize = (nSize + 1) / 2 {
		//进行两两分组
		//i是左侧分支节点的索引。因为两个一组哈希 所以 i+=2
		for i := 0; i < nSize; i += 2 {
			//ii是跟i配套，凑成一组右侧分支节点的索引
			ii := min(i+1, nSize-1)
			node := NewMerkleNode(&nodes[j+i], &nodes[j+ii], nil)
			nodes = append(nodes, *node)
		}
		j += nSize
	}
	return &MerkleTree{&(nodes[len(nodes)-1]), nodes[:len(dataList)]}
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := new(MerkleNode)
	mNode.LeftNode = left
	mNode.RightNode = right
	//叶子节点
	if left != nil && right != nil {
		//	对左右两侧分支节点的hash进行双哈希
		hashValue := append(left.Data, right.Data...)
		hashData := sha256.Sum256(hashValue)
		mNode.Data = hashData[:]
		mNode.Flag = left.Flag + right.Flag
	} else {
		mNode.Data = data
	}
	return mNode
}

// 遍历
func (m *MerkleNode) Traverse(res map[string]*MerkleNode) {
	res[m.Flag] = m
	if m.LeftNode != nil {
		m.LeftNode.Traverse(res)
	}
	if m.RightNode != nil {
		m.RightNode.Traverse(res)
	}
}

// 默克尔路径
func (m *MerkleTree) MerklePath() {
	var res = make(map[string]*MerkleNode)
	m.Node.Traverse(res)
	for _, v := range m.LeafNodes {
		root := m.Node
		fmt.Printf("[%s]的路径：\n", v.Flag)
		var strList []string
		for root.LeftNode != nil && root.RightNode != nil {
			left := root.LeftNode
			right := root.RightNode
			if strings.Contains(left.Flag, v.Flag) {
				root = left
				strList = append(strList, "[%s,\""+hex.EncodeToString(res[right.Flag].Data)+"\"]")
			} else {
				root = right
				strList = append(strList, "[\""+hex.EncodeToString(res[left.Flag].Data)+"\",%s]")
			}
		}
		str := "\"" + hex.EncodeToString(res[root.Flag].Data) + "\""
		for i := len(strList); i > 0; i-- {
			str = fmt.Sprintf(strList[i-1], str)
		}
		fmt.Println(str)
	}
}
