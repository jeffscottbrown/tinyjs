package parser

import (
	"fmt"

	"example.com/tinyjs/internal/ast"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var tinyLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Keyword", Pattern: `\bprint\b`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Int", Pattern: `[0-9]+`},
	{Name: "Operator", Pattern: `[+\-*/]`},
	{Name: "Punct", Pattern: `[=();]`},
	{Name: "Whitespace", Pattern: `\s+`},
})

type Parser struct {
	parser *participle.Parser[ast.Program]
}

func New() (*Parser, error) {
	p, err := participle.Build[ast.Program](
		participle.Lexer(tinyLexer),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, fmt.Errorf("build parser: %w", err)
	}

	return &Parser{parser: p}, nil
}

func MustNew() *Parser {
	p, err := New()
	if err != nil {
		panic(err)
	}
	return p
}

func (p *Parser) ParseString(filename, input string) (*ast.Program, error) {
	program, err := p.parser.ParseString(filename, input)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", filename, err)
	}
	return program, nil
}
