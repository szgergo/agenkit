# Agenkit — Copilot Instructions

A Go CLI tool that manages AI agent configs (MCP servers, hooks, instructions) from a single YAML manifest and applies them to every AI provider (Claude, Copilot, Cursor, Codex, etc.).

## Developer Context

The developer is **experienced in programming but new to Go**. When writing any Go code, always explain:
- What each Go-specific construct does and why it's used (e.g., pointer receivers, struct tags, `init()`, `:=`)
- Any Go idioms or patterns that differ from common languages (e.g., error handling, interfaces, package visibility)
- Why a particular approach is idiomatic Go vs alternatives

Do not over-explain general programming concepts (loops, conditionals, data structures) — focus explanations on Go-specific behavior.

## Project Docs

- **Design & architecture**: `implementation/DESIGN.md`
- **Go, Cobra, Bubbletea, YAML conventions**: `implementation/shared-context.md` — read this before writing any code
- **Per-feature implementation specs**: `implementation/step-XX-*.md`

## Directory Layout

```
cmd/agenkit/main.go      # Entry point — tiny, wires cli.Execute() only
internal/
  cli/                   # Cobra command definitions
  manifest/              # YAML data model + parser
  provider/              # Provider interface + per-provider generators
  importer/              # Import from existing provider configs
  registry/              # MCP registry client
  lock/                  # Lock file management
  sync/                  # Git sync operations
  tui/                   # Bubbletea TUI components
  template/              # Variable interpolation engine
testdata/                # Test fixtures (manifests/, provider_configs/)
```

## Build & Test

```bash
go run ./cmd/agenkit          # Run during development
go build -o agenkit ./cmd/agenkit
go test ./...                 # Run all tests
go test -v -run TestXxx ./internal/manifest/
go mod tidy                   # After adding/removing imports
```

## Code Conventions

**Cobra commands:** Always use `RunE` (not `Run`) — errors must propagate. Persistent flags on root, local flags on each command.

**Error handling:** Wrap with context using `fmt.Errorf("doing X: %w", err)`. Sentinel errors (`ErrManifestNotFound`) for conditions callers check; custom types for structured errors.

**Receivers:** Pointer receivers (`*T`) when method mutates or struct is large. Value receivers (`T`) for small read-only methods. Be consistent per type.

**Package structure per package:**
1. `<pkg>.go` — struct definitions + constructors
2. Logic files (`parser.go`, `validator.go`, etc.)
3. `<file>_test.go` — table-driven tests next to the code they test

**Naming:**
- Packages: lowercase single word (`manifest`, `provider`)
- Files: lowercase with underscores (`claude_desktop.go`)
- Exported: PascalCase (`LoadManifest`), unexported: camelCase (`expandVars`)
- Error messages: lowercase, no trailing punctuation, include context

**Tests:** Table-driven pattern always. Use `testdata/` for fixture files. Use `t.Helper()` in test helpers.

## Key Interfaces

The `Provider` interface (`internal/provider/`) is the core abstraction — every AI tool implements `Name()`, `Detect()`, `Generate(*Manifest)`, `Read(path)`. New providers must satisfy this interface.

## Library Docs

**Always use Context7 for up-to-date library documentation.** Key library IDs:
- Go spec: `/websites/go_dev_ref_spec`
- Cobra: resolve via `github.com/spf13/cobra`
