# Step 21 — Phase 2: Hooks

> **Depends on:** the completed Phase 1 foundation

## Outcome

agenkit can model, validate, inspect, and project lifecycle hooks across providers without breaking the clarity of the Phase 1 MCP foundation.

## Why this step exists

Hooks are the first serious test of whether the original architecture can grow. If Phase 1 was designed too narrowly, this step will expose it immediately.

## Decisions you should make

- What is the shared hook model, and what should remain provider-specific?
- Which event vocabulary belongs in the manifest contract?
- How should provider capability differences be represented?
- When a provider does not support a hook concept, should agenkit warn, skip, or fail?
- How should hooks appear in validation, inspection commands, and lock state?

## Suggested work order

1. Define the shared hook behavior from the product perspective before choosing types.
2. Map manifest-level hook concepts to provider capabilities explicitly.
3. Extend validation and generation in a way that preserves the Phase 1 boundaries.
4. Add inspection and apply behavior so hooks are not a hidden side feature.
5. Cover supported, partially supported, and unsupported cases with fixtures.

## Go learning focus

- extending an existing model without destabilizing it
- capability mapping between a shared contract and provider-specific behavior
- backward-compatible growth of internal types and serialized config

## Learn more

- Credible sources: [Effective Go](https://go.dev/doc/effective_go), [Go error handling guidance](https://go.dev/blog/error-handling-and-go), [go.yaml.in/yaml/v4 docs](https://pkg.go.dev/go.yaml.in/yaml/v4)
- Search keywords: "go backward compatible struct changes", "go capability mapping", "go warnings vs errors", "yaml struct evolution go"

## Watch for

- forcing all providers into one fake “universal” hook model
- making unsupported behavior silent when users need to know about it
- adding extension points so early that the feature becomes harder to understand than it needs to be

## Definition of done

- [ ] Hooks fit the shared manifest model cleanly.
- [ ] Provider support differences are handled intentionally.
- [ ] Validation, inspection, and apply behavior all understand hooks.
- [ ] The Phase 1 architecture still feels coherent after this extension.
