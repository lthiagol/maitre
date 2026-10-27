package ui

import "os"

// IsTerminal reports whether f is a character device, i.e. a terminal.
//
// The console refuses to start without one, so an agent can never launch it by
// accident — the same guarantee v1's separate binaries gave, with one binary.
func IsTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
