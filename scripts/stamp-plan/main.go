// Command stamp-plan records a plan's id migration in its journal.
//
// A plan whose ids were renamed cannot have its journal rewritten — the journal is
// append-only — so the migration is *added* to it, as a control event carrying the
// version the plan now lives under.
//
// The event is written through the domain's operation, not by editing the file: the
// plan's stamp moves because the tool wrote it, not because a hand touched the JSON.
//
//	go run ./scripts/stamp-plan <plan-dir> <old-id-form> <new-id-form>
//
// For this repository's own plan:
//
//	go run ./scripts/stamp-plan kitchen R01 RCP01
package main

import (
	"fmt"
	"os"

	"github.com/lthiagol/maitre/internal/core"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: stamp-plan <plan-dir> <old-id-form> <new-id-form>")
		os.Exit(2)
	}
	dir, from, to := os.Args[1], os.Args[2], os.Args[3]

	store := core.Open(dir)
	ops := core.NewOperations(store, core.AuditContext{
		Actor: "runner", Session: "s-migrate", Role: "runner",
	})
	if err := ops.RecordMigration(from, to); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("the migration %s -> %s is recorded in %s\n", from, to, dir)
}
