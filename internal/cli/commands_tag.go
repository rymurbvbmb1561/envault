package cli

func init() {
	registerHelp("tag", tagHelp)
}

const tagHelp = `tag — manage metadata tags on secret keys

Usage:
  envault tag add KEY tag1 [tag2 ...]   attach one or more tags to KEY
  envault tag remove KEY tag1 [tag2 ...]  remove tags from KEY
  envault tag list KEY                  list tags attached to KEY

Tags are stored alongside secrets in the encrypted vault and are useful
for grouping or annotating keys (e.g. "production", "ci", "deprecated").

Examples:
  envault tag add DATABASE_URL production
  envault tag add API_KEY production ci
  envault tag remove API_KEY ci
  envault tag list DATABASE_URL
`
