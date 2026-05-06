// Package config manages the per-project envault configuration file
// (.envault.json). It provides helpers to write, load, and locate the
// configuration by walking up the directory tree from the current working
// directory, mirroring the behaviour of tools such as git.
//
// Typical usage:
//
//	// On 'envault init'
//	cfg := config.DefaultConfig()
//	cfg.CreatedAt = time.Now().UTC().Format(time.RFC3339)
//	config.Write(".envault.json", cfg)
//
//	// On any subsequent command
//	cfg, root, err := config.LoadFromCwd()
package config
