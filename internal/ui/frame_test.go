package ui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The frame forks no process.
//
// The console's own rule — in `crew.go` — says **a frame must never wait on a
// process**, and the render path broke it: the snapshot asked the runtime for
// the crew, and the runtime is two spawns. The witnesses here are **process
// level**: a fake `herdr` sits on the PATH and records every call, so the
// count is the claim.

// fakeHerdr is a `herdr` on the PATH that records what it was asked.
//
// The script answers with a fixed document and appends each call to a log, so
// a test can count the processes the console forked.
type fakeHerdr struct {
	// dir is where the binary sits — the directory the PATH points at.
	dir string
	// log is the file the calls are appended to.
	log string
}

// herdrOnPath puts a `herdr` on the PATH.
//
// The answer is what the runtime would say: the agents and the tabs, in the
// shape `herdr api snapshot` prints.
func herdrOnPath(t *testing.T, answer any) *fakeHerdr {
	t.Helper()
	dir := t.TempDir()
	log := filepath.Join(dir, "calls.log")
	script := "#!/bin/sh\n" +
		"printf '%s\\n' \"$*\" >> '" + log + "'\n" +
		"cat '" + filepath.Join(dir, "answer.json") + "'\n"
	if err := os.WriteFile(filepath.Join(dir, "herdr"), []byte(script), 0o755); err != nil {
		t.Fatalf("write the fake herdr: %v", err)
	}
	raw, err := json.Marshal(answer)
	if err != nil {
		t.Fatalf("write the fake herdr's answer: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "answer.json"), raw, 0o644); err != nil {
		t.Fatalf("write the fake herdr's answer: %v", err)
	}
	return &fakeHerdr{dir: dir, log: log}
}

// calls are the calls the fake herdr recorded.
func (h *fakeHerdr) calls(t *testing.T) []string {
	t.Helper()
	raw, err := os.ReadFile(h.log)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		t.Fatalf("read the calls: %v", err)
	}
	var out []string
	for _, line := range strings.Split(strings.TrimRight(string(raw), "\n"), "\n") {
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

// herdrAnswer is the document the runtime answers with.
//
// It is the shape `herdr api snapshot` prints: the agents and the tabs, in one
// document.
func herdrAnswer(root string) map[string]any {
	return map[string]any{
		"result": map[string]any{
			"snapshot": map[string]any{
				"agents": []map[string]any{
					{
						"pane_id":        "w1:p1",
						"tab_id":         "w1:t1",
						"cwd":            root,
						"foreground_cwd": root,
						"agent_status":   "working",
					},
				},
				"tabs": []map[string]any{
					{"tab_id": "w1:t1", "label": "maitre: orchestrator"},
				},
			},
		},
	}
}

// TestASnapshotSpawnsNothing covers the render path.
//
// A frame draws the screen; it does not fork a process to draw it. The runtime
// is **on the PATH** here, and it is never called: the snapshot asks nothing.
func TestASnapshotSpawnsNothing(t *testing.T) {
	store := goldenStore(t)
	herdr := herdrOnPath(t, herdrAnswer(store.ProjectRoot()))
	t.Setenv("PATH", herdr.dir)

	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatalf("the snapshot: %v", err)
	}
	m.tab = TabSousChef
	// The frame, drawn twice: a frame that asked would ask on every frame.
	_ = plain(m)
	_ = plain(m)

	if calls := herdr.calls(t); len(calls) != 0 {
		t.Errorf("the render path forked %d times, want none: %v", len(calls), calls)
	}
}

// TestTheTabRendersWithHerdrOffThePath covers the tab without a runtime.
//
// **herdr is not on the PATH**: the tab still renders, and it says what it
// knows. It does not say the runtime is not here — the console never asked —
// and it does not say the automation is not running, which is a guess about a
// crew it has not heard about.
func TestTheTabRendersWithHerdrOffThePath(t *testing.T) {
	t.Setenv("PATH", "")
	store := goldenStore(t)

	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatalf("the snapshot: %v", err)
	}
	m.tab = TabSousChef
	out := plain(m)

	if !strings.Contains(out, "THE STAMP") {
		t.Fatalf("the tab has no stamp:\n%s", out)
	}
	if strings.Contains(out, "herdr is not here") {
		t.Errorf("the tab says the runtime is not here, and it never asked:\n%s", out)
	}
	if strings.Contains(out, "the automation is not running") {
		t.Errorf("the tab says the automation is not running, and it does not know:\n%s", out)
	}
	if !strings.Contains(out, "the crew has not been asked") {
		t.Errorf("the tab should say it has not asked the runtime:\n%s", out)
	}
}

// TestTheGoldensRenderWithHerdrOffThePath covers the goldens.
//
// A golden that needs a binary on the machine is a golden that fails on
// someone else's. Every golden is rendered here with **herdr off the PATH**,
// and each one is the file in the tree — byte for byte.
func TestTheGoldensRenderWithHerdrOffThePath(t *testing.T) {
	t.Setenv("PATH", "")
	home := t.TempDir()
	t.Setenv("HOME", home)

	goldens := []struct {
		path   string
		render func(t *testing.T) string
	}{
		{goldenPath, func(t *testing.T) string {
			m, err := Snapshot(goldenStore(t), 100)
			if err != nil {
				t.Fatalf("the snapshot: %v", err)
			}
			m.tab = TabRecipes
			return StripANSI(m.View().Content)
		}},
		{homeGoldenPath, func(t *testing.T) string {
			store := copyTemplate(t, planTemplate, home, "demo")
			m, err := Snapshot(store, 80)
			if err != nil {
				t.Fatalf("the snapshot: %v", err)
			}
			return clockPattern.ReplaceAllString(StripANSI(m.View().Content), "HH:MM:SS")
		}},
		{narrowBoardGoldenPath, func(t *testing.T) string {
			m, err := Snapshot(goldenStore(t), 100)
			if err != nil {
				t.Fatalf("the snapshot: %v", err)
			}
			m.tab = TabBoard
			return StripANSI(m.View().Content)
		}},
		{filterGoldenPath, func(t *testing.T) string {
			m, err := Snapshot(goldenStore(t), 100)
			if err != nil {
				t.Fatalf("the snapshot: %v", err)
			}
			m.tab = TabRecipes
			if _, err := Drive(m, []string{"/", "text:glossary"}); err != nil {
				t.Fatalf("drive the filter: %v", err)
			}
			return StripANSI(m.View().Content)
		}},
		{stressGoldenPath, func(t *testing.T) string {
			m, err := Snapshot(stressStore(t), 100)
			if err != nil {
				t.Fatalf("the snapshot: %v", err)
			}
			m.tab = TabRecipes
			m.height = 40
			return StripANSI(m.View().Content)
		}},
	}

	for _, golden := range goldens {
		want, err := os.ReadFile(golden.path)
		if err != nil {
			t.Fatalf("read %s: %v", golden.path, err)
		}
		if got := golden.render(t); got != string(want) {
			t.Errorf("%s does not render with herdr off the PATH.\n--- got ---\n%s\n--- want ---\n%s",
				golden.path, got, want)
		}
	}
}

// TestTheSnapshotDoesNotAskTheRuntime is the seam behind the process.
//
// The runtime is a seam, and the snapshot does not touch it: the console asks
// in a command, and a frame reads what it already knows.
func TestTheSnapshotDoesNotAskTheRuntime(t *testing.T) {
	store := goldenStore(t)
	m, err := Snapshot(store, 100)
	if err != nil {
		t.Fatalf("the snapshot: %v", err)
	}
	if m.crewAsked {
		t.Error("the snapshot asked the runtime for the crew")
	}
	if len(m.crew) != 0 {
		t.Errorf("the snapshot holds %d crew members, want none", len(m.crew))
	}
	if m.crewErr != nil {
		t.Errorf("the snapshot holds an error it never asked for: %v", m.crewErr)
	}
}
