package compiler

import (
	"fmt"

	"example.com/tinyjs/internal/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (c *Compiler) GenerateIR(program *ast.Program) (string, error) {
	if program == nil {
		return "", fmt.Errorf("program cannot be nil")
	}

	mod := ir.NewModule()

	// @.fmt = private unnamed_addr constant [5 x i8] c"%ld\0A\00"
	fmtText := "%ld\n"
	fmtConst := constant.NewCharArrayFromString(fmtText + "\x00")
	fmtGlobal := mod.NewGlobalDef(".fmt", fmtConst)
	fmtGlobal.Immutable = true
	fmtArrayType := fmtConst.Typ

	// declare i32 @printf(i8*, ...)
	printf := mod.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true

	// define i32 @main()
	mainFn := mod.NewFunc("main", types.I32)
	entry := mainFn.NewBlock("entry")

	symbols := map[string]*ir.InstAlloca{}

	for _, stmt := range program.Statements {
		if stmt == nil {
			return "", fmt.Errorf("program contains nil statement")
		}

		switch {
		case stmt.Assignment != nil:
			s := stmt.Assignment
			if s.Name == "" {
				return "", fmt.Errorf("assignment name cannot be empty")
			}

			ptr, ok := symbols[s.Name]
			if !ok {
				ptr = entry.NewAlloca(types.I64)
				ptr.LocalName = s.Name
				symbols[s.Name] = ptr
			}

			val, err := c.compileExpression(entry, symbols, s.Value)
			if err != nil {
				return "", fmt.Errorf("compile assignment value for %q: %w", s.Name, err)
			}

			entry.NewStore(val, ptr)

		case stmt.Print != nil:
			s := stmt.Print
			if s.Arg == nil {
				return "", fmt.Errorf("print argument cannot be nil")
			}

			val, err := c.compileExpression(entry, symbols, s.Arg)
			if err != nil {
				return "", fmt.Errorf("compile print argument: %w", err)
			}

			fmtPtr := entry.NewGetElementPtr(
				fmtArrayType,
				fmtGlobal,
				constant.NewInt(types.I64, 0),
				constant.NewInt(types.I64, 0),
			)

			entry.NewCall(printf, fmtPtr, val)

		default:
			return "", fmt.Errorf("statement has no supported variant set")
		}
	}

	entry.NewRet(constant.NewInt(types.I32, 0))

	return mod.String(), nil
}

func (c *Compiler) compileExpression(
	block *ir.Block,
	symbols map[string]*ir.InstAlloca,
	expr *ast.Expression,
) (value.Value, error) {
	if expr == nil {
		return nil, fmt.Errorf("expression cannot be nil")
	}

	if expr.Value != nil {
		return c.compileValueExpr(block, symbols, expr.Value)
	}

	if expr.Binary != nil {
		left, err := c.compileValueExpr(block, symbols, expr.Binary.Left)
		if err != nil {
			return nil, fmt.Errorf("compile left operand: %w", err)
		}

		right, err := c.compileValueExpr(block, symbols, expr.Binary.Right)
		if err != nil {
			return nil, fmt.Errorf("compile right operand: %w", err)
		}

		switch expr.Binary.Op {
		case "+":
			return block.NewAdd(left, right), nil
		case "-":
			return block.NewSub(left, right), nil
		case "*":
			return block.NewMul(left, right), nil
		case "/":
			return block.NewSDiv(left, right), nil
		default:
			return nil, fmt.Errorf("unsupported operator %q", expr.Binary.Op)
		}
	}

	return nil, fmt.Errorf("expression has no supported variant set")
}

func (c *Compiler) compileValueExpr(
	block *ir.Block,
	symbols map[string]*ir.InstAlloca,
	expr *ast.ValueExpr,
) (value.Value, error) {
	if expr == nil {
		return nil, fmt.Errorf("value expression cannot be nil")
	}

	if expr.Int != nil {
		return constant.NewInt(types.I64, *expr.Int), nil
	}

	if expr.Ident != nil {
		name := *expr.Ident
		ptr, ok := symbols[name]
		if !ok {
			return nil, fmt.Errorf("reference to undefined variable %q", name)
		}
		return block.NewLoad(types.I64, ptr), nil
	}

	return nil, fmt.Errorf("value expression has no supported variant set")
}
