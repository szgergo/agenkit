# Step 04 — Manifest Resolve

> **Depends on:** Steps 01–03
> **Can be worked on alongside:** Step 05 once merging and resolution rules are clear
>
> **Pipeline position:** Resolution is stage 3 in the four-stage manifest pipeline: **Load → Merge → Resolve → Generate/Write**. `LoadGlobalManifest` and `LoadProjectManifest` return raw, unresolved `*Manifest`. The Merge step (step 03, `manifestMerger.go`) combines global + project manifests first. Only then does every command that needs resolved values call the resolver explicitly. This ordering means:
> - `vars:` values from the project manifest are available when resolving global entries.
> - Raw values stay available for `config view --raw` (which uses Merge output, bypassing Resolve).
> - Resolution failures are visible at command time rather than silently at load time.

## Outcome

You have a clear strategy for resolving manifest placeholders such as environment-driven values and runtime-derived values without turning the template system into a mini language.

## Why this step exists

Template resolution touches many user-visible behaviors: portability, secrets, provider paths, and generated config stability. It needs to be predictable long before the tool becomes feature-rich.

## Decisions you should make

- Which placeholder forms are truly part of the supported contract?
- At what point in the pipeline should resolution happen?
- What should happen when a referenced value is missing or malformed?
- How much of the manifest should stay raw for inspection versus resolved for execution?
- How will you keep template behavior deterministic and easy to test?

## Suggested work order

1. List the placeholder kinds required by the design.
2. Decide what inputs the resolver needs and what outputs it returns.
3. Separate resolution rules from manifest parsing so both remain understandable.
4. Cover edge cases such as missing environment values and invalid placeholder syntax.
5. Re-check how this step will affect `config view`, `diff`, and `apply`.

## Go learning focus

- string processing and validation boundaries
- error-returning helper functions and where to wrap context
- keeping logic pure where possible so testing stays cheap
- using standard library facilities before inventing abstractions

## Learn more

- Credible sources: [Effective Go](https://go.dev/doc/effective_go), [Go error handling guidance](https://go.dev/blog/error-handling-and-go), [Go os package](https://pkg.go.dev/os)
- Search keywords: "go environment variable lookup", "go string replacement patterns", "go error wrapping", "go validate template placeholders"

## Watch for

- silently falling back when a placeholder cannot be resolved
- coupling resolution logic too tightly to one command
- designing a more powerful template language than the product actually needs

## Definition of done

- [ ] Supported placeholder forms are explicit.
- [ ] Resolution behavior is deterministic and testable.
- [ ] Failure behavior is clear rather than silent.
- [ ] The design still leaves room for raw versus resolved views where needed.
