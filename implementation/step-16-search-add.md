# Step 16 — Search and Add Commands

> **Depends on:** Steps 02 and 15

## Outcome

Users can discover MCP servers from configured registries and add them to the manifest in a way that feels deliberate rather than magical.

## Why this step exists

Search and add are the bridge between external discovery and local configuration. If the UX is sloppy here, users will either stop trusting registry results or stop trusting what gets written into their manifest.

## Decisions you should make

- What information should search results expose by default?
- When should `add` be fully non-interactive versus prompt-driven?
- How should duplicate or near-duplicate additions be handled?
- What normalization should happen on insert, and what should stay close to the registry data?
- How will users understand where an added entry came from?

## Suggested work order

1. Decide the user flow from discovery to manifest update.
2. Make the search output useful before optimizing the add path.
3. Define how additions become manifest entries, including naming and defaults.
4. Handle ambiguity intentionally instead of silently picking one interpretation.
5. Re-check that the resulting manifest still looks like something a human would maintain.

## Go learning focus

- command UX design around real user decisions
- mutating structured config safely
- keeping interactive and scriptable behaviors from collapsing into one messy path

## Learn more

- Credible sources: [Cobra package docs](https://pkg.go.dev/github.com/spf13/cobra), [Cobra guides](https://cobra.dev/), [Go error handling guidance](https://go.dev/blog/error-handling-and-go)
- Search keywords: "cobra flag validation", "cobra interactive selection", "go update yaml file", "go duplicate detection"

## Watch for

- auto-filling too much and creating opaque manifest entries
- mixing registry transport concerns with manifest mutation logic
- making ambiguity resolution inconsistent across search and add

## Definition of done

- [ ] Search results are useful enough to choose from confidently.
- [ ] Add behavior is intentional in both interactive and direct modes.
- [ ] Duplicate handling is explicit.
- [ ] Inserted entries remain maintainable by a human.
