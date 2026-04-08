package compiler_test

import (
	"testing"

	"example.com/tinyjs/internal/ast"
	"example.com/tinyjs/internal/compiler"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CompilerSuite struct {
	suite.Suite
}

func TestCompilerSuite(t *testing.T) {
	suite.Run(t, new(CompilerSuite))
}

func (s *CompilerSuite) TestCompileNilProgramFails() {
	out, err := compiler.Compile(nil)
	s.Empty(out)
	s.Error(err)
	s.ErrorContains(err, "program is nil")
}

func (s *CompilerSuite) TestCompileEmptyProgramProducesHeaderOnly() {
	prog := &ast.Program{}

	out, err := compiler.Compile(prog)
	s.Require().NoError(err)
	s.Equal("; tinyjs output\n; accepted syntax: <ident> = <integer>;\n\n", out)
}

func (s *CompilerSuite) TestCompileSingleAssignment() {
	prog := &ast.Program{
		Assignments: []*ast.Assignment{{Name: "x", Eq: "=", Value: 1, Semi: ";"}},
	}

	out, err := compiler.Compile(prog)
	s.Require().NoError(err)
	s.Equal("; tinyjs output\n; accepted syntax: <ident> = <integer>;\n\nvar x = 1\n", out)
}

func (s *CompilerSuite) TestCompileMultipleAssignmentsPreservesOrder() {
	prog := &ast.Program{
		Assignments: []*ast.Assignment{
			{Name: "x", Eq: "=", Value: 1, Semi: ";"},
			{Name: "y", Eq: "=", Value: 20, Semi: ";"},
			{Name: "z", Eq: "=", Value: 300, Semi: ";"},
		},
	}

	out, err := compiler.Compile(prog)
	s.Require().NoError(err)
	s.Contains(out, "var x = 1\nvar y = 20\nvar z = 300\n")
}

func (s *CompilerSuite) TestCompileRejectsNilAssignment() {
	prog := &ast.Program{Assignments: []*ast.Assignment{nil}}

	out, err := compiler.Compile(prog)
	s.Empty(out)
	s.Error(err)
	s.ErrorContains(err, "assignment is nil")
}

func (s *CompilerSuite) TestCompileRejectsEmptyName() {
	prog := &ast.Program{
		Assignments: []*ast.Assignment{{Name: "", Eq: "=", Value: 1, Semi: ";"}},
	}

	out, err := compiler.Compile(prog)
	s.Empty(out)
	s.Error(err)
	s.ErrorContains(err, "empty variable name")
}

func TestCompile_TableDrivenValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		statement *ast.Assignment
		wantLine  string
	}{
		{name: "zero", statement: &ast.Assignment{Name: "x", Value: 0}, wantLine: "var x = 0\n"},
		{name: "positive", statement: &ast.Assignment{Name: "y", Value: 42}, wantLine: "var y = 42\n"},
		{name: "negative", statement: &ast.Assignment{Name: "z", Value: -7}, wantLine: "var z = -7\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			prog := &ast.Program{Assignments: []*ast.Assignment{tt.statement}}
			out, err := compiler.Compile(prog)
			require.NoError(t, err)
			require.Contains(t, out, tt.wantLine)
		})
	}
}
