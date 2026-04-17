# Step 12 — Validate Command

> **Depends on:** Step 02
> **Can be worked on alongside:** Steps 10, 11, and 13 once the manifest model is stable

## Outcome

Users can ask agenkit whether configuration is structurally and semantically sound before applying it, and they get results that are specific enough to act on.

## Why this step exists

Validation is where you turn hidden assumptions into explicit rules. It also forces you to decide what the parser guarantees, what later commands require, and what should merely produce warnings.

## Decisions you should make

- Which checks belong to parsing and which belong to validation?
- Which problems should be errors versus warnings?
- Should a missing project `mode:` produce a warning and default to merge, or should it be a hard error?
- Should validation run on raw, merged, resolved, or multiple views of the config?
- How much guidance should validation output include?
- How will you report multiple issues without making the output unreadable?

## Suggested work order

1. Write down the classes of problems the tool should detect.
2. Split structural checks from semantic checks intentionally.
3. Decide how validation will surface multiple findings.
4. Add fixtures for both ordinary mistakes and design-specific edge cases.
5. Re-check that validation output supports later commands instead of duplicating them.

## Go learning focus

- aggregated error reporting
- table-driven testing for rule combinations
- designing user-facing diagnostics that still map back to plain Go logic

## Learn more

- Credible sources: [Go error handling guidance](https://go.dev/blog/error-handling-and-go), [Go testing package](https://pkg.go.dev/testing), [Table-driven tests](https://go.dev/wiki/TableDrivenTests)
- Search keywords: "go aggregated errors", "go validation errors", "go table driven tests", "go warnings vs errors"

## Watch for

- duplicating parsing logic inside validation
- turning every odd situation into a hard error
- making diagnostics so verbose that users stop reading them
- treating project-level missing `mode:` as silent behavior instead of a user-facing warning

## Definition of done

- [ ] Validation rules are intentionally scoped.
- [ ] Error versus warning behavior is consistent.
- [ ] Multiple findings can be surfaced in one run.
- [ ] The command helps users fix problems instead of merely announcing failure.
