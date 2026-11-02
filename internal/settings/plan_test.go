package settings_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/core"
	"github.com/lthiagol/maitre/internal/settings"
)

// TestTheSettingsFileLivesInThePlanDirectory covers where the file is.
//
// The settings travel with the plan: the file is in the plan's directory, beside
// the plan's records — and beside them, not among them.
func TestTheSettingsFileLivesInThePlanDirectory(t *testing.T) {
	root := filepath.Join(t.TempDir(), core.PlanDirName)

	path := settings.Path(root)
	if filepath.Base(path) != "settings.json" {
		t.Fatalf("the settings file is %q, want settings.json", filepath.Base(path))
	}
	if filepath.Base(filepath.Dir(path)) != core.PlanDirName {
		t.Fatalf("the settings file is in %q, want the plan's directory %q",
			filepath.Base(filepath.Dir(path)), core.PlanDirName)
	}
}

// TestThePlanFormatDoesNotMove covers the version fence of the plan.
//
// The settings carry no rules, so the plan's format does not move. The domain's
// version is the plan's, and it stays 7.
func TestThePlanFormatDoesNotMove(t *testing.T) {
	if core.DomainVersion != 7 {
		t.Fatalf("the domain's version is %d, want 7", core.DomainVersion)
	}
}

// TestTheDomainDoesNotReadTheSettings covers the ownership.
//
// The file is the tool's, not the domain's: no file in the domain mentions it.
func TestTheDomainDoesNotReadTheSettings(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("..", "core"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join("..", "core", e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		text := string(raw)
		for _, mention := range []string{"settings.json", "internal/settings"} {
			if strings.Contains(text, mention) {
				t.Errorf("%s mentions %q — the domain does not know the file", e.Name(), mention)
			}
		}
	}
}

// TestThePlansReadersIgnoreTheFile covers a file in the plan's directory that is
// not a plan record.
//
// A settings file beside the plan's records does not change what the plan is.
func TestThePlansReadersIgnoreTheFile(t *testing.T) {
	root := filepath.Join(t.TempDir(), core.PlanDirName)
	store := core.Open(root)
	if err := core.Seed(store); err != nil {
		t.Fatalf("seed: %v", err)
	}
	before, err := store.ListRecipes()
	if err != nil {
		t.Fatalf("list recipes: %v", err)
	}

	if err := settings.Write(settings.Path(root), settings.Defaults()); err != nil {
		t.Fatalf("write the settings: %v", err)
	}

	after, err := store.ListRecipes()
	if err != nil {
		t.Fatalf("list recipes with a settings file in the plan's directory: %v", err)
	}
	if len(after) != len(before) {
		t.Fatalf("the plan holds %d recipes with a settings file, want %d", len(after), len(before))
	}
	for i, rec := range after {
		if rec.ID != before[i].ID {
			t.Fatalf("recipe %d is %q with a settings file, want %q", i, rec.ID, before[i].ID)
		}
	}
}
