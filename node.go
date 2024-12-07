package treesittergo

import (
	"context"
	"fmt"
)

type Node struct {
	t Treesitter
	n uint64
}

func newNode(t Treesitter, n uint64) Node {
	return Node{t, n}
}

func (n Node) Kind(ctx context.Context) (string, error) {
	nodeTypeStrPtr, err := n.t.nodeType.Call(ctx, n.n)
	if err != nil {
		return "", fmt.Errorf("getting node type: %w", err)
	}
	return n.t.readString(ctx, nodeTypeStrPtr[0])
}

func (n Node) Child(ctx context.Context, index uint64) (Node, error) {
	// allocate tsnode 24 bytes
	nodePtr, err := n.t.malloc.Call(ctx, uint64(24))
	if err != nil {
		return Node{}, fmt.Errorf("allocating node: %w", err)
	}

	_, err = n.t.nodeChild.Call(ctx, nodePtr[0], n.n, index)
	if err != nil {
		return Node{}, fmt.Errorf("getting node child: %w", err)
	}
	return newNode(n.t, nodePtr[0]), nil
}

func (n Node) ChildCount(ctx context.Context) (uint64, error) {
	res, err := n.t.nodeChildCount.Call(ctx, n.n)
	if err != nil {
		return 0, fmt.Errorf("getting node child count: %w", err)
	}
	return res[0], nil
}

func (n Node) String(ctx context.Context) (string, error) {
	strPtr, err := n.t.nodeString.Call(ctx, n.n)
	if err != nil {
		return "", fmt.Errorf("getting node string: %w", err)
	}
	return n.t.readString(ctx, strPtr[0])
}
