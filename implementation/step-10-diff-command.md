# Step 10 — Diff Command

> **Depends on:** Step 09
> **Can be worked on alongside:** Steps 11-13 once the core model and apply pipeline are stable

## Outcome

Users can preview what `apply` would change without writing anything. The preview should feel trustworthy because it reflects the same underlying decisions as `apply`, not a second interpretation of the system.

## Why this step exists

Preview is one of the main trust-building features in a tool that edits other tools’ config files. If preview and apply diverge, users will stop trusting both.

## Decisions you should make

- What should be the comparison unit: whole file text, logical sections, or both?
- How closely should `diff` share logic with `apply`?
- What should “no changes” look like to a user?
- How much formatting and color is helpful before output becomes noisy?
- Should `apply --dry-run` be a thin alias to the same behavior or a separate path?

## Suggested work order

1. Reuse the same preparation pipeline as `apply` up to the point just before writing.
2. Decide what “current” versus “proposed” means for each provider.
3. Design output that helps review changes quickly.
4. Cover additions, removals, modifications, and no-op cases with fixtures.
5. Verify that preview semantics still match apply semantics after refactoring.

## Go learning focus

- reusing orchestration without duplicating behavior
- output formatting as a product concern
- keeping user-facing comparison logic understandable rather than clever

## Learn more

- Credible sources: [Cobra guides](https://cobra.dev/), [Go testing package](https://pkg.go.dev/testing), [Table-driven tests](https://go.dev/wiki/TableDrivenTests)
- Search keywords: "go diff output testing", "cobra dry run command", "go golden tests", "go compare generated files"

## Watch for

- building a preview path that slowly drifts away from `apply`
- over-investing in diff presentation before you trust the underlying comparison
- hiding important changes inside overly compact output

## Definition of done

- [ ] `diff` and `apply` agree on what would change.
- [ ] Additions, removals, and modifications are all visible.
- [ ] No-op output is clear and not ambiguous.
- [ ] The implementation does not duplicate the core preparation pipeline unnecessarily.
