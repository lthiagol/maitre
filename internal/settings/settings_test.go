package settings

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// planDir makes a directory to hold the settings file.
func planDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

// TestAMissingFileYieldsTheDefaults covers the read of a plan nobody has set up.
//
// A missing file is not an error: the chef has not chosen anything yet, and the
// defaults are the answer until they do.
func TestAMissingFileYieldsTheDefaults(t *testing.T) {
	path := Path(planDir(t))

	rec, err := Read(path)
	if err != nil {
		t.Fatalf("a missing file should not be an error: %v", err)
	}
	want := Defaults()
	if rec != want {
		t.Fatalf("the defaults are %+v, want %+v", rec, want)
	}
	if rec.View.Theme != DefaultTheme || rec.View.Order != DefaultOrder || !rec.View.ShowServed {
		t.Fatalf("the defaults are %+v, want the theme %q, the order %q and served shown",
			rec.View, DefaultTheme, DefaultOrder)
	}
}

// TestTheFileHasItsOwnVersion covers the shape of the file.
//
// The settings carry a version of their own, the way the plan's records do. It is
// the schema's version, not the plan's — the plan's format does not move.
func TestTheFileHasItsOwnVersion(t *testing.T) {
	dir := planDir(t)
	path := Path(dir)

	if err := Write(path, Defaults()); err != nil {
		t.Fatalf("write: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("the file is not a JSON document: %v\n%s", err, raw)
	}
	if doc["version"] != float64(Version) {
		t.Fatalf("the file's version is %v, want %d", doc["version"], Version)
	}
	view, ok := doc["view"].(map[string]any)
	if !ok {
		t.Fatalf("the file has no view section: %s", raw)
	}
	for _, key := range []string{"theme", "order", "showServed"} {
		if _, ok := view[key]; !ok {
			t.Errorf("the view section has no %q: %s", key, raw)
		}
	}
}

// TestTheWriterStampsItsOwnVersion covers the version on the way out.
//
// The writer knows the schema it writes, so the version on the file is this
// build's — a record that carries another version is not written as it was.
func TestTheWriterStampsItsOwnVersion(t *testing.T) {
	path := Path(planDir(t))

	rec := Defaults()
	rec.Version = 99
	if err := Write(path, rec); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.Version != Version {
		t.Fatalf("the file's version is %d, want %d", got.Version, Version)
	}
}

// TestAnUnknownVersionIsRefusedByName covers the version fence.
//
// The tool does not guess what it cannot read: a safe unknown value disabled a
// whole feature with green tests. The refusal names the version it refused.
func TestAnUnknownVersionIsRefusedByName(t *testing.T) {
	path := Path(planDir(t))
	raw := `{"version": 99, "view": {"theme": "Kitchen", "order": "triage", "showServed": true}}`
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Read(path)
	if err == nil {
		t.Fatal("an unknown version should be refused")
	}
	if !strings.Contains(err.Error(), "99") {
		t.Fatalf("the refusal should name the version it refused: %v", err)
	}
	if !strings.Contains(err.Error(), "not one this build has") {
		t.Fatalf("the refusal should say why: %v", err)
	}
}

// TestAMissingVersionIsRefused covers a file that names no version.
//
// The file's version is part of its schema. A file without one is a file this
// build cannot read — and a safe guess would be the silent fallback the house
// learned to refuse.
func TestAMissingVersionIsRefused(t *testing.T) {
	path := Path(planDir(t))
	if err := os.WriteFile(path, []byte(`{"view": {}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Read(path); err == nil {
		t.Fatal("a file with no version should be refused")
	}
}

// TestUnknownKeysAreIgnored covers the forward compatibility.
//
// A file written by a newer build, with keys this build does not know, still
// reads: the unknown keys are ignored, and the ones this build knows are read.
func TestUnknownKeysAreIgnored(t *testing.T) {
	path := Path(planDir(t))
	raw := `{
	  "version": 1,
	  "crew": {"harness": "codex", "count": 3},
	  "view": {"theme": "Dracula", "order": "lane", "showServed": false, "density": "compact"}
	}`
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	rec, err := Read(path)
	if err != nil {
		t.Fatalf("a file with unknown keys should read: %v", err)
	}
	if rec.View.Theme != "Dracula" || rec.View.Order != "lane" || rec.View.ShowServed {
		t.Fatalf("the record is %+v, want the theme Dracula, the order lane and served hidden", rec.View)
	}
}

// TestAWriteLeavesNothingBehind covers the tidiness of a write.
//
// A write leaves the plan's folder holding **the record and nothing else** —
// the successful one and the failed one alike. It is not only tidiness: the
// plan's fingerprint hashes every file in that folder, so a leftover is the
// plan's own scan to answer.
//
// **The witness bites.** A write that does not remove its temporary file leaves
// one behind, and this test says so.
func TestAWriteLeavesNothingBehind(t *testing.T) {
	dir := planDir(t)
	path := Path(dir)

	// A successful write.
	if err := Write(path, Defaults()); err != nil {
		t.Fatalf("write: %v", err)
	}
	assertOnlyTheRecord(t, dir)

	// A failed one: the bytes go half way into the temporary file, and the
	// write fails. The temporary file goes with it.
	restore := writeTemp
	writeTemp = func(f *os.File, data []byte) error {
		if len(data) > 1 {
			if _, err := f.Write(data[:len(data)/2]); err != nil {
				return err
			}
		}
		return errors.New("the disk is full")
	}
	defer func() { writeTemp = restore }()

	rec := Defaults()
	rec.View.Theme = "Dracula"
	if err := Write(path, rec); err == nil {
		t.Fatal("the write should have failed — the fault was injected into it")
	}
	assertOnlyTheRecord(t, dir)
}

// TestAWriteIntoAFolderThatIsNotThereIsRefused covers the refusal.
//
// The record lives in the plan's folder, and the write does not make a folder
// that is not there: a write that cannot happen is refused, and the file that
// was there before is left alone.
//
// **The witness bites.** A write that creates the folder it was asked for
// succeeds, and this test says so.
func TestAWriteIntoAFolderThatIsNotThereIsRefused(t *testing.T) {
	dir := planDir(t)
	path := Path(dir)
	whole := writeRecord(t, path, Defaults())

	rec := Defaults()
	rec.View.Theme = "Dracula"
	missing := Path(filepath.Join(dir, "nowhere", FileName))
	if err := Write(missing, rec); err == nil {
		t.Fatal("a write into a directory that is not there should fail")
	}

	// And the file that was there before is as it was: the old bytes are still
	// the whole file, not half of it.
	assertWhole(t, path, whole)
}

// TestAHalfWrittenFileCannotLoseTheChefsChoices is the witness for the atomic write.
//
// A write that stops halfway — a full disk, a pulled cable — must leave the file
// that was there before **whole**: not truncated, not empty, not a partial record.
// The fault is injected into the temporary file's write: the bytes go half way and
// the write fails.
//
// **The witness bites.** A `Write` that overwrites the old file in place has no
// temporary file, so the fault is never injected, the write succeeds, and this
// test fails on the line that says the write should have failed.
func TestAHalfWrittenFileCannotLoseTheChefsChoices(t *testing.T) {
	dir := planDir(t)
	path := Path(dir)
	whole := writeRecord(t, path, Defaults())

	// The fault: half the bytes, and then a failure.
	restore := writeTemp
	writeTemp = func(f *os.File, data []byte) error {
		if len(data) > 1 {
			if _, err := f.Write(data[:len(data)/2]); err != nil {
				return err
			}
		}
		return errors.New("the disk is full")
	}
	defer func() { writeTemp = restore }()

	rec := Defaults()
	rec.View.Theme = "Dracula"
	if err := Write(path, rec); err == nil {
		t.Fatal("the write should have failed — the fault was injected into it")
	}

	// The chef's file is whole.
	assertWhole(t, path, whole)
}

// TestARenameThatFailsCannotLoseTheChefsChoices is the second half of the witness.
//
// The bytes are staged and the **last** step fails — the rename. The old file must
// be whole here too: a failure at any point of the write is a failure of the write.
//
// **The witness bites.** A `Write` that replaces the old file by writing over it
// has no rename, so the fault is never injected, the write succeeds, and this test
// fails on the line that says the write should have failed.
func TestARenameThatFailsCannotLoseTheChefsChoices(t *testing.T) {
	dir := planDir(t)
	path := Path(dir)
	whole := writeRecord(t, path, Defaults())

	// The fault: the rename says no.
	restore := renameTemp
	renameTemp = func(oldpath, newpath string) error {
		return errors.New("the rename failed")
	}
	defer func() { renameTemp = restore }()

	rec := Defaults()
	rec.View.Theme = "Dracula"
	if err := Write(path, rec); err == nil {
		t.Fatal("the write should have failed — the fault was injected into it")
	}

	// The chef's file is whole.
	assertWhole(t, path, whole)
}

// writeRecord writes the record and returns the bytes that are on the disk.
func writeRecord(t *testing.T, path string, rec Record) []byte {
	t.Helper()
	if err := Write(path, rec); err != nil {
		t.Fatalf("write: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// assertWhole requires the file at path to be byte for byte what it was.
func assertWhole(t *testing.T, path string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the file that was there before is gone: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("the file that was there before is not whole:\n%s", got)
	}
}

// assertOnlyTheRecord requires the folder to hold the record and nothing else.
func assertOnlyTheRecord(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != FileName {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("the folder holds %v, want only %s", names, FileName)
	}
}

// TestTheRoundTrip covers the record's own life: written, read back, the same.
func TestTheRoundTrip(t *testing.T) {
	path := Path(planDir(t))

	rec := Record{View: View{Theme: "Alucard", Order: "plain", ShowServed: false}}
	if err := Write(path, rec); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.View != rec.View {
		t.Fatalf("the record is %+v, want %+v", got.View, rec.View)
	}
}
