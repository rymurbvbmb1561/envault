# envault

A lightweight local secrets manager that encrypts `.env` files using [age](https://github.com/FiloSottile/age) encryption for safe project storage.

---

## Installation

```bash
go install github.com/yourusername/envault@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envault.git
cd envault && go build -o envault .
```

---

## Usage

**Encrypt a `.env` file:**

```bash
envault encrypt .env
```

This produces a `.env.age` file that is safe to commit to version control.

**Decrypt when needed:**

```bash
envault decrypt .env.age -o .env
```

**Run a command with secrets injected directly (no plaintext file written):**

```bash
envault run .env.age -- go run main.go
```

On first use, envault generates a key stored at `~/.config/envault/key`. Keep this key backed up — it is required to decrypt your secrets.

---

## Why envault?

- No external services or accounts required
- Uses [age](https://age-encryption.org/) — a modern, audited encryption format
- Works offline and integrates into any workflow
- Tiny binary with zero runtime dependencies

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)