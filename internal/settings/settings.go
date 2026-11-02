// Package settings holds the tool's settings record: its schema, the read, the
// atomic write and the version check.
//
// The settings are the **tool's** file, not a plan record. A record is the
// machine's record about the work; a theme, an order and a served toggle are the
// chef's preferences, and none of them carries a rule the plan was written under.
// So the plan's format does not move, the domain does not read the file, and the
// file carries a version of its own.
package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// Version is the schema version this build reads and writes.
	//
	// The file's own version, the way the plan's records carry one and the way the
	// skill does. A file written by a newer build is refused, by name — the tool
	// does not guess what it cannot read.
	Version = 1

	// FileName is the settings file, in the plan's directory.
	FileName = "settings.json"

	// DefaultTheme is the theme a console nobody told otherwise uses.
	DefaultTheme = "Kitchen"

	// DefaultOrder is the order the recipe list uses by default.
	DefaultOrder = "triage"
)

// View is what the chef sees: the theme, the order of the list, and whether the
// served recipes are shown.
//
// The values are names, not enums. The console owns the vocabulary — it is what
// renders them — and it refuses a name it does not know, by name.
type View struct {
	Theme      string `json:"theme"`
	Order      string `json:"order"`
	ShowServed bool   `json:"showServed"`
}

// Record is the whole file.
type Record struct {
	Version int  `json:"version"`
	View    View `json:"view"`
}

// Defaults is the record a plan with no settings file has.
//
// A missing file is not an error: the chef's choices have not been made yet, and
// the defaults are the answer until they are.
func Defaults() Record {
	return Record{
		Version: Version,
		View: View{
			Theme:      DefaultTheme,
			Order:      DefaultOrder,
			ShowServed: true,
		},
	}
}

// Path is the settings file inside a plan directory.
//
// The plan's directory is where the file lives, so the settings travel with the
// plan: a checkout that wants its own keeps the file out of git.
func Path(planRoot string) string {
	return filepath.Join(planRoot, FileName)
}

// Read reads the record at path.
//
// A missing file yields the defaults. An unknown version is refused, by name —
// the house's *silent enum fallback* lesson: a safe unknown value disabled a
// whole feature with green tests. Unknown keys are ignored, so a file that
// carries more than this build knows still reads.
func Read(path string) (Record, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Defaults(), nil
		}
		return Record{}, fmt.Errorf("read %s: %w", path, err)
	}
	var rec Record
	if err := json.Unmarshal(raw, &rec); err != nil {
		return Record{}, fmt.Errorf("read %s: %w", path, err)
	}
	if rec.Version != Version {
		return Record{}, fmt.Errorf(
			"the settings file has version %d, which is not one this build has (version %d)",
			rec.Version, Version)
	}
	return rec, nil
}

// The seams of the atomic write.
//
// The write is atomic: the bytes go to a temporary file in the record's folder and
// are renamed over the old one. These two are variables so a test can inject a
// fault at **either** stage — a write that stops halfway, and a rename that fails —
// and assert the property the write is for: **the file that was there before is
// whole**.
//
// A write that overwrites the old file in place has neither stage, so the fault
// cannot be injected at all: the write succeeds and the witness goes red. That is
// what the witness is for — the rule is not the rename, the rule is that the
// chef's file survives a write that fails.
var (
	// writeTemp writes the record's bytes to the temporary file.
	writeTemp = func(f *os.File, data []byte) error {
		_, err := f.Write(data)
		return err
	}

	// renameTemp renames the temporary file over the old one.
	renameTemp = func(oldpath, newpath string) error {
		return os.Rename(oldpath, newpath)
	}
)

// Write writes the record at path, atomically.
//
// The bytes go to a temporary file and are renamed over the old one, so a
// half-written file cannot lose the chef's choices. The version stamped is this
// build's: the writer knows the schema it writes.
func Write(path string, rec Record) error {
	rec.Version = Version
	raw, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("encode the settings: %w", err)
	}
	raw = append(raw, '\n')

	dir := filepath.Dir(path)
	// The temporary file is not a .json file: the plan's fingerprint hashes the
	// .json files, and a half-written one must not be read as the settings.
	tmp, err := os.CreateTemp(dir, ".settings-*.tmp")
	if err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := writeTemp(tmp, raw); err != nil {
		tmp.Close()
		return fmt.Errorf("write %s: %w", tmpName, err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("sync %s: %w", tmpName, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close %s: %w", tmpName, err)
	}
	if err := renameTemp(tmpName, path); err != nil {
		return fmt.Errorf("rename %s: %w", tmpName, err)
	}
	return nil
}
