package runner

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RunIR(ir string) (string, error) {
	cmd := exec.Command("lli")
	cmd.Stdin = strings.NewReader(ir)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("lli failed: %w\nstderr:\n%s", err, stderr.String())
	}

	return stdout.String(), nil
}
