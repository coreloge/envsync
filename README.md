# envsync

> A utility to reconcile and audit environment variable drift between local, staging, and production configs.

---

## Installation

```bash
go install github.com/yourname/envsync@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envsync.git && cd envsync && go build -o envsync .
```

---

## Usage

Compare environment files across environments and surface any missing or mismatched keys:

```bash
# Audit drift between local and production
envsync audit --base .env.local --target .env.production

# Reconcile staging against local (dry run)
envsync reconcile --base .env.local --target .env.staging --dry-run

# Output a diff report to a file
envsync audit --base .env.local --target .env.production --output report.txt
```

**Example output:**

```
[MISSING]  DATABASE_URL        found in .env.production, missing in .env.local
[MISMATCH] LOG_LEVEL           local=debug | production=error
[OK]       APP_PORT            consistent across both configs
```

---

## Commands

| Command       | Description                                      |
|---------------|--------------------------------------------------|
| `audit`       | Report on key drift between two config files     |
| `reconcile`   | Sync missing keys from base into target          |
| `validate`    | Check all required keys are present in a file    |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)