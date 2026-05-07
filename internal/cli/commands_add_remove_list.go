package cli

// This file wires the add, remove, and list sub-commands into the Runner's
// command dispatch table defined in commands.go.
//
// Each command is registered via the same pattern used by init/lock/unlock so
// that run.go can route to them without additional changes.

func init() {
	// Registration happens through the Runner.Dispatch method; the actual
	// routing is done in run.go using a switch statement.  This file exists
	// to document the grouping and to hold any future shared helpers for
	// the three mutating vault commands.
}

// dispatchMutating routes add / remove / list commands.
// Called from Runner.Dispatch when the primary switch falls through.
func (r *Runner) dispatchMutating(cmd string, args []string) (bool, error) {
	switch cmd {
	case "add":
		return true, r.addCmd(args)
	case "remove", "rm":
		return true, r.removeCmd(args)
	case "list", "ls":
		return true, r.listCmd(args)
	}
	return false, nil
}
