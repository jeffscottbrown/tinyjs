package parser_test

import (
	"testing"

	"example.com/tinyjs/internal/parser"
	"github.com/stretchr/testify/require"
)

func TestParseString_AssignmentAndPrint(t *testing.T) {
	t.Parallel()

	p := parser.MustNew()

	program, err := p.ParseString("test.tinyjs", `
x = 42;
print(x);
`)
	require.NoError(t, err)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 2)

	require.NotNil(t, program.Statements[0])
	require.NotNil(t, program.Statements[0].Assignment)
	require.Equal(t, "x", program.Statements[0].Assignment.Name)
	require.EqualValues(t, 42, *program.Statements[0].Assignment.Value.Value.Int)
	require.Nil(t, program.Statements[0].Print)

	require.NotNil(t, program.Statements[1])
	require.NotNil(t, program.Statements[1].Print)
	require.Equal(t, "x", *program.Statements[1].Print.Arg.Value.Ident)
	require.Nil(t, program.Statements[1].Assignment)
}

func TestParseString_MultipleAssignmentsAndPrints(t *testing.T) {
	t.Parallel()

	p := parser.MustNew()

	program, err := p.ParseString("test.tinyjs", `
x = 1;
print(x);
y = 99;
print(y);
`)
	require.NoError(t, err)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 4)

	require.Equal(t, "x", program.Statements[0].Assignment.Name)
	require.EqualValues(t, 1, *program.Statements[0].Assignment.Value.Value.Int)

	require.Equal(t, "x", *program.Statements[1].Print.Arg.Value.Ident)

	require.Equal(t, "y", program.Statements[2].Assignment.Name)
	require.EqualValues(t, 99, *program.Statements[2].Assignment.Value.Value.Int)

	require.Equal(t, "y", *program.Statements[3].Print.Arg.Value.Ident)
}

func TestParseString_InvalidSource_ReturnsError(t *testing.T) {
	t.Parallel()

	p := parser.MustNew()

	program, err := p.ParseString("test.tinyjs", `
x = ;
print(x);
`)
	require.Error(t, err)
	require.Nil(t, program)
}

func TestParseString_UnknownStatement_ReturnsError(t *testing.T) {
	t.Parallel()

	p := parser.MustNew()

	program, err := p.ParseString("test.tinyjs", `
dance(x);
`)
	require.Error(t, err)
	require.Nil(t, program)
}
