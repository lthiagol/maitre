// Package cli implements the agent-facing command surface: the `maitre-cli` binary.
//
// Every command is a thin adapter over an operation in the core package. The CLI owns no
// domain rules and never touches the plan files directly, so it cannot drift from the
// console — which calls the same operations.
//
// See docs/decisions/0003-surface-split.md and docs/decisions/0007-operation-layer.md.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lthiagol/maitre/internal/core"
)

// Exit codes.
const (
	ExitOK    = 0
	ExitError = 1
	ExitGate  = 2
)

// Runner executes one command.
type Runner struct {
	Args   []string
	Stdout io.Writer
	Stderr io.Writer

	// The identity given before the command, in the documented global form:
	// `maitre-cli --actor line-cook recipe claim R01`.
	actor   string
	session string
	role    string
}

// Run dispatches a command and returns the process exit code.
func (r *Runner) Run() int {
	if len(r.Args) == 0 {
		r.usage()
		return ExitError
	}
	args, id, err := splitIdentity(r.Args)
	if err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	r.actor, r.session, r.role = id.Actor, id.Session, id.Role
	if len(args) == 0 {
		r.usage()
		return ExitError
	}
	switch args[0] {
	case "help", "-h", "--help":
		r.usage()
		return ExitOK
	case "init":
		return r.cmdInit(args[1:])
	case "demo":
		return r.cmdDemo(args[1:])
	case "gates":
		return r.cmdGates(args[1:])
	case "plan":
		return r.cmdPlan(args[1:])
	case "path":
		return r.cmdPath(args[1:])
	case "log":
		return r.cmdLog(args[1:])
	case "recipe":
		return r.dispatch(args[1:], "recipe", map[string]func([]string) int{
			"create":          r.cmdRecipeCreate,
			"list":            r.cmdRecipeList,
			"show":            r.cmdRecipeShow,
			"update":          r.cmdRecipeUpdate,
			"advance":         r.cmdRecipeAdvance,
			"approve":         r.cmdRecipeApprove,
			"request-changes": r.cmdRecipeRequestChanges,
			"claim":           r.cmdRecipeClaim,
			"document":        r.cmdRecipeDocument,
			"review":          r.cmdRecipeReview,
			"comment":         r.cmdRecipeComment,
			"block":           r.cmdRecipeBlock,
			"unblock":         r.cmdRecipeUnblock,
			"defer":           r.cmdRecipeDefer,
			"resume":          r.cmdRecipeResume,
			"cancel":          r.cmdRecipeCancel,
		})
	case "ac":
		return r.dispatch(args[1:], "ac", map[string]func([]string) int{
			"add":     r.cmdACAdd,
			"verify":  r.cmdACVerify,
			"signoff": r.cmdACSignOff,
			"defer":   r.cmdACDefer,
		})
	case "wp":
		return r.dispatch(args[1:], "wp", map[string]func([]string) int{
			"add": r.cmdWPAdd,
		})
	case "step":
		return r.dispatch(args[1:], "step", map[string]func([]string) int{
			"add":     r.cmdStepAdd,
			"start":   r.cmdStepStart,
			"done":    r.cmdStepDone,
			"skip":    r.cmdStepSkip,
			"fail":    r.cmdStepFail,
			"deviate": r.cmdStepDeviate,
		})
	case "finding":
		return r.dispatch(args[1:], "finding", map[string]func([]string) int{
			"add":     r.cmdFindingAdd,
			"resolve": r.cmdFindingResolve,
			"dismiss": r.cmdFindingDismiss,
			"reopen":  r.cmdFindingReopen,
		})
	case "special":
		return r.dispatch(args[1:], "special", map[string]func([]string) int{
			"list":    r.cmdSpecialList,
			"create":  r.cmdSpecialCreate,
			"update":  r.cmdSpecialUpdate,
			"dismiss": r.cmdSpecialDismiss,
			"promote": r.cmdSpecialPromote,
		})
	case "prepnote":
		return r.dispatch(args[1:], "prepnote", map[string]func([]string) int{
			"list":    r.cmdPrepNoteList,
			"create":  r.cmdPrepNoteCreate,
			"update":  r.cmdPrepNoteUpdate,
			"promote": r.cmdPrepNotePromote,
		})
	default:
		r.fail("unknown command", args[0])
		return ExitError
	}
}

// dispatch routes a subcommand of a group.
func (r *Runner) dispatch(args []string, group string, verbs map[string]func([]string) int) int {
	if len(args) == 0 {
		r.fail(group+" needs a verb", strings.Join(sortedKeys(verbs), ", "))
		return ExitError
	}
	verb, ok := verbs[args[0]]
	if !ok {
		r.fail("unknown "+group+" verb", args[0],
			"known verbs: "+strings.Join(sortedKeys(verbs), ", "))
		return ExitError
	}
	return verb(args[1:])
}

// sortedKeys gives a stable order for error messages.
func sortedKeys(m map[string]func([]string) int) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sortStrings(out)
	return out
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

func (r *Runner) usage() {
	fmt.Fprint(r.Stdout, `maitre-cli — the maitre agent CLI

Usage:
  maitre-cli init                create the plan directory
  maitre-cli demo [--chaos]      seed a plan to look at
  maitre-cli plan                the whole plan, one document
  maitre-cli path                the execution order
  maitre-cli log [--limit N]     the journal
  maitre-cli gates               the gate table
  maitre-cli recipe <verb>       the recipe's own state
  maitre-cli ac <verb>           acceptance criteria
  maitre-cli step <verb>         guidance steps
  maitre-cli finding <verb>      review findings
  maitre-cli special <verb>      the ladder's rung 0
  maitre-cli prepnote <verb>     the ladder's rung 1

Identity (who is acting, and in which session):
  --actor <name>    default: $SCALE_ACTOR, else "agent"
  --session <id>    default: $SCALE_SESSION, else a fresh id
  --role <name>     default: $SCALE_ROLE

Every command prints one JSON document: on stdout when it succeeds, on stderr when it
fails. This help is the one exception — it prints text.
`)
}

// ---------------------------------------------------------------------------
// Identity
// ---------------------------------------------------------------------------

// identity is who is acting. Every write is stamped with it.
type identity struct {
	Actor   string
	Session string
	Role    string
}

// splitIdentity pulls the identity flags out of the front of the argument list.
//
// The documented form puts them before the command — `maitre-cli --actor line-cook recipe
// claim R01` — and they may also come after the verb, where the verb's own flag set
// reads them. The first argument that is not an identity flag ends the scan.
func splitIdentity(args []string) ([]string, identity, error) {
	var id identity
	i := 0
	for i < len(args) {
		a := args[i]
		if !strings.HasPrefix(a, "--") {
			break
		}
		name := strings.TrimPrefix(a, "--")
		value := ""
		hasValue := false
		if eq := strings.Index(name, "="); eq >= 0 {
			name, value, hasValue = name[:eq], name[eq+1:], true
		}

		var target *string
		switch name {
		case "actor":
			target = &id.Actor
		case "session":
			target = &id.Session
		case "role":
			target = &id.Role
		default:
			return args[i:], id, nil
		}
		if !hasValue {
			if i+1 >= len(args) {
				return nil, id, fmt.Errorf("--%s needs a value", name)
			}
			i++
			value = args[i]
		}
		*target = value
		i++
	}
	return args[i:], id, nil
}

// resolveIdentity reads the actor, session and role from the flags or the environment.
//
// The flags come first: the verb's own flag set, then the global form, then the
// environment.
//
// The session must come from the caller: CPD's L5 law is about *sessions*, and a
// session that spans several commands is what makes the check meaningful. An agent that
// runs ten commands in one task should set SCALE_SESSION once.
func resolveIdentity(fs *flagSet, global identity) identity {
	id := identity{
		Actor:   firstNonEmpty(fs.StringValue("actor"), global.Actor, os.Getenv("SCALE_ACTOR"), "agent"),
		Session: firstNonEmpty(fs.StringValue("session"), global.Session, os.Getenv("SCALE_SESSION"), newSessionID()),
		Role:    firstNonEmpty(fs.StringValue("role"), global.Role, os.Getenv("SCALE_ROLE")),
	}
	return id
}

// newSessionID makes a session id for a caller who did not supply one.
func newSessionID() string {
	return fmt.Sprintf("s-%d-%d", time.Now().UnixNano(), os.Getpid())
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// Output
// ---------------------------------------------------------------------------

// emit prints one JSON document.
func (r *Runner) emit(v any) int {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(r.Stderr, "marshal output: %v\n", err)
		return ExitError
	}
	fmt.Fprintln(r.Stdout, string(raw))
	return ExitOK
}

// fail prints an error document.
func (r *Runner) fail(msg string, details ...string) {
	payload := map[string]any{"ok": false, "error": msg}
	if len(details) > 0 {
		payload["details"] = details
	}
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Fprintf(r.Stderr, "%s\n", msg)
		return
	}
	fmt.Fprintln(r.Stderr, string(raw))
}

// failOp prints the error from an operation, and the exit code it deserves.
//
// A gate failure is exit 2 — the command was understood, the plan said no. Anything else
// is exit 1.
func (r *Runner) failOp(err error) int {
	var ge *core.GateError
	if errorsAs(err, &ge) {
		gates := make([]string, 0, len(ge.Gates))
		for _, g := range ge.Gates {
			gates = append(gates, g)
		}
		raw, _ := json.MarshalIndent(map[string]any{
			"ok":    false,
			"error": err.Error(),
			"event": string(ge.Event),
			"gates": gates,
		}, "", "  ")
		fmt.Fprintln(r.Stderr, string(raw))
		return ExitGate
	}
	r.fail("failed", err.Error())
	return ExitError
}

// ---------------------------------------------------------------------------
// Store access
// ---------------------------------------------------------------------------

// openStore finds the plan and returns it with the operations bound to the identity.
func (r *Runner) openStore(fs *flagSet) (*core.Store, *core.Operations, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}
	root, err := core.FindRoot(cwd)
	if err != nil {
		return nil, nil, fmt.Errorf("%w — run `maitre-cli init` first", err)
	}
	store := core.Open(root)
	id := resolveIdentity(fs, identity{Actor: r.actor, Session: r.session, Role: r.role})
	return store, core.NewOperations(store, core.AuditContext{
		Actor: id.Actor, Session: id.Session, Role: id.Role,
	}), nil
}

// absDir resolves a path to an absolute directory.
func absDir(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

// ---------------------------------------------------------------------------
// Workspace commands
// ---------------------------------------------------------------------------

func (r *Runner) cmdInit(args []string) int {
	fs := newFlagSet("init")
	dir := fs.String("dir", ".", "directory to create the plan in")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	store := core.Open(filepath.Join(absDir(*dir), core.PlanDirName))
	if err := store.Init(); err != nil {
		r.fail("init failed", err.Error())
		return ExitError
	}
	return r.emit(map[string]any{"ok": true, "root": store.Root})
}

func (r *Runner) cmdDemo(args []string) int {
	fs := newFlagSet("demo")
	dir := fs.String("dir", ".", "directory to seed the plan in")
	chaos := fs.Bool("chaos", false, "also seed the pathological recipes")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	store := core.Open(filepath.Join(absDir(*dir), core.PlanDirName))
	seed := core.Seed
	if *chaos {
		seed = core.SeedStress
	}
	if err := seed(store); err != nil {
		r.fail("seed failed", err.Error())
		return ExitError
	}
	return r.emit(map[string]any{"ok": true, "root": store.Root, "next": "maitre"})
}

func (r *Runner) cmdGates(args []string) int {
	fs := newFlagSet("gates")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	return r.emit(map[string]any{"ok": true, "gates": core.GateNames()})
}

// cmdPlan prints the whole plan in one document — the cheapest read for an agent.
func (r *Runner) cmdPlan(args []string) int {
	fs := newFlagSet("plan")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	store, _, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	recipes, err := store.ListRecipes()
	if err != nil {
		r.fail("list recipes", err.Error())
		return ExitError
	}
	specials, _ := store.ListSpecials()
	notes, _ := store.ListPrepNotes()
	project, _ := store.LoadProject()
	journal, _ := store.LoadJournal()

	return r.emit(map[string]any{
		"ok":        true,
		"project":   project,
		"recipes":   recipes,
		"specials":  specials,
		"prepNotes": notes,
		"events":    len(journal.Events),
	})
}

// cmdPath prints the execution order.
func (r *Runner) cmdPath(args []string) int {
	fs := newFlagSet("path")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	store, _, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	recipes, err := store.ListRecipes()
	if err != nil {
		r.fail("list recipes", err.Error())
		return ExitError
	}
	ready, waiting := core.Path(recipes)
	return r.emit(map[string]any{"ok": true, "ready": ready, "waiting": waiting})
}

// cmdLog prints the journal, newest first.
func (r *Runner) cmdLog(args []string) int {
	fs := newFlagSet("log")
	limit := fs.Int("limit", 50, "how many events")
	since := fs.Int("since", 0, "only events after this sequence number")
	if err := fs.parse(args); err != nil {
		r.fail("bad flags", err.Error())
		return ExitError
	}
	store, _, err := r.openStore(fs)
	if err != nil {
		r.fail("no plan", err.Error())
		return ExitError
	}
	journal, err := store.LoadJournal()
	if err != nil {
		r.fail("load journal", err.Error())
		return ExitError
	}
	if *since > 0 {
		return r.emit(map[string]any{"ok": true, "events": journal.Since(*since)})
	}
	return r.emit(map[string]any{"ok": true, "events": journal.Recent(*limit)})
}

// atoi parses an int argument, reporting a clean error.
func atoi(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%q is not a number", s)
	}
	return n, nil
}
