package cli

// doctor command registration is handled via init() in doctor.go.
// This file documents the command for reference.
//
// Usage:
//
//	envault doctor
//
// Runs a series of environment checks and reports whether the project is
// correctly initialised, the key and vault files are present, an editor
// is configured, and a clipboard tool is available.
//
// Exit codes:
//
//	0 — all checks passed
//	1 — one or more checks failed
