# homefim

`homefim` contains two independent Go command-line applications for generating
and checking a SHA-256 file-integrity baseline:

- `generate-baseline` records hashes for the current user's shell startup files
  (`.bash_profile`, `.bash_login`, `.profile`, `.bashrc`, and `.bash_logout`)
  and files under `~/.ssh`. It skips `known_hosts` and `known_hosts.old`.
- `secure_fim` checks the paths and hashes in the embedded baseline. It reports
  missing files, hash mismatches, and file read errors.

The baseline contains file paths and hashes, not file contents. Review the
generated `internal/baseline/baseline.json` before building or distributing the
scanner.

## Requirements

- Go 1.22 or newer
- Run baseline generation from the repository root

## Generate a baseline

Generate a baseline for the current user:

```sh
go run ./cmd/generate-baseline
```

This writes `internal/baseline/baseline.json`. The scanner embeds that file at
build time, so regenerate the scanner binary after changing the baseline.

## Build

Build each application independently:

```sh
mkdir -p bin
go build -o bin/generate-baseline ./cmd/generate-baseline
go build -o bin/secure_fim ./cmd/secure-fim
```

Or generate a baseline and then build both applications with the default Make
target:

```sh
make
```

The Makefile also provides `make generate`, `make build`,
`make build-generator`, `make build-scanner`, and `make clean`.

## Run the scanner

By default, `secure_fim` performs one check and exits, making it suitable for
scheduled execution such as cron:

```sh
./bin/secure_fim
```

To run continuously and check every 30 seconds, pass `--loop`:

```sh
./bin/secure_fim --loop
```

The scanner uses the baseline embedded in its binary; it does not load the
JSON file at runtime. When it detects a mismatch or a missing path, it writes
a diagnostic to the log. File read errors are also logged.

## Test

Run the Go package tests from the repository root:

```sh
go test ./...
```
