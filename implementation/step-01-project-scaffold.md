# Step 01 — Project Scaffold

> **Depends on:** nothing

## Outcome

A runnable CLI foundation exists. You should be able to build and run the binary, see meaningful root help output, and feel that the repository shape is ready for the next few steps without already being over-designed.

## Why this step exists

This step gives you a stable place to learn from. It is not about finishing the architecture up front. It is about creating just enough shape that later decisions have somewhere sensible to live.

## Decisions you should make

- What is the smallest project structure that supports the next few steps well?
- Which responsibilities already deserve a separate package, and which ones do not yet?
- What should the root command expose immediately versus later?
- What minimal local workflow do you want for build, run, and test?

## Suggested work order

1. Establish the module and executable entry point.
2. Create only the initial project structure that the next steps clearly justify.
3. Add the root command and baseline help behavior.
4. Add a simple developer workflow for building, running, and testing.
5. Sanity-check that the project still feels small and understandable.

## Go learning focus

- what the program entry point is and why it should stay tiny
- how modules affect imports and project identity
- package visibility through naming and compiler-enforced internal boundaries
- Cobra’s command model and why command errors should flow upward cleanly

## Learn more
- Credible sources: [Go modules reference](https://go.dev/ref/mod), [Organizing a Go module](https://go.dev/doc/modules/layout), [Cobra package docs](https://pkg.go.dev/github.com/spf13/cobra), [Cobra guides](https://cobra.dev/)
- Search keywords: "go module layout", "go package visibility export", "cobra RunE", "cobra persistent flags"

## Watch for

- creating too many packages before you have real pressure for them
- letting command handlers become the place where all logic lives
- treating help text as filler instead of part of the product

## Definition of done

- [X] The root command builds and runs.
- [X] The root help text explains the tool clearly enough for a new reader.
- [X] The project layout supports the next steps without feeling inflated.
- [X] You can explain why each top-level area exists.
