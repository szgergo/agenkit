# Step 11 — Get, Describe, and Delete Commands

> **Depends on:** Step 02
> **Can be worked on alongside:** Steps 10, 12, and 13 after the manifest model is stable

## Outcome

Users can inspect what agenkit knows about resources and remove manifest-defined resources intentionally, with command behavior that feels consistent and unsurprising.

## Why this step exists

Tools become much easier to trust when users can inspect state before and after changes. Deletion is especially important because it tests whether your resource model is actually clear.

## Decisions you should make

- What counts as a resource in the CLI surface?
- Which outputs belong in list views versus detail views?
- What exactly does delete change: manifest state, generated state, or both?
- How should global versus project context affect what users see and remove?
- What level of confirmation or warning is appropriate for destructive actions?

## Suggested work order

1. Define the resource vocabulary and command behavior before writing handlers.
2. Make listing and detail views clear before worrying about output variations.
3. Decide how delete interacts with layered config and follow-up apply behavior.
4. Add tests around resource lookup, missing resources, and destructive intent.
5. Re-check that the command names and outputs still feel aligned with the `kubectl`-style direction.

## Go learning focus

- modeling command behavior around domain concepts instead of internal implementation details
- separating presentation logic from resource mutation logic
- making mutation paths explicit and testable

## Learn more

- Credible sources: [Cobra package docs](https://pkg.go.dev/github.com/spf13/cobra), [Cobra guides](https://cobra.dev/), [Go error handling guidance](https://go.dev/blog/error-handling-and-go)
- Search keywords: "cobra subcommands", "cobra exact args", "cobra help text", "go cli resource commands"

## Watch for

- exposing internal implementation terms as user-facing resources by accident
- making delete semantics implicit or surprising
- mixing inspection output with mutation side effects

## Definition of done

- [ ] Users can list resources in a predictable way.
- [ ] Detail views answer more than the list views without becoming noisy.
- [ ] Delete semantics are explicit and consistent with the rest of the workflow.
- [ ] Missing-resource and wrong-scope behavior is clear.
