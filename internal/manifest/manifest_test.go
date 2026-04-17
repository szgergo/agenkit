package manifest

import (
	"os"
	"path/filepath"
	"testing"

	"go.yaml.in/yaml/v4"
)

func loadTestManifest(t *testing.T, filename string) Manifest {
	t.Helper()

	var testManifest Manifest
	path := filepath.Join(".", "..", "..", "testdata", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", filename, err)
	}

	if err := yaml.Load(data, &testManifest, yaml.WithKnownFields()); err != nil {
		t.Fatalf("parse %s: %v", filename, err)
	}

	return testManifest
}

func TestParseFullManifest(t *testing.T) {
	testManifest := loadTestManifest(t, "full_manifest_correct.yaml")

	// Test Version
	if testManifest.Version != 1 {
		t.Errorf("Version should be 1, it is %d", testManifest.Version)
		t.FailNow()
	}

	// Test Mode (zero-value expected for current fixture)
	if testManifest.Mode != "" {
		t.Errorf("Mode should be empty string, got %q", testManifest.Mode)
		t.FailNow()
	}

	// Test Vars
	if testManifest.Vars == nil {
		t.Errorf("Vars should not be nil")
		t.FailNow()
	}

	if testManifest.Vars["mcp_dir"] != "{home}/mcp-servers" {
		t.Errorf("mcp_dir should be {home}/mcp-servers, got %s", testManifest.Vars["mcp_dir"])
		t.FailNow()
	}

	if testManifest.Vars["work_dir"] != "{home}/work" {
		t.Errorf("work_dir should be {home}/work, got %s", testManifest.Vars["work_dir"])
		t.FailNow()
	}

	// Test McpServers
	if testManifest.McpServers == nil {
		t.Errorf("McpServers should not be nil")
		t.FailNow()
	}

	if len(testManifest.McpServers) != 4 {
		t.Errorf("Should have 4 MCP servers, got %d", len(testManifest.McpServers))
		t.FailNow()
	}

	// Test filesystem server
	filesystem, exists := testManifest.McpServers["filesystem"]
	if !exists {
		t.Errorf("filesystem server should exist")
		t.FailNow()
	}

	if filesystem.Command != "node" {
		t.Errorf("filesystem command should be 'node', got %s", filesystem.Command)
		t.FailNow()
	}

	if len(filesystem.Args) != 2 {
		t.Errorf("filesystem should have 2 args, got %d", len(filesystem.Args))
		t.FailNow()
	}

	if filesystem.Env["NODE_ENV"] != "production" {
		t.Errorf("NODE_ENV should be 'production', got %s", filesystem.Env["NODE_ENV"])
		t.FailNow()
	}

	// Test postgres server with providers filter
	postgres, exists := testManifest.McpServers["postgres"]
	if !exists {
		t.Errorf("postgres server should exist")
		t.FailNow()
	}

	if len(postgres.Providers) != 2 {
		t.Errorf("postgres should have 2 providers, got %d", len(postgres.Providers))
		t.FailNow()
	}

	if postgres.Providers[0] != "claude" || postgres.Providers[1] != "cursor" {
		t.Errorf("providers should be [claude, cursor], got %v", postgres.Providers)
		t.FailNow()
	}

	// Test context7 server with provider overrides
	context7, exists := testManifest.McpServers["context7"]
	if !exists {
		t.Errorf("context7 server should exist")
		t.FailNow()
	}

	if len(context7.ProviderOverrides) != 2 {
		t.Errorf("context7 should have 2 provider overrides, got %d", len(context7.ProviderOverrides))
		t.FailNow()
	}

	cursor, hasOverride := context7.ProviderOverrides["cursor"]
	if !hasOverride {
		t.Errorf("cursor override should exist")
		t.FailNow()
	}

	if len(cursor.Args) != 4 {
		t.Errorf("cursor override should have 4 args, got %d", len(cursor.Args))
		t.FailNow()
	}

	if cursor.Args[2] != "--profile" {
		t.Errorf("3rd arg should be '--profile', got %s", cursor.Args[2])
		t.FailNow()
	}

	// Test Providers
	if testManifest.Providers == nil {
		t.Errorf("Providers should not be nil")
		t.FailNow()
	}

	if len(testManifest.Providers) != 3 {
		t.Errorf("Should have 3 providers, got %d", len(testManifest.Providers))
		t.FailNow()
	}

	claudeDesktop, exists := testManifest.Providers["claude-desktop"]
	if !exists {
		t.Errorf("claude-desktop provider should exist")
		t.FailNow()
	}

	if claudeDesktop.ConfigPath != "/custom/path/to/claude_desktop_config.json" {
		t.Errorf("claude-desktop config_path should be '/custom/path/to/claude_desktop_config.json', got %s", claudeDesktop.ConfigPath)
		t.FailNow()
	}

	codex, exists := testManifest.Providers["codex"]
	if !exists {
		t.Errorf("codex provider should exist")
		t.FailNow()
	}

	if *codex.Enabled != false {
		t.Errorf("codex Enabled should be false, got %v", codex.Enabled)
		t.FailNow()
	}

}

func TestParseEmptyManifest(t *testing.T) {
	testManifest := loadTestManifest(t, "empty_manifest.yaml")

	if testManifest.Version != 0 {
		t.Errorf("Version should be 0, got %d", testManifest.Version)
	}

	// Mode should be zero-value for empty manifest
	if testManifest.Mode != "" {
		t.Errorf("Mode should be empty string for empty manifest, got %q", testManifest.Mode)
	}

	if testManifest.Vars != nil {
		t.Errorf("Vars should be nil, got %v", testManifest.Vars)
	}

	if testManifest.McpServers != nil {
		t.Errorf("McpServers should be nil, got %v", testManifest.McpServers)
	}

	if testManifest.Providers != nil {
		t.Errorf("Providers should be nil, got %v", testManifest.Providers)
	}
}

func TestParseVersionOnly(t *testing.T) {
	testManifest := loadTestManifest(t, "version_only.yaml")

	if testManifest.Version != 1 {
		t.Errorf("Version should be 1, got %d", testManifest.Version)
	}

	// Mode should be zero-value for version-only manifest
	if testManifest.Mode != "" {
		t.Errorf("Mode should be empty string for version-only manifest, got %q", testManifest.Mode)
	}

	if testManifest.Vars != nil {
		t.Errorf("Vars should be nil, got %v", testManifest.Vars)
	}

	if testManifest.McpServers != nil {
		t.Errorf("McpServers should be nil, got %v", testManifest.McpServers)
	}

	if testManifest.Providers != nil {
		t.Errorf("Providers should be nil, got %v", testManifest.Providers)
	}
}

func TestParseManifestWithModeMerge(t *testing.T) {
	testManifest := loadTestManifest(t, "manifest_with_mode_merge.yaml")

	if testManifest.Version != 1 {
		t.Errorf("Version should be 1, got %d", testManifest.Version)
	}

	if testManifest.Mode != ModeMerge {
		t.Errorf("Mode should be 'merge', got %q", testManifest.Mode)
	}
}

func TestParseManifestWithModeReplace(t *testing.T) {
	testManifest := loadTestManifest(t, "manifest_with_mode_replace.yaml")

	if testManifest.Version != 1 {
		t.Errorf("Version should be 1, got %d", testManifest.Version)
	}

	if testManifest.Mode != ModeReplace {
		t.Errorf("Mode should be 'replace', got %q", testManifest.Mode)
	}
}

func TestParseManifestWithModeExtend(t *testing.T) {
	testManifest := loadTestManifest(t, "manifest_with_mode_extend.yaml")

	if testManifest.Version != 1 {
		t.Errorf("Version should be 1, got %d", testManifest.Version)
	}

	if testManifest.Mode != ModeExtend {
		t.Errorf("Mode should be 'extend', got %q", testManifest.Mode)
	}
}

func TestParseManifestWithInvalidMode(t *testing.T) {
	path := filepath.Join(".", "..", "..", "testdata", "manifest_with_mode_invalid.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read manifest_with_mode_invalid.yaml: %v", err)
	}

	var testManifest Manifest
	err = yaml.Load(data, &testManifest, yaml.WithKnownFields())
	if err == nil {
		t.Fatalf("expected parsing to fail for invalid mode, but succeeded")
	}

	// Check that the error mentions mode validation
	if !contains(err.Error(), "mode") && !contains(err.Error(), "Invalid") {
		t.Errorf("expected error about mode validation, got: %v", err)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestParseNoVars(t *testing.T) {
	testManifest := loadTestManifest(t, "no_vars.yaml")

	if testManifest.Vars != nil {
		t.Errorf("Vars should be nil, got %v", testManifest.Vars)
	}

	if testManifest.McpServers == nil {
		t.Errorf("McpServers should not be nil")
		t.FailNow()
	}

	if len(testManifest.McpServers) != 1 {
		t.Errorf("Should have 1 MCP server, got %d", len(testManifest.McpServers))
	}

	if testManifest.Providers == nil {
		t.Errorf("Providers should not be nil")
		t.FailNow()
	}

	if len(testManifest.Providers) != 1 {
		t.Errorf("Should have 1 provider, got %d", len(testManifest.Providers))
	}
}

func TestParseNoMcpServers(t *testing.T) {
	testManifest := loadTestManifest(t, "no_mcp_servers.yaml")

	if testManifest.Vars == nil {
		t.Errorf("Vars should not be nil")
		t.FailNow()
	}

	if testManifest.McpServers != nil {
		t.Errorf("McpServers should be nil, got %v", testManifest.McpServers)
	}

	if testManifest.Providers == nil {
		t.Errorf("Providers should not be nil")
		t.FailNow()
	}
}

func TestParseNoProviders(t *testing.T) {
	testManifest := loadTestManifest(t, "no_providers.yaml")

	if testManifest.Vars == nil {
		t.Errorf("Vars should not be nil")
		t.FailNow()
	}

	if testManifest.McpServers == nil {
		t.Errorf("McpServers should not be nil")
		t.FailNow()
	}

	if testManifest.Providers != nil {
		t.Errorf("Providers should be nil, got %v", testManifest.Providers)
	}
}

func TestParseMcpServerMinimal(t *testing.T) {
	testManifest := loadTestManifest(t, "mcp_server_minimal.yaml")

	simple, exists := testManifest.McpServers["simple"]
	if !exists {
		t.Errorf("simple server should exist")
		t.FailNow()
	}

	if simple.Command != "npx" {
		t.Errorf("Command should be 'npx', got %s", simple.Command)
	}

	if simple.Args != nil {
		t.Errorf("Args should be nil, got %v", simple.Args)
	}

	if simple.Env != nil {
		t.Errorf("Env should be nil, got %v", simple.Env)
	}

	if simple.Providers != nil {
		t.Errorf("Providers should be nil, got %v", simple.Providers)
	}

	if simple.ProviderOverrides != nil {
		t.Errorf("ProviderOverrides should be nil, got %v", simple.ProviderOverrides)
	}
}

func TestParseMcpServerNoArgs(t *testing.T) {
	testManifest := loadTestManifest(t, "mcp_server_no_args.yaml")

	github, exists := testManifest.McpServers["github"]
	if !exists {
		t.Errorf("github server should exist")
		t.FailNow()
	}

	if github.Command != "npx" {
		t.Errorf("Command should be 'npx', got %s", github.Command)
	}

	if github.Args != nil {
		t.Errorf("Args should be nil, got %v", github.Args)
	}

	if github.Env == nil {
		t.Errorf("Env should not be nil")
		t.FailNow()
	}

	if github.Env["GITHUB_TOKEN"] != "{env.GITHUB_TOKEN}" {
		t.Errorf("GITHUB_TOKEN should be '{env.GITHUB_TOKEN}', got %s", github.Env["GITHUB_TOKEN"])
	}
}

func TestParseMcpServerNoEnv(t *testing.T) {
	testManifest := loadTestManifest(t, "mcp_server_no_env.yaml")

	filesystem, exists := testManifest.McpServers["filesystem"]
	if !exists {
		t.Errorf("filesystem server should exist")
		t.FailNow()
	}

	if filesystem.Command != "node" {
		t.Errorf("Command should be 'node', got %s", filesystem.Command)
	}

	if len(filesystem.Args) != 2 {
		t.Errorf("Should have 2 args, got %d", len(filesystem.Args))
	}

	if filesystem.Env != nil {
		t.Errorf("Env should be nil, got %v", filesystem.Env)
	}
}

func TestParseMcpServerNoProvidersFilter(t *testing.T) {
	testManifest := loadTestManifest(t, "mcp_server_no_providers_filter.yaml")

	github, exists := testManifest.McpServers["github"]
	if !exists {
		t.Errorf("github server should exist")
		t.FailNow()
	}

	if github.Providers != nil {
		t.Errorf("Providers should be nil (meaning all providers), got %v", github.Providers)
	}
}

func TestParseMcpServerNoOverrides(t *testing.T) {
	testManifest := loadTestManifest(t, "mcp_server_no_overrides.yaml")

	github, exists := testManifest.McpServers["github"]
	if !exists {
		t.Errorf("github server should exist")
		t.FailNow()
	}

	if github.ProviderOverrides != nil {
		t.Errorf("ProviderOverrides should be nil, got %v", github.ProviderOverrides)
	}

	if len(github.Providers) != 2 {
		t.Errorf("Should have 2 providers, got %d", len(github.Providers))
	}
}

func TestParseProviderNoConfigPath(t *testing.T) {
	testManifest := loadTestManifest(t, "provider_no_config_path.yaml")

	codex, exists := testManifest.Providers["codex"]
	if !exists {
		t.Errorf("codex provider should exist")
		t.FailNow()
	}

	if codex.ConfigPath != "" {
		t.Errorf("ConfigPath should be empty, got %s", codex.ConfigPath)
	}

	if *codex.Enabled != true {
		t.Errorf("Enabled should be true, got %v", codex.Enabled)
	}
}

func TestParseProviderNoEnabled(t *testing.T) {
	testManifest := loadTestManifest(t, "provider_no_enabled.yaml")

	cursor, exists := testManifest.Providers["cursor"]
	if !exists {
		t.Errorf("cursor provider should exist")
		t.FailNow()
	}

	if cursor.ConfigPath != "{home}/.config/cursor/mcp.json" {
		t.Errorf("ConfigPath should be '{home}/.config/cursor/mcp.json', got %s", cursor.ConfigPath)
	}

	if cursor.Enabled != nil {
		t.Errorf("Enabled should default to false, got %v", cursor.Enabled)
	}
}
