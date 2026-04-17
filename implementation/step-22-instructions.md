# Step 22 — Phase 3: Instructions

> **Depends on:** the completed Phase 1 foundation
> **Best tackled after:** Step 21, once you have already extended the core model once

## Outcome

agenkit can manage instruction or policy-style configuration from the same source-of-truth model while respecting that providers may express these concepts very differently.

## Why this step exists

Instructions are the most abstraction-sensitive feature in the roadmap. They test whether the project can stay simple while still projecting richer behavior into multiple tools.

## Decisions you should make

- What belongs in a shared instruction model versus provider-specific translation?
- How should global, project, and provider-level precedence work for instructions?
- Which instruction concepts are stable enough to normalize, and which should remain provider-specific?
- How will users inspect and validate instructions before applying them?
- What boundaries keep this phase from overcomplicating the rest of the system?

## Suggested work order

1. Define the user-visible problem this feature solves before designing a schema.
2. Decide the smallest shared instruction shape that is still useful.
3. Design provider-specific projection rules and identify unsupported or lossy areas explicitly.
4. Extend validation, inspection, diff, and apply behavior where instructions need to participate.
5. Re-check that the resulting model is still something humans can maintain over time.

## Go learning focus

- schema evolution without breaking earlier decisions
- composition over framework-like abstraction
- backward-compatible thinking for configuration-heavy systems

## Learn more

- Credible sources: [Effective Go](https://go.dev/doc/effective_go), [go.yaml.in/yaml/v4 docs](https://pkg.go.dev/go.yaml.in/yaml/v4), [Go error handling guidance](https://go.dev/blog/error-handling-and-go)
- Search keywords: "go schema evolution config", "go configuration precedence", "go provider specific translation", "yaml backward compatibility"

## Watch for

- inventing an overly general instruction system for future possibilities instead of current needs
- conflating instructions with hooks because both look like “extra config”
- letting provider-specific translation rules leak everywhere in the codebase

## Definition of done

- [ ] Instructions have a clear place in the overall model.
- [ ] Precedence and provider-translation rules are explicit.
- [ ] Users can inspect and validate instruction behavior before apply.
- [ ] The final system still feels like one coherent tool rather than several overlapping subsystems.
