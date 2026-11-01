package cli

import (
	"fmt"
	"strings"
)

// A textCommand is a command that prints text instead of a JSON document.
type textCommand struct {
	// Name is how the docs and the help name it.
	Name string
	// Args are the arguments it needs, when it needs any.
	Args []string
}

// argv is the command line that runs it.
func (c textCommand) argv() []string {
	fields := strings.Fields(c.Name)
	if len(fields) > 0 && strings.HasPrefix(fields[len(fields)-1], "<") {
		// The last field is a placeholder — `skill load <name>`.
		fields = fields[:len(fields)-1]
	}
	return append(append([]string{}, fields...), c.Args...)
}

// textCommands are the commands that print text instead of a JSON document.
//
// The list is the single source of truth: the CLI's own help is built from it, and
// the two documents that state the rule are checked against it. Three copies, one
// list — so the tool cannot miscount its own surface in the place an agent reads
// first.
var textCommands = []textCommand{
	{Name: "help"},
	{Name: "skill help"},
	{Name: "skill load <name>", Args: []string{"orchestrator-role"}},
}

// textCommandNames are the text commands, as the documents name them.
func textCommandNames() []string {
	out := make([]string, 0, len(textCommands))
	for _, c := range textCommands {
		out = append(out, c.Name)
	}
	return out
}

// textCommandsSentence is the sentence the help and the documents share.
//
// It is generated from the list, so the binary cannot say one thing and do another.
func textCommandsSentence() string {
	return fmt.Sprintf("%s commands print text instead: %s.",
		numberWord(len(textCommands)), strings.Join(textCommandNames(), ", "))
}

// numberWord spells a small count.
func numberWord(n int) string {
	words := []string{
		"zero", "one", "two", "three", "four", "five",
		"six", "seven", "eight", "nine", "ten",
	}
	if n < 0 || n >= len(words) {
		return fmt.Sprintf("%d", n)
	}
	return words[n]
}
