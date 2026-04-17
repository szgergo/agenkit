package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGatherConfigFromPath(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name             string
		setupPath        func(t *testing.T) string
		wantNilManifest  bool
		wantErrContains  string
		expectErr        bool
		expectedVersion  int
	}

	testCases := []testCase{
		{
			name: "missing file returns nil manifest and nil error",
			setupPath: func(t *testing.T) string {
				t.Helper()
				return filepath.Join(t.TempDir(), ".agenkit", "agenkit.yml")
			},
			wantNilManifest: true,
		},
		{
			name: "non-not-exist read error is returned",
			setupPath: func(t *testing.T) string {
				t.Helper()
				return t.TempDir()
			},
			expectErr:       true,
			wantErrContains: "reading file",
		},
		{
			name: "parse error is returned",
			setupPath: func(t *testing.T) string {
				t.Helper()

				path := filepath.Join(t.TempDir(), "agenkit.yml")
				if err := os.WriteFile(path, []byte("version: ["), 0o644); err != nil {
					t.Fatalf("write invalid manifest: %v", err)
				}
				return path
			},
			expectErr:       true,
			wantErrContains: "invalid config file",
		},
		{
			name: "valid manifest is loaded",
			setupPath: func(t *testing.T) string {
				t.Helper()

				path := filepath.Join(t.TempDir(), "agenkit.yml")
				if err := os.WriteFile(path, []byte("version: 1\n"), 0o644); err != nil {
					t.Fatalf("write valid manifest: %v", err)
				}
				return path
			},
			expectedVersion: 1,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			path := tc.setupPath(t)
			gotManifest, err := gatherConfigFromPath(path)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tc.wantErrContains != "" && !strings.Contains(err.Error(), tc.wantErrContains) {
					t.Fatalf("expected error to contain %q, got %q", tc.wantErrContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if tc.wantNilManifest {
				if gotManifest != nil {
					t.Fatalf("expected nil manifest, got %+v", gotManifest)
				}
				return
			}

			if gotManifest == nil {
				t.Fatalf("expected manifest, got nil")
			}
			if gotManifest.Version != tc.expectedVersion {
				t.Fatalf("expected version %d, got %d", tc.expectedVersion, gotManifest.Version)
			}
		})
	}
}
