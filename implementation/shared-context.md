# Agenkit — Shared Implementation Context

> This document is the reference companion for the implementation steps.
> It explains what you should understand before or during the work, but it intentionally does not tell you exactly what code to write.

---

## 1. How to use this guide

Use the project docs in this order:

1. `DESIGN.md` for product behavior, command semantics, and architectural constraints.
2. This file for Go concepts, implementation judgment, and research direction.
3. The step files for delivery order, scope, and definition of done.

This file is intentionally code-free. Its job is to help you reason well, not to become something you copy from line by line.

When a step feels hard, do not immediately search for a pattern to mimic. First ask:

- What behavior is fixed by the design?
- What is still my implementation decision?
- What is the smallest design that solves the current step well?
- Which choice will make the next two steps easier instead of only making this one step pass?

---

## 2. Study map by step range

| Step range | Main concern | Research before starting | Go understanding that matters most |
|------------|--------------|--------------------------|------------------------------------|
| 01-04 | Build the project skeleton, model the manifest, and implement the Load/Merge/Resolve pipeline | Go modules, package layout, Cobra basics, YAML modeling | `main`, packages, exported vs unexported names, structs, struct tags, zero values, slices, maps, basic error handling |
| 05-09 | Define provider boundaries, generate output, write safely, and wire the first full flow | Go interfaces, JSON/TOML encoding, deterministic output, filesystem writing | interfaces, method sets, pointer vs value receivers, marshaling, file I/O, error wrapping |
| 10-14 | Make the tool inspectable and trustworthy | CLI output patterns, validation design, reconciliation workflows, fixture-driven testing | table-driven tests, comparison logic, formatting, parsing external config safely |
| 15-20 | Add networked, interactive, and operational features | HTTP clients, caching, prompt design, process execution, release workflow design | `net/http`, `context`, time, `os/exec`, environment and path discovery, portability |
| 21-22 | Extend the model without overcomplicating the core | capability mapping, extensible schema design, compatibility strategy | composition, package boundaries, backward-compatible config changes, deliberate abstraction |

If you want to prepare ahead, the most useful research order is:

1. Go package layout and module basics.
2. Structs, zero values, and struct tags.
3. Interfaces and method receivers.
4. Error handling and testing style.
5. File and path handling across operating systems.
6. HTTP, context, and external process execution.

For concrete sources and better search terms, use the topic guide in Section 10.

---

## 3. The Go mental model that matters here

Go will feel easier if you treat it as a language that prefers clarity over cleverness.

Important mindset shifts:

- Go is package-oriented, not class-oriented. You design responsibilities around packages and types, not inheritance trees.
- Go treats error handling as normal control flow. Most failures should be visible in function signatures, not hidden in exceptions.
- Small interfaces are usually better than central “everything” abstractions. In Go, interfaces are best when they describe a real boundary in behavior.
- Zero values matter more than in many other languages. A type should make sense even before every field is explicitly initialized.
- The standard library solves more than you may expect. Reaching for another dependency should be a deliberate choice.

For this project, that means:

- keep the command layer thin
- keep configuration modeling explicit
- prefer straightforward types over abstract frameworks
- add interfaces where multiple implementations genuinely need the same behavior
- let the data flow be obvious from manifest input to provider output

---

## 4. Go concepts worth understanding early

### Packages and visibility

Go uses package names and capitalization as part of its visibility model.

- Names starting with an uppercase letter are exported from the package.
- Names starting with a lowercase letter stay package-private.
- The `internal/` directory adds a stronger boundary: only code inside the parent tree can import those packages.

Why this matters for agenkit:

- you want clear internal boundaries between command wiring, manifest handling, provider behavior, and operational features
- you do not need to export types just because multiple files use them; package-level organization is often enough

### Structs and struct tags

Structs are the main way you model configuration data in Go. Tags describe how that data is serialized.

What to understand:

- a struct describes a stable shape
- a tag such as `yaml`, `json`, or `toml` defines how a field maps to external data
- `omitempty` changes output behavior, which can affect whether generated configs stay clean or become noisy

Why this matters for agenkit:

- manifest data is stable, structured, and central to the whole application
- provider configs are serialized boundaries, so tags are part of the contract

### Zero values

Every Go type has a default zero value. This is not just a language detail; it influences how good your config model feels.

Examples of decisions zero values affect:

- whether an empty slice means “none” or “not set”
- whether an empty string is a valid default or an invalid state
- whether a boolean needs a pointer because “unset” differs from `false`

For config-heavy tools, poor zero-value choices create confusing merge and validation behavior.

### Maps vs slices

Go gives you both, but they express different intent.

- A map is best when identity is keyed by name.
- A slice is best when order matters or duplicates are meaningful.

When designing agenkit data, ask:

- Is this concept identified by a stable key?
- Does ordering matter to users or to generated output?
- Do I need deterministic output even if the internal structure is unordered?

### Pointer vs value receivers

Methods can receive either a copy of a value or access to the original through a pointer.

What to understand:

- pointer receivers are usually right when methods mutate state or when copying the value would be wasteful
- value receivers are often fine for small, read-only types
- consistency matters more than clever micro-optimization

This is an important Go-specific design choice because it affects method sets, interface satisfaction, and how mutation flows through your code.

### Interfaces

Go interfaces describe behavior, not hierarchy.

Good interface questions:

- Is there a real boundary with more than one implementation?
- Does the interface describe a behavior the caller actually needs?
- Is the interface small enough to stay easy to satisfy and easy to understand?

For agenkit, provider behavior is a natural interface candidate. Many internal helpers are not.

### Error handling

Go uses explicit error returns. That means you should decide:

- where errors gain context
- where validation errors should be aggregated instead of returned one by one
- where warnings are user-facing output rather than actual errors

Idiomatic Go generally favors:

- returning errors instead of logging deep in the stack
- wrapping errors with context at boundaries
- keeping error messages lowercase and descriptive

### `io.Reader` and `io.Writer`

These two interfaces show up everywhere in Go because they are reusable boundaries across files, buffers, network calls, and tests.

Even if you do not use them everywhere, understanding them helps you design code that can:

- read from disk or memory with the same logic
- write output to stdout, a buffer, or a file
- become easier to test without inventing custom abstractions

### `context.Context`

`context.Context` matters once network requests or external processes enter the project.

You should understand:

- cancellation
- deadlines and timeouts
- how to pass context across call chains without storing it in structs

It is most relevant for registry access and possibly long-running operational commands.

### Paths and portability

Cross-platform tools live or die on path handling.

Use Go’s path and filesystem thinking, not string concatenation thinking:

- use path utilities that understand operating-system differences
- separate path discovery from path use
- think about home directories, environment variables, and provider-specific defaults as explicit behavior

---

## 5. How to think about the CLI layer

The CLI layer is not the application. It is the user-facing entry point into the application.

For Cobra specifically, understand these ideas:

- the command tree is your public surface area
- local flags belong to one command; persistent flags affect a whole subtree
- `RunE` matters because it lets errors propagate cleanly
- help text is part of the product, not an afterthought
- completion is useful only when it reflects real command semantics

Good CLI design habits for agenkit:

- keep parsing and user interaction close to the command layer
- move reusable domain logic out of command handlers
- avoid letting command files become the only place where business logic lives
- decide early which commands are interactive and which must stay scriptable

Because the project borrows heavily from `kubectl`-style ergonomics, consistency matters:

- resource-oriented verbs should feel predictable
- output formats should be stable
- destructive behavior should be explicit
- preview and apply flows should agree on meaning

---

## 6. How to think about serialized config data

agenkit sits between one manifest model and several provider-specific config formats. That means serialization is not an implementation detail; it is part of the core design.

Things to reason about:

- what belongs in the shared manifest model versus what belongs in provider-specific projection logic
- which fields are required semantically versus only required by one provider
- how omission of empty fields affects readability and diff stability
- how to make generated output deterministic so `diff` and `apply` feel trustworthy

Data boundaries to keep separate in your thinking:

- raw manifest input
- validated manifest data
- resolved or expanded manifest state
- provider-specific generated representation
- lock state describing what agenkit previously managed

If these boundaries blur together, later features become much harder.

---

## 7. Filesystem, process, and network concerns

This project interacts with the local machine a lot. That means operational behavior is as important as in-memory design.

### Filesystem

Things worth thinking through:

- how you will avoid clobbering user-managed config
- what “safe write” means for generated files
- when temporary files or atomic replacement are appropriate
- how lock files stay in sync with actual writes

### External processes

Git sync is intentionally designed around shelling out to `git`.

What you should understand before implementing that:

- command execution should surface stderr/stdout clearly enough to debug failures
- working directory and environment are part of the behavior
- user intent matters: clone, pull, push, and apply should not blur together

### Networking

Registry access introduces a different failure model:

- timeouts
- offline behavior
- partial success when multiple registries are queried
- cache validity versus freshness

Networked features are usually where `context`, timeout policy, and retry/fallback decisions begin to matter.

---

## 8. Testing strategy for this project

Go testing tends to work best when it is simple and local.

Useful patterns for agenkit:

- keep tests next to the package they exercise
- use table-driven tests for merge logic, parsing edge cases, and validation matrices
- use `testdata/` for fixture manifests and provider config examples
- test behavior and invariants, not just function existence

Areas where tests will matter disproportionately:

- manifest parsing and merge semantics
- template expansion
- provider generation output
- lock file behavior
- import deduplication and conflict handling
- validation rules
- command output for user-trust features such as `diff` and `config view`

When testing generated config, ask:

- am I asserting behavior or just asserting one exact formatting artifact?
- if exact output matters, is it because users will diff or inspect it directly?
- do I need fixture-based expectations to keep the tests readable?

---

## 9. Agenkit-specific constraints to keep in view

These are product constraints that should stay visible throughout implementation:

- The manifest is the source of truth; provider files are generated artifacts.
- The tool is non-destructive by design, so lock state is central rather than optional.
- The project uses two-layer config with `merge`, `replace`, and `extend` semantics.
- Provider detection follows an explicit priority order, rather than ad hoc discovery.
- Interactive onboarding should stay simple and plain-text instead of depending on a heavier interactive UI stack.
- `sync clone` and `apply` are separate actions; chaining should be explicit.
- Registry support is configurable and should support cache, live fetch, and fallback behavior.
- Hooks and instructions are later phases, but the Phase 1 design should leave room for them without overfitting to them too early.

Whenever you are unsure about an implementation decision, check whether it preserves these constraints.

---

## 10. Recommended reading and search keywords

Keep this list as the main source hub so the step files can stay lighter.

| Topic | Best use here | Credible source | Useful search keywords |
|-------|---------------|-----------------|------------------------|
| Go language tour | First contact with core syntax and language feel | https://go.dev/tour/ | `go tour basics`, `go package main`, `go short variable declaration` |
| Effective Go | Idioms and style decisions | https://go.dev/doc/effective_go | `effective go interfaces`, `effective go errors`, `effective go methods` |
| Go modules reference | Module behavior and dependency rules | https://go.dev/ref/mod | `go modules tutorial`, `go mod tidy`, `go module path imports` |
| Organizing a Go module | Project layout decisions | https://go.dev/doc/modules/layout | `go module layout`, `go cmd internal layout`, `go project structure official` |
| `internal/` package boundary | Package privacy and architecture boundaries | https://go.dev/doc/go1.4#internalpackages | `go internal package`, `go package visibility export`, `go internal directory rule` |
| Go `io` package | Reader/writer boundaries for reusable logic | https://pkg.go.dev/io | `go io reader writer`, `go reader writer interface`, `go bytes buffer io writer` |
| Go `testing` package | Test structure, subtests, helpers | https://pkg.go.dev/testing | `go testing package`, `go subtests`, `go t helper` |
| Table-driven tests | The dominant Go testing style | https://go.dev/wiki/TableDrivenTests | `go table driven tests`, `go test fixtures testdata`, `go test cases slice struct` |
| Go error handling guidance | Idiomatic error flow and wrapping | https://go.dev/blog/error-handling-and-go | `go error wrapping`, `fmt errorf percent w`, `errors is as go` |
| Go `context` package | Timeouts, cancellation, request lifecycles | https://pkg.go.dev/context | `go context timeout`, `context cancellation go`, `withtimeout go` |
| Go `net/http` package | Registry client and HTTP request design | https://pkg.go.dev/net/http | `go http client timeout`, `newrequestwithcontext go`, `go http client headers` |
| Go `path/filepath` package | Cross-platform file and path logic | https://pkg.go.dev/path/filepath | `go filepath join`, `filepath vs path go`, `go expand home path` |
| Go `os/exec` package | Shelling out to git and external tools | https://pkg.go.dev/os/exec | `go exec command`, `exec command context go`, `go combinedoutput stderr stdout` |
| Go `encoding/json` package | JSON provider config reading and writing | https://pkg.go.dev/encoding/json | `go json tags`, `marshalindent go`, `omitempty json go` |
| Cobra package docs | Core command tree, flags, and execution model | https://pkg.go.dev/github.com/spf13/cobra | `cobra RunE`, `cobra persistent flags`, `cobra subcommands go` |
| Cobra guides | Practical CLI patterns and completions | https://cobra.dev/ | `cobra shell completion`, `cobra help examples`, `cobra flags guide` |
| `go.yaml.in/yaml/v4` docs | Manifest modeling and YAML tags | https://pkg.go.dev/go.yaml.in/yaml/v4 | `yaml v4 struct tags go`, `omitempty yaml go`, `yaml unmarshal go struct` |
| BurntSushi TOML docs | TOML handling for provider config | https://pkg.go.dev/github.com/BurntSushi/toml | `burntsushi toml encode`, `go toml struct tags`, `toml decode file go` |

Search advice:

- Prefer official docs and package docs first, then maintainer-authored guides.
- When searching, include the package or project name, not just the concept. `cobra persistent flags` is better than `go cli flags`.
- For Go language behavior, adding `site:go.dev` or `pkg.go.dev` usually improves results.
- For libraries, search using the real import path or maintainer name when possible.
