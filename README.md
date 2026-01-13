# linode-objl+

**linode-objl+** is a Go CLI that streamlines Linode Object Storage operations. It aims to improve basic object manipulation with **recursive** capabilities for `cp` (copy) and `rm` (remove), along with a fast `ls` (list) command.

This tool focuses on practical day-to-day tasks: listing bucket contents, copying objects (including whole prefixes), and removing objects safely and efficiently.

---

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Global Flags](#global-flags)
  - [`ls` — List objects](#ls--list-objects)
  - [`cp` — Copy objects (recursive supported)](#cp--copy-objects-recursive-supported)
  - [`rm` — Remove objects (recursive supported)](#rm--remove-objects-recursive-supported)
- [Examples](#examples)
- [Performance & Tuning](#performance--tuning)
- [Completion Scripts](#completion-scripts)
- [Output & Logging](#output--logging)
- [Exit Codes](#exit-codes)
- [Development](#development)
- [Contributing](#contributing)
- [Roadmap](#roadmap)
- [License](#license)

---

## Features

- **S3-compatible** operations against Linode Object Storage.
- **Recursive `cp` and `rm`** for bulk operations on prefixes.
- **Human-friendly `ls`** with size, modified time, and prefix filtering.
- **Safe-by-default deletes** (prompt/confirmation and `--dry-run`).
- **Parallel transfers** with configurable concurrency.
- **Consistent, script-friendly output** (JSON and plain).

---

## Quick Start

```bash
# 1) Configure the tool with your credentials
linode-objl+ configure  --profile

# 2) Install the CLI
go install github.com/your-org/linode-objlplus@latest

# 3) Try listing a bucket
linode-objl+ ls objs://my-bucket/photos/

## Installation

### Prebuilt binaries

Download the appropriate binary for your platform from Releases, then place it on your `PATH` as `linode-objl+`.

### From source (Go >= 1.21)

```bash
git clone https://github.com/your-org/linode-objlplus.git
cd linode-objlplus
make build
# or
go build -o linode-objl+
```

---

## Configuration

You can configure credentials and defaults via environment variables or a config file.

### Config file

By default, the CLI looks for a file at:
- `~/.linodeobjim/credentials.ini` (Linux/macOS)
- `%USERPROFILE%/.linodeobjim/credentials.init` (Windows)

Example `credentials.init`:

```ini
[default]
token  = YOUR_LINODE_TOKEN
region = YOUR_OBJ_CLUSTER_ID
```

---

## Usage

```text
linode-objl+ [command] [args] [flags]
```

### Global Flags

- `-p, --profile <name>` — use a specific profile
- `--endpoint <url>` — override S3 endpoint
- `--region <name>` — override region
- `-v, --verbose` — increase logging verbosity
- `--json` — output results in JSON
- `--dry-run` — show what would happen without making changes

### `ls` — List objects

List objects in a bucket or prefix.

```bash
linode-objl+ ls objs://bucket[/prefix] [--long] [--recursive] [--delimiter "/"] [--max-items N]
```

**Flags**
- `--long` — show size and modified time
- `--recursive` — traverse all child prefixes
- `--delimiter` — group by delimiter (like folders)
- `--max-items` — limit number of listed items

### `cp` — Copy objects (recursive supported)

Copy single objects or entire prefixes. Supports cross-bucket copying.

```bash
linode-objl+ cp SOURCE TARGET [--recursive] [--include PATTERN] [--exclude PATTERN] \
  [--acl private|public-read] [--storage-class SC] [--concurrency N] [--dry-run]
```

**Examples**
- Copy one object: `linode-objl+ cp objs://a/img.jpg objs://b/img.jpg`
- Copy a prefix: `linode-objl+ cp objs://a/photos/ objs://b/photos/ --recursive`

**Flags**
- `--recursive` — copy all objects under the source prefix
- `--include` / `--exclude` — glob-like filters
- `--acl` — object ACL (e.g., `public-read`)
- `--storage-class` — storage class (if supported)
- `--concurrency` — number of parallel transfers
- `--dry-run` — show planned operations

### `rm` — Remove objects (recursive supported)

Safely delete single objects or entire prefixes.

```bash
linode-objl+ rm objs://bucket/key | objs://bucket/prefix/ [--recursive] [--force] [--include PATTERN] [--exclude PATTERN] [--batch N] [--dry-run]
```

**Flags**
- `--recursive` — delete all objects under the prefix
- `--force` — skip confirmation prompts
- `--include` / `--exclude` — filter objects
- `--batch` — batch size for multi-object delete API
- `--dry-run` — show planned deletions without executing

---

## Examples

List a bucket root with details:

```bash
linode-objl+ ls objs://my-bucket --long
```

List recursively under a prefix:

```bash
linode-objl+ ls objs://my-bucket/logs/ --recursive
```

Copy all PNGs from one prefix to another, 16 concurrent transfers:

```bash
linode-objl+ cp objs://src/assets/ objs://dst/assets/ --recursive --include "*.png" --concurrency 16
```

Dry-run a destructive delete to verify scope:

```bash
linode-objl+ rm objs://my-bucket/tmp/ --recursive --dry-run
```

Force delete a specific object:

```bash
linode-objl+ rm objs://my-bucket/old/report.csv --force
```

---

## Performance & Tuning

- **Concurrency**: Increase `--concurrency` for faster copies; tune based on network and API limits.
- **Batch deletes**: Use `--batch` for efficient multi-object deletion (where supported by API).
- **Pagination**: Large listings can be paged with `--max-items`; combine with `--json` for scripts.
- **Retries**: The CLI performs exponential backoff on transient failures.

---

## Completion Scripts

Generate shell completions:

```bash
# Bash
linode-objl+ completion bash > /etc/bash_completion.d/linode-objl+

# Zsh
linode-objl+ completion zsh > ~/.zsh/completions/_linode-objl+

# Fish
linode-objl+ completion fish > ~/.config/fish/completions/linode-objl+.fish
```

---

## Output & Logging

- Default output is human-readable.
- Use `--json` for structured output suited for automation.
- Increase verbosity with `-v` / `--verbose`.

---

## Exit Codes

- `0` — Success
- `1` — General error
- `2` — Invalid arguments or configuration
- `3` — Authentication or authorization failure
- `4` — Network or endpoint errors

---

## Development

Requires Go >= 1.21.

```bash
git clone https://github.com/your-org/linode-objlplus.git
cd linode-objlplus
make test
make build
```

---

## Contributing

Issues and PRs are welcome! Please:

1. Open an issue describing the feature/bug.
2. Write tests for new functionality.
3. Follow Go formatting and linting (`gofmt`, `golangci-lint`).

---

## Roadmap

- Implement basic commands options and flags
- Implement tunning for uploads and downloads
- Multipart copy with automatic part sizing.
- Sync mode (`sync` subcommand) for bi-directional reconciliation.
- Server-side encryption options and KMS integration.
- Enhanced progress bars with throughput stats.
- More robust include/exclude with regex.

---

## License

MIT (see `LICENSE` file).

---

## Disclaimer

This tool is not an official Linode product. Use at your own risk. Validate destructive operations with `--dry-run` before proceeding.
