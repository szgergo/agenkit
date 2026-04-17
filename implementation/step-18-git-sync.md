# Step 18 — Git Sync

> **Depends on:** Step 01
> **Can be worked on alongside:** Steps 10-17 once the base command structure exists

## Outcome

Users can manage git-backed syncing of agenkit state with command behavior that is explicit about what mutates the repository and what does not.

## Why this step exists

Sync is a core product promise, but it is easy to make it feel dangerous if commands blur together or silently trigger more actions than users expected.

## Decisions you should make

- Which sync actions should exist as separate commands, and which should be chainable by explicit choice?
- What repo location and initialization model best match the design?
- When should commits happen automatically, and how visible should that be?
- What should clone do by itself, and what should it intentionally not do?
- How should git failures be surfaced so users can recover confidently?

## Suggested work order

1. Define the command semantics before implementing any `git` execution.
2. Map each command to the exact repository mutation it is allowed to perform.
3. Implement the shell-out path with clear error reporting.
4. Cover clone, init, status, pull, and push behaviors separately.
5. Re-check that chaining remains explicit rather than hidden.

## Go learning focus

- `os/exec` and the lifecycle of external commands
- process output handling and debugging ergonomics
- making side-effect-heavy commands understandable

## Learn more

- Credible sources: [Go os/exec package](https://pkg.go.dev/os/exec), [Go context package](https://pkg.go.dev/context), [Go error handling guidance](https://go.dev/blog/error-handling-and-go)
- Search keywords: "go exec command", "exec command context go", "go capture stderr stdout", "go git command wrapper"

## Watch for

- hiding git behavior behind overly friendly abstractions
- collapsing clone, apply, pull, and push into one surprising workflow
- assuming a clean repository state when the real world is messier

## Definition of done

- [ ] Each sync command has a clear, bounded effect.
- [ ] Clone and apply remain separate unless the user explicitly chains them.
- [ ] Git failures are visible enough to debug.
- [ ] Automatic commit behavior is intentional and reviewable.
