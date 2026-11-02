package cli

import (
	"github.com/lthiagol/maitre/internal/settings"
)

// cmdSettingsShow prints the tool's settings record.
//
// The record is the tool's file, not the plan's: the console writes it — it is the
// one that knows the vocabulary — and the CLI reads it. An agent that wants to
// know how the plan is set up reads it here.
func (r *Runner) cmdSettingsShow(args []string) int {
	fs := newFlagSet("settings show")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	root, err := planRoot()
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	rec, err := settings.Read(settings.Path(root))
	if err != nil {
		r.fail("read the settings", err.Error())
		return ExitError
	}
	return r.emit(map[string]any{"ok": true, "settings": rec})
}
