// compiler/compiler.go
package compiler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"example.com/tinyjs/internal/ast"
	"example.com/tinyjs/internal/parser"
	"example.com/tinyjs/internal/runner"
)

type Compiler struct {
	parser *parser.Parser
}

func New() (*Compiler, error) {
	p, err := parser.New()
	if err != nil {
		return nil, fmt.Errorf("create compiler parser: %w", err)
	}

	return &Compiler{
		parser: p,
	}, nil
}

func MustNew() *Compiler {
	c, err := New()
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Compiler) Parse(filename, input string) (*ast.Program, error) {
	return c.parser.ParseString(filename, input)
}

func (c *Compiler) CompileProgram(program *ast.Program) (string, error) {
	ir, err := c.GenerateIR(program)
	if err != nil {
		return "", fmt.Errorf("compile program: %w", err)
	}
	return ir, nil
}

func (c *Compiler) CompileString(filename, input string) (string, error) {
	program, err := c.Parse(filename, input)
	if err != nil {
		return "", err
	}

	return c.CompileProgram(program)
}

func (c *Compiler) RunProgram(program *ast.Program) (string, error) {
	ir, err := c.CompileProgram(program)
	if err != nil {
		return "", err
	}

	out, err := runner.RunIR(ir)
	if err != nil {
		return "", fmt.Errorf("run IR: %w", err)
	}

	return out, nil
}

func (c *Compiler) RunString(filename, input string) (string, error) {
	program, err := c.Parse(filename, input)
	if err != nil {
		return "", err
	}

	return c.RunProgram(program)
}

// Build parses source, generates IR, and uses clang to create an executable at destPath.
func (c *Compiler) Build(ir string, destPath string) error {

	// Create a temporary directory for the intermediate .ll file
	tmpDir, err := os.MkdirTemp("", "tinyjs-build-*")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	irFile := filepath.Join(tmpDir, "out.ll")
	if err := os.WriteFile(irFile, []byte(ir), 0644); err != nil {
		return fmt.Errorf("write IR: %w", err)
	}

	// Run clang to produce the binary at the requested destination
	clangCmd := exec.Command("clang", clangArgs(irFile, destPath)...)
	var clangStderr bytes.Buffer
	clangCmd.Stderr = &clangStderr

	if err := clangCmd.Run(); err != nil {
		return fmt.Errorf("compile: %w\n%s", err, clangStderr.String())
	}

	return nil
}

func clangArgs(irFile, binFile string) []string {
	args := []string{"-Wno-override-module"}
	if runtime.GOOS == "darwin" {
		if out, err := exec.Command("xcrun", "--show-sdk-path").Output(); err == nil {
			args = append(args, "-isysroot", strings.TrimSpace(string(out)))
		}
	}
	return append(args, irFile, "-o", binFile)
}
