# Step 08 — Lock State Management

> **Depends on:** Step 02

## Outcome

You have a design for tracking what agenkit manages so that generated changes remain non-destructive and later cleanup operations can distinguish tool-managed state from user-managed state.

## Why this step exists

Lock state is central to user trust. Without it, `apply`, `delete`, and later reconciliation features will either overwrite too much or become too conservative to be useful.

## Decisions you should make

- What identity should a lock entry record so later operations can reason about managed state correctly?
- How much historical or structural detail should the lock keep?
- How should lock versioning work so future changes remain survivable?
- What should happen if lock state is missing, stale, or inconsistent with what is on disk?
- Which behaviors depend on lock state and which should never depend on it?

## Suggested work order

1. Define the specific problems lock state must solve.
2. Design the minimal persisted shape that solves those problems.
3. Think through drift scenarios before writing the first serializer.
4. Add tests that compare previous managed state with proposed managed state.
5. Check that later commands such as `diff`, `delete`, and `apply` can all rely on the same lock story.

## Go learning focus

- serializing internal metadata cleanly
- file I/O and failure handling for internal state files
- versioned data formats and compatibility thinking
- designing internal persistence that is useful without becoming the real source of truth

## Learn more

- Credible sources: [Go os package](https://pkg.go.dev/os), [Go path/filepath package](https://pkg.go.dev/path/filepath), [Go error handling guidance](https://go.dev/blog/error-handling-and-go)
- Search keywords: "go write file safely", "go filepath join", "go json file persistence", "go error wrapping"

## Watch for

- storing so much state that the lock becomes a second config system
- storing so little state that you cannot make safe decisions later
- treating the lock as authoritative when the manifest should still lead

## Definition of done

- [ ] Lock state has a clear purpose and a clear boundary.
- [ ] The stored shape is minimal but sufficient for non-destructive behavior.
- [ ] Drift scenarios are thought through, not left accidental.
- [ ] Later commands can rely on the lock design without special cases everywhere.
