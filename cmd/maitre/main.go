// Command maitre is the Kitchen Brigade human console.
//
// It refuses to start the interactive view without a terminal, so an agent can never
// launch it by accident. `--render` prints a static snapshot instead, which needs no
// terminal and works in CI. `--keys` drives the console from a script, so its behaviour
// can be checked without a terminal either.
//
// See docs/05-human-interface.md and docs/decisions/0003-surface-split.md.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"charm.land/bubbletea/v2"

	"github.com/lthiagol/maitre/internal/core"
	"github.com/lthiagol/maitre/internal/ui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "maitre: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fs := flag.NewFlagSet("maitre", flag.ContinueOnError)
	render := fs.Bool("render", false, "print a static snapshot of the console and exit")
	plain := fs.Bool("plain", false, "strip colours from the output")
	width := fs.Int("width", 100, "render width")
	show := fs.String("show", "", "render one recipe by id instead of the list")
	keys := fs.String("keys", "", "drive the console with these keys before rendering, comma separated")
	all := fs.Bool("all", false, "render every tab instead of the current one")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	root, err := core.FindRoot(cwd)
	if err != nil {
		return fmt.Errorf("%w — run `scale init` first", err)
	}
	store := core.Open(root)

	if *keys != "" || *all || *show != "" || *render {
		// A driven console writes under its own session, so the journal can show
		// that a verdict was scripted rather than given.
		m, err := ui.SnapshotScripted(store, *width)
		if err != nil {
			return err
		}
		if *show != "" {
			out, err := ui.RenderDetail(store, *show, *width)
			if err != nil {
				return err
			}
			fmt.Print(finish(out, *plain))
			return nil
		}
		if *keys != "" {
			out, err := ui.Drive(m, strings.Split(*keys, ","))
			if err != nil {
				return err
			}
			fmt.Print(finish(out, *plain))
			return nil
		}
		out, err := ui.Render(store, *width)
		if err != nil {
			return err
		}
		fmt.Print(finish(out, *plain))
		return nil
	}

	if !ui.IsTerminal(os.Stdin) {
		return fmt.Errorf("the console needs a terminal — use `maitre --render` for a snapshot, or `scale` for scripted access")
	}

	// The interactive console watches the disk: the chef sees a change on disk
	// without pressing `r`.
	p := tea.NewProgram(ui.NewLive(store))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

// finish strips the colours when asked.
func finish(out string, plain bool) string {
	if plain {
		return ui.StripANSI(out)
	}
	return out
}
