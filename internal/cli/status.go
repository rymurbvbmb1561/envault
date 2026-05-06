package cli

import "fmt"

// Status prints the current envault state for the project: whether a key and
// vault exist, and whether the vault is currently unlocked (plaintext present).
func (r *Runner) Status() error {
	r.printf("envault status\n")
	r.printf("  config  : %s\n", r.cfg.ConfigPath)

	if r.keys.Exists() {
		r.printf("  key     : found (%s)\n", r.cfg.KeyPath)
	} else {
		r.printf("  key     : not initialised — run 'envault init'\n")
	}

	if r.vault.Exists() {
		r.printf("  vault   : locked (%s)\n", r.cfg.VaultPath)
	} else {
		r.printf("  vault   : no vault found — run 'envault lock' to create one\n")
	}

	if r.vault.PlaintextExists() {
		r.printf("  .env    : unlocked (plaintext present)\n")
	} else {
		r.printf("  .env    : locked / absent\n")
	}

	return nil
}

// Version prints the application version string.
func (r *Runner) Version(version string) error {
	if version == "" {
		version = "dev"
	}
	r.printf("envault %s\n", version)
	return nil
}

// notInitialisedError is returned when a command requires initialisation.
type notInitialisedError struct{ msg string }

func (e *notInitialisedError) Error() string { return e.msg }

func errNotInitialised() error {
	return &notInitialisedError{msg: "envault is not initialised in this project — run 'envault init'"}
}

// requireInit returns an error if the keystore has not been set up yet.
func (r *Runner) requireInit() error {
	if !r.keys.Exists() {
		return errNotInitialised()
	}
	return nil
}

// requireVault returns an error if no encrypted vault exists yet.
func (r *Runner) requireVault() error {
	if err := r.requireInit(); err != nil {
		return err
	}
	if !r.vault.Exists() {
		return fmt.Errorf("no vault found — run 'envault lock' to encrypt your .env file first")
	}
	return nil
}
