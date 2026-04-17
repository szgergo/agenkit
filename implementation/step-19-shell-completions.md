# Step 19 — Shell Completions

> **Depends on:** Step 01

## Outcome

The CLI can generate or expose shell completion behavior that reflects the actual command surface and makes common commands faster to use.

## Why this step exists

Completion is a quality-of-life feature, but it only helps if it stays aligned with the real command tree and real argument expectations.

## Decisions you should make

- Which commands and flags benefit from completion enough to justify the work?
- How much should be static versus dynamically generated?
- When is the command surface stable enough that completion work will not churn constantly?
- How should installation guidance be presented to users?

## Suggested work order

1. Wait until the command surface is recognizable enough to stabilize.
2. Add basic completion support first.
3. Layer in dynamic completions only where they materially help.
4. Check that completion behavior matches help text and argument semantics.
5. Add user guidance for installation and verification.

## Go learning focus

- Cobra’s completion model
- command metadata as part of the public interface
- keeping shell integration aligned with real command behavior

## Learn more

- Credible sources: [Cobra package docs](https://pkg.go.dev/github.com/spf13/cobra), [Cobra guides](https://cobra.dev/)
- Search keywords: "cobra shell completion", "cobra valid args function", "cobra dynamic completion", "cobra completion install"

## Watch for

- implementing completion too early and rewriting it repeatedly
- adding dynamic completion where static suggestions would be simpler and clearer
- letting completion promise commands or values the CLI does not actually support

## Definition of done

- [ ] Completion support exists for the current command surface.
- [ ] Dynamic completion is used only where it clearly helps.
- [ ] User guidance is enough to enable the feature.
- [ ] Completion behavior and help text agree.
