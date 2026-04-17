# Step 06 — Provider Config Generation

> **Depends on:** Steps 02 and 05
>
> **Note:** Generators receive a manifest that has already been resolved by step 04 — they never see raw `{home}` or `{vars.x}` placeholders. Resolution is the caller's responsibility (step 09 and similar). Safe file writes are handled by the writer (step 07); generators return content in memory only. This is why step 03 is not listed as a hard dependency here: the resolver and generators are pipeline peers, not a direct call chain within this step.

## Outcome

Each supported provider can be generated from the shared manifest model in a way that respects provider-specific format differences without losing the shared semantics agenkit is supposed to preserve.

## Why this step exists

This is where the “one source, many targets” promise becomes real. It is also where weak abstraction decisions show up quickly.

## Decisions you should make

- What should stay shared across providers, and what should remain intentionally provider-specific?
- How will you preserve deterministic output so generated files diff cleanly?
- How should provider overrides be represented during generation?
- What should happen when a manifest concept has no direct equivalent in a provider?
- How much read/modify/write behavior is necessary versus how much can be generated from a clean projection?

## Suggested work order

1. Pick one provider and prove the projection flow end to end.
2. Identify what logic is genuinely shared after implementing a concrete provider, not before.
3. Expand to the remaining providers while preserving their differences.
4. Validate generated output against representative fixtures.
5. Re-check whether your shared helpers still deserve to be shared.

## Go learning focus

- serialization boundaries across YAML, JSON, and TOML
- deterministic data shaping even when internal data structures are unordered
- composition of shared helpers without forcing everything into one abstraction
- keeping provider-specific logic obvious rather than hidden behind generic layers

## Learn more

- Credible sources: [Go encoding/json package](https://pkg.go.dev/encoding/json), [BurntSushi TOML docs](https://pkg.go.dev/github.com/BurntSushi/toml), [go.yaml.in/yaml/v4 docs](https://pkg.go.dev/go.yaml.in/yaml/v4)
- Search keywords: "go json tags", "MarshalIndent go", "burntsushi toml encode", "go deterministic output map order"

## Watch for

- flattening provider differences in the name of reuse
- introducing shared helpers before you have seen repetition clearly
- letting output stability depend accidentally on map iteration order

## Definition of done

- [ ] Each target provider can be generated from the manifest model.
- [ ] Generated output is stable enough for review and diffing.
- [ ] Provider-specific differences remain understandable.
- [ ] Shared logic exists only where multiple providers genuinely need it.
