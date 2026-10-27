package core

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// readDoc reads one of the project's design docs.
func readDoc(t *testing.T, name string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "docs", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

// sectionOf returns the part of a doc between a heading and the next one.
func sectionOf(doc, heading string) string {
	rest := ""
	if i := strings.Index(doc, heading); i >= 0 {
		rest = doc[i+len(heading):]
	}
	if i := strings.Index(rest, "\n## "); i >= 0 {
		rest = rest[:i]
	}
	return rest
}

// subsectionOf returns the part of a doc between a sub-heading and the next heading
// of any level.
//
// A sub-heading's section ends at the next `###` as well as the next `##` — without
// that, a check reads the tables of the sections that follow it.
func subsectionOf(doc, heading string) string {
	rest := sectionOf(doc, heading)
	if i := strings.Index(rest, "\n### "); i >= 0 {
		rest = rest[:i]
	}
	return rest
}

// checkDerivedViews compares the doc's table of derived views with the code's.
//
// The table is how a reader finds a view, and it drifted silently: the chef's moves
// gained a fourth one, the tree got no row, and two rows claimed views the domain
// does not compute. The check runs both ways — a view the code has and the table does
// not is a view a reader cannot find, and a row the table has and the code does not
// is a view that never happens.
//
// Each row names the function that computes it, and the name is checked against the
// code's own: a function that was renamed leaves the table naming something that is
// not there.
func checkDerivedViews(doc string) error {
	section := sectionOf(doc, "## 8. Derived views")
	if section == "" {
		return errNoSection("8. Derived views")
	}
	rows := map[string]docRow{}
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		cells := strings.Split(line, "|")
		if len(cells) < 4 {
			continue
		}
		rows[firstBackticked(cells[1])] = docRow{by: strings.TrimSpace(cells[2]), text: cells[3]}
	}
	if len(rows) == 0 {
		return errNoSection("8. Derived views")
	}

	var missing []string
	var extra []string
	for _, v := range DerivedViews {
		row, ok := rows[v.Name]
		if !ok {
			missing = append(missing, v.Name)
			continue
		}
		// The row says who computes it, and it says what the code says.
		if v.Planned {
			if row.by != "planned" {
				extra = append(extra, v.Name+": the code does not compute it, the doc says "+row.by)
			}
			continue
		}
		want := funcName(v.Function)
		if want == "" {
			missing = append(missing, v.Name+": it names no function the code has")
			continue
		}
		if row.by != "`"+want+"`" {
			missing = append(missing, v.Name+": "+want)
			continue
		}
		if len(v.Values) == 0 {
			continue
		}
		// The row's enumeration: the values the view can take.
		listed := map[string]bool{}
		for _, name := range allBackticked(row.text) {
			listed[name] = true
		}
		if len(listed) == 0 {
			return errNoSection("the " + v.Name + " row's enumeration")
		}
		for _, value := range v.Values {
			if value == "" {
				// The empty value has no name; the doc writes "none".
				continue
			}
			if !listed[value] {
				missing = append(missing, v.Name+": "+value)
			}
		}
		for name := range listed {
			if name == "none" {
				// The doc's name for the empty value.
				continue
			}
			if !containsString(v.Values, name) {
				extra = append(extra, v.Name+": "+name)
			}
		}
	}
	for name, row := range rows {
		v := derivedView(name)
		if v.Name == "" {
			extra = append(extra, name)
			continue
		}
		// A row that claims to be computed must name the function the code uses.
		if !v.Planned && row.by != "`"+funcName(v.Function)+"`" {
			extra = append(extra, name+": "+row.by)
		}
	}
	if len(missing) > 0 || len(extra) > 0 {
		return errMismatch{what: "derived views", missing: missing, extra: extra}
	}
	return nil
}

// docRow is one row of the doc's table of derived views.
type docRow struct {
	// by is the cell that names who computes the view.
	by string
	// text is the cell that defines the view.
	text string
}

// derivedView is the declared view by that name.
func derivedView(name string) DerivedView {
	for _, v := range DerivedViews {
		if v.Name == name {
			return v
		}
	}
	return DerivedView{}
}

// checkOpenQuestions checks that section 9 is findable and whole.
//
// Two of the questions have proposals, and a reader has to be able to reach them from
// the question. The questions that were there stay: an open question that disappears
// from the list is an answer nobody gave.
func checkOpenQuestions(doc, proposals string) error {
	section := sectionOf(doc, "## 9. Open questions")
	if section == "" {
		return errNoSection("9. Open questions")
	}

	// The links to the proposals, and the headings they point at.
	for _, want := range []string{"#1-courses", "#2-design_decisions"} {
		if !strings.Contains(section, "06-proposals.md"+want) {
			return errMismatch{what: "open questions", missing: []string{want}}
		}
		if !proposalsHaveAnchor(proposals, want) {
			return errMismatch{what: "the proposals", extra: []string{want}}
		}
	}

	// The question about questions.
	if !strings.Contains(section, "Should the plan have questions at all?") {
		return errMismatch{what: "open questions", missing: []string{"the question about questions"}}
	}

	// The questions that were there are still there.
	for _, q := range []string{
		"Should `ready` and `approved` both exist?",
		"Does `epic` need its own file type",
		"Does `program` (rung 4) need an identity",
		"Can the L5 session boundary be *enforced*",
		"Is the guidance `done_when` field worth its cost?",
		"What is the workspace directory called?",
		"How much of the charter (goals / non-goals / principles) belongs in `project.json`",
	} {
		if !strings.Contains(section, q) {
			return errMismatch{what: "open questions", missing: []string{q}}
		}
	}
	return nil
}

// proposalsHaveAnchor reports whether the proposals doc has a heading that makes the
// given anchor.
func proposalsHaveAnchor(doc, anchor string) bool {
	for _, line := range strings.Split(doc, "\n") {
		if !strings.HasPrefix(line, "#") {
			continue
		}
		if headingAnchor(strings.TrimLeft(line, "# ")) == anchor {
			return true
		}
	}
	return false
}

// headingAnchor is the link target a heading makes.
//
// A reader's editor builds it the same way: lowercase, the punctuation gone, the
// spaces dashes.
func headingAnchor(heading string) string {
	text := strings.ToLower(heading)
	for _, r := range []string{"`", ".", ",", ":", "?", "!", "(", ")", "'", `"`} {
		text = strings.ReplaceAll(text, r, "")
	}
	text = strings.ReplaceAll(text, " ", "-")
	return "#" + text
}

// checkViewCoverage compares the package's exported functions with the registry.
//
// The wire that checks the specification's table fires when the *doc* is edited. This
// one fires when a *view* is added: a function nobody registered computes something the
// specification does not know about, and the table is how a reader finds a view.
//
// Every exported function is either a view — a function that carries the `// view`
// marker — or a declared helper (`core.NonViews`). A function that is neither fails,
// and so does an exemption that no longer names a function.
//
// The whole package is read, not one file: the domain does not keep a file boundary
// for its views — one of them already lives outside views.go.
func checkViewCoverage(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	found := map[string]bool{}
	marked := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		for name, f := range exportedFuncs(string(raw)) {
			found[name] = true
			if f.view {
				marked[name] = true
			}
		}
	}

	// The registry: the views the specification's table is checked against.
	registered := map[string]bool{}
	for _, v := range DerivedViews {
		if v.Planned {
			continue
		}
		registered[funcName(v.Function)] = true
	}

	var unregistered []string
	var unregisteredView []string
	var unmarkedView []string
	var stale []string
	for name := range found {
		exempt := exemptFromView(name)
		isView := marked[name]
		isRegistered := registered[name]
		switch {
		case (isView || isRegistered) && exempt:
			// A view that is also a helper is judged twice, and the two answers
			// cannot both be right.
			unregistered = append(unregistered, name+" (a view and a helper at once)")
		case isView && !isRegistered:
			// A view the specification's table does not name.
			unregisteredView = append(unregisteredView, name)
		case isRegistered && !isView:
			// A registered view whose marker was deleted: the same judgement
			// made twice, the other way round. A function cannot be reclassified
			// by deleting a word.
			unmarkedView = append(unmarkedView, name)
		case !isView && !exempt:
			unregistered = append(unregistered, name)
		}
	}
	for _, name := range NonViews {
		if !familyMatchesAny(found, name) {
			stale = append(stale, name)
		}
	}
	if len(unregistered) > 0 || len(unregisteredView) > 0 || len(unmarkedView) > 0 || len(stale) > 0 {
		return errCoverage{
			unregistered:     unregistered,
			unregisteredView: unregisteredView,
			unmarkedView:     unmarkedView,
			stale:            stale,
		}
	}
	return nil
}

// exemptFromView reports whether a name is a declared helper.
//
// An exemption may name a function exactly, or a family of them: "Operations.*" is
// every method of the operations, "*.Label" is every method called Label.
func exemptFromView(name string) bool {
	for _, entry := range NonViews {
		if entry == name {
			return true
		}
		if !strings.Contains(entry, "*") {
			continue
		}
		if strings.HasPrefix(entry, "*") {
			if strings.HasSuffix(name, strings.TrimPrefix(entry, "*")) {
				return true
			}
			continue
		}
		if strings.HasPrefix(name, strings.TrimSuffix(entry, "*")) {
			return true
		}
	}
	return false
}

// familyMatchesAny reports whether an exemption still names a function.
//
// A family matches when any function it covers is there.
func familyMatchesAny(found map[string]bool, entry string) bool {
	if !strings.Contains(entry, "*") {
		return found[entry]
	}
	if strings.HasPrefix(entry, "*") {
		suffix := strings.TrimPrefix(entry, "*")
		for name := range found {
			if strings.HasSuffix(name, suffix) {
				return true
			}
		}
		return false
	}
	prefix := strings.TrimSuffix(entry, "*")
	for name := range found {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// exportedFuncs finds the exported functions and methods in a Go source file, and
// whether each carries the view marker.
//
// A method is qualified by its receiver's type — "Criterion.Status" — so it cannot be
// confused with a function of the same name.
func exportedFuncs(source string) map[string]scannedFunc {
	out := map[string]scannedFunc{}
	comment := ""
	for _, line := range strings.Split(source, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") {
			comment = trimmed
			continue
		}
		if !strings.HasPrefix(trimmed, "func ") {
			if trimmed != "" {
				comment = ""
			}
			continue
		}
		name, ok := declaredName(trimmed)
		if ok {
			out[name] = scannedFunc{view: comment == "// view"}
		}
		comment = ""
	}
	return out
}

// scannedFunc is one exported function the scan found.
type scannedFunc struct {
	// view is true when the function carries the `// view` marker.
	view bool
}

// declaredName is the name a declaration introduces, and whether it is exported.
//
// A method is qualified by its receiver's type — "Criterion.Status" — so it cannot be
// confused with a function of the same name.
func declaredName(line string) (string, bool) {
	rest := line[len("func "):]
	if strings.HasPrefix(rest, "(") {
		// A method: "func (r Recipe) Name(".
		end := strings.Index(rest, ")")
		if end < 0 {
			return "", false
		}
		receiver := strings.TrimSpace(rest[1:end])
		name := declaredIdentifier(rest[end+1:])
		if !isOperationName(name) {
			return "", false
		}
		typ := receiver
		if i := strings.LastIndex(typ, " "); i >= 0 {
			typ = typ[i+1:]
		}
		return strings.TrimPrefix(typ, "*") + "." + name, true
	}
	name := declaredIdentifier(rest)
	if !isOperationName(name) {
		return "", false
	}
	return name, true
}

// declaredIdentifier is the name a declaration introduces.
func declaredIdentifier(rest string) string {
	if i := strings.IndexAny(rest, "(["); i >= 0 {
		rest = rest[:i]
	}
	return strings.TrimSpace(rest)
}

// errCoverage is the failure of the coverage check.
type errCoverage struct {
	// unregistered are the functions that are neither a view nor a declared helper.
	unregistered []string
	// unregisteredView are the views the specification's table does not name.
	unregisteredView []string
	// unmarkedView are the registered views whose marker was deleted.
	unmarkedView []string
	// stale are the declared helpers that no longer name a function.
	stale []string
}

func (e errCoverage) Error() string {
	var b strings.Builder
	b.WriteString("the package's exported functions and the view registry disagree")
	if len(e.unregistered) > 0 {
		b.WriteString(": " + strings.Join(e.unregistered, ", ") +
			" — no view registers them and no exemption names them")
	}
	if len(e.unregisteredView) > 0 {
		if len(e.unregistered) > 0 {
			b.WriteString("; ")
		}
		b.WriteString(strings.Join(e.unregisteredView, ", ") + " — views the specification's table does not name")
	}
	if len(e.unmarkedView) > 0 {
		if len(e.unregistered) > 0 || len(e.unregisteredView) > 0 {
			b.WriteString("; ")
		}
		b.WriteString(strings.Join(e.unmarkedView, ", ") + " — registered views with no marker")
	}
	if len(e.stale) > 0 {
		if len(e.unregistered) > 0 || len(e.unregisteredView) > 0 || len(e.unmarkedView) > 0 {
			b.WriteString("; ")
		}
		b.WriteString(strings.Join(e.stale, ", ") + " — declared helpers that no longer name a function")
	}
	return b.String()
}

// checkLifecycles compares the doc's state table with the code's lifecycles.
//
// The states are a closed set, and a reader who cannot enumerate them cannot tell
// where a recipe is.
func checkLifecycles(doc string) error {
	section := subsectionOf(doc, "### 3.1 States")
	if section == "" {
		return errNoSection("3.1 States")
	}
	listed := map[string]bool{}
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		listed[firstBackticked(strings.Split(line, "|")[1])] = true
	}
	if len(listed) == 0 {
		return errNoSection("3.1 States")
	}

	var missing []string
	var extra []string
	for _, l := range AllLifecycles {
		if !listed[string(l)] {
			missing = append(missing, string(l))
		}
	}
	for name := range listed {
		if !containsString(lifecycleNames(), name) {
			extra = append(extra, name)
		}
	}
	if len(missing) > 0 || len(extra) > 0 {
		return errMismatch{what: "lifecycles", missing: missing, extra: extra}
	}
	return nil
}

// lifecycleNames are the lifecycles, as names.
func lifecycleNames() []string {
	out := make([]string, 0, len(AllLifecycles))
	for _, l := range AllLifecycles {
		out = append(out, string(l))
	}
	return out
}

// checkGateTable compares the doc's gate table with the code's gate table.
//
// The doc says "one table, one place, one definition each". A gate the doc does not
// list is a gate a reader cannot find, and a gate the doc lists that the code does not
// have is a gate that never runs.
func checkGateTable(doc string) error {
	section := sectionOf(doc, "## 4. Gates")
	listed := map[string]bool{}
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		// Only the first column: the gate's name. The rest of the row is prose.
		listed[firstBackticked(strings.Split(line, "|")[1])] = true
	}
	if len(listed) == 0 {
		return errNoSection("4. Gates")
	}

	var missing []string
	for _, g := range GateNames() {
		if !listed[g] {
			missing = append(missing, g)
		}
	}
	var extra []string
	for name := range listed {
		if !containsString(GateNames(), name) {
			extra = append(extra, name)
		}
	}
	if len(missing) > 0 || len(extra) > 0 {
		return errMismatch{what: "gates", missing: missing, extra: extra}
	}
	return nil
}

// checkJournalKinds compares the doc's closed-enum sentence with the code's kinds.
//
// The doc says an unknown kind is an error. A reader cannot enumerate the enum from
// the doc if the doc's list is short.
func checkJournalKinds(doc string) error {
	listed := map[string]bool{}
	for _, para := range strings.Split(doc, "\n\n") {
		if !strings.Contains(para, "closed enum**:") {
			continue
		}
		// Only the sentence that declares the enum. The words in the sentences
		// around it are prose, not kinds.
		for _, name := range allBackticked(sentenceWith(para, "closed enum")) {
			if name == "kind" {
				continue
			}
			listed[name] = true
		}
	}
	if len(listed) == 0 {
		return errNoSection("the closed enum of journal kinds")
	}

	var missing []string
	for _, k := range JournalKinds {
		if !listed[k] {
			missing = append(missing, k)
		}
	}
	var extra []string
	for name := range listed {
		if !containsString(JournalKinds, name) {
			extra = append(extra, name)
		}
	}
	if len(missing) > 0 || len(extra) > 0 {
		return errMismatch{what: "journal kinds", missing: missing, extra: extra}
	}
	return nil
}

// checkOperationMarks compares the doc's operation marks with the core's operations.
//
// The doc's tables mark each operation ✅ (built) or ⏳ (not built). A reader who
// trusts the mark must not be misled.
func checkOperationMarks(doc string) error {
	section := sectionOf(doc, "## The operations")
	rows := 0
	built := operationMethods()
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| ") {
			continue
		}
		mark := ""
		switch {
		case strings.Contains(line, "✅"):
			mark = "built"
		case strings.Contains(line, "⏳"):
			mark = "planned"
		default:
			continue
		}
		rows++
		for _, name := range operationNames(strings.Split(line, "|")[1]) {
			switch mark {
			case "built":
				if !built[name] {
					return errMismatch{what: "operations", extra: []string{name}}
				}
			case "planned":
				if built[name] {
					return errMismatch{what: "operations", missing: []string{name}}
				}
			}
		}
	}
	if rows == 0 {
		return errNoSection("The operations")
	}
	return nil
}

// operationNames reads the operation names in a table cell.
//
// Only the names: a cell that also carries prose ("`GetRecipe(id)` — the id is the
// id") names one operation, not two.
func operationNames(cell string) []string {
	var out []string
	for _, name := range allBackticked(cell) {
		name = strings.TrimSpace(strings.SplitN(name, "(", 2)[0])
		if !isOperationName(name) {
			continue
		}
		out = append(out, name)
	}
	return out
}

// isOperationName reports whether the word looks like an operation's name.
//
// The core's operations are exported Go identifiers, so a name starts with a capital
// letter. A stray word in the cell is not an operation, and reading it as one would
// make the tripwire cry wolf on a reworded sentence.
func isOperationName(name string) bool {
	if name == "" {
		return false
	}
	return name[0] >= 'A' && name[0] <= 'Z'
}

// sentenceWith returns the sentence that holds the given text.
func sentenceWith(text, needle string) string {
	i := strings.Index(text, needle)
	if i < 0 {
		return ""
	}
	start := strings.LastIndex(text[:i], ".") + 1
	end := strings.Index(text[i:], ".")
	if end < 0 {
		return text[start:]
	}
	return text[start : i+end]
}

// operationMethods is the set of operations the core exposes.
func operationMethods() map[string]bool {
	out := map[string]bool{}
	t := reflect.TypeOf(&Operations{})
	for i := 0; i < t.NumMethod(); i++ {
		out[t.Method(i).Name] = true
	}
	return out
}

// firstBackticked is the first backticked word in a text.
func firstBackticked(text string) string {
	names := allBackticked(text)
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

// allBackticked is every backticked word in a text.
func allBackticked(text string) []string {
	var out []string
	for {
		i := strings.Index(text, "`")
		if i < 0 {
			return out
		}
		text = text[i+1:]
		j := strings.Index(text, "`")
		if j < 0 {
			return out
		}
		out = append(out, text[:j])
		text = text[j+1:]
	}
}

// containsString reports whether the list holds the value.
func containsString(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// errNoSection is the failure of a tripwire that found nothing to read.
type errNoSection string

func (e errNoSection) Error() string {
	return "the doc has no " + string(e) + " section — the tripwire would check nothing"
}

// errMismatch is the failure of a tripwire that found drift.
type errMismatch struct {
	what    string
	missing []string
	extra   []string
}

func (e errMismatch) Error() string {
	var b strings.Builder
	b.WriteString("the doc and the code disagree about the " + e.what)
	if len(e.missing) > 0 {
		b.WriteString(": the doc does not list " + strings.Join(e.missing, ", "))
	}
	if len(e.extra) > 0 {
		if len(e.missing) > 0 {
			b.WriteString("; ")
		}
		b.WriteString("the doc lists " + strings.Join(e.extra, ", ") + ", which the code does not have")
	}
	return b.String()
}

// TestTheDocListsTheGatesTheCodeHas keeps the gate table honest.
func TestTheDocListsTheGatesTheCodeHas(t *testing.T) {
	if err := checkGateTable(readDoc(t, "03-domain-model.md")); err != nil {
		t.Fatal(err)
	}
}

// TestTheDocListsTheDerivedViews keeps the table of derived views honest.
//
// It drifted once already: the chef's moves gained a fourth one and the tree got no
// row, and nothing said so.
func TestTheDocListsTheDerivedViews(t *testing.T) {
	if err := checkDerivedViews(readDoc(t, "03-domain-model.md")); err != nil {
		t.Fatal(err)
	}
}

// TestTheDocListsTheLifecycles keeps the state table honest.
func TestTheDocListsTheLifecycles(t *testing.T) {
	if err := checkLifecycles(readDoc(t, "03-domain-model.md")); err != nil {
		t.Fatal(err)
	}
}

// TestTheInterfaceNamesTheTriage keeps the footer paragraph honest.
//
// The console has a fourth chef's move — the triage — and the interface doc listed
// three.
func TestTheInterfaceNamesTheTriage(t *testing.T) {
	doc := readDoc(t, "05-human-interface.md")
	for _, para := range strings.Split(doc, "\n\n") {
		if !strings.Contains(para, "footer is pinned") {
			continue
		}
		// The paragraph wraps over several lines.
		if strings.Contains(para, "triage") {
			return
		}
	}
	t.Fatal("the footer paragraph should name the triage")
}

// TestTheDerivedViewsTripwireFailsOnRealDrift keeps the tripwire from being turned
// off.
//
// A tripwire that cannot fail is not a tripwire. The drift this one guards against
// happened once: the chef's moves gained a fourth one and the tree got no row, and
// the spec contradicted itself.
func TestTheDerivedViewsTripwireFailsOnRealDrift(t *testing.T) {
	doc := readDoc(t, "03-domain-model.md")

	cases := []struct {
		name  string
		doc   string
		check func(string) error
	}{
		{"a chef move removed", strings.Replace(doc, " · `triage`", "", 1), checkDerivedViews},
		{"a chef move invented", strings.Replace(doc, " · `serve`", " · `serve` · `invented`", 1), checkDerivedViews},
		{"a lane removed", strings.Replace(doc, "`grooming` · `approval`", "`approval`", 1), checkDerivedViews},
		{"a lane invented", strings.Replace(doc, "`blocked` · `done`", "`blocked` · `invented` · `done`", 1), checkDerivedViews},
		{"a review state removed", strings.Replace(doc, " · `remediated`", "", 1), checkDerivedViews},
		{"a review state invented", strings.Replace(doc, " · `remediated`", " · `remediated` · `invented`", 1), checkDerivedViews},
		{"the tree's row deleted", strings.Replace(doc, "| `tree` | `TreeOf` | The plan as a dependency tree |\n", "", 1), checkDerivedViews},
		{"a view invented", strings.Replace(doc, "| `tree` |", "| `invented` |", 1), checkDerivedViews},
		{"a view's function mistyped", strings.Replace(doc, "| `lane` | `LaneOf` |", "| `lane` | `LaneOff` |", 1), checkDerivedViews},
		{"the triage's function mistyped", strings.Replace(doc, "| `triage` | `TriageOf` |", "| `triage` | `TriageOff` |", 1), checkDerivedViews},
		{"a triage value dropped", strings.Replace(doc, " · `On hold`", "", 1), checkDerivedViews},
		{"a triage value invented", strings.Replace(doc, " · `Served`", " · `Served` · `invented`", 1), checkDerivedViews},
		{"a function that does not exist", strings.Replace(doc, "| `lane` | `LaneOf` |", "| `lane` | `NoSuchFunction` |", 1), checkDerivedViews},
		{"a planned view relabelled as computed", strings.Replace(doc, "| `coverage` | planned |", "| `coverage` | `Coverage` |", 1), checkDerivedViews},
		{"a computed view relabelled as planned", strings.Replace(doc, "| `tree` | `TreeOf` |", "| `tree` | planned |", 1), checkDerivedViews},
		{"a lifecycle removed", strings.Replace(doc, "| `remediation` | 9 |", "| `invented` | 9 |", 1), checkLifecycles},
		{"a lifecycle invented", strings.Replace(doc, "| `reviewed` | 10 |", "| `invented` | 10 |", 1), checkLifecycles},
		{"an event removed", strings.Replace(doc, "draft --> groomed : groom", "draft --> groomed : nothing", 1), checkDiagram},
		{"a journal kind removed", strings.Replace(doc, " · `step`", "", 1), checkJournalKinds},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.doc == doc {
				t.Fatal("the scratch copy is the same as the real doc — the corruption did not apply")
			}
			if err := c.check(c.doc); err == nil {
				t.Fatal("the tripwire should have tripped")
			}
		})
	}
}

// TestTheDerivedViewsTripwireFailsOnAnEmptySection keeps a renamed heading from
// silencing the tripwire.
//
// A renamed heading used to leave a tripwire with nothing to read and green: a safe
// unknown value disabling a whole feature, with the tests passing.
func TestTheDerivedViewsTripwireFailsOnAnEmptySection(t *testing.T) {
	doc := readDoc(t, "03-domain-model.md")
	renamed := strings.Replace(doc, "## 8. Derived views", "## 8. Views and things", 1)
	if err := checkDerivedViews(renamed); err == nil {
		t.Fatal("the tripwire should fail when the section is gone")
	}
}

// TestADeclaredViewThatNamesNoFunctionFails covers the declaration itself.
//
// A view that points at a function the code does not have computes nothing, and the
// wire says so — the doc's row cannot be right either, because there is no function
// behind it.
func TestADeclaredViewThatNamesNoFunctionFails(t *testing.T) {
	doc := readDoc(t, "03-domain-model.md")
	orig := DerivedViews
	t.Cleanup(func() { DerivedViews = orig })

	views := append([]DerivedView{}, orig...)
	views[2].Function = "NoSuchFunction"
	DerivedViews = views

	if err := checkDerivedViews(doc); err == nil {
		t.Fatal("the wire should fail when a view names a function that is not there")
	}
}

// TestTheOpenQuestionsLinkToTheProposals keeps section 9 findable and whole.
//
// Two of the questions have proposals, and the questions themselves must survive: an
// open question that vanishes is an answer nobody gave.
func TestTheOpenQuestionsLinkToTheProposals(t *testing.T) {
	if err := checkOpenQuestions(readDoc(t, "03-domain-model.md"), readDoc(t, "06-proposals.md")); err != nil {
		t.Fatal(err)
	}
}

// TestTheOpenQuestionsTripwireFailsOnRealDrift keeps the section 9 wire honest.
func TestTheOpenQuestionsTripwireFailsOnRealDrift(t *testing.T) {
	doc := readDoc(t, "03-domain-model.md")
	proposals := readDoc(t, "06-proposals.md")

	cases := []struct {
		name string
		doc  string
	}{
		{"a proposal link removed", strings.Replace(doc, "06-proposals.md#1-courses", "06-proposals.md", 1)},
		{"a proposal link pointing nowhere", strings.Replace(doc, "06-proposals.md#2-design_decisions", "06-proposals.md#9-nowhere", 1)},
		{"the question about questions removed", strings.Replace(doc, "Should the plan have questions at all?", "Should the plan have questions at all", 1)},
		{"an open question removed", strings.Replace(doc, "Does `epic` need its own file type", "Does the rung need its own file type", 1)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.doc == doc {
				t.Fatal("the scratch copy is the same as the real doc — the corruption did not apply")
			}
			if err := checkOpenQuestions(c.doc, proposals); err == nil {
				t.Fatal("the tripwire should have tripped")
			}
		})
	}

	t.Run("a proposal heading renamed", func(t *testing.T) {
		broken := strings.Replace(proposals, "## 1. Courses", "## 1. Coursez", 1)
		if broken == proposals {
			t.Fatal("the corruption did not apply")
		}
		if err := checkOpenQuestions(doc, broken); err == nil {
			t.Fatal("the tripwire should have tripped when the anchor no longer resolves")
		}
	})
}

// TestEveryExportedFunctionIsRegisteredOrDeclared covers the hole F-63 left.
//
// The wire that checks the specification's table fires when the doc is edited. This
// one fires when a *view* is added: a function nobody registered computes something
// the specification does not know about.
func TestEveryExportedFunctionIsRegisteredOrDeclared(t *testing.T) {
	if err := checkViewCoverage("."); err != nil {
		t.Fatal(err)
	}
}

// scratchPackage copies the package's non-test sources into a scratch directory, so a
// test can change the code without changing the code.
func scratchPackage(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		raw, err := os.ReadFile(e.Name())
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, e.Name()), raw, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// TestAViewNobodyRegisteredFails is the reviewer's demonstration.
//
// A view added to the package and registered nowhere must fail: the hole is not that
// the doc was edited, but that a view exists and no list knows it.
func TestAViewNobodyRegisteredFails(t *testing.T) {
	for _, file := range []string{"model.go", "operations.go", "store.go", "lifecycle.go"} {
		t.Run("added to "+file, func(t *testing.T) {
			scratch := scratchPackage(t)
			path := filepath.Join(scratch, file)
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			added := []string{
				"",
				"// A view nobody registered.",
				"func AViewNobodyRegistered(r Recipe) string {",
				"\treturn \"nothing\"",
				"}",
				"",
			}
			if err := os.WriteFile(path, append(append([]byte{}, raw...), []byte(strings.Join(added, "\n"))...), 0o644); err != nil {
				t.Fatal(err)
			}
			err = checkViewCoverage(scratch)
			if err == nil {
				t.Fatal("the check should fail when a view is added and registered nowhere")
			}
			if !strings.Contains(err.Error(), "AViewNobodyRegistered") {
				t.Errorf("the failure should name the function: %v", err)
			}
		})
	}
}

// TestAMarkedViewThatIsAlsoExemptFails covers the contradiction.
//
// A function that says it is a view and is also declared a helper is judged twice, and
// the two answers cannot both be right. It must not pass quietly.
func TestAMarkedViewThatIsAlsoExemptFails(t *testing.T) {
	scratch := scratchPackage(t)
	path := filepath.Join(scratch, "views.go")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// A view that is also on the exemption list.
	added := []string{
		"",
		"// view",
		"func LaneNames() []string {",
		"\treturn nil",
		"}",
		"",
	}
	if err := os.WriteFile(path, append(append([]byte{}, raw...), []byte(strings.Join(added, "\n"))...), 0o644); err != nil {
		t.Fatal(err)
	}
	err = checkViewCoverage(scratch)
	if err == nil {
		t.Fatal("the check should fail when a function is a view and a helper at once")
	}
	if !strings.Contains(err.Error(), "a view and a helper at once") {
		t.Errorf("the failure should say the contradiction: %v", err)
	}
}

// TestAStaleExemptionFails keeps the exemption list honest.
//
// An exemption for a function that is gone is a hole in the other direction: the
// list says a name is a helper, and no such function exists. A family that matches
// nothing is the same hole, wider.
func TestAStaleExemptionFails(t *testing.T) {
	t.Run("a stale name", func(t *testing.T) {
		scratch := scratchPackage(t)
		path := filepath.Join(scratch, "views.go")
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		// A helper removed from the file, and its name left in the exemption list.
		stripped := strings.Replace(string(raw), "func KnownLane(name string) bool {", "func knownLane(name string) bool {", 1)
		if stripped == string(raw) {
			t.Fatal("the corruption did not apply")
		}
		if err := os.WriteFile(path, []byte(stripped), 0o644); err != nil {
			t.Fatal(err)
		}
		err = checkViewCoverage(scratch)
		if err == nil {
			t.Fatal("the check should fail when an exemption no longer names a function")
		}
		if !strings.Contains(err.Error(), "KnownLane") {
			t.Errorf("the failure should name the stale exemption: %v", err)
		}
	})

	t.Run("a family that matches nothing", func(t *testing.T) {
		orig := NonViews
		t.Cleanup(func() { NonViews = orig })
		NonViews = append(append([]string{}, orig...), "Nothing.*")

		err := checkViewCoverage(".")
		if err == nil {
			t.Fatal("the check should fail when a family matches nothing")
		}
		if !strings.Contains(err.Error(), "Nothing.*") {
			t.Errorf("the failure should name the stale family: %v", err)
		}
	})
}

// TestAViewAsAMethodOfRecipeFails is the demonstration to beat.
//
// A view added as a method of Recipe, with no marker, must fail: the read-side
// methods are named one by one, so a family cannot admit a view nobody classified.
func TestAViewAsAMethodOfRecipeFails(t *testing.T) {
	scratch := scratchPackage(t)
	path := filepath.Join(scratch, "model.go")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	added := []string{
		"",
		"// A view nobody registered.",
		"func (r Recipe) Coverage() string {",
		"\treturn \"nothing\"",
		"}",
		"",
	}
	if err := os.WriteFile(path, append(append([]byte{}, raw...), []byte(strings.Join(added, "\n"))...), 0o644); err != nil {
		t.Fatal(err)
	}
	err = checkViewCoverage(scratch)
	if err == nil {
		t.Fatal("the check should fail when a view is added as a method and registered nowhere")
	}
	if !strings.Contains(err.Error(), "Recipe.Coverage") {
		t.Errorf("the failure should name the method: %v", err)
	}
}

// TestTheJudgementIsSymmetric covers the four cases.
//
// A function is a view if and only if it is marked and registered. The judgement is
// made twice — once by the marker, once by the registry — and the two answers cannot
// differ in either direction.
func TestTheJudgementIsSymmetric(t *testing.T) {
	cases := []struct {
		name string
		// change is the edit the scratch copy gets.
		change func(t *testing.T, dir string)
		want   string
	}{
		{
			name: "marked and not registered",
			change: func(t *testing.T, dir string) {
				addFunction(t, dir, "views.go", "// view\nfunc AViewNobodyRegistered(r Recipe) string {")
			},
			want: "the specification's table does not name",
		},
		{
			name: "marked and exempt",
			change: func(t *testing.T, dir string) {
				markFunction(t, dir, "views.go", "func LaneNames() []string {")
			},
			want: "a view and a helper at once",
		},
		{
			name: "registered and not marked",
			change: func(t *testing.T, dir string) {
				unmarkFunction(t, dir, "views.go", "func LaneOf(r Recipe) Lane {")
			},
			want: "registered views with no marker",
		},
		{
			name: "registered and exempt",
			change: func(t *testing.T, dir string) {
				// The exemption list is the code's own, so it is changed there
				// and restored after.
				orig := NonViews
				t.Cleanup(func() { NonViews = orig })
				NonViews = append(append([]string{}, orig...), "LaneOf")
			},
			want: "a view and a helper at once",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			scratch := scratchPackage(t)
			c.change(t, scratch)
			err := checkViewCoverage(scratch)
			if err == nil {
				t.Fatal("the check should fail when the judgement is not symmetric")
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("the failure should say %q: %v", c.want, err)
			}
		})
	}
}

// addFunction adds a function to a file in the scratch copy.
func addFunction(t *testing.T, dir, file, decl string) {
	t.Helper()
	path := filepath.Join(dir, file)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	added := []string{"", decl, "}", ""}
	if err := os.WriteFile(path, append(append([]byte{}, raw...), []byte(strings.Join(added, "\n"))...), 0o644); err != nil {
		t.Fatal(err)
	}
}

// markFunction gives a function the view marker.
func markFunction(t *testing.T, dir, file, decl string) {
	t.Helper()
	path := filepath.Join(dir, file)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	patched := strings.Replace(string(raw), decl, "// view\n"+decl, 1)
	if patched == string(raw) {
		t.Fatal("the declaration was not found")
	}
	if err := os.WriteFile(path, []byte(patched), 0o644); err != nil {
		t.Fatal(err)
	}
}

// unmarkFunction takes the view marker off a function.
func unmarkFunction(t *testing.T, dir, file, decl string) {
	t.Helper()
	path := filepath.Join(dir, file)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	patched := strings.Replace(string(raw), "// view\n"+decl, decl, 1)
	if patched == string(raw) {
		t.Fatal("the marker was not found")
	}
	if err := os.WriteFile(path, []byte(patched), 0o644); err != nil {
		t.Fatal(err)
	}
}

// exemptFunction adds a name to the exemption list.
func exemptFunction(t *testing.T, dir, file, name string) {
	t.Helper()
	path := filepath.Join(dir, file)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	patched := strings.Replace(string(raw), "var NonViews = []string{", "var NonViews = []string{\n\t\""+name+"\",", 1)
	if patched == string(raw) {
		t.Fatal("the exemption list was not found")
	}
	if err := os.WriteFile(path, []byte(patched), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestTheDocListsTheJournalKinds keeps the closed enum honest.
func TestTheDocListsTheJournalKinds(t *testing.T) {
	if err := checkJournalKinds(readDoc(t, "03-domain-model.md")); err != nil {
		t.Fatal(err)
	}
}

// TestTheDocsOperationMarksMatchTheCode keeps the operation contract honest.
func TestTheDocsOperationMarksMatchTheCode(t *testing.T) {
	if err := checkOperationMarks(readDoc(t, "04-operations.md")); err != nil {
		t.Fatal(err)
	}
}

// TestTheTripwiresFailOnRealDrift keeps the tripwires from being turned off.
//
// A tripwire that cannot fail is not a tripwire.
func TestTheTripwiresFailOnRealDrift(t *testing.T) {
	gates := readDoc(t, "03-domain-model.md")
	ops := readDoc(t, "04-operations.md")

	cases := []struct {
		name  string
		doc   string
		check func(string) error
	}{
		{"a gate renamed", strings.Replace(gates, "`intent.present`", "`intent.there`", 1), checkGateTable},
		{"a gate invented", strings.Replace(gates, "| `intent.present` |", "| `intent.invented` |", 1), checkGateTable},
		{"a journal kind removed", strings.Replace(gates, " · `step`", "", 1), checkJournalKinds},
		{"a journal kind invented", strings.Replace(gates, " · `step`", " · `step` · `invented`", 1), checkJournalKinds},
		{"an operation mark flipped", strings.Replace(ops, "✅ `CreateSpecial`", "⏳ `CreateSpecial`", 1), checkOperationMarks},
		{"an operation marked built that does not exist",
			strings.Replace(ops, "✅ `CreateSpecial`", "✅ `CreateSpecial` · `Invented`", 1), checkOperationMarks},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.check(c.doc); err == nil {
				t.Fatal("the tripwire should have tripped")
			}
		})
	}
}

// TestTheTripwiresFailOnAnEmptySection keeps a renamed heading from silencing a
// tripwire.
//
// Renaming the heading used to leave the tripwire with nothing to read and green:
// a safe unknown value disabling a whole feature, with the tests passing.
func TestTheTripwiresFailOnAnEmptySection(t *testing.T) {
	ops := readDoc(t, "04-operations.md")
	renamed := strings.Replace(ops, "## The operations", "## The verbs", 1)
	if err := checkOperationMarks(renamed); err == nil {
		t.Fatal("the tripwire should fail when the section is gone")
	}
}

// TestTheTripwiresSurviveARewordedSentence keeps a tripwire from crying wolf.
//
// The kinds are read from the sentence that declares the enum, and the operations from
// the names in a cell. A word that is neither is prose, and prose must not trip the
// wire.
func TestTheTripwiresSurviveARewordedSentence(t *testing.T) {
	gates := readDoc(t, "03-domain-model.md")
	ops := readDoc(t, "04-operations.md")

	// A backticked word in the sentence beside the enum.
	reworded := strings.Replace(gates,
		"Unknown kinds are an error (C11), never silently ignored.",
		"Unknown kinds are an error (`null` is not one).", 1)
	if err := checkJournalKinds(reworded); err != nil {
		t.Fatalf("a reworded sentence should not trip the tripwire: %v", err)
	}

	// A backticked word in a row's first column.
	marked := strings.Replace(ops,
		"| ✅ `GetRecipe(id)` |",
		"| ✅ `GetRecipe(id)` (`id` is the id) |", 1)
	if err := checkOperationMarks(marked); err != nil {
		t.Fatalf("a reworded cell should not trip the tripwire: %v", err)
	}
}

// TestTheDiagramAgreesWithTheTransitionTable keeps the lifecycle's graph honest.
//
// The diagram is the spec's statement of the lifecycle: every edge is labelled with
// the event and the gates that guard it. A gate the edge leaves out reads as a
// transition nobody guards.
func TestTheDiagramAgreesWithTheTransitionTable(t *testing.T) {
	if err := checkDiagram(readDoc(t, "03-domain-model.md")); err != nil {
		t.Fatal(err)
	}
}

// checkDiagram compares the lifecycle diagram with the transition table.
//
// Both ways: every event in the code has an edge with the same from, the same to and
// the same gates, and every edge in the diagram is an event in the code.
func checkDiagram(doc string) error {
	edges := diagramEdges(doc)
	if len(edges) == 0 {
		return errNoSection("the lifecycle diagram")
	}

	// Every event in the code has an edge.
	for _, ev := range AllEvents {
		want := transitions[ev]
		for _, from := range want.from {
			found := false
			for _, e := range edges {
				if e.event != string(ev) || e.from != string(from) || e.to != string(want.to) {
					continue
				}
				found = true
				if !sameStrings(e.gates, want.gates) {
					return errMismatch{what: "diagram edge " + string(ev) + " from " + string(from),
						extra: []string{"gates " + strings.Join(e.gates, ", ")}}
				}
			}
			if !found {
				return errMismatch{what: "diagram edges",
					missing: []string{string(ev) + " from " + string(from)}}
			}
		}
	}

	// Every edge in the diagram is an event in the code.
	for _, e := range edges {
		ev, ok := eventByName(e.event)
		if !ok {
			return errMismatch{what: "diagram edges", extra: []string{e.event}}
		}
		want := transitions[ev]
		if e.to != string(want.to) {
			return errMismatch{what: "diagram edge " + e.event, extra: []string{"to " + e.to}}
		}
		if !containsLifecycle(want.from, Lifecycle(e.from)) {
			return errMismatch{what: "diagram edge " + e.event, extra: []string{"from " + e.from}}
		}
		if !sameStrings(e.gates, want.gates) {
			return errMismatch{what: "diagram edge " + e.event, extra: []string{"gates " + strings.Join(e.gates, ", ")}}
		}
	}
	return nil
}

// diagramEdge is one edge of the lifecycle diagram.
type diagramEdge struct {
	event string
	from  string
	to    string
	gates []string
}

// diagramEdges reads the lifecycle diagram's edges.
func diagramEdges(doc string) []diagramEdge {
	var out []diagramEdge
	edge := regexp.MustCompile(`(\w+) --> (\w+) : (\w+)(<br/>(.*))?`)
	for _, m := range edge.FindAllStringSubmatch(sectionOf(doc, "### 3.2 Transitions"), -1) {
		var gates []string
		for _, g := range strings.Split(strings.ReplaceAll(m[5], "<br/>", " · "), " · ") {
			g = strings.TrimSpace(g)
			if g != "" {
				gates = append(gates, g)
			}
		}
		out = append(out, diagramEdge{event: m[3], from: m[1], to: m[2], gates: gates})
	}
	return out
}

// eventByName is the event with the given name.
func eventByName(name string) (Event, bool) {
	for _, ev := range AllEvents {
		if string(ev) == name {
			return ev, true
		}
	}
	return "", false
}

// sameStrings reports whether two lists hold the same values in the same order.
func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestTheDiagramCheckTripsOnDrift keeps the diagram check honest.
//
// A check that cannot fail is not a check: the diagram loses an edge, or an edge loses
// a gate, and the spec drifts from the code without a word.
func TestTheDiagramCheckTripsOnDrift(t *testing.T) {
	doc := readDoc(t, "03-domain-model.md")

	cases := []struct {
		name string
		doc  string
	}{
		{"an edge removed", strings.Replace(doc,
			"ready --> groomed : request_changes<br/>changes.recorded\n", "", 1)},
		{"a gate dropped from an edge", strings.Replace(doc,
			"review.verdict_ok · review.reviewer_independent · review.clean", "review.clean", 1)},
		{"an edge invented", strings.Replace(doc,
			"    ready --> approved : approve<br/>approval.by_human\n",
			"    ready --> approved : approve<br/>approval.by_human\n    draft --> closed : close<br/>review.clean\n", 1)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := checkDiagram(c.doc); err == nil {
				t.Fatal("the diagram check should have tripped")
			}
		})
	}
}

// TestTheSpecAgreesWithTheCodeAboutReopen keeps one event from being described three
// ways.
//
// The diagram, the gate table and the transition table must say the same thing about
// the gates the reopen runs.
func TestTheSpecAgreesWithTheCodeAboutReopen(t *testing.T) {
	doc := readDoc(t, "03-domain-model.md")

	// The code's answer.
	want := transitions[EventReopen].gates
	if len(want) != 0 {
		t.Fatalf("the code runs %v on reopen; the doc must say the same", want)
	}

	// The diagram's account.
	gates := diagramGates(doc, string(EventReopen))
	if len(gates) != 0 {
		t.Errorf("the diagram says the reopen runs %v; the code runs none", gates)
	}

	// The gate table's account: no gate may claim the reopen.
	for _, gate := range GateNames() {
		for _, event := range appliesTo(doc, gate) {
			if event == string(EventReopen) {
				t.Errorf("the doc says %q applies to the reopen; the code runs it on no gate", gate)
			}
		}
	}
}

// diagramGates reads the gates on an event's edge in the lifecycle diagram.
func diagramGates(doc, event string) []string {
	edge := regexp.MustCompile(`--> \w+ : ` + regexp.QuoteMeta(event) + `(<br/>(.*))?`)
	for _, m := range edge.FindAllStringSubmatch(sectionOf(doc, "### 3.2 Transitions"), -1) {
		var gates []string
		for _, g := range strings.Split(strings.ReplaceAll(m[2], "<br/>", " · "), " · ") {
			g = strings.TrimSpace(g)
			if g != "" {
				gates = append(gates, g)
			}
		}
		return gates
	}
	return nil
}

// appliesTo reads the gate's "Applies to" column from the doc's gate table.
func appliesTo(doc, gate string) []string {
	section := sectionOf(doc, "## 4. Gates")
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| `"+gate+"` |") {
			continue
		}
		cells := strings.Split(line, "|")
		if len(cells) < 4 {
			return nil
		}
		var out []string
		for _, e := range strings.Split(cells[3], ",") {
			e = strings.Trim(strings.TrimSpace(e), "`")
			if e != "" {
				out = append(out, e)
			}
		}
		return out
	}
	return nil
}

// TestTheDocDoesNotOverstateTheDenyList keeps the spec honest about what the
// permission check can do.
//
// The list matches the actor's name, so an agent that calls itself `chef` passes it.
func TestTheDocDoesNotOverstateTheDenyList(t *testing.T) {
	doc := readDoc(t, "03-domain-model.md")
	for _, claim := range []string{
		"cannot sign a recipe off at all",
		"cannot be satisfied by an agent",
	} {
		if strings.Contains(doc, claim) {
			t.Errorf("the doc claims %q — the deny-list matches a name", claim)
		}
	}
	if strings.Contains(doc, "cancel.reason_present") {
		t.Error("the doc names cancel.reason_present as a gate — the reason is checked inline")
	}
}

// checkDenyList checks the spec against the deny-list the code holds.
//
// It is a function of the doc, so the tripwire can be given a broken doc and be seen
// to fail.
func checkDenyList(doc string) error {
	for _, role := range AgentRoles {
		// The spec writes a role the way it is spoken: "sous chef", not "sous-chef".
		spoken := strings.ReplaceAll(role, "-", " ")
		if !strings.Contains(doc, role) && !strings.Contains(doc, spoken) {
			return fmt.Errorf("the doc does not name %q — the code refuses it as an agent", role)
		}
	}
	for name := range retiredAgentRoles {
		if !strings.Contains(doc, "`"+name+"`") {
			return fmt.Errorf("the doc does not name the retired role %q — a plan written before the rename still meets it", name)
		}
	}
	return nil
}

// TestTheDocNamesTheRolesTheCodeRefuses keeps the deny-list and the spec together.
//
// The list is a deny-list of names, and the spec names them in prose. Nothing said
// the two had to agree, so a role renamed in the code left the spec telling the
// reader about a role the tool no longer refuses — and a reader who believes the
// spec would sign off on the wrong name.
func TestTheDocNamesTheRolesTheCodeRefuses(t *testing.T) {
	if err := checkDenyList(readDoc(t, "03-domain-model.md")); err != nil {
		t.Fatal(err)
	}
}

// TestTheDenyListTripwireFailsOnRealDrift keeps the tripwire from being turned off.
//
// A tripwire that cannot fail is not a tripwire.
func TestTheDenyListTripwireFailsOnRealDrift(t *testing.T) {
	doc := readDoc(t, "03-domain-model.md")
	// The spec without the runner: the code refuses it, the spec does not name it.
	broken := strings.ReplaceAll(doc, "runner", "line cook")
	if broken == doc {
		t.Fatal("the doc does not hold the word the drift removes")
	}
	if err := checkDenyList(broken); err == nil {
		t.Error("the tripwire should fail on a spec that does not name the runner")
	}
}
