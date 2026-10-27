package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Directory names inside the plan root.
const (
	DirRecipes   = "recipes"
	DirCourses   = "courses"
	DirPrepNotes = "prep-notes"
	DirSpecials  = "specials"

	FileProject   = "project.json"
	FileActivity  = "activity.json"
	FileDecisions = "decisions.json"
)

// PlanDirName is the plan directory, relative to the project root.
const PlanDirName = "kitchen"

// ErrNotFound is returned when an artifact does not exist.
var ErrNotFound = errors.New("not found")

// Store reads and writes the plan directory. It is the only thing that touches disk.
type Store struct {
	// Root is the plan directory, e.g. <project>/kitchen.
	Root string
}

// Open returns a store rooted at the given plan directory.
func Open(root string) *Store { return &Store{Root: root} }

// FindRoot walks up from start looking for a plan directory.
func FindRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, PlanDirName)
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no %s/ directory found above %s", PlanDirName, start)
		}
		dir = parent
	}
}

// Init creates the plan directory and its subdirectories.
func (s *Store) Init() error {
	for _, d := range []string{"", DirRecipes, DirCourses, DirPrepNotes, DirSpecials} {
		if err := os.MkdirAll(filepath.Join(s.Root, d), 0o755); err != nil {
			return fmt.Errorf("create %s: %w", d, err)
		}
	}
	project := filepath.Join(s.Root, FileProject)
	if _, err := os.Stat(project); os.IsNotExist(err) {
		if err := writeJSONAtomic(project, map[string]any{
			"name":     filepath.Base(filepath.Dir(s.Root)),
			"created":  Now(),
			"charter":  map[string]any{"goals": []string{}, "non_goals": []string{}, "principles": []string{}},
			"settings": map[string]any{},
		}); err != nil {
			return err
		}
	}
	return nil
}

// LoadProject reads the plan-level metadata.
func (s *Store) LoadProject() (Project, error) {
	raw, err := os.ReadFile(filepath.Join(s.Root, FileProject))
	if err != nil {
		if os.IsNotExist(err) {
			return Project{Name: filepath.Base(filepath.Dir(s.Root))}, nil
		}
		return Project{}, fmt.Errorf("read project: %w", err)
	}
	var doc struct {
		Name    string `json:"name"`
		Created string `json:"created"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return Project{}, fmt.Errorf("parse project: %w", err)
	}
	return Project{Name: doc.Name, Created: doc.Created}, nil
}

// filepathJoin is filepath.Join, named short because it is used a lot.
func filepathJoin(parts ...string) string { return filepath.Join(parts...) }

// readDirSafe reads a directory, treating a missing one as empty.
func readDirSafe(dir string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return entries, nil
}

// readFile reads a file.
func readFile(path string) ([]byte, error) { return os.ReadFile(path) }

// unmarshal decodes JSON into v.
func unmarshal(raw []byte, v any) error { return json.Unmarshal(raw, v) }

// sortByID sorts a slice in place by the id at the given index.
func sortByID[T any](items []T, id func(int) string) {
	sort.Slice(items, func(i, j int) bool { return recipeLess(id(i), id(j)) })
}

// RecipeDir is the directory holding recipes.
func (s *Store) RecipeDir() string { return filepath.Join(s.Root, DirRecipes) }

// ProjectRoot is the directory that contains the plan directory.
//
// Verifications run in the project, not in the plan — `go test ./...` has to see the
// project root.
func (s *Store) ProjectRoot() string { return filepath.Dir(s.Root) }

// SaveRecipe writes a recipe atomically.
//
// The file is named `<id>-<slug>.json`, and the slug follows the title — so a rename
// writes a new file and the old one would stay. The directory is the plan's index
// (C2): one artifact, one file. So the write removes the copy the old name left
// behind — and a removal that fails does not fail the write: the artifact landed,
// and the next write cleans the leftover up.
func (s *Store) SaveRecipe(r *Recipe) error {
	if r.ID == "" {
		return fmt.Errorf("recipe has no id")
	}
	if r.Slug == "" {
		r.Slug = Slugify(r.Title)
	}
	name := fmt.Sprintf("%s-%s.json", r.ID, r.Slug)
	dir := s.RecipeDir()
	if err := writeJSONAtomic(filepath.Join(dir, name), r); err != nil {
		return err
	}
	return removeOtherCopies(dir, name, r.ID)
}

// removeOtherCopies deletes every file in dir that carries the same id as the one
// just written.
//
// A copy left behind is one artifact twice, and every reader picks one arbitrarily:
// `show` reads the old title, `list` counts the recipe twice, the board draws it
// twice, and the gates read whichever copy they found.
//
// The match is on the **whole id**, not the number behind it. The ids are numbered
// per kind, so "C01" and "R01" share a number — and a file of another kind in this
// directory is not a copy of this artifact. It is somebody else's.
//
// The write came first, so the artifact is there whatever happens here. A removal
// that fails is not a failure of the write: the leftover stays, the next write
// removes it, and the reader still picks the current copy.
func removeOtherCopies(dir, keep, id string) error {
	paths, err := copiesNamed(dir, id)
	if err != nil {
		return err
	}
	for _, p := range paths {
		name := filepath.Base(p)
		if name == keep {
			continue
		}
		// A leftover is not worth failing a write over: the artifact landed, the
		// reader picks the current copy, and the next write cleans this up.
		_ = os.Remove(p)
	}
	return nil
}

// copiesNamed lists the files in dir whose name carries exactly this id.
//
// The whole id, not the number: the file name always carries the canonical id, and
// the ids are numbered per kind.
func copiesNamed(dir, id string) ([]string, error) {
	entries, err := readDirSafe(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		base := strings.SplitN(e.Name(), "-", 2)[0]
		if !strings.EqualFold(base, id) {
			continue
		}
		out = append(out, filepath.Join(dir, e.Name()))
	}
	return out, nil
}

// copiesWithID lists the files in dir that carry the given id.
//
// The reader is tolerant: "r1", "R01" and "1" all name the same artifact, so the
// number behind the id is what is compared.
func copiesWithID(dir, id string) ([]string, error) {
	entries, err := readDirSafe(dir)
	if err != nil {
		return nil, err
	}
	want := normaliseID(id)
	var out []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		if normaliseID(strings.SplitN(e.Name(), "-", 2)[0]) != want {
			continue
		}
		out = append(out, filepath.Join(dir, e.Name()))
	}
	return out, nil
}

// currentCopy picks the copy of an artifact that the plan means.
//
// A plan written before the writer removed stale copies can hold two files for one
// id. The current copy is the one the artifact itself says is newest — its own
// `updated` stamp — because that is the plan's record of when it last changed.
//
// The stamps have second resolution, so a rename and the copy it leaves behind
// usually carry the same one. The file's own time breaks that tie: the copy written
// last is the current one. The name order is the last resort, so the answer never
// depends on the order the filesystem happens to hand back.
func currentCopy(paths []string) string {
	if len(paths) == 1 {
		return paths[0]
	}
	best := ""
	bestStamp := ""
	bestMod := time.Time{}
	for _, p := range paths {
		stamp := updatedStamp(p)
		mod := modifiedAt(p)
		if best == "" {
			best, bestStamp, bestMod = p, stamp, mod
			continue
		}
		// The stamp decides. The file's own time breaks the tie the second
		// resolution leaves. The first file in name order is the last resort, so
		// the answer never depends on the order the filesystem hands back.
		if stamp > bestStamp || (stamp == bestStamp && mod.After(bestMod)) {
			best, bestStamp, bestMod = p, stamp, mod
		}
	}
	return best
}

// updatedStamp reads the `updated` stamp of an artifact file.
//
// A file that cannot be read or parsed has no stamp: it loses to a copy that has
// one, and among copies that have none the file time decides.
func updatedStamp(path string) string {
	raw, err := readFile(path)
	if err != nil {
		return ""
	}
	var probe struct {
		Updated string `json:"updated"`
	}
	if err := unmarshal(raw, &probe); err != nil {
		return ""
	}
	return probe.Updated
}

// modifiedAt is when a file was last written.
func modifiedAt(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

// dedupeByID keeps one entry per id: the copy with the newest `updated` stamp.
//
// The directory is the plan's index, and an index that holds the same artifact twice
// has drifted — `list` would count it twice and the board would draw it twice.
func dedupeByID[T any](items []T, id, updated func(T) string) []T {
	at := map[string]int{}
	out := make([]T, 0, len(items))
	for _, it := range items {
		k := id(it)
		if i, ok := at[k]; ok {
			if updated(it) > updated(out[i]) {
				out[i] = it
			}
			continue
		}
		at[k] = len(out)
		out = append(out, it)
	}
	return out
}

// DiskFingerprint summarises the plan's files on disk.
//
// It is what the console compares between polls: the same fingerprint means nothing
// changed. The scan stats the plan's directories and the files in them — it opens
// none of them, and it reads no plan data.
//
// A directory's own time changes when a file is added or removed, but not when a
// file inside it is rewritten. So both are watched: the directories for the files
// that come and go, the files for the ones that change.
//
// The scan tolerates a plan that is being written to: a file that disappears
// between the listing and the stat is skipped, and one that appears after the
// listing is seen by the next scan.
func (s *Store) DiskFingerprint() (string, error) {
	var b strings.Builder
	dirs := []string{
		s.Root,
		filepath.Join(s.Root, DirRecipes),
		filepath.Join(s.Root, DirPrepNotes),
		filepath.Join(s.Root, DirSpecials),
		filepath.Join(s.Root, DirCourses),
	}
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				// A directory that is not there yet — the courses are not built.
				continue
			}
			return "", err
		}
		fmt.Fprintf(&b, "d %s %d %d\n", dir, info.ModTime().UnixNano(), info.Size())
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		for _, e := range entries {
			if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
				continue
			}
			info, err := e.Info()
			if err != nil {
				if os.IsNotExist(err) {
					// It went away between the listing and the stat.
					continue
				}
				return "", err
			}
			fmt.Fprintf(&b, "f %s %d %d\n", filepath.Join(dir, e.Name()), info.ModTime().UnixNano(), info.Size())
		}
	}
	return b.String(), nil
}

// LoadRecipe reads one recipe by id.
func (s *Store) LoadRecipe(id string) (*Recipe, error) {
	path, err := s.findRecipePath(id)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var r Recipe
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	// Before the normalisation, so the legacy digest can still be read.
	r.restampLegacyApproval()
	r.normaliseLists()
	return &r, nil
}

// ListRecipes returns every recipe, sorted by id.
//
// The list is never nil: an empty plan is an empty list, and the CLI marshals an
// empty list as `[]` and a nil slice as `null`. A client cannot iterate `null`.
func (s *Store) ListRecipes() ([]Recipe, error) {
	entries, err := os.ReadDir(s.RecipeDir())
	if err != nil {
		if os.IsNotExist(err) {
			return []Recipe{}, nil
		}
		return nil, fmt.Errorf("read %s: %w", s.RecipeDir(), err)
	}
	out := make([]Recipe, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(s.RecipeDir(), e.Name()))
		if err != nil {
			return nil, err
		}
		var r Recipe
		if err := json.Unmarshal(raw, &r); err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		r.restampLegacyApproval()
		r.normaliseLists()
		out = append(out, r)
	}
	// One artifact, one entry: a copy left behind by a rename must not make the
	// recipe count twice.
	out = dedupeByID(out,
		func(r Recipe) string { return r.ID },
		func(r Recipe) string { return r.Updated },
	)
	sort.Slice(out, func(i, j int) bool { return recipeLess(out[i].ID, out[j].ID) })
	return out, nil
}

// NextRecipeID allocates the next recipe id by scanning what exists, never by counting
// (constraint P9).
//
// It also scans what is *referenced*: an id a prep note points at is taken, even
// though no recipe carries it. A promotion records where it went before the target is
// written, so a crash can leave a pointer to an id that does not exist — and that
// pointer must resolve to nothing, never to a stranger.
func (s *Store) NextRecipeID() (string, error) {
	reserved, err := s.reservedRecipeIDs()
	if err != nil {
		return "", err
	}
	return nextIDIn(s.RecipeDir(), "R", reserved...)
}

// reservedRecipeIDs are the recipe ids a prep note points at.
func (s *Store) reservedRecipeIDs() ([]string, error) {
	notes, err := s.ListPrepNotes()
	if err != nil {
		return nil, err
	}
	promoted := make([]string, 0, len(notes))
	for _, n := range notes {
		promoted = append(promoted, n.PromotedTo)
	}
	return referencedIDs(promoted, "R"), nil
}

// parseNumericID extracts the number from an id like "R07" or "07".
func parseNumericID(id string) (int, bool) {
	s := id
	if len(s) > 0 && s[0] >= 'A' && s[0] <= 'Z' {
		s = s[1:]
	}
	if s == "" {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}

// nextID scans ids for the highest number and returns the next one.
//
// The ids are allocated from a scan of what exists, never from a count — v1's
// `len()+1` produced [AC-01, AC-03, AC-03] after a removal. The width is the zero
// padding: 2 for "AC-01", 0 for "G1".
func nextID(prefix string, ids []string, width int) string {
	max := 0
	for _, id := range ids {
		n, ok := parseNumericID(strings.TrimPrefix(id, prefix))
		if !ok {
			continue
		}
		if n > max {
			max = n
		}
	}
	if width > 0 {
		return fmt.Sprintf("%s%0*d", prefix, width, max+1)
	}
	return fmt.Sprintf("%s%d", prefix, max+1)
}

// findRecipePath resolves an id to its file, tolerating a missing or zero-padded id.
//
// A plan can hold two files for one id — a copy left behind by a rename. The current
// copy wins, so the reader never picks one arbitrarily.
func (s *Store) findRecipePath(id string) (string, error) {
	paths, err := copiesWithID(s.RecipeDir(), id)
	if err != nil {
		return "", err
	}
	if len(paths) == 0 {
		return "", fmt.Errorf("recipe %s: %w", id, ErrNotFound)
	}
	return currentCopy(paths), nil
}

// normaliseID strips a leading kind prefix and zero padding so "r1", "R01" and "1"
// match.
//
// The prefix is stripped whatever it is — "R", "N", "S" — so every artifact kind
// answers to the same forms. The ids are unambiguous by directory, so the kind
// letter carries no information once the number is known.
func normaliseID(id string) string {
	s := strings.ToUpper(strings.TrimSpace(id))
	if len(s) > 0 && s[0] >= 'A' && s[0] <= 'Z' {
		s = s[1:]
	}
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return "0"
	}
	return s
}

// recipeLess orders ids numerically, so R10 sorts after R2.
func recipeLess(a, b string) bool {
	as := idSegments(a)
	bs := idSegments(b)
	for i := 0; i < len(as) && i < len(bs); i++ {
		an, aerr := strconv.Atoi(as[i])
		bn, berr := strconv.Atoi(bs[i])
		if aerr == nil && berr == nil {
			if an != bn {
				return an < bn
			}
			continue
		}
		if as[i] != bs[i] {
			return as[i] < bs[i]
		}
	}
	return len(as) < len(bs)
}

// idSegments splits an id into its numeric segments, dropping the kind prefix.
func idSegments(id string) []string {
	s := id
	if len(s) > 0 && s[0] >= 'A' && s[0] <= 'Z' {
		s = s[1:]
	}
	return strings.Split(s, ".")
}

// writeJSONAtomic writes JSON via a temp file, fsync and rename, so a crash leaves the
// previous or the new document, never a hybrid.
func writeJSONAtomic(path string, v any) error {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	raw = append(raw, '\n')

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("temp file in %s: %w", dir, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(raw); err != nil {
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
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("rename %s: %w", tmpName, err)
	}
	// fsync the directory so the rename itself is durable.
	if d, err := os.Open(dir); err == nil {
		_ = d.Sync()
		_ = d.Close()
	}
	return nil
}
