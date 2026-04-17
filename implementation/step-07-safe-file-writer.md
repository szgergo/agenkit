# Step 07 — Safe File Writer

> **Depends on:** Steps 05–06
>
> **Pipeline position:** Write is stage 4 in the four-stage manifest pipeline: **Load → Merge → Resolve → Generate/Write**. Generators (step 06) produce provider config content in memory. This step is responsible for getting that content onto disk safely — without corrupting an existing file if something goes wrong mid-write, and without touching files that have not changed.

## Outcome

You have a safe, reusable write primitive that all provider output goes through. It handles atomic writes, optional backup of existing files, and a dry-run mode that returns what *would* be written without touching the filesystem. The `diff` command (step 10) and `apply` command (step 09) both depend on this.

## Why this step exists

Writing generated config directly with `os.WriteFile` is unsafe: a crash or error mid-write leaves a partially written, corrupt provider config. Provider configs (Claude Desktop, Cursor, Copilot) are often the user's only working config — corrupting them silently would be a serious trust failure. A dedicated Write step makes the safety guarantee explicit and testable before the apply command is built.

## Decisions you should make

- What is the atomic write pattern in Go, and when is it safe to use `os.Rename`?
- Should backup files be created automatically, or only on request?
- How does dry-run mode fit into the writer's interface — a flag, a separate function, or a mode enum?
- What happens if the destination directory does not exist?
- Should the writer compare content before writing to avoid unnecessary writes (and avoid unnecessary git diffs)?

## Suggested work order

1. Implement the atomic write (temp file + rename) and test it with a simple fixture.
2. Add content comparison — skip the write if the new content is byte-for-byte identical.
3. Add backup behavior — write a `.bak` copy of the existing file before replacing it.
4. Add dry-run support — return the content that would be written, without touching disk.
5. Wire the writer into a stub call from the apply command scaffolding so the interface is validated early.

## Go learning focus

- `os.CreateTemp` — creates a temp file in the same directory as the target (important for `os.Rename` to be atomic across filesystems)
- `os.Rename` — atomic on most OSes when source and destination are on the same filesystem
- `defer` for cleanup — closing and removing temp files if an error occurs before rename
- `bytes.Equal` for content comparison before writing
- `os.MkdirAll` for ensuring destination directories exist

## Learn more

- Credible sources: [Go os package](https://pkg.go.dev/os), [Effective Go](https://go.dev/doc/effective_go), [Go file I/O patterns](https://pkg.go.dev/io)
- Search keywords: "go atomic file write temp rename", "os.CreateTemp go", "go safe file write pattern", "go deferred cleanup temp file"

## Watch for

- creating temp files in `/tmp` rather than the destination directory — `os.Rename` across filesystems is not atomic
- not closing the temp file handle before calling `os.Rename`
- forgetting to remove the temp file when an error occurs after creation but before rename
- silently skipping writes when permissions are wrong — surface the error

## Definition of done

- [ ] Atomic write (temp + rename) is implemented and tested.
- [ ] Content comparison skips unnecessary writes.
- [ ] Backup behavior is implemented and opt-in.
- [ ] Dry-run mode returns content without writing.
- [ ] Error cases (missing dir, bad permissions) surface errors rather than silently failing.
- [ ] At least one provider generator (step 06) uses the writer.
