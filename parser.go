package treesittergo

import (
	"context"
	_ "embed"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/emscripten"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed ts-combined-sql.wasm
var tsWasm []byte

type Parser struct {
	m          api.Module
	p          uint64
	Runtime    wazero.Runtime
	WASMModule api.Module
}

// Create a new parser.
func NewParser() *Parser {
	return NewParserCtx(context.Background())
}

func NewParserCtx(ctx context.Context) *Parser {
	r := wazero.NewRuntime(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	compiled, err := r.CompileModule(ctx, tsWasm)
	if err != nil {
		panic(err)
	}

	_, err = emscripten.InstantiateForModule(ctx, r, compiled)
	if err != nil {
		panic(err)
	}

	mod, err := r.InstantiateModule(ctx, compiled, wazero.NewModuleConfig().WithName("ts"))
	if err != nil {
		panic(err)
	}

	newParser := mod.ExportedFunction("ts_parser_new")
	p, err := newParser.Call(ctx)
	if err != nil {
		panic(err)
	}

	return &Parser{
		m:          mod,
		p:          p[0],
		Runtime:    r,
		WASMModule: mod,
	}
}

func (p *Parser) Close() {
	p.CloseCtx(context.Background())
}

func (p *Parser) Delete() {
	delete := p.WASMModule.ExportedFunction("ts_parser_delete")
	_, err := delete.Call(context.Background(), p.p)
	if err != nil {
		panic(err)
	}
}

func (p *Parser) CloseCtx(ctx context.Context) {
	deleteParser := p.m.ExportedFunction("ts_parser_delete")
	_, err := deleteParser.Call(ctx, p.p)
	if err != nil {
		panic(err)
	}
	p.Runtime.Close(ctx)
}

func (p *Parser) SetLanguageCtx(ctx context.Context, l *Language) error {
	setLanguage := p.m.ExportedFunction("ts_parser_set_language")
	ok, err := setLanguage.Call(ctx, p.p, l.l)
	if err != nil {
		panic(err)
	}
	if ok[0] == 0 {
		v, err := p.GetLanguageVersionCtx(ctx, l)
		if err != nil {
			panic(err)
		}
		return LanguageError{v}
	}

	return nil
}

func (p *Parser) GetLanguageVersionCtx(ctx context.Context, l *Language) (uint64, error) {
	languageVersion := p.m.ExportedFunction("ts_language_version")
	v, err := languageVersion.Call(ctx, l.l)
	if err != nil {
		panic(err)
	}
	return v[0], nil
}

func (p *Parser) ParseStringCtx(ctx context.Context, str string) (*Tree, error) {
	malloc := p.m.ExportedFunction("malloc")
	free := p.m.ExportedFunction("free")
	parseString := p.m.ExportedFunction("ts_parser_parse_string")

	strByte := []byte(str)
	strSize := uint64(len(strByte))
	strPtr, err := malloc.Call(ctx, strSize)
	if err != nil {
		return nil, err
	}
	defer free.Call(ctx, strPtr[0])

	if !p.m.Memory().Write(uint32(strPtr[0]), strByte) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
			strPtr[0], strSize, p.m.Memory().Size())
	}

	tree, err := parseString.Call(ctx, p.p, uint64(0), strPtr[0], strSize)
	if err != nil {
		panic(err)
	}
	return newTree(tree[0], p.m), nil
}

// func main() {
// 	log.Printf("exported functions: %+v\n", mod.ExportedFunctionDefinitions())
//
// 	parserTimeoutMicros := mod.ExportedFunction("ts_parser_timeout_micros")
// 	t, err := parserTimeoutMicros.Call(ctx, p[0])
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("t: %+v\n", t)
//
// 	setParserTimeoutMicros := mod.ExportedFunction("ts_parser_set_timeout_micros")
// 	_, err = setParserTimeoutMicros.Call(ctx, p[0], uint64(1232))
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	t, err = parserTimeoutMicros.Call(ctx, p[0])
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("t2: %+v\n", t)
// }
