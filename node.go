package treesittergo

import (
	"context"
	"log"

	"github.com/tetratelabs/wazero/api"
)

type Node struct {
	n uint64
	m api.Module
}

func newNode(n uint64, m api.Module) *Node {
	return &Node{n, m}
}

func (n *Node) KindCtx(ctx context.Context) string {
	nodeType := n.m.ExportedFunction("ts_node_type")
	nodeTypeStrPtr, err := nodeType.Call(ctx, n.n)
	if err != nil {
		panic(err)
	}
	strlen := n.m.ExportedFunction("strlen")
	strSize, err := strlen.Call(ctx, nodeTypeStrPtr[0])
	if err != nil {
		panic(err)
	}
	strBytes, ok := n.m.Memory().Read(uint32(nodeTypeStrPtr[0]), uint32(strSize[0]))
	if !ok {
		log.Panicf("Memory.Read(%d, %d) invalid memory size %d",
			nodeTypeStrPtr[0], strSize[0], n.m.Memory().Size())
	}
	return string(strBytes)
}

func (n *Node) ChildCtx(ctx context.Context, index uint64) *Node {
	// allocate tsnode 24 bytes
	malloc := n.m.ExportedFunction("malloc")
	nodePtr, err := malloc.Call(ctx, uint64(24))
	if err != nil {
		panic(err)
	}

	child := n.m.ExportedFunction("ts_node_child")
	_, err = child.Call(ctx, nodePtr[0], n.n, index)
	if err != nil {
		panic(err)
	}
	return newNode(nodePtr[0], n.m)
}

func (n *Node) ChildCountCtx(ctx context.Context) uint64 {
	childCount := n.m.ExportedFunction("ts_node_child_count")
	res, err := childCount.Call(ctx, n.n)
	if err != nil {
		panic(err)
	}
	return res[0]
}

func (n *Node) StringCtx(ctx context.Context) string {
	str := n.m.ExportedFunction("ts_node_string")
	strPtr, err := str.Call(ctx, n.n)
	if err != nil {
		panic(err)
	}
	strlen := n.m.ExportedFunction("strlen")
	strSize, err := strlen.Call(ctx, strPtr[0])
	if err != nil {
		panic(err)
	}
	strBytes, ok := n.m.Memory().Read(uint32(strPtr[0]), uint32(strSize[0]))
	if !ok {
		log.Panicf("Memory.Read(%d, %d) invalid memory size %d",
			strPtr[0], strSize[0], n.m.Memory().Size())
	}
	return string(strBytes)
}
