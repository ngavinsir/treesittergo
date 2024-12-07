package treesittergo

import (
	"context"
	"fmt"
)

type Query struct {
	t Treesitter
	q uint64
}

func (t Treesitter) NewQuery(ctx context.Context, pattern string, l Language) (Query, error) {
	patternPtr, patternSize, freePattern, err := t.allocateString(ctx, pattern)
	if err != nil {
		return Query{}, fmt.Errorf("allocating pattern string: %w", err)
	}
	defer freePattern()
	queryPtr, err := t.queryNew.Call(ctx, l.l, patternPtr, patternSize, 0, 0)
	if err != nil {
		return Query{}, fmt.Errorf("creating query: %w", err)
	}
	return Query{t, queryPtr[0]}, nil
}
