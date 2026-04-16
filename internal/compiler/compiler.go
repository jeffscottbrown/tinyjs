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

// Build generates an executable at destPath from the provided LLVM IR.
// On macOS it produces a universal (fat) binary that runs on both Intel
// and Apple Silicon, so only one macOS runner is needed in CI.
func (c *Compiler) Build(ir string, destPath string) error {
	tmpDir, err := os.MkdirTemp("", "tinyjs-build-*")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	irFile := filepath.Join(tmpDir, "out.ll")
	if err := os.WriteFile(irFile, []byte(ir), 0644); err != nil {
		return fmt.Errorf("write IR: %w", err)
	}

	if runtime.GOOS == "darwin" {
		return buildDarwinUniversal(irFile, destPath, tmpDir)
	}
	return runClang([]string{"-Wno-override-module", irFile, "-o", destPath})
}

// buildDarwinUniversal compiles irFile for both x86_64 and arm64, then
// merges the two slices into a universal binary at destPath using lipo.
func buildDarwinUniversal(irFile, destPath, tmpDir string) error {
	var sdkArgs []string
	if out, err := exec.Command("xcrun", "--show-sdk-path").Output(); err == nil {
		sdkArgs = []string{"-isysroot", strings.TrimSpace(string(out))}
	}

	slices := []struct {
		target string
		out    string
	}{
		{"x86_64-apple-macosx10.15", filepath.Join(tmpDir, "out_x86_64")},
		{"arm64-apple-macosx11.0", filepath.Join(tmpDir, "out_arm64")},
	}

	for _, s := range slices {
		args := append([]string{"-Wno-override-module", "--target=" + s.target}, sdkArgs...)
		args = append(args, irFile, "-o", s.out)
		if err := runClang(args); err != nil {
			return err
		}
	}

	var stderr bytes.Buffer
	cmd := exec.Command("lipo", "-create", "-output", destPath, slices[0].out, slices[1].out)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("lipo: %w\n%s", err, stderr.String())
	}
	return nil
}

func runClang(args []string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("clang", args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clang: %w\n%s", err, stderr.String())
	}
	return nil
}
