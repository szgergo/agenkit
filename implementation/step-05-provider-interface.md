# Step 05 — Provider Boundary and Detection

> **Depends on:** Steps 01-02

## Outcome

You have a provider boundary that fits the actual product responsibilities and a detection strategy that follows the intended precedence rules without scattering provider-specific logic across the whole project.

## Why this step exists

Providers are the main abstraction boundary in agenkit. If this boundary is too large, every provider becomes hard to implement. If it is too small, later orchestration becomes awkward and repetitive.

## Decisions you should make

- What is the smallest shared provider contract that supports the current product?
- Which responsibilities belong to the shared provider boundary, and which should stay provider-specific?
- How should detection results be represented so callers can reason about installed, missing, and explicitly configured providers?
- Where should detection precedence be enforced so it stays consistent everywhere?
- Do you need capability metadata already, or can that wait until later phases?

## Suggested work order

1. List the provider behaviors that multiple commands genuinely need.
2. Define the shared contract around those behaviors only.
3. Design a detection approach that respects explicit config, defaults, and environment-derived paths in the right order.
4. Check that the boundary supports current providers without assuming future providers must look the same.
5. Add tests around detection precedence and provider registration behavior.

## Go learning focus

- interfaces as small behavioral contracts
- pointer versus value receivers and how they affect method sets
- package boundaries that keep provider details from leaking into unrelated areas
- choosing concrete types first and interfaces where polymorphism is real

## Learn more

- Credible sources: [Go language tour](https://go.dev/tour/), [Effective Go](https://go.dev/doc/effective_go), [The internal package boundary](https://go.dev/doc/go1.4#internalpackages)
- Search keywords: "go interfaces small interfaces", "go method sets", "pointer vs value receiver go", "go internal package"

## Watch for

- building one oversized interface because it feels “clean”
- letting provider detection rules spread across commands and helpers
- designing for hypothetical providers that do not exist yet

## Definition of done

- [ ] The provider boundary is small, understandable, and sufficient for current needs.
- [ ] Detection precedence is explicit and testable.
- [ ] Callers do not need provider-specific conditionals for ordinary flows.
- [ ] You can explain why each shared provider method belongs in the shared contract.
