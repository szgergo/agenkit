# Known Issues and TODOs

This document tracks issues and code improvements discovered during implementation audits that are too specific to fix out of order. They are tracked here to be addressed during the relevant implementation step.

## Code Issues

### 1. `manifestLoader.go` — missing file should return `(nil, nil)` 

**Location:** `internal/manifest/manifestLoader.go`

**Issue:** `LoadGlobalManifest()` and `LoadProjectManifest()` currently propagate `os.ReadFile`'s error when a config file doesn't exist, returning `(nil, error)`. However, a missing manifest is a valid state — users may have no global config, or be in a directory with no project config.

**Impact:** The Merge function (step 03) must be able to receive `nil` inputs without error. Until this is fixed, the caller cannot distinguish between "file not found" (which should succeed with nil) and "file exists but is invalid" (which should fail with error).

**Fix:** In `gatherConfigFromPath()`, detect `os.IsNotExist(err)` and return `(nil, nil)` instead of propagating the error. Other errors (permission denied, disk read failure) should still propagate.

**Tracked against:** **Step 03 — Manifest Merge**

---

### 2. `manifest_test.go` — hand-rolled `contains()` should use `strings.Contains`

**Location:** `internal/manifest/manifest_test.go`

**Issue:** The test file has a manual reimplementation of `strings.Contains` at the bottom of the file:

```go
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
```

This is unnecessary — Go's standard library provides `strings.Contains()`.

**Fix:** 
1. Add `"strings"` to the imports if not already present
2. Remove the `contains()` helper function
3. Replace the one call site `contains(err.Error(), "mode")` with `strings.Contains(err.Error(), "mode")`

**Tracked against:** **Step 02 — Manifest Data Model** (test quality improvement)
