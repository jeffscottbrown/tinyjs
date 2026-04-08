package ast

import "github.com/alecthomas/participle/v2/lexer"

type Program struct {
	Statements []*Statement `@@*`
}

type Statement struct {
	Pos        lexer.Position
	Assignment *Assignment `  @@`
	Print      *Print      `| @@`
}

type Assignment struct {
	Pos   lexer.Position
	Name  string      `@Ident`
	Value *Expression `"=" @@ ";"`
}

type Print struct {
	Pos lexer.Position
	Arg *Expression `"print" "(" @@ ")" ";"`
}

type Expression struct {
	Pos    lexer.Position
	Binary *BinaryExpr `  @@`
	Value  *ValueExpr  `| @@`
}

type BinaryExpr struct {
	Pos   lexer.Position
	Left  *ValueExpr `@@`
	Op    string     `@("+" | "-" | "*" | "/")`
	Right *ValueExpr `@@`
}

type ValueExpr struct {
	Pos   lexer.Position
	Int   *int64  `  @Int`
	Ident *string `| @Ident`
}
