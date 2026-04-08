package runner_test

import (
	"strings"
	"testing"

	"example.com/tinyjs/internal/compiler"
	"example.com/tinyjs/internal/parser"
	"github.com/stretchr/testify/require"
)

func TestGenerateIRFromParsedSource_AssignmentThenPrint(t *testing.T) {
	t.Parallel()

	p := parser.MustNew()

	program, err := p.ParseString("test.tinyjs", `
x = 42;
print(x);
`)
	require.NoError(t, err)

	compiler := compiler.MustNew()

	ir, err := compiler.GenerateIR(program)
	require.NoError(t, err)

	require.Contains(t, ir, "define i32 @main()")
	require.Contains(t, ir, "%x = alloca i64")
	require.Contains(t, ir, "store i64 42, i64* %x")
	require.Contains(t, ir, "load i64, i64* %x")
	require.Contains(t, ir, "@printf")
}

func TestGenerateIRFromParsedSource_ReassignSameVariable_AllocaOnce(t *testing.T) {
	t.Parallel()

	p := parser.MustNew()

	program, err := p.ParseString("test.tinyjs", `
x = 1;
x = 2;
print(x);
`)
	require.NoError(t, err)

	compiler := compiler.MustNew()

	ir, err := compiler.GenerateIR(program)
	require.NoError(t, err)

	require.Equal(t, 1, strings.Count(ir, "%x = alloca i64"))
	require.Contains(t, ir, "store i64 1, i64* %x")
	require.Contains(t, ir, "store i64 2, i64* %x")
}

func TestGenerateIRFromParsedSource_PrintUndefinedVariable_ReturnsError(t *testing.T) {
	t.Parallel()

	p := parser.MustNew()

	program, err := p.ParseString("test.tinyjs", `
print(missing);
`)
	require.NoError(t, err)

	compiler := compiler.MustNew()
	ir, err := compiler.GenerateIR(program)
	require.Error(t, err)
	require.Empty(t, ir)
	require.Contains(t, err.Error(), `compile print argument: reference to undefined variable "missing"`)
}
