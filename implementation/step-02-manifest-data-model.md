# Step 02 — Manifest Data Model and Parsing

> **Depends on:** Step 01
> **Can be worked on alongside:** Step 03 once the basic project scaffold exists

## Outcome

You have a clear in-memory representation of **agenkit's own manifest file** and a parsing boundary that turns YAML into that representation without yet blurring parsing, validation, expansion, and merge behavior together.

Specifically, this step is about the data model and parser for `agenkit.yml`:

- the global manifest
- the project manifest
- the shared shape both of those files use

It is **not** about provider-specific config files such as JSON or TOML files for Copilot, Claude, Cursor, or Codex.

## Why this step exists

Almost every later feature depends on this model. If agenkit's own manifest shape is confused, every command and provider feature will inherit that confusion.

## Decisions you should make

- Which concepts in `agenkit.yml` deserve first-class types and which should stay simple fields?
- Where should you use keyed collections versus ordered collections?
- What belongs in parsing, and what belongs in later validation?
- How will you represent the fact that both the global and project `agenkit.yml` files share one schema, while still supporting later merge behavior?
- How will you represent a project-only field like `mode:` so global manifests remain valid and meaningfully different?
- Which zero values are safe, and which ones would create ambiguity?

## Example Manifest

This is what you're building a data model for. The parser should turn this YAML into Go structs. This example covers Phase 1 only:

```yaml
version: 1

# Path variables for cross-platform/cross-machine portability
vars:
  mcp_dir: "{home}/mcp-servers"
  work_dir: "{home}/work"

# MCP Server configurations
mcp_servers:
  filesystem:
    command: node
    args:
      - "{vars.mcp_dir}/filesystem/dist/index.js"
      - "{vars.work_dir}"
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

# Provider configuration paths (optional overrides)
providers:
  claude-desktop:
    config_path: "/custom/path/to/claude_desktop_config.json"
  cursor:
    config_path: "{home}/.config/cursor/mcp.json"
  codex:
    enabled: false
```

## Suggested work order

1. Re-read the manifest behavior in `DESIGN.md` and list the concepts you must preserve.
2. Design the model around agenkit's own manifest semantics, not around one sample file.
3. Decide what successful parsing guarantees and what it intentionally leaves for validation.
4. Add fixtures that cover both ordinary and malformed inputs.
5. Check that the model still leaves room for later phases such as hooks and instructions.

## Go learning focus

- structs as stable data shapes
- struct tags as a serialized-data boundary
- maps, slices, and zero values in configuration modeling
- how package-private helpers can keep parsing code readable without exporting everything

## Learn more

- Credible sources: [go.yaml.in/yaml/v4 docs](https://pkg.go.dev/go.yaml.in/yaml/v4), [Effective Go](https://go.dev/doc/effective_go), [Go language tour](https://go.dev/tour/)
- Search keywords: "yaml v4 struct tags go", "yaml unmarshal go struct", "go zero values config", "go map vs slice"

## Watch for

- mixing parsing with validation or template expansion too early
- modeling only today’s examples instead of the design constraints
- choosing types whose zero values make merge semantics confusing later

## Definition of done

- [X] The `agenkit.yml` model captures the Phase 1 behavior cleanly.
- [X] Parsing has a clear success boundary and clear failure modes.
- [X] Fixtures cover both expected and malformed inputs.
- [X] You can explain why each major manifest collection is a map, slice, or scalar.
