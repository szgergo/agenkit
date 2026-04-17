# Step 09 — Apply Orchestration

> **Depends on:** Steps 02–08
>
> **Pipeline position:** `apply` is where the full four-stage manifest pipeline runs end to end: **Load → Merge → Resolve → Generate/Write/Lock**. The apply command owns the orchestration — it calls each stage explicitly in order, validates between stages where necessary, and keeps side effects (writes, lock updates) at the end. The `diff` command (step 10) reuses the same pipeline up to and including Generate, then stops before Write.

## Outcome

You have the first full end-to-end flow: load configuration, resolve it, detect targets, generate provider output, write safely, and update internal state in a way that is understandable and trustworthy.

## Why this step exists

This is the first point where the system behaves like a product instead of a set of building blocks. It also exposes whether your earlier boundaries are actually usable.

## Decisions you should make

- What are the exact stages of the apply pipeline?
- Where should validation happen, and what must happen before any writes occur?
- How should failures across multiple providers behave?
- What user-facing output is necessary for trust without becoming noisy?
- How tightly should `apply`, `diff`, and internal generation logic be coupled?

## Suggested work order

1. Write down the pipeline in stages before implementing it.
2. Decide which stages are pure transformation and which ones touch the outside world.
3. Make the provider loop explicit and predictable.
4. Ensure write behavior and lock updates happen in a defensible order.
5. Revisit the user-facing output once the flow works end to end.

## Go learning focus

- orchestration versus domain logic boundaries
- explicit error propagation across multi-stage flows
- keeping side effects near the edges of the system
- designing command behavior that remains readable as features accumulate

## Learn more

- Credible sources: [Cobra package docs](https://pkg.go.dev/github.com/spf13/cobra), [Go error handling guidance](https://go.dev/blog/error-handling-and-go), [Effective Go](https://go.dev/doc/effective_go)
- Search keywords: "cobra command lifecycle", "go orchestration error propagation", "cobra RunE", "go multi stage pipeline design"

## Watch for

- letting the command layer learn too much about every internal detail
- performing destructive writes before all prerequisite checks are done
- making warning, error, and info output indistinguishable to users

## Definition of done

- [ ] The apply pipeline has clear, ordered stages.
- [ ] Writes happen only after the necessary checks succeed.
- [ ] Provider failures are handled intentionally rather than accidentally.
- [ ] Lock state and provider files stay consistent with the applied result.
