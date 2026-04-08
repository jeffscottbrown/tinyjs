package compiler_test

import (
	"os/exec"
	"strings"
	"testing"

	"example.com/tinyjs/internal/ast"
	"example.com/tinyjs/internal/compiler"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CompilerSuite struct {
	suite.Suite
	compiler *compiler.Compiler
}

func TestCompilerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CompilerSuite))
}

func (s *CompilerSuite) SetupTest() {
	s.compiler = compiler.MustNew()
}

func (s *CompilerSuite) TestCompileProgram_AssignmentThenPrint() {
	varValue := int64(42)
	varName := "x"
	program := &ast.Program{
		Statements: []*ast.Statement{
			{
				Assignment: &ast.Assignment{
					Name: varName,
					Value: &ast.Expression{
						Value: &ast.ValueExpr{
							Int: &varValue,
						},
					},
				},
			}, {
				Print: &ast.Print{
					Arg: &ast.Expression{
						Value: &ast.ValueExpr{
							Ident: &varName,
						},
					},
				},
			},
		},
	}

	ir, err := s.compiler.CompileProgram(program)
	s.Require().NoError(err)

	s.Contains(ir, "define i32 @main()")
	s.Contains(ir, "%x = alloca i64")
	s.Contains(ir, "store i64 42, i64* %x")
	s.Contains(ir, "load i64, i64* %x")
	s.Contains(ir, "@printf")
}

func (s *CompilerSuite) TestCompileString_AssignmentThenPrint() {
	ir, err := s.compiler.CompileString("test.tinyjs", `
x = 42;
print(x);
`)
	s.Require().NoError(err)

	s.Contains(ir, "define i32 @main()")
	s.Contains(ir, "%x = alloca i64")
	s.Contains(ir, "store i64 42, i64* %x")
	s.Contains(ir, "load i64, i64* %x")
	s.Contains(ir, "@printf")
}

func (s *CompilerSuite) TestCompileString_ReassignSameVariable_AllocaOnce() {
	ir, err := s.compiler.CompileString("test.tinyjs", `
x = 1;
x = 2;
print(x);
`)
	s.Require().NoError(err)

	s.Equal(1, strings.Count(ir, "%x = alloca i64"))
	s.Contains(ir, "store i64 1, i64* %x")
	s.Contains(ir, "store i64 2, i64* %x")
}

func (s *CompilerSuite) TestRunString_AssignmentThenPrint() {
	s.requireLLI()

	out, err := s.compiler.RunString("test.tinyjs", `
x = 42;
print(x);
`)
	s.Require().NoError(err)
	s.Equal("42\n", normalizeNewlines(out))
}

func (s *CompilerSuite) TestRunString_ReassignmentPrintsLatestValue() {
	s.requireLLI()

	out, err := s.compiler.RunString("test.tinyjs", `
x = 1;
x = 99;
print(x);
`)
	s.Require().NoError(err)
	s.Equal("99\n", normalizeNewlines(out))
}

func (s *CompilerSuite) TestRunString_MultiplePrints() {
	s.requireLLI()

	out, err := s.compiler.RunString("test.tinyjs", `
x = 5;
print(x);
y = 8;
print(y);
`)
	s.Require().NoError(err)
	s.Equal("5\n8\n", normalizeNewlines(out))
}

func (s *CompilerSuite) TestRunString_PrintUndefinedVariable_ReturnsError() {
	s.requireLLI()

	out, err := s.compiler.RunString("test.tinyjs", `
print(x);
`)
	s.Require().Error(err)
	s.Empty(out)
	s.Contains(err.Error(), `compile program: compile print argument: reference to undefined variable "x"`)
}

func (s *CompilerSuite) TestCompileString_InvalidSource_ReturnsError() {
	ir, err := s.compiler.CompileString("test.tinyjs", `
x = ;
print(x);
`)
	s.Require().Error(err)
	s.Empty(ir)
}

func (s *CompilerSuite) requireLLI() {
	s.T().Helper()

	_, err := exec.LookPath("lli")
	if err != nil {
		s.T().Skip("lli not found in PATH; skipping execution test")
	}
}

func normalizeNewlines(v string) string {
	return strings.ReplaceAll(v, "\r\n", "\n")
}

func TestNormalizeNewlines(t *testing.T) {
	t.Parallel()
	require.Equal(t, "a\nb\n", normalizeNewlines("a\r\nb\r\n"))
}
