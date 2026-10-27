//go:build !unix

package core

import "os/exec"

// startGroup is a no-op: a process group is a Unix idea.
func startGroup(cmd *exec.Cmd) {}

// killGroup kills the command itself.
func killGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}

// processAlive reports whether a pid is still running.
//
// Without a signal to probe with, the only honest answer is "we cannot tell".
func processAlive(pid int) bool { return true }
