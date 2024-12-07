package treesittergo

import (
	"context"
	"fmt"
)

func (t Treesitter) LanguageSQL(ctx context.Context) (Language, error) {
	sqlLangPtr, err := t.languageSQL.Call(ctx)
	if err != nil {
		return Language{}, fmt.Errorf("initiating sql language: %w", err)
	}
	return NewLanguage(sqlLangPtr[0], t), nil
}
