# Step 17 — Init Wizard

> **Depends on:** Steps 02, 05, and 15
> **Can be worked on alongside:** Step 18 once provider detection and registry access are stable

## Outcome

New users can answer a short sequence of plain-text prompts and end up with a sensible starting manifest instead of a blank file and a long setup burden.

## Why this step exists

The product becomes much easier to adopt when users can bootstrap a usable configuration without first understanding the full manifest format.

## Decisions you should make

- What is the minimum set of questions needed to create a useful first manifest?
- Which values should be detected automatically, and which ones should be explicitly confirmed?
- How much import and registry discovery belongs in the init flow?
- What should happen in non-interactive environments?
- How will users review what init is about to write?

## Suggested work order

1. Define the shortest successful onboarding path.
2. Decide which optional branches genuinely help first-time users.
3. Keep the prompt layer thin and focused on gathering decisions.
4. Make the generated output reviewable before it is saved.
5. Re-check that init teaches users the model without forcing them to learn everything at once.

## Go learning focus

- input and output boundaries for command-line interaction
- keeping interactive flows readable without introducing framework complexity
- deciding where prompting belongs versus where domain logic belongs

## Learn more

- Credible sources: [Go bufio package](https://pkg.go.dev/bufio), [Go os package](https://pkg.go.dev/os), [Cobra guides](https://cobra.dev/)
- Search keywords: "go bufio scanner stdin", "go prompt user input", "go detect terminal stdin", "cobra interactive command"

## Watch for

- turning init into a catch-all workflow for unrelated features
- asking users questions the tool could answer safely itself
- hiding too much of the generated result from the user

## Definition of done

- [ ] A new user can create a useful starting manifest through plain-text interaction.
- [ ] Auto-detection and prompting work together instead of fighting each other.
- [ ] The final write is reviewable and intentional.
- [ ] Non-interactive behavior is explicit and predictable.
