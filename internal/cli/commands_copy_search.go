package cli

// This file ensures copy and search commands are registered via their
// respective init() functions in copy.go and search.go.
//
// Both commands are read-only operations on the vault:
//
//   copy  <KEY>            - copy a secret value to the system clipboard
//   search <QUERY> [--values] - search keys (and optionally values) by substring
//
// They follow the same Runner-based pattern as other commands and require
// the vault to be initialised and locked before use.
