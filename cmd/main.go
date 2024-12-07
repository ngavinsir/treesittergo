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
	ts, err := treesittergo.New(ctx)
	if err != nil {
		panic(err)
	}
	p, err := ts.NewParser(ctx)
	if err != nil {
		panic(err)
	}
	// defer p.Close(ctx)

	// p.Delete()

	sqlLang, err := ts.LanguageSQL(ctx)
	if err != nil {
		panic(err)
	}
	v, err := p.GetLanguageVersion(ctx, sqlLang)
	if err != nil {
		panic(err)
	}
	log.Printf("sql lang version: %+v\n", v)
	sqlLangName, err := sqlLang.Name(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("lang name: %+v\n", sqlLangName)

	err = p.SetLanguage(ctx, sqlLang)
	if err != nil {
		panic(err)
	}

	tree, err := p.ParseString(ctx, "select 1;")
	if err != nil {
		panic(err)
	}
	root, err := tree.RootNode(ctx)
	if err != nil {
		panic(err)
	}
	rootKind, err := root.Kind(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("root node kind: %+v\n", rootKind)
	rootString, err := root.String(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("root node string: %+v\n", rootString)
	rootChildCount, err := root.ChildCount(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("root node child count: %+v\n", rootChildCount)
	child1, err := root.Child(ctx, 0)
	if err != nil {
		panic(err)
	}
	child1Kind, err := child1.Kind(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 1 kind: %+v\n", child1Kind)
	child1String, err := child1.String(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 1 string: %+v\n", child1String)
	child1Count, err := child1.ChildCount(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 1 child count: %+v\n", child1Count)
	child1child1, err := child1.Child(ctx, 0)
	if err != nil {
		panic(err)
	}
	child1child1Kind, err := child1child1.Kind(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 1 child 1 kind: %+v\n", child1child1Kind)
	child1child1String, err := child1child1.String(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 1 child 1 string: %+v\n", child1child1String)
	child2, err := root.Child(ctx, 1)
	if err != nil {
		panic(err)
	}
	child2Kind, err := child2.Kind(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 2 kind: %+v\n", child2Kind)
	child2String, err := child2.String(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 2 string: %+v\n", child2String)
}
