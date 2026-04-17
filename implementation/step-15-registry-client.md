# Step 15 — Registry Client

> **Depends on:** Step 02
> **Can be worked on alongside:** Steps 10-14 once the manifest model is stable

## Outcome

agenkit can query one or more MCP registries, cache useful results, and surface normalized information that later commands can depend on without knowing registry-specific quirks.

## Why this step exists

Registry support is where the project starts interacting with outside systems that can be slow, unavailable, or inconsistent. The design here will shape reliability and user trust.

## Decisions you should make

- What is the internal result shape that search and add commands should consume?
- What cache behavior is worth implementing now?
- What are the timeout and fallback expectations?
- How should multi-registry results be merged and labeled?
- How should offline and partial-failure behavior be surfaced to users?

## Suggested work order

1. Define the client boundary and result shape first.
2. Decide on timeout, caching, and fallback rules before adding concurrency.
3. Implement the simplest reliable request path.
4. Add multi-registry merging and source labeling.
5. Cover online, cached, offline, and partial-failure scenarios with focused tests.

## Go learning focus

- `net/http` and request lifecycle basics
- `context` for timeout and cancellation control
- time-based cache reasoning
- keeping network code simple until complexity is justified

## Learn more

- Credible sources: [Go net/http package](https://pkg.go.dev/net/http), [Go context package](https://pkg.go.dev/context), [Go time package](https://pkg.go.dev/time)
- Search keywords: "go http client timeout", "newrequestwithcontext go", "go caching ttl", "go partial failure handling"

## Watch for

- building concurrency first and behavior second
- hiding registry failures so thoroughly that debugging becomes hard
- letting cache behavior become implicit or surprising

## Definition of done

- [ ] Registry results have a stable internal shape.
- [ ] Cache, live fetch, and fallback behavior are explicit.
- [ ] Multi-registry results remain understandable.
- [ ] Later commands can consume the registry client without caring about transport details.
