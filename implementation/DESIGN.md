# Agenkit — Design Document

## Table of Contents

- [Problem Statement](#problem-statement)
- [Core Insight](#core-insight)
- [Design Principles](#design-principles)
- [Manifest](#manifest)
  - [Manifest Pipeline](#manifest-pipeline)
  - [Manifest Storage & Git Sync](#manifest-storage--git-sync)
    - [Storage Locations & Config Layering](#storage-locations-config-layering)
    - [Git Sync Setup](#git-sync-setup)
    - [Auto-commit on Change](#auto-commit-on-change)
    - [CLI Commands](#cli-commands)
    - [New Machine Bootstrap](#new-machine-bootstrap)
    - [What Goes in the Git Repo](#what-goes-in-the-git-repo)
    - [Lock Files — Tracking Managed Entries](#lock-files-tracking-managed-entries)
    - [Security Considerations](#security-considerations)
    - [Conflict Resolution](#conflict-resolution)
    - [Phase 1: MCP Server Management](#phase-1-mcp-server-management)
    - [Config Creation UX (Zero-to-Manifest)](#config-creation-ux-zero-to-manifest)
    - [Phase 2: Agent Lifecycle Hooks](#phase-2-agent-lifecycle-hooks)
    - [Phase 3: Custom Instructions & Rules](#phase-3-custom-instructions-rules)
- [CLI Design — kubectl-Style](#cli-design-kubectl-style)
  - [Concept Mapping](#concept-mapping)
  - [Resource-Oriented Commands](#resource-oriented-commands)
  - [Output Formats (`-o` flag)](#output-formats-o-flag)
  - [Global Flags (consistent across all commands)](#global-flags-consistent-across-all-commands)
  - [CLI Help Standards](#cli-help-standards)
  - [Shell Completion](#shell-completion)
- [Technology Choice](#technology-choice)
- [File Structure (Go Project)](#file-structure-go-project)
- [Resolved Design Decisions](#resolved-design-decisions)
- [Evolution of This Idea](#evolution-of-this-idea)


> A CLI tool that manages AI agent configurations from a single source of truth.
> Define your MCP servers, hooks, and instructions once — apply them to every AI provider.

---

## Problem Statement

Developers using multiple AI coding agents (GitHub Copilot, Claude Code, Cursor, Codex, etc.) face a fragmented configuration landscape. Each tool stores MCP server configs, lifecycle hooks, and custom instructions in different locations, different formats, and with different semantics. Setting up a new machine or keeping configs in sync across tools is manual, error-prone, and tedious.

## Core Insight

These configurations share the same *intent* — "use this MCP server," "run this script on session start," "follow these coding rules" — but differ in *format*. A single manifest can be the source of truth, with provider-specific configs generated automatically.

## Design Principles

1. **No proprietary cloud** — sync is handled by git (GitHub, GitLab, Gitea, anything). You own the repo and the data.
2. **Declarative manifest** — one YAML file describes your desired state.
3. **Non-destructive** — generates provider configs alongside existing ones; never deletes what it didn't create.
4. **Graceful degradation** — if a provider doesn't support a feature (e.g., hooks), skip it and warn.
5. **Minimal abstraction** — don't invent a universal config language. Map concepts 1:1 where possible, note gaps clearly.
6. **Git-native** — the manifest is version-controlled; changes are committed and optionally auto-pushed.

---

## Manifest

### Manifest Pipeline

Every command that needs usable config runs the manifest through a fixed four-stage pipeline. Stages are always executed in this order; no stage may be skipped (except where noted):

```
LoadGlobalManifest()  ──┐
                        ├──→  Merge()  ──→  Resolve()  ──→  Generate() / Write() / Display()
LoadProjectManifest() ──┘
```

| Stage | Responsibility | Code location |
|-------|---------------|---------------|
| **Load** | Parse YAML into raw `*Manifest`; reject unknown fields and invalid enums | `internal/manifest/manifestLoader.go` |
| **Merge** | Combine global + project manifests according to project's `mode:` field (`merge` / `replace` / `extend`) | `internal/manifest/manifestMerger.go` |
| **Resolve** | Expand `{home}`, `{env.X}`, `{vars.X}` placeholders; error on unknown or missing references | `internal/manifest/manifestResolver.go` |
| **Generate / Write / Display** | Provider-specific config generation, safe file writes, lock file updates | `internal/provider/`, `internal/lock/`, `internal/cli/` |

**Key invariants:**
- Loaders return raw `*Manifest` — never merged or resolved.
- Merge runs before Resolve; `vars:` values may themselves contain placeholders that resolve in the next stage.
- Resolve runs before any provider sees the manifest — generators never handle raw `{…}` strings.
- Commands that expose unresolved state (e.g. `config view --raw`) use the output of Merge, bypassing Resolve.
- If no project manifest is found, Merge is a no-op (returns the global manifest unchanged).

---

### Manifest Storage & Git Sync

#### Storage Locations & Config Layering

Agenkit uses a **two-layer config system**, applied in order (later layers override earlier):

| Layer | Location | Scope | Who controls it |
|-------|----------|-------|----------------|
| **Global** | `~/.agenkit/agenkit.yml` | Machine-wide defaults | You — your personal AI tooling baseline |
| **Project** | `./.agenkit/agenkit.yml` (project root) | Project-specific | Per-repo, committed to git |

The `~/.agenkit/` dotdir is itself a git repo — synced across machines as described below.

**Merge semantics:** The project manifest declares how it relates to the global config via a top-level `mode:` field. The global manifest has no merge target, so `mode:` is only meaningful in a project manifest.

If a project manifest omits `mode:`, agenkit warns and defaults to `merge`; the user can suppress the warning by adding `mode: merge` explicitly.

| Mode | Behavior | When to use |
|------|----------|-------------|
| `merge` **(default)** | Global + project combined; project wins on key conflicts | Standard daily use |
| `replace` | Project replaces global entirely; no inheritance | Isolated project needing full control |
| `extend` | Project can only *add* new entries; cannot override or remove global keys | Org-managed configs — platform team MCPs stay always-on |

```yaml
# ./.agenkit/agenkit.yml (project manifest)
version: 1

# Default is "merge" — omit to get standard behavior
mode: merge

mcp_servers:
  # Added on top of global servers (merge mode)
  project-api:
    command: npx
    args: ["-y", "@myorg/project-api-mcp"]
```

```yaml
# Isolated project — ignores global config entirely
mode: replace

mcp_servers:
  local-only:
    command: node
    args: ["./tools/mcp-server.js"]
```

```yaml
# Org-enforced baseline — project can only add, never remove or override
mode: extend

mcp_servers:
  extra-tool:
    command: uvx
    args: ["extra-mcp-tool"]
```

> **Note on `extend`:** If a project manifest in `extend` mode tries to override a key present in the global config, agenkit emits a warning and ignores the override. Use `agenkit validate` to check for these conflicts before applying.

```
~/.agenkit/
├── agenkit.yml             # Global manifest
├── hooks/                  # Hook scripts
│   ├── on-session-start.sh
│   └── audit-tool.sh
├── instructions.md         # Global instructions source (Phase 3)
└── cache/                  # Registry response cache
```

> **Future consideration: Profiles.** A multi-profile system (named configs like `work`, `oss`, `data-science` — switchable via `agenkit config use-profile`) is a natural future extension. For now, the global + project model covers the primary use cases without added complexity.

---

#### Git Sync Setup

```
$ agenkit sync init
  ? Where should your global manifest be stored?
  ❯ New private GitHub repo (recommended)
    New private GitLab repo
    Existing git repo URL
    Local only (no remote sync)

  Creating private repo 'agenkit-config' on GitHub...
  ✓ Initialized ~/.agenkit as a git repo
  ✓ Added remote: git@github.com:gergo/agenkit-config.git
  ✓ Initial commit and push

  On a new machine, run:
    agenkit sync clone git@github.com:gergo/agenkit-config.git
```

#### Auto-commit on Change

When git sync is configured, any operation that modifies the global manifest **automatically commits the change**:

```
$ agenkit add mcp postgres
  ✓ Added 'postgres' MCP server to agenkit.yml
  ✓ Committed: "agenkit: add postgres MCP server"
  ✓ Pushed to origin/main
```

If push fails (e.g., offline), agenkit commits locally and warns:

```
  ✓ Committed: "agenkit: add postgres MCP server"
  ⚠ Push failed (offline?) — will retry on next `agenkit sync push`
```

#### CLI Commands

```
agenkit sync init              # Set up git repo for global manifest
agenkit sync clone <url>       # Clone manifest repo on a new machine (clone only)
agenkit sync clone <url> --apply  # Clone + immediately apply to detected providers
agenkit sync push              # Manually push pending commits
agenkit sync pull              # Pull latest manifest from remote
agenkit sync pull --apply      # Pull + re-apply to providers
agenkit sync status            # Show local vs remote diff
```

#### New Machine Bootstrap

This is the primary cross-machine value proposition — **two explicit commands to fully configure a new machine**:

```
$ agenkit sync clone git@github.com:gergo/agenkit-config.git
  Cloning manifest repo...
  ✓ Cloned to ~/.agenkit/

  Run `agenkit apply` to apply your config to detected providers.
  Or re-run with --apply to do both in one step.

$ agenkit apply

  🔍 Detecting installed AI providers...
    ✓ Claude Code
    ✓ GitHub Copilot (VS Code)
    ✗ Cursor (not installed — skipping)

  Applying 7 MCP servers to 2 providers...
    ✓ filesystem   → Claude Code, Copilot
    ✓ github       → Claude Code, Copilot
    ✓ postgres     → Claude Code only
    ...

  ✅ Done. Your AI agent setup is fully configured.
```

> **Why separate commands?** Keeping clone and apply distinct lets you inspect the manifest before it makes any changes to your local provider configs. `--apply` on `sync clone` is a convenience shortcut when you trust the remote state.

#### What Goes in the Git Repo

```
~/.agenkit/              # Git root
├── agenkit.yml                 # The manifest (committed)
├── hooks/                      # Hook scripts (committed)
│   ├── on-session-start.sh
│   └── audit-tool.sh
├── instructions.md             # Global instructions source (Phase 3)
└── locks/                      # Lock files (committed — tracks agenkit ownership)
    ├── claude-desktop.lock.json
    ├── claude-code.lock.json
    ├── cursor.lock.json
    └── copilot.lock.json
```

#### Lock Files — Tracking Managed Entries

Agenkit must know which entries in a provider's config *it* created, so it can safely update or remove them without touching anything the user added manually. Provider JSON configs don't support comments, so injecting metadata keys (e.g. `_agenkit_managed: true`) into provider configs would risk breaking providers that validate their schema strictly.

Instead, agenkit maintains a **parallel lock file per provider** in `~/.agenkit/locks/`:

```json
// ~/.agenkit/locks/claude-desktop.lock.json
{
  "managed_servers": ["filesystem", "github", "postgres"],
  "last_applied": "2025-04-08T10:30:00Z",
  "manifest_sha": "abc123"
}
```

- `apply` writes to the provider config AND updates the lock file in one operation.
- On next `apply`, agenkit removes only the servers listed in the lock (not user-added ones).
- Lock files are committed to git — every `apply` produces a new commit.
- **Rollback:** `git revert <sha>` + `agenkit apply` to restore any previous state.

#### Security Considerations

**Never commit secrets.** The manifest uses `{env.VAR_NAME}` references for anything sensitive (API keys, tokens). The env var values live in your shell profile (`~/.zshrc`, `~/.bashrc`) or a secrets manager — never in the manifest.

Agenkit enforces this:

```
$ agenkit validate
  ⚠ agenkit.yml line 14: 'GITHUB_TOKEN' looks like a literal secret.
    Use {env.GITHUB_TOKEN} instead and set it in your shell profile.
```

**Private repo recommended.** Even without literal secrets, your manifest reveals which tools and services you use. `agenkit sync init` creates a **private** repo by default.

#### Conflict Resolution

If the manifest is edited on two machines before syncing:

```
$ agenkit sync pull
  ⚠ Merge conflict in agenkit.yml (diverged from remote)

  Options:
  ❯ Open in editor to resolve manually
    Take remote version (discard local changes)
    Keep local version (overwrite remote on next push)
```

For the common case (adding a server on machine A while machine B adds a different server), git's automatic merge handles it correctly — MCP server entries are independent YAML keys.

---



#### Phase 1: MCP Server Management

**Goal:** Define MCP servers once, generate correct configs for all installed AI providers.

**Manifest format** (`agenkit.yml`):

```yaml
version: 1

# (Phase 2) Split large configs across files — resolved and merged before parsing
# includes:
#   - ./work-mcps.yml
#   - ./personal-mcps.yml

# Path variables for cross-platform/cross-machine portability
vars:
  mcp_dir: "{home}/mcp-servers"

mcp_servers:
  filesystem:
    command: node
    args:
      - "{vars.mcp_dir}/filesystem/dist/index.js"
      - "{home}/Projects"
    env:
      NODE_ENV: production

  github:
    command: npx
    args:
      - "-y"
      - "@modelcontextprotocol/server-github"
    env:
      GITHUB_TOKEN: "{env.GITHUB_TOKEN}"

  postgres:
    command: uvx
    args:
      - "mcp-server-postgres"
      - "{env.DATABASE_URL}"
    # Only generate for specific providers
    providers: [claude, cursor]

  context7:
    command: npx
    args: ["-y", "@upstash/context7-mcp"]
    # Provider-specific arg overrides — same server, different flags per tool
    provider_overrides:
      cursor:
        args: ["-y", "@upstash/context7-mcp", "--profile", "cursor-profile"]
      claude:
        args: ["-y", "@upstash/context7-mcp", "--profile", "claude-profile"]
```

**Provider config generation targets:**

| Provider | Config Location | Format |
|----------|----------------|--------|
| Claude Desktop | `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) | JSON |
| Claude Code | `.claude/settings.json` or `~/.claude/settings.json` | JSON |
| Cursor | `.cursor/mcp.json` (project) or global settings | JSON |
| GitHub Copilot (VS Code) | `.vscode/mcp.json` or VS Code `settings.json` | JSON |
| GitHub Copilot CLI | `~/.copilot/mcp-config.json` | JSON |
| Codex (OpenAI) | `~/.codex/config.toml` (global) or `.codex/config.toml` (project) | TOML |

**Provider auto-detection — three-tier lookup:**

Agenkit detects which providers are installed using a priority-ordered lookup. Each tier can override the previous:

1. **`providers.<name>.config_path` in `agenkit.yml`** — explicit user override; highest priority
2. **OS-specific known default path** — compiled into agenkit per-provider; used if no override set
3. **ENV var documented by the provider** — last fallback for non-standard installs

> Why ENV vars are last: GUI applications (e.g. Claude Desktop) are launched outside the terminal and don't inherit your shell's environment. ENV vars set in `.zshrc` are invisible to agenkit when running as a tray app or from a launcher. The known default path is more reliable.

**Known ENV vars (where providers document them):**

| Provider | ENV var | Notes |
|----------|---------|-------|
| Claude Code | `CLAUDE_CONFIG_DIR` | Overrides `~/.claude/` |
| Codex | `CODEX_HOME` | Overrides `~/.codex/` |
| GitHub Copilot | *(none documented)* | Relies on VS Code extension path detection |
| Claude Desktop | *(none documented)* | Uses OS-specific default path only |
| Cursor | *(none documented)* | Uses OS-specific default path only |

**Overriding provider path in `agenkit.yml`:**

```yaml
# agenkit.yml
providers:
  claude-desktop:
    config_path: "/custom/path/to/claude_desktop_config.json"
  cursor:
    config_path: "{home}/.config/cursor/mcp.json"
  # Disable a provider entirely
  codex:
    enabled: false
```

**CLI commands** (see [CLI Design](#cli-design--kubectl-style) for the full kubectl-style command reference):

```
agenkit init                   # Interactive wizard — scan + build manifest
agenkit import                 # Auto-detect existing configs → generate manifest
agenkit search <query>           # Search across all configured registries
agenkit registry list            # Show configured registries + connectivity
agenkit registry add <name> <url> # Add a registry
agenkit registry remove <name>   # Remove a registry
agenkit registry test <name>     # Test connectivity + auth
agenkit add mcp <name>         # Add an MCP server from registry or manually
agenkit apply                  # Generate all provider configs from manifest
agenkit apply --provider claude # Generate for specific provider only
agenkit apply --dry-run        # Show what would be generated without writing
agenkit status                 # Show which providers are configured, which are stale
agenkit validate               # Validate manifest syntax and check for issues
```

---

#### Config Creation UX (Zero-to-Manifest)

**Problem:** If the first step is "open a text editor and write YAML," most users won't bother. The tool must meet users where they are — with configs they already have.

#### Strategy 1: `agenkit import` — Reverse-Engineer From Existing Configs

This is the **primary onboarding path.** Most users already have MCP servers configured in at least one provider. Import reads those configs and generates `agenkit.yml` automatically.

**Flow:**

```
$ agenkit import

🔍 Scanning for AI agent configurations...

Found:
  ✓ Claude Code     ~/.claude/settings.json          3 MCP servers
  ✓ Cursor          ~/.cursor/mcp.json               5 MCP servers
  ✓ GitHub Copilot  ~/.vscode/settings.json           2 MCP servers
  ✗ Claude Desktop  (not found)
  ✗ Codex           (not found)

Merging 7 unique MCP servers (3 shared across providers)...

? Review the generated manifest? [Y/n]

  # Shows the YAML with all discovered servers, paths already templated

? Where to save? [agenkit.yml]

✅ Created agenkit.yml with 7 MCP servers
   Run `agenkit apply --dry-run` to preview what would be generated.
```

**Key behaviors:**

- **Deduplication:** If the same MCP server (same command + args pattern) appears in multiple providers, it's merged into one entry.
- **Path abstraction:** Absolute paths like `/Users/gergo/...` are automatically converted to `{home}/...` for portability.
- **Env var detection:** If an arg looks like a secret/token, suggest using `{env.VAR_NAME}` instead of the literal value.
- **Provider tagging:** If a server only exists in one provider's config, add a `providers: [cursor]` tag so `apply` doesn't push it everywhere.

**Conflict resolution during import:**

When the same server name appears across multiple providers with different configurations, agenkit handles three cases:

| Case | Handling |
|------|----------|
| Same name, identical config | Silent dedup — one entry, no prompt |
| Same name, different args/env | Ask: keep provider A's version / keep B's / keep both (prompt for a name suffix) |
| Same name, provider-specific intentional differences | Merge as one entry with `provider_overrides:` block — not treated as a conflict |

```
⚠ Conflict: 'context7' found in Cursor and Claude Code with different args.

  Cursor version:  npx -y @upstash/context7-mcp --profile cursor-profile
  Claude version:  npx -y @upstash/context7-mcp --profile claude-profile

  Options:
  ❯ Merge as one entry with provider-specific overrides (recommended)
    Use Cursor's version for all providers
    Use Claude's version for all providers
    Keep both (name the second one: context7-claude)
```

#### Strategy 2: `agenkit init` — Interactive Wizard

For new setups or users who want a guided experience.

**Flow:**

```
$ agenkit init

? What scope is this manifest for?
  ❯ Global (all projects — saved to ~/.agenkit/agenkit.yml)
    Project (this project only — saved to ./.agenkit/agenkit.yml)

? Would you like to import existing configs from installed AI tools?
  ❯ Yes, scan and import
    No, start from scratch

  🔍 Found 3 MCP servers in Claude Code, 2 in Copilot.
  ✓ Imported 4 unique servers.

? Add more MCP servers?
  ❯ Browse popular servers
    Add by command manually
    Done

  Popular MCP servers:
  ❯ ◻ filesystem — Local file system access
    ◻ github — GitHub API (issues, PRs, repos)
    ◻ postgres — PostgreSQL database
    ◻ sqlite — SQLite database
    ◻ brave-search — Web search via Brave
    ◻ puppeteer — Browser automation
    (↑↓ navigate, space to select, enter to confirm)

  ? GITHUB_TOKEN is required for the github server.
    How should it be configured?
    ❯ Reference env var: {env.GITHUB_TOKEN}
      Enter value now (will be stored as env reference)
      Skip for now

? Which providers should agenkit manage?
  ❯ ◉ Claude Code
    ◉ GitHub Copilot
    ◉ Cursor
    ◻ Claude Desktop
    ◻ Codex

✅ Created agenkit.yml
```

#### Strategy 3: `agenkit add` — Incremental Addition

For adding servers one at a time after initial setup.

```
# Add from a curated registry of popular servers
$ agenkit add mcp github
  Added 'github' MCP server (npx @modelcontextprotocol/server-github)
  Requires: GITHUB_TOKEN env var

# Add by specifying command directly
$ agenkit add mcp custom-server --command "node" --args "./my-server/index.js"
  Added 'custom-server' MCP server

# Interactive mode
$ agenkit add mcp
  ? Server name: my-api
  ? Command: npx
  ? Arguments: -y @my-org/mcp-api-server
  ? Environment variables (key=value, empty to skip): API_KEY={env.API_KEY}
  ? Limit to specific providers? (empty = all): claude,copilot
  Added 'my-api' MCP server
```

#### MCP Server Registry Integration

To power `agenkit init`, `agenkit add`, and `agenkit search`, agenkit queries one or more **MCP registries**. The official MCP Registry is the default, but users can configure additional or alternative registries — including private enterprise registries.

**Default registry:** `https://registry.modelcontextprotocol.io`

All registries must implement the standard MCP Registry API spec:

```
GET /v0/servers                          # List/search all servers (paginated)
GET /v0/servers/{name}/versions          # List versions of a server
GET /v0/servers/{name}/versions/latest   # Get install config for latest version
```

**Configuring registries in `agenkit.yml`:**

```yaml
registries:
  # Official MCP Registry — always included by default, can be disabled
  official:
    url: https://registry.modelcontextprotocol.io
    enabled: true

  # Private enterprise registry (e.g. Azure API Center, self-hosted)
  acme-corp:
    url: https://mcp-registry.acme.corp/v0
    auth:
      type: bearer
      token: "{env.ACME_REGISTRY_TOKEN}"
    # Search this registry first
    priority: high

  # Another community registry
  community:
    url: https://mcp.community-hub.dev
    enabled: true
```

**Auth options for private registries:**

| Type | Config |
|------|--------|
| None (public) | `auth: none` or omit |
| Bearer token | `auth: { type: bearer, token: "{env.TOKEN}" }` |
| Basic auth | `auth: { type: basic, username: "x", password: "{env.PASS}" }` |
| Header | `auth: { type: header, name: "X-API-Key", value: "{env.KEY}" }` |

**CLI registry management:**

```
agenkit registry list                    # Show configured registries + status
agenkit registry add <name> <url>        # Add a registry interactively
agenkit registry remove <name>           # Remove a registry
agenkit registry test <name>             # Check connectivity + auth
agenkit registry update                  # Refresh cache from all enabled registries
```

**Multi-registry search behavior:**

When searching across multiple registries, results are merged and labelled by source:

```
$ agenkit search postgres

  [official]    io.github.modelcontextprotocol/server-postgres  — PostgreSQL via MCP
  [official]    io.github.supabase/mcp-supabase                 — Supabase (Postgres + storage)
  [acme-corp]   acme/internal-db-server                         — ACME internal database server
```

- Results from higher-priority registries appear first.
- If the same server exists in multiple registries, both are shown with their source label.
- `agenkit add mcp acme/internal-db-server --registry acme-corp` to be explicit.

**Offline fallback — three-tier chain:**

1. **Cache-first:** Serve from `~/.agenkit/cache/<registry-name>.json` (24h TTL). This is the primary offline story for established installs.
2. **Live fetch:** On cache miss or expired cache, query the registry.
3. **Seed list (bootstrap only):** A minimal list of ~50 well-known servers is bundled via Go's `embed.FS`. Used only on first run when no cache exists yet. Results labelled `[seed]` in output to make clear they may be stale.

> No large embedded registry ships with the binary. The seed list is a last-resort bootstrap aid for first-run offline scenarios, not a maintained catalog.

**Caching:** Registry responses are cached per-registry (`~/.agenkit/cache/<registry-name>.json`) with a configurable TTL (default: 24h). Override with `--no-cache` or `registries.<name>.cache_ttl: 0`.

**Design rationale — why not hardcode-only:**

| Approach | Pros | Cons |
|----------|------|------|
| Hardcoded only | Works offline, no API dependency | Stale on release day, no private registries |
| Single official registry | Always up-to-date | Can't use private/enterprise registries |
| **Configurable registries (chosen)** | Personal, enterprise, community all supported | Slightly more complex |

#### Config File Watcher (Future Enhancement)

A possible future feature: `agenkit watch` monitors provider config files for changes and offers to import new servers back into the manifest. This closes the loop — if you add a server via Cursor's UI, agenkit detects it and asks if you want to sync it.

**Built-in variables for path templating:**

| Variable | Example Value |
|----------|---------------|
| `{home}` | `/Users/gergo` or `C:\Users\gergo` |
| `{os}` | `darwin`, `linux`, `windows` |
| `{arch}` | `arm64`, `x64` |
| `{hostname}` | `gergo-macbook` |
| `{env.VAR_NAME}` | Value of environment variable |
| `{vars.name}` | User-defined variable from manifest |

**Key design decisions:**

- **Merge strategy:** By default, agenkit-managed servers are *added* to existing provider configs. Existing manually-added servers are preserved. Agenkit marks its entries with a comment/metadata field so it can update them on re-apply.
- **Install support:** Phase 1 does NOT auto-install MCP servers (no `npm install`). It only generates configs pointing to commands. The user is responsible for having `npx`, `uvx`, etc. available. Auto-install is a future enhancement.
- **Scope:** Supports both global configs and project-level configs. The manifest location determines scope — `.agenkit/agenkit.yml` in a project root generates project-level configs; `~/.agenkit/agenkit.yml` generates global configs.

---

#### Phase 2: Agent Lifecycle Hooks

**Goal:** Define lifecycle hooks (bash scripts triggered on agent events) once, distribute to providers that support them.

**Event model:**

Hooks are bash commands triggered on agent lifecycle events. Different providers use different event names for similar concepts.

**Manifest extension:**

```yaml
hooks:
  session_start:
    command: "echo 'Session started at $(date)' >> ~/.agenkit/session.log"
    # Or reference a script
    # command: "{home}/.agenkit/hooks/on-session-start.sh"

  before_tool_use:
    command: "{home}/.agenkit/hooks/audit-tool.sh"
    # Optionally restrict to specific tools
    tools: ["bash", "write_file"]

  after_tool_use:
    command: "{home}/.agenkit/hooks/post-tool-check.sh"

  session_end:
    command: "{home}/.agenkit/hooks/cleanup.sh"
```

**Event mapping across providers:**

| Agenkit Event | Copilot Hook | Claude Code Hook |
|---------------|-------------|-----------------|
| `session_start` | `copilot-agent:session-start` | `SessionStart` |
| `session_end` | `copilot-agent:session-end` | `SessionEnd` |
| `before_tool_use` | *(not supported — skip)* | `PreToolUse` |
| `after_tool_use` | *(not supported — skip)* | `PostToolUse` |
| `user_prompt` | *(TBD)* | *(TBD)* |

> **Note:** This mapping will evolve as providers add/change hook support. The tool should warn when a hook is defined but not supported by a target provider.

**Key design decisions:**

- Scripts are referenced by path or inline command — agenkit does NOT manage/copy script files, just the config pointing to them.
- Windows support: hooks are bash commands. On Windows, they require WSL, Git Bash, or similar. Agenkit should detect and warn.
- Security: hooks are user-authored, local-only. No community/shared hooks in scope.

---

#### Phase 3: Custom Instructions & Rules

**Goal:** Maintain coding instructions/rules in one place, generate provider-specific instruction files.

**Manifest extension:**

```yaml
instructions:
  global:
    # A single source file for global coding instructions
    source: "{home}/.agenkit/instructions.md"

  project:
    # Project-level instructions (when agenkit.yml is in a project)
    source: "./coding-rules.md"

    # Provider-specific overrides (optional)
    cursor:
      # Cursor supports structured rules with globs
      rules:
        - name: "API conventions"
          glob: "src/api/**"
          source: "./docs/api-rules.md"
        - name: "Test patterns"
          glob: "**/*.test.*"
          source: "./docs/test-rules.md"
```

**Generation targets:**

| Provider | Output | Format |
|----------|--------|--------|
| GitHub Copilot | `.github/copilot-instructions.md` | Markdown |
| Claude Code | `CLAUDE.md` (project), `~/.claude/CLAUDE.md` (global) | Markdown |
| Cursor | `.cursor/rules/*.mdc` | MDC (Markdown + frontmatter) |
| Codex | `AGENTS.md` (project), `~/.codex/AGENTS.md` (global) | Markdown |

**Key design decisions:**

- **Lossy by design:** The tool acknowledges that Cursor's glob-scoped rules can't be expressed in Copilot's flat markdown. It generates the best approximation and documents the gaps.
- **Source format is plain markdown.** No invented format. Provider-specific features (like Cursor globs) are expressed in the manifest YAML, not in the markdown source.
- **Additive, not replacing:** Generated instruction files include a header comment indicating they're managed by agenkit. Manual content in those files is preserved in a separate section.

---

## CLI Design — kubectl-Style

Agenkit's CLI is intentionally modeled after `kubectl` patterns. If you know kubectl, you already know agenkit. This reduces cognitive load for the target audience (developers already fluent in Kubernetes-style tooling).

### Concept Mapping

| kubectl Concept | agenkit Equivalent |
|----------------|-------------------|
| resource type | `mcp`, `hook`, `instruction`, `registry`, `provider` |
| kubeconfig | `~/.agenkit/agenkit.yml` (global) |
| namespace (sort of) | project scope — determined by cwd having `./.agenkit/agenkit.yml` |
| `kubectl apply -f` | `agenkit apply -f` |
| `kubectl get pods` | `agenkit get mcps` |
| `kubectl describe pod X` | `agenkit describe mcp X` |
| `kubectl diff` | `agenkit diff` |
| `kubectl config view` | `agenkit config view` |

### Resource-Oriented Commands

Like kubectl, all data manipulation follows the verb-resource pattern:

```
agenkit <verb> <resource> [name] [flags]
```

#### Core Verbs

| Verb | kubectl Analogy | Description |
|------|----------------|-------------|
| `get` | `kubectl get` | List resources (table format by default) |
| `describe` | `kubectl describe` | Show detailed info about a single resource |
| `apply` | `kubectl apply` | Apply manifest to providers (reconcile desired → actual state) |
| `delete` | `kubectl delete` | Remove a resource from the manifest |
| `diff` | `kubectl diff` | Show what `apply` would change without changing it |
| `edit` | `kubectl edit` | Open resource in `$EDITOR` |

#### Full Command Reference

**Resources (get/describe/delete):**

```bash
# MCP Servers
agenkit get mcps                          # List all MCP servers (global + project merged)
agenkit get mcps -o yaml                  # Output as YAML
agenkit get mcps -o json                  # Output as JSON
agenkit get mcps -o wide                  # Wide table with extra columns
agenkit describe mcp github               # Detailed view: config, providers, status
agenkit delete mcp postgres               # Remove from manifest

# Hooks
agenkit get hooks                         # List all hooks
agenkit describe hook session-start       # Show hook details + provider mapping

# Registries
agenkit get registries                    # List configured registries + health
agenkit describe registry acme-corp       # Registry details, server count, last sync

# Providers
agenkit get providers                     # List detected AI providers + status
agenkit describe provider claude-code     # Provider config path, managed servers, etc.
```

**Apply & Diff:**

```bash
agenkit apply                             # Apply global + project manifest to all providers
agenkit apply -f custom.yml               # Apply a specific manifest file
agenkit apply --provider claude-code      # Apply to one provider only
agenkit apply --dry-run                   # Same as: agenkit diff
agenkit diff                              # Show what apply would change (colored diff)
```

**Config (mirrors `kubectl config`):**

```bash
agenkit config view                       # Show merged config (global + project, resolved)
agenkit config view --raw                 # Show unresolved (with {var} placeholders)
agenkit config view --global              # Show only global config
agenkit config view --project             # Show only project config
agenkit config view --minify              # Compact output
```

**Init & Import:**

```bash
agenkit init                              # Interactive wizard — scan + build manifest
agenkit import                            # Detect existing provider configs → generate manifest
```

**Search (registry interaction):**

```bash
agenkit search github                     # Search across all registries
agenkit search github --registry official # Search specific registry
```

**Add (shorthand for common mutations):**

```bash
agenkit add mcp github                    # Add MCP server from registry
agenkit add mcp custom --command "node" --args "./server.js"  # Add manually
agenkit add registry acme https://mcp.acme.corp/v0            # Add a registry
```

**Sync (git operations):**

```bash
agenkit sync init                            # Set up git repo for config dir
agenkit sync clone <url>                     # Clone config on new machine (clone only)
agenkit sync clone <url> --apply             # Clone + immediately apply to providers
agenkit sync push                            # Push pending commits
agenkit sync pull                            # Pull latest from remote
agenkit sync pull --apply                    # Pull + re-apply to providers
agenkit sync status                          # Local vs remote diff
```

**Validate:**

```bash
agenkit validate                          # Check manifest syntax, warn on secrets, etc.
```

### Output Formats (`-o` flag)

Consistent across all `get` commands, just like kubectl:

| Flag | Output |
|------|--------|
| *(default)* | Human-readable table |
| `-o wide` | Extended table with more columns |
| `-o yaml` | YAML (copy-pasteable into manifest) |
| `-o json` | JSON |
| `-o name` | Just resource names (for scripting) |

```bash
$ agenkit get mcps
NAME         COMMAND   PROVIDERS          STATUS
github       npx       claude,copilot     applied
postgres     uvx       claude             applied
filesystem   node      claude,copilot     stale

$ agenkit get mcps -o wide
NAME         COMMAND   ARGS                                    ENV              PROVIDERS          STATUS
github       npx       -y @modelcontextprotocol/server-github  GITHUB_TOKEN     claude,copilot     applied
...

$ agenkit get mcps -o yaml
mcp_servers:
  github:
    command: npx
    args: ["-y", "@modelcontextprotocol/server-github"]
    ...
```

### Global Flags (consistent across all commands)

```
--global                 # Target global config only (like --all-namespaces, sort of)
--project <path>         # Override project manifest location
--no-cache               # Skip registry cache
-v, --verbose            # Verbose output
-q, --quiet              # Minimal output
```

### CLI Help Standards

Every command in agenkit **must** be self-documenting via `--help`. The target standard: a developer new to agenkit should be able to learn the full workflow from `agenkit --help` + `agenkit help <command>` alone, without reading external docs.

**Required for every command (cobra fields):**

| Field | Purpose | Example |
|-------|---------|---------|
| `Short` | One-line summary shown in `agenkit --help` | `"Apply manifest to AI provider configs"` |
| `Long` | Full description shown in `agenkit apply --help` | Explains behaviour, defaults, scope resolution |
| `Example` | Runnable examples, shown in `--help` | See below |
| `SeeAlso` | Related commands | `["diff", "validate", "sync push"]` |

**Example: `agenkit apply --help`**

```
Apply the agenkit manifest to all detected AI provider configs.

Reads the current manifest (global + project layers merged) and writes provider-specific
config files for each installed AI tool. Agenkit tracks which entries it manages via lock
files so it never touches manually-added configuration.

By default, applies to all detected providers. Use --provider to target one.
Use --dry-run (or `agenkit diff`) to preview changes without writing.

Usage:
  agenkit apply [flags]

Examples:
  # Apply to all detected providers
  agenkit apply

  # Preview what would change without writing
  agenkit apply --dry-run

  # Apply only to Claude Code
  agenkit apply --provider claude-code

  # Apply a specific manifest file
  agenkit apply -f ./staging.yml

Flags:
  --provider string   Target a specific provider (e.g. claude-code, cursor, copilot)
  --dry-run           Show what would change without writing (same as: agenkit diff)
  -f, --file string   Apply a specific manifest file instead of the default

See also:
  agenkit diff        Preview changes before applying
  agenkit validate    Check manifest for errors before applying
  agenkit sync push   Push the resulting commit to the remote
```

### Shell Completion

Auto-generated by cobra, works like kubectl completions:

```bash
# Bash
source <(agenkit completion bash)

# Zsh
source <(agenkit completion zsh)

# Fish
agenkit completion fish | source
```

Completions include resource names (e.g., `agenkit describe mcp <TAB>` lists server names).

---

## Technology Choice

**Language: Go**

Go was chosen for:
- Single binary distribution — no runtime dependencies
- Excellent cross-platform support (cross-compile with `GOOS`/`GOARCH`)
- Strong CLI ecosystem (`cobra` for commands)
- Native JSON/YAML handling (`encoding/json`, `go.yaml.in/yaml/v4`)
- Fast enough for a config tool — startup time is negligible
- Familiar to the target audience (devtools developers)

**Key Go libraries:**

| Library | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI command structure + rich help text |
| `go.yaml.in/yaml/v4` | YAML parsing |

**Distribution:**

- `go install github.com/USER/agenkit@latest`
- Homebrew tap (`brew install agenkit`)
- GitHub Releases (prebuilt binaries for linux/darwin/windows × amd64/arm64)
- Optional: npm wrapper (`npx agenkit`) for discoverability

---

## File Structure (Go Project)

```
agenkit/
├── cmd/
│   └── agenkit/
│       └── main.go                 # Entry point
├── internal/
│   ├── cli/                        # Cobra command definitions
│   │   ├── root.go                 # Root command + global flags
│   │   ├── get.go                  # `agenkit get <resource>` commands
│   │   ├── describe.go             # `agenkit describe <resource>` commands
│   │   ├── apply.go                # `agenkit apply` generator
│   │   ├── diff.go                 # `agenkit diff` (apply --dry-run)
│   │   ├── delete.go               # `agenkit delete <resource>` commands
│   │   ├── init.go                 # `agenkit init` wizard
│   │   ├── import.go               # `agenkit import` scanner
│   │   ├── add.go                  # `agenkit add mcp|registry` commands
│   │   ├── search.go               # `agenkit search` registry query
│   │   ├── config.go               # `agenkit config view` commands
│   │   ├── sync.go                 # `agenkit sync` git operations
│   │   └── validate.go             # `agenkit validate` checker
│   ├── manifest/                   # Manifest parsing, loading & resolution
│   │   ├── manifest.go             # Struct definitions + Mode type
│   │   ├── manifestLoader.go       # YAML loading (returns raw, unresolved Manifest)
│   │   ├── manifestMerger.go       # Merge global + project manifests (respects mode:)
│   │   └── manifestResolver.go     # {home}, {env.X}, {vars.X} placeholder expansion
│   ├── provider/                   # Provider-specific generators
│   │   ├── provider.go             # Provider interface
│   │   ├── detect.go               # Auto-detect installed providers
│   │   ├── claude_desktop.go
│   │   ├── claude_code.go
│   │   ├── copilot.go
│   │   ├── cursor.go
│   │   └── codex.go
│   ├── importer/                   # Reverse-engineer existing configs
│   │   ├── importer.go             # Import orchestrator
│   │   ├── scanner.go              # Scan for installed provider configs
│   │   └── dedup.go                # Deduplicate servers across providers
│   ├── registry/                   # MCP Registry client (multi-registry)
│   │   ├── client.go               # Per-registry HTTP client with auth support
│   │   ├── manager.go              # Multi-registry fan-out, merging, priority
│   │   ├── cache.go                # Per-registry response caching (configurable TTL)
│   │   └── seed.go                 # Embedded seed list (~50 servers via embed.FS, bootstrap only)
│   ├── lock/                       # Lock file management
│   │   └── lock.go                 # Track agenkit-owned entries per provider config
│   ├── sync/                       # Git sync operations
│   │   ├── git.go                  # Git repo init, clone, commit, push, pull
│   │   └── github.go               # GitHub API client (create private repo)
├── testdata/
│   ├── manifests/                  # Sample agenkit.yml files
│   └── provider_configs/           # Expected provider config outputs
├── go.mod
├── go.sum
├── Makefile                        # Build, test, lint, release targets
├── goreleaser.yml                  # Cross-platform release automation
├── agenkit.yml                     # Dogfood: agenkit's own manifest
├── README.md
├── LICENSE
└── implementation/
    ├── DESIGN.md                       # This file
    ├── shared-context.md               # Go concepts + research guide for implementers
    └── step-NN-*.md                    # Per-step implementation guides
```

---

## Resolved Design Decisions

Previously tracked as "open questions" — now closed.

1. **Provider auto-detection:** Three-tier lookup: `agenkit.yml providers.X.config_path` → OS-default path → provider ENV var. ENV is last (GUI apps don't inherit shell env). See [Provider Detection](#provider-auto-detection--three-tier-lookup) section.

2. **Project vs global merge:** Controlled by `mode:` in project manifest. `merge` (default), `replace`, `extend`. See [Storage Locations & Config Layering](#storage-locations--config-layering).

3. **Manifest includes:** Use top-level `includes:` key (not YAML `!include` tag). Agenkit resolves includes before YAML parsing. Deferred to **Phase 2**.

4. **Registry offline fallback:** Cache-first (24h TTL) → live fetch → embedded seed (~50 servers, `embed.FS`, labelled `[seed]`). No large embedded registry. `agenkit registry update` refreshes cache manually.

5. **Import conflict handling:** Three cases — silent dedup / ask user / merge as `provider_overrides:`. See [agenkit import](#strategy-1-agenkit-import--reverse-engineer-from-existing-configs).

6. **Marking managed entries in provider configs:** Use a parallel lock file (`~/.agenkit/locks/<provider>.lock.json`) instead of injecting metadata keys into provider configs. Every `apply` produces a new git commit for easy rollback.

7. **git implementation:** Shell out to `git`. `go-git` has too many edge cases with SSH keys and credential helpers. `git` is a hard dependency — agenkit errors on `sync init`/`sync clone` if `git` is not found in `$PATH`.

8. **Auto-push behavior:** Auto-commit always on any manifest change. Auto-push is opt-in via `--push` flag or `sync.autopush: true` in global config. Push failures are non-fatal — agenkit warns and queues for `agenkit sync push`.

---

## Evolution of This Idea

This design emerged from iterative critique of the original concept:

1. **Original idea:** Tray app that syncs all AI agent data to cloud with E2E encryption.
   - *Rejected:* Agent data is mostly ephemeral/context-bound. Cloud sync + accounts adds massive complexity for little value. Existing tools (Syncthing, git) already solve generic file sync.

2. **Refinement:** Sync MCP configs and agent settings via cloud profiles.
   - *Rejected:* MCP configs aren't portable as-is (absolute paths, local deps). "Profile + login" = SaaS infrastructure costs. TAM is too small for a paid product.

3. **Pivot:** Declarative CLI tool that generates provider-specific configs from a single manifest.
   - *Accepted:* No cloud, no accounts, no servers. Unix philosophy. Solves the real pain (config fragmentation) with minimal complexity.

4. **Additions:** Hooks (Phase 2) and instructions (Phase 3) follow the same pattern — one source, many targets — and extend naturally from the MCP foundation.
