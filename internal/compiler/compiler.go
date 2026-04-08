package compiler

import (
	"fmt"
	"strings"

	"example.com/tinyjs/internal/ast"
)

const (
	headerLine1 = "; tinyjs output\n"
	headerLine2 = "; accepted syntax: <ident> = <integer>;\n\n"
)

// Compile emits a deliberately tiny pseudo output representation.
//
// Later you can replace this with a richer IR or LLVM backend without having
// to change the frontend structure very much.
func Compile(p *ast.Program) (string, error) {
	if p == nil {
		return "", fmt.Errorf("program is nil")
	}

	var b strings.Builder
	b.WriteString(headerLine1)
	b.WriteString(headerLine2)

	for _, a := range p.Assignments {
		if a == nil {
			return "", fmt.Errorf("assignment is nil")
		}
		if a.Name == "" {
			return "", fmt.Errorf("assignment has empty variable name")
		}

		fmt.Fprintf(&b, "var %s = %d\n", a.Name, a.Value)
	}

	return b.String(), nil
}
