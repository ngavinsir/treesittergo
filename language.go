package treesittergo

import (
	"context"
	"fmt"
	"log"

	"github.com/tetratelabs/wazero/api"
)

type (
	Language struct {
		l uint64
		m api.Module
	}

	LanguageError struct {
		version uint64
	}
)

func (l LanguageError) Error() string {
	return fmt.Sprintf("Incompatible language version %d", l.version)
}

func NewLanguage(l uint64, m api.Module) *Language {
	return &Language{l, m}
}

func (l *Language) Name() string {
	langName := l.m.ExportedFunction("ts_language_name")
	langNamePtr, err := langName.Call(context.Background(), l.l)
	if err != nil {
		panic(err)
	}
	strlen := l.m.ExportedFunction("strlen")
	strSize, err := strlen.Call(context.Background(), langNamePtr[0])
	if err != nil {
		panic(err)
	}
	strBytes, ok := l.m.Memory().Read(uint32(langNamePtr[0]), uint32(strSize[0]))
	if !ok {
		log.Panicf("Memory.Read(%d, %d) invalid memory size %d",
			langNamePtr[0], strSize[0], l.m.Memory().Size())
	}
	return string(strBytes)
}
