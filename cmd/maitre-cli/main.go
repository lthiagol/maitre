// Command maitre-cli is the maitre agent CLI.
//
// Agents drive the plan through it. The human-facing console is the separate `maitre`
// binary. See docs/decisions/0003-surface-split.md.
package main

import (
	"os"

	"github.com/lthiagol/maitre/internal/cli"
)

func main() {
	r := &cli.Runner{Args: os.Args[1:], Stdout: os.Stdout, Stderr: os.Stderr}
	os.Exit(r.Run())
}
