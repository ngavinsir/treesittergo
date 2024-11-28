package main

import (
	"context"
	_ "embed"
	"log"

	"github.com/ngavinsir/treesittergo"
)

//go:embed ts-sql.wasm
var sqlLanguageWasm []byte

func main() {
	ctx := context.Background()
	p := treesittergo.NewParserCtx(ctx)
	// defer p.CloseCtx(ctx)

	// p.Delete()

	initSQLLangWasm := p.WASMModule.ExportedFunction("tree_sitter_sql")
	sqlLangPtr, err := initSQLLangWasm.Call(ctx)
	if err != nil {
		panic(err)
	}

	sqlLang := treesittergo.NewLanguage(sqlLangPtr[0], p.WASMModule)
	v, err := p.GetLanguageVersionCtx(ctx, sqlLang)
	if err != nil {
		panic(err)
	}
	log.Printf("sql lang version: %+v\n", v)
	log.Printf("lang name: %+v\n", sqlLang.Name())

	err = p.SetLanguageCtx(ctx, sqlLang)
	if err != nil {
		panic(err)
	}

	tree, err := p.ParseStringCtx(ctx, "select 1;")
	if err != nil {
		panic(err)
	}
	root := tree.RootNodeCtx(ctx)
	log.Printf("root node kind: %+v\n", root.KindCtx(ctx))
	log.Printf("root node string: %+v\n", root.StringCtx(ctx))
	childCount := root.ChildCountCtx(ctx)
	log.Printf("root node child count: %+v\n", childCount)
	child1 := root.ChildCtx(ctx, 0)
	log.Printf("child 1 kind: %+v\n", child1.KindCtx(ctx))
	log.Printf("child 1 string: %+v\n", child1.StringCtx(ctx))
	child1Count := child1.ChildCountCtx(ctx)
	log.Printf("child 1 child count: %+v\n", child1Count)
	child1child1 := child1.ChildCtx(ctx, 0)
	log.Printf("child 1 child 1 kind: %+v\n", child1child1.KindCtx(ctx))
	log.Printf("child 1 child 1 string: %+v\n", child1child1.StringCtx(ctx))
	child2 := root.ChildCtx(ctx, 1)
	log.Printf("child 2 kind: %+v\n", child2.KindCtx(ctx))
	log.Printf("child 2 string: %+v\n", child2.StringCtx(ctx))
}
