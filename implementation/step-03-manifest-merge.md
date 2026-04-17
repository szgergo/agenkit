# Step 03 — Manifest Merge

> **Depends on:** Steps 01–02
>
> **Pipeline position:** Merge is stage 2 in the four-stage manifest pipeline: **Load → Merge → Resolve → Generate/Write**. `LoadGlobalManifest` and `LoadProjectManifest` (step 02) return two independent raw `*Manifest` values. Merge combines them into one according to the project manifest's `mode:` field. The resulting merged `*Manifest` is what every downstream stage (Resolve, Generate, Write) works with.

## Outcome

You have a function that takes a global and a project `*Manifest` and returns a single merged `*Manifest` that correctly reflects the `mode:` semantics (`merge`, `replace`, `extend`). Either manifest may be nil (no global config, or no project config).

## Why this step exists

The two-layer config system is a core product promise — global defaults with per-project overrides. Without an explicit Merge stage, each downstream consumer would have to implement its own layering logic, which leads to inconsistency and untestable combinations.

## Decisions you should make

- What does "merge" mean field-by-field? Maps, slices, and scalar fields each need a rule.
- What does `replace` mean — drop the global entirely, or only drop conflicting keys?
- What does `extend` mean — project can only add, never override or remove?
- What happens when either manifest is nil? (No global config found, or no project config found.)
- Where should the merged result live — should Merge return a new value or mutate one of its inputs?

## Suggested work order

1. Write down the expected output for each `mode:` value before writing any code — use your `testdata/` fixtures.
2. Implement `merge` mode first; it is the most common and exposes the most field-by-field decisions.
3. Add `replace` and `extend` once `merge` is proven.
4. Handle nil inputs explicitly so callers never need to guard against missing configs.
5. Cover edge cases: one side has an `mcp_server` key the other does not; both sides have the same key; `vars:` on both sides.

## Go learning focus

- passing and returning pointer types (`*Manifest`) — when nil is a valid value vs when it signals an error
- map merge in Go: iterating over one map and writing into another
- struct-by-struct field decisions: you cannot `merge` a struct blindly, each field needs an explicit rule
- keeping Merge a pure function (no I/O, no side effects) so it is trivially testable

## Learn more

- Credible sources: [Effective Go](https://go.dev/doc/effective_go), [Go spec — composite literals](https://go.dev/ref/spec#Composite_literals), [Go maps in action](https://go.dev/blog/maps)
- Search keywords: "go merge two structs", "go map merge pattern", "go nil pointer receiver", "go deep copy struct"

## Watch for

- mutating an input manifest — callers may still hold references to it; return a new value
- silently dropping fields from one side without an explicit decision
- encoding `mode:` logic in multiple places — it belongs only in `manifestMerger.go`
- assuming the project manifest always exists; a user may only have a global config

## Definition of done

- [ ] `Merge(global, project *Manifest) (*Manifest, error)` (or similar signature) is implemented in `internal/manifest/manifestMerger.go`.
- [ ] All three modes (`merge`, `replace`, `extend`) have explicit behavior.
- [ ] Nil inputs are handled without panicking.
- [ ] Table-driven tests in `manifest_test.go` cover all three modes and the nil cases.
- [ ] The output of Merge is the correct input to Resolve (step 04).
