package runner_test

import (
	"os/exec"
	"strings"
	"testing"

	"example.com/tinyjs/internal/compiler"
	"example.com/tinyjs/internal/runner"
	"github.com/stretchr/testify/suite"
)

func (s *RunnerSuite) TestPrintingBinaryExpressionWithLiterals() {

	tjsSource := `
	print(14 + 2);
	`

	ir, err := s.compiler.CompileString("", tjsSource)

	out, err := runner.RunIR(ir)
	s.NoError(err)
	s.Equal("16\n", s.normalizeNewlines(out))
}

func (s *RunnerSuite) TestRunIR_PrintsValue() {

	ir := `@.fmt = private unnamed_addr constant [5 x i8] c"%ld\0A\00"

declare i32 @printf(ptr noundef, ...)

define i32 @main() {
entry:
  %x = alloca i64
  store i64 42, ptr %x
  %t1 = load i64, ptr %x
  %t2 = getelementptr inbounds [5 x i8], ptr @.fmt, i64 0, i64 0
  %t3 = call i32 (ptr, ...) @printf(ptr %t2, i64 %t1)
  ret i32 0
}
`

	out, err := runner.RunIR(ir)
	s.Require().NoError(err)
	s.Require().Equal("42\n", s.normalizeNewlines(out))
}

func (s *RunnerSuite) TestRunIR_InvalidIR_ReturnsError() {

	out, err := runner.RunIR(`this is not valid llvm ir`)
	s.Require().Error(err)
	s.Require().Empty(out)
	s.Require().Contains(err.Error(), "lli failed")
}

func requireLLI(t *testing.T) {
	t.Helper()

	_, err := exec.LookPath("lli")
	if err != nil {
		t.Skip("lli not found in PATH; skipping runner execution tests")
	}
}

func (s *RunnerSuite) normalizeNewlines(v string) string {
	return strings.ReplaceAll(v, "\r\n", "\n")
}

type RunnerSuite struct {
	suite.Suite
	compiler *compiler.Compiler
}

func TestCompilerSuite(t *testing.T) {
	requireLLI(t)

	suite.Run(t, new(RunnerSuite))
}

func (s *RunnerSuite) SetupTest() {
	s.compiler = compiler.MustNew()
}
