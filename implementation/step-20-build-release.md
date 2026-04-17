# Step 20 — Build and Release

> **Depends on:** Step 01
> **Can be worked on alongside:** Most other steps, but it is best finalized once the command surface is mostly stable

## Outcome

You have a reproducible way to build, package, and release the tool without relying on one specific machine or one developer remembering the right commands.

## Why this step exists

Release work is not just packaging. It forces you to decide what counts as the product, how versions are represented, and which automation should live in the repository versus external systems.

## Decisions you should make

- What build targets and artifacts matter for the project’s current stage?
- How should version information enter the binary?
- Which steps should be automated locally, and which belong in CI?
- What release process is worth supporting before the project is feature-complete?
- How much packaging complexity is justified right now?

## Suggested work order

1. Define the minimum reproducible build story.
2. Decide how versioning and artifact naming should work.
3. Add automation that supports repeatability before adding convenience flourishes.
4. Align local build behavior with CI expectations.
5. Re-check that release automation still matches the actual supported platforms and workflows.

## Go learning focus

- module-aware builds and reproducibility
- separating build concerns from runtime concerns
- thinking about delivery as part of engineering design, not just as an afterthought

## Learn more

- Credible sources: [Go modules reference](https://go.dev/ref/mod), [Organizing a Go module](https://go.dev/doc/modules/layout), [Go command docs](https://pkg.go.dev/cmd/go)
- Search keywords: "go build ldflags version", "go install command module", "go cross compile GOOS GOARCH", "go release automation"

## Watch for

- building an elaborate release pipeline before the product surface is stable
- letting CI become the only place where builds reliably work
- supporting more packaging variants than the project can maintain

## Definition of done

- [ ] Another machine or CI environment can build the project predictably.
- [ ] Versioning and artifact expectations are explicit.
- [ ] Build automation is useful without being over-engineered.
- [ ] Release behavior matches the project’s actual scope.
