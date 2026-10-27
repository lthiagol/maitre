//go:build unix

package core

import (
	"os/exec"
	"syscall"
)

// startGroup puts the command in its own process group.
//
// A verification is a shell, and the shell has children: the test suite it starts.
// Killing the shell leaves them running, holding the project directory.
func startGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killGroup kills the command's whole process group.
func killGroup(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	// The negative pid is the group: the command and everything it started.
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}

// processAlive reports whether a pid is still running.
func processAlive(pid int) bool {
	return syscall.Kill(pid, 0) == nil
}
