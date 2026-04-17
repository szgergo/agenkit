# Step 13 — Config View Command

> **Depends on:** Steps 02 and 04
> **Can be worked on alongside:** Steps 10-12 once merging and resolution rules are clear

## Outcome

Users can inspect what agenkit believes the configuration is at the points where that inspection is genuinely useful, such as after layering and after value resolution.

## Why this step exists

This command becomes the main tool for understanding agenkit’s internal interpretation without reading code or guessing from generated provider files.

## Decisions you should make

- Which representations are worth exposing to users?
- How should raw, merged, and resolved views differ?
- Which output formats actually help users inspect state?
- How much internal detail is too implementation-specific to expose as a public command?
- Should this command explain provenance or just show the resulting state?

## Suggested work order

1. Decide the user questions this command is meant to answer.
2. Define the supported views around those questions.
3. Make output deterministic and readable before optimizing for more formats.
4. Add cases that prove layering and resolution behavior are inspectable.
5. Check that this command complements `validate`, `diff`, and `apply` rather than overlapping confusingly.

## Go learning focus

- separating internal representation from presentation
- formatting structured output without leaking accidental implementation details
- keeping one command useful across both debugging and normal use

## Learn more

- Credible sources: [Cobra guides](https://cobra.dev/), [Go encoding/json package](https://pkg.go.dev/encoding/json), [go.yaml.in/yaml/v4 docs](https://pkg.go.dev/go.yaml.in/yaml/v4)
- Search keywords: "go pretty print json", "yaml marshal go", "cobra output format flags", "go inspect config state"

## Watch for

- exposing every internal intermediate state just because it exists
- creating output modes with no clear user purpose
- turning this command into a substitute for proper validation

## Definition of done

- [ ] The command answers clear user questions about configuration state.
- [ ] Supported views are intentional and understandable.
- [ ] Output is stable enough for inspection and automation where appropriate.
- [ ] Layering and resolution behavior can be checked without guesswork.
