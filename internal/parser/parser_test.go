package parser_test

import (
	"testing"

	"example.com/tinyjs/internal/parser"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ParserSuite struct {
	suite.Suite
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}

func (s *ParserSuite) TestParseSingleAssignment() {
	prog, err := parser.Parse("x = 1;")
	s.Require().NoError(err)
	s.Require().NotNil(prog)
	s.Require().Len(prog.Assignments, 1)

	assign := prog.Assignments[0]
	s.Equal("x", assign.Name)
	s.Equal("=", assign.Eq)
	s.Equal(int64(1), assign.Value)
	s.Equal(";", assign.Semi)
}

func (s *ParserSuite) TestParseMultipleAssignments() {
	src := "x = 1;\ny = 2;\nz = 300;"

	prog, err := parser.Parse(src)
	s.Require().NoError(err)
	s.Require().Len(prog.Assignments, 3)

	s.Equal("x", prog.Assignments[0].Name)
	s.Equal(int64(1), prog.Assignments[0].Value)
	s.Equal("y", prog.Assignments[1].Name)
	s.Equal(int64(2), prog.Assignments[1].Value)
	s.Equal("z", prog.Assignments[2].Name)
	s.Equal(int64(300), prog.Assignments[2].Value)
}

func (s *ParserSuite) TestParseIgnoresWhitespaceAndComments() {
	src := `
		// assign a value
		alpha = 10;

		// assign another value
		beta = 20;
	`

	prog, err := parser.Parse(src)
	s.Require().NoError(err)
	s.Require().Len(prog.Assignments, 2)

	s.Equal("alpha", prog.Assignments[0].Name)
	s.Equal(int64(10), prog.Assignments[0].Value)
	s.Equal("beta", prog.Assignments[1].Name)
	s.Equal(int64(20), prog.Assignments[1].Value)
}

func (s *ParserSuite) TestParseRejectsMissingSemicolon() {
	prog, err := parser.Parse("x = 1")
	s.Nil(prog)
	s.Error(err)
	s.ErrorContains(err, ";")
}

func (s *ParserSuite) TestParseRejectsIdentifierOnRightHandSide() {
	prog, err := parser.Parse("x = y;")
	s.Nil(prog)
	s.Error(err)
	s.ErrorContains(err, "Int")
}

func (s *ParserSuite) TestParseRejectsKeywordLikeLetBecauseGrammarDoesNotSupportIt() {
	prog, err := parser.Parse("let x = 1;")
	s.Nil(prog)
	s.Error(err)
	s.ErrorContains(err, "=")
}

func (s *ParserSuite) TestParseRejectsPrintCall() {
	prog, err := parser.Parse("print(x);")
	s.Nil(prog)
	s.Error(err)
	s.ErrorContains(err, "=")
}

func (s *ParserSuite) TestParseRejectsUnknownCharacter() {
	prog, err := parser.Parse("x = 1$;")
	s.Nil(prog)
	s.Error(err)
	s.ErrorContains(err, "$")
}

func (s *ParserSuite) TestParseRejectsEmptyInputOnlyAsEmptyProgramIsAllowed() {
	prog, err := parser.Parse("")
	s.Require().NoError(err)
	s.Require().NotNil(prog)
	s.Empty(prog.Assignments)
}

func (s *ParserSuite) TestParseVeryLargeIntegerWithinInt64Range() {
	prog, err := parser.Parse("big = 9223372036854775807;")
	s.Require().NoError(err)
	s.Require().Len(prog.Assignments, 1)
	s.Equal(int64(9223372036854775807), prog.Assignments[0].Value)
}

func TestParse_TableDrivenInvalidExamples(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		want string
	}{
		{name: "missing identifier", src: "= 1;", want: "Ident"},
		{name: "missing number", src: "x = ;", want: "Int"},
		{name: "double equals", src: "x == 1;", want: "Int"},
		{name: "no assignment operator", src: "x 1;", want: "="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			prog, err := parser.Parse(tt.src)
			require.Nil(t, prog)
			require.Error(t, err)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
