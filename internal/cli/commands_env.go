package cli

const envUsage = `Usage: envault env [PREFIX] [--export]

Print decrypted secrets as shell-evaluable KEY=VALUE lines.

Options:
  PREFIX      Optional key prefix filter (e.g. AWS prints only AWS_* keys)
  --export    Prefix each line with "export " so variables are exported
              to child processes when the output is eval'd

Examples:
  # Load all secrets into the current shell
  eval $(envault env)

  # Load only AWS_* secrets, exported
  eval $(envault env AWS --export)
`

func init() {
	registerUsage("env", envUsage)
}
