package parser

import (
	"fmt"

	"example.com/tinyjs/internal/ast"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var tinyLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Comment", Pattern: `//[^\n\r]*`},
	{Name: "Whitespace", Pattern: `[ \t\r\n]+`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Int", Pattern: `\d+`},
	{Name: "Punct", Pattern: `[=;]`},
})

var tinyParser = participle.MustBuild[ast.Program](
	participle.Lexer(tinyLexer),
	participle.Elide("Whitespace", "Comment"),
)

// Parse parses the full source into an AST.
func Parse(src string) (*ast.Program, error) {
	prog, err := tinyParser.ParseString("", src)
	if err != nil {
		return nil, fmt.Errorf("parse source: %w", err)
	}
	return prog, nil
}
