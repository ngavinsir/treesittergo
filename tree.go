package treesittergo

import (
	"context"

	"github.com/tetratelabs/wazero/api"
)

type Tree struct {
	t uint64
	m api.Module
}

// Create a new tree from a raw pointer.
func newTree(t uint64, m api.Module) *Tree {
	return &Tree{t, m}
}

func (t *Tree) RootNodeCtx(ctx context.Context) *Node {
	// allocate tsnode 24 bytes
	malloc := t.m.ExportedFunction("malloc")
	nodePtr, err := malloc.Call(ctx, uint64(24))
	if err != nil {
		panic(err)
	}

	treeRootNode := t.m.ExportedFunction("ts_tree_root_node")
	_, err = treeRootNode.Call(ctx, nodePtr[0], t.t)
	if err != nil {
		panic(err)
	}
	return newNode(nodePtr[0], t.m)
}
