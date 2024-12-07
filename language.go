package treesittergo

import (
	"context"
	"fmt"
)

type (
	Language struct {
		t Treesitter
		l uint64
	}

	LanguageError struct {
		version uint64
	}
)

func (l LanguageError) Error() string {
	return fmt.Sprintf("Incompatible language version %d", l.version)
}

func NewLanguage(l uint64, t Treesitter) Language {
	return Language{l: l, t: t}
}

func (l Language) Name(ctx context.Context) (string, error) {
	langNamePtr, err := l.t.languageName.Call(context.Background(), l.l)
	if err != nil {
		return "", fmt.Errorf("getting language name: %w", err)
	}
	return l.t.readString(ctx, langNamePtr[0])
}
