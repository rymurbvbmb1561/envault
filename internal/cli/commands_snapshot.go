package cli

// commands_snapshot.go ensures the snapshot, snapshot-restore, and
// snapshot-list commands are registered via their respective init() calls
// in snapshot.go and snapshot_restore.go.
//
// No additional wiring is required; registration happens automatically
// through the init() functions in those files.
