# Step 14 — Import Command

> **Depends on:** Steps 02, 05, and 06
> **Can be worked on alongside:** Steps 10-13 once provider reading and manifest modeling are reliable

## Outcome

Users can reverse-engineer existing provider config into a manifest that is reviewable, understandable, and as lossless as the product design allows.

## Why this step exists

Import is one of the strongest onboarding features. It is also one of the riskiest because users will assume the imported result preserves intent, not just syntax.

## Decisions you should make

- How will you identify equivalent servers across providers?
- When should differences be merged, duplicated, or surfaced for a user decision?
- What belongs in provider overrides versus separate entries?
- How much should import normalize names and shapes?
- Where should interactive conflict handling begin and end?

## Suggested work order

1. Write down the import conflict cases from the design.
2. Define the deduplication and conflict rules before implementing UI around them.
3. Start with fixture-based imports from a small set of representative provider configs.
4. Add the user-decision path only where silent resolution would lose intent.
5. Re-check that the generated manifest is something a human would actually want to maintain.

## Go learning focus

- reconciling multiple external data shapes into one internal model
- making transformation logic reviewable with fixtures
- keeping ambiguity explicit instead of silently “fixing” it away

## Learn more

- Credible sources: [Go encoding/json package](https://pkg.go.dev/encoding/json), [BurntSushi TOML docs](https://pkg.go.dev/github.com/BurntSushi/toml), [Table-driven tests](https://go.dev/wiki/TableDrivenTests)
- Search keywords: "go json unmarshal struct", "toml decode file go", "go reconcile configs", "go fixtures testdata"

## Watch for

- aggressive deduplication that destroys meaningful differences
- naming imported resources in ways users cannot reason about later
- burying important choices in import heuristics that users never see

## Definition of done

- [ ] Import preserves intent well enough to become a maintainable manifest.
- [ ] Conflict handling is explicit where it needs to be.
- [ ] Equivalent entries and genuinely different entries are treated differently.
- [ ] Users can understand the imported result without reading the importer code.
