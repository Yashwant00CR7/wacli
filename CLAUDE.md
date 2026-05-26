# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

- **Build**: `pnpm build` (builds to `dist/wacli`) or `CGO_ENABLED=1 CGO_CFLAGS="-Wno-error=missing-braces" go build -tags sqlite_fts5 -o dist/wacli ./cmd/wacli`
- **Test all**: `pnpm test`
- **Test Go only**: `pnpm test:go` (uses `go test ./...`)
- **Test FTS**: `pnpm test:fts` (uses `go test -tags sqlite_fts5 ./...`)
- **Test single test**: `go test -tags sqlite_fts5 -v ./internal/store -run TestSearchMessages`
- **Lint**: `pnpm lint` (uses `go vet ./...`)
- **Format**: `pnpm format` (runs `gofmt -w .`)
- **Format check**: `pnpm format:check` (runs `gofmt -l .`)
- **Database generate**: `pnpm generate:sqlc` (regenerates sqlc models/queries in `internal/store/storedb`)

## Store Statistics

Current store stats:
- **Messages**: 15,102
- **Chats**: 444
- **Groups**: 140

Run `wacli store stats` to view current statistics.

## Codebase Architecture

- **CLI Layer (`cmd/wacli/`)**: CLI endpoints built with Cobra. Handles parsing flags/arguments, initializing App, and outputting JSON/human formatting.
- **Coordination Layer (`internal/app/`)**: Coordinates syncing, events dispatching, media download queues, webhooks, and core logic.
- **WhatsApp Client (`internal/wa/`)**: Wraps `whatsmeow` library. Manages QR/phone auth pairing, connection state, decryption, and media transport.
- **Store Layer (`internal/store/`)**: SQLite store. Uses `sqlc` for typed query models. Includes schema (`schema.sql`), queries (`sqlc/queries.sql`), and Full-Text Search (FTS5).
- **Concurrency Lock (`internal/lock/`)**: Manages SQLite/Store multi-process locks to prevent concurrent write corruption.

## Store Location

Default store: `D:\Users\yashw\.wacli` (moved from `C:\Users\yashw\.wacli`)

The `WACLI_STORE_DIR` environment variable is set to point to this location. All `wacli` commands use this directory automatically.

To change store location:
```bash
# Move existing store
Move-Item "D:\Users\yashw\.wacli" "NEW_PATH"

# Set new location
[System.Environment]::SetEnvironmentVariable("WACLI_STORE_DIR", "NEW_PATH", "User")
```

## Sync Behavior

`wacli sync` syncs from the last sync point, not from the beginning. It connects to WhatsApp's servers and receives new messages in real-time.

### Sync Modes

| Mode | Flag | Behavior |
|------|------|----------|
| **Once** (default) | none or `--once` | Syncs history + new messages, exits after 30s idle |
| **Follow** | `--follow` | Continuous sync - only new messages after connection |

### Key Notes

- **History sync**: The "Processing history sync (X conversations)..." message shows LID (Local ID) migration for existing conversations - it's NOT downloading all message history
- **New messages only**: Use `wacli sync --follow` for continuous live syncing of new messages
- **Storage limits**: Use `--max-messages` or `--max-db-size` to bound local history growth
- **On-demand backfill**: Use `wacli history` command to fetch specific chat history

### Common Commands

```bash
# Sync new messages continuously
wacli sync --follow

# Sync once and exit (with history migration)
wacli sync

# Sync with storage limits
wacli sync --max-messages 10000 --max-db-size 500MB
```
