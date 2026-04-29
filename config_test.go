package sentinelone

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTestConfig creates a temp file with the given content and returns its path.
func writeTestConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "s1cfg*.ini")
	if err != nil {
		t.Fatalf("create temp config: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	f.Close()
	return f.Name()
}

const multiProfileConfig = `# comment line
; also a comment

[default]
url   = https://default.sentinelone.net
token = default-token

[production]
url   = https://prod.sentinelone.net
token = prod-token

[staging]
url   = https://staging.sentinelone.net
token = staging-token
`

// --- NewClientFromEnv ---

func TestNewClientFromEnv_Success(t *testing.T) {
	t.Setenv(EnvURL, "https://test.sentinelone.net")
	t.Setenv(EnvToken, "test-token")

	cli, err := NewClientFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://test.sentinelone.net" {
		t.Errorf("baseURL = %q", cli.baseURL)
	}
	if cli.apiToken != "test-token" {
		t.Errorf("apiToken = %q", cli.apiToken)
	}
}

func TestNewClientFromEnv_MissingURL(t *testing.T) {
	t.Setenv(EnvURL, "")
	t.Setenv(EnvToken, "test-token")

	_, err := NewClientFromEnv()
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
	if !strings.Contains(err.Error(), EnvURL) {
		t.Errorf("error should mention %s: %v", EnvURL, err)
	}
}

func TestNewClientFromEnv_MissingToken(t *testing.T) {
	t.Setenv(EnvURL, "https://test.sentinelone.net")
	t.Setenv(EnvToken, "")

	_, err := NewClientFromEnv()
	if err == nil {
		t.Fatal("expected error for missing token")
	}
	if !strings.Contains(err.Error(), EnvToken) {
		t.Errorf("error should mention %s: %v", EnvToken, err)
	}
}

func TestNewClientFromEnv_WithOptions(t *testing.T) {
	t.Setenv(EnvURL, "https://test.sentinelone.net/")
	t.Setenv(EnvToken, "tok")

	cli, err := NewClientFromEnv(WithRateLimiting(false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.rateLimits != nil {
		t.Error("expected rate limiting to be disabled")
	}
	// trailing slash should be stripped by NewClient
	if cli.baseURL != "https://test.sentinelone.net" {
		t.Errorf("baseURL = %q", cli.baseURL)
	}
}

// --- NewClientFromConfig ---

func TestNewClientFromConfig_DefaultProfile(t *testing.T) {
	t.Setenv(EnvConfig, writeTestConfig(t, multiProfileConfig))

	cli, err := NewClientFromConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://default.sentinelone.net" {
		t.Errorf("baseURL = %q", cli.baseURL)
	}
	if cli.apiToken != "default-token" {
		t.Errorf("apiToken = %q", cli.apiToken)
	}
}

func TestNewClientFromConfig_NamedProfile(t *testing.T) {
	t.Setenv(EnvConfig, writeTestConfig(t, multiProfileConfig))

	cli, err := NewClientFromConfig(WithProfile("production"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://prod.sentinelone.net" {
		t.Errorf("baseURL = %q", cli.baseURL)
	}
	if cli.apiToken != "prod-token" {
		t.Errorf("apiToken = %q", cli.apiToken)
	}
}

func TestNewClientFromConfig_ProfileNotFound(t *testing.T) {
	t.Setenv(EnvConfig, writeTestConfig(t, multiProfileConfig))

	_, err := NewClientFromConfig(WithProfile("nonexistent"))
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention profile name: %v", err)
	}
}

func TestNewClientFromConfig_MissingURL(t *testing.T) {
	t.Setenv(EnvConfig, writeTestConfig(t, "[default]\ntoken = tok\n"))

	_, err := NewClientFromConfig()
	if err == nil {
		t.Fatal("expected error for missing url")
	}
	if !strings.Contains(err.Error(), configKeyURL) {
		t.Errorf("error should mention 'url': %v", err)
	}
}

func TestNewClientFromConfig_MissingToken(t *testing.T) {
	t.Setenv(EnvConfig, writeTestConfig(t, "[default]\nurl = https://test.sentinelone.net\n"))

	_, err := NewClientFromConfig()
	if err == nil {
		t.Fatal("expected error for missing token")
	}
	if !strings.Contains(err.Error(), configKeyToken) {
		t.Errorf("error should mention 'token': %v", err)
	}
}

func TestNewClientFromConfig_FileNotFound(t *testing.T) {
	t.Setenv(EnvConfig, "/nonexistent/path/credentials")

	_, err := NewClientFromConfig()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "open config file") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewClientFromConfig_WithOptions(t *testing.T) {
	t.Setenv(EnvConfig, writeTestConfig(t, multiProfileConfig))

	cli, err := NewClientFromConfig(WithProfile("staging"), WithRateLimiting(false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.rateLimits != nil {
		t.Error("expected rate limiting to be disabled")
	}
	if cli.baseURL != "https://staging.sentinelone.net" {
		t.Errorf("baseURL = %q", cli.baseURL)
	}
}

func TestNewClientFromConfig_WithConfigFile(t *testing.T) {
	// SENTINELONE_CONFIG is not set; path comes exclusively from WithConfigFile.
	t.Setenv(EnvConfig, "")

	path := writeTestConfig(t, multiProfileConfig)

	cli, err := NewClientFromConfig(WithConfigFile(path))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://default.sentinelone.net" {
		t.Errorf("baseURL = %q", cli.baseURL)
	}
	if cli.apiToken != "default-token" {
		t.Errorf("apiToken = %q", cli.apiToken)
	}
}

func TestNewClientFromConfig_WithConfigFileAndProfile(t *testing.T) {
	t.Setenv(EnvConfig, "")

	path := writeTestConfig(t, multiProfileConfig)

	cli, err := NewClientFromConfig(WithConfigFile(path), WithProfile("production"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://prod.sentinelone.net" {
		t.Errorf("baseURL = %q", cli.baseURL)
	}
	if cli.apiToken != "prod-token" {
		t.Errorf("apiToken = %q", cli.apiToken)
	}
}

// --- NewDefaultClient ---

func TestNewDefaultClient_EnvVarsPriority(t *testing.T) {
	t.Setenv(EnvConfig, writeTestConfig(t, multiProfileConfig))
	t.Setenv(EnvURL, "https://env.sentinelone.net")
	t.Setenv(EnvToken, "env-token")

	cli, err := NewDefaultClient(WithProfile("production"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://env.sentinelone.net" {
		t.Errorf("baseURL = %q, want env value", cli.baseURL)
	}
	if cli.apiToken != "env-token" {
		t.Errorf("apiToken = %q, want env value", cli.apiToken)
	}
}

func TestNewDefaultClient_FallbackToConfig(t *testing.T) {
	t.Setenv(EnvConfig, writeTestConfig(t, multiProfileConfig))
	t.Setenv(EnvURL, "")
	t.Setenv(EnvToken, "")

	cli, err := NewDefaultClient(WithProfile("staging"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://staging.sentinelone.net" {
		t.Errorf("baseURL = %q", cli.baseURL)
	}
	if cli.apiToken != "staging-token" {
		t.Errorf("apiToken = %q", cli.apiToken)
	}
}

func TestNewDefaultClient_EmptyProfileUsesEnvProfile(t *testing.T) {
	t.Setenv(EnvConfig, writeTestConfig(t, multiProfileConfig))
	t.Setenv(EnvURL, "")
	t.Setenv(EnvToken, "")
	t.Setenv(EnvProfile, "production")

	cli, err := NewDefaultClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://prod.sentinelone.net" {
		t.Errorf("baseURL = %q", cli.baseURL)
	}
}

func TestNewDefaultClient_EmptyProfileFallsBackToDefault(t *testing.T) {
	t.Setenv(EnvConfig, writeTestConfig(t, multiProfileConfig))
	t.Setenv(EnvURL, "")
	t.Setenv(EnvToken, "")
	t.Setenv(EnvProfile, "")

	cli, err := NewDefaultClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://default.sentinelone.net" {
		t.Errorf("baseURL = %q", cli.baseURL)
	}
}

func TestNewDefaultClient_OnlyURLSet(t *testing.T) {
	// Only URL set → env condition not met → falls back to config file.
	t.Setenv(EnvConfig, writeTestConfig(t, multiProfileConfig))
	t.Setenv(EnvURL, "https://env.sentinelone.net")
	t.Setenv(EnvToken, "")
	t.Setenv(EnvProfile, "")

	cli, err := NewDefaultClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://default.sentinelone.net" {
		t.Errorf("baseURL = %q, want config file default", cli.baseURL)
	}
}

func TestNewDefaultClient_OnlyTokenSet(t *testing.T) {
	// Only token set → env condition not met → falls back to config file.
	t.Setenv(EnvConfig, writeTestConfig(t, multiProfileConfig))
	t.Setenv(EnvURL, "")
	t.Setenv(EnvToken, "env-token")
	t.Setenv(EnvProfile, "")

	cli, err := NewDefaultClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://default.sentinelone.net" {
		t.Errorf("baseURL = %q, want config file default", cli.baseURL)
	}
}

func TestNewDefaultClient_WithConfigFile(t *testing.T) {
	t.Setenv(EnvConfig, "")
	t.Setenv(EnvURL, "")
	t.Setenv(EnvToken, "")

	path := writeTestConfig(t, multiProfileConfig)

	cli, err := NewDefaultClient(WithConfigFile(path), WithProfile("staging"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli.baseURL != "https://staging.sentinelone.net" {
		t.Errorf("baseURL = %q", cli.baseURL)
	}
}

func TestNewDefaultClient_ConfigError(t *testing.T) {
	t.Setenv(EnvConfig, "")
	t.Setenv(EnvURL, "")
	t.Setenv(EnvToken, "")

	orig := userConfigDirFn
	userConfigDirFn = func() (string, error) { return "", errors.New("no home directory") }
	defer func() { userConfigDirFn = orig }()

	_, err := NewDefaultClient()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "resolve config directory") {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- parseConfigFile ---

func TestParseConfigFile_MultipleProfiles(t *testing.T) {
	path := writeTestConfig(t, multiProfileConfig)

	profiles, err := parseConfigFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 3 {
		t.Errorf("expected 3 profiles, got %d", len(profiles))
	}
	for _, name := range []string{"default", "production", "staging"} {
		if _, ok := profiles[name]; !ok {
			t.Errorf("missing profile %q", name)
		}
	}
}

func TestParseConfigFile_ColonSeparator(t *testing.T) {
	path := writeTestConfig(t, "[default]\nurl: https://test.sentinelone.net\ntoken: tok\n")

	profiles, err := parseConfigFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := profiles["default"]
	if p.URL != "https://test.sentinelone.net" {
		t.Errorf("URL = %q", p.URL)
	}
	if p.Token != "tok" {
		t.Errorf("Token = %q", p.Token)
	}
}

func TestParseConfigFile_LinesBeforeSection(t *testing.T) {
	path := writeTestConfig(t, "orphan = value\n[default]\nurl = https://test.sentinelone.net\ntoken = tok\n")

	profiles, err := parseConfigFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Errorf("expected 1 profile (orphan ignored), got %d", len(profiles))
	}
}

func TestParseConfigFile_NoSeparator(t *testing.T) {
	path := writeTestConfig(t, "[default]\nbadline\nurl = https://test.sentinelone.net\ntoken = tok\n")

	profiles, err := parseConfigFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profiles["default"].URL != "https://test.sentinelone.net" {
		t.Errorf("URL = %q", profiles["default"].URL)
	}
}

func TestParseConfigFile_UnknownKey(t *testing.T) {
	path := writeTestConfig(t, "[default]\nurl = https://test.sentinelone.net\ntoken = tok\nunknown = value\n")

	profiles, err := parseConfigFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Errorf("expected 1 profile, got %d", len(profiles))
	}
}

func TestParseConfigFile_EmptyLines(t *testing.T) {
	path := writeTestConfig(t, "\n\n[default]\n\nurl = https://test.sentinelone.net\n\ntoken = tok\n\n")

	profiles, err := parseConfigFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := profiles["default"]
	if p.URL != "https://test.sentinelone.net" || p.Token != "tok" {
		t.Errorf("unexpected profile values: %+v", p)
	}
}

func TestParseConfigFile_FileNotFound(t *testing.T) {
	_, err := parseConfigFile("/nonexistent/path/credentials")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "open config file") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseConfigFile_ScannerError(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "s1cfg*.ini")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	// Write a line longer than bufio.MaxScanTokenSize to trigger a scanner error.
	_, _ = f.WriteString("[default]\nurl = " + strings.Repeat("x", bufio.MaxScanTokenSize+1) + "\n")
	_ = f.Close()

	_, err = parseConfigFile(f.Name())
	if err == nil {
		t.Fatal("expected scanner error")
	}
	if !strings.Contains(err.Error(), "read config file") {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- configFilePath ---

func TestConfigFilePath_EnvOverride(t *testing.T) {
	t.Setenv(EnvConfig, "/custom/path/credentials")

	path, err := configFilePath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/custom/path/credentials" {
		t.Errorf("path = %q", path)
	}
}

func TestConfigFilePath_Default(t *testing.T) {
	t.Setenv(EnvConfig, "")

	path, err := configFilePath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(path, "sentinelone") {
		t.Errorf("path should contain 'sentinelone': %q", path)
	}
	if !strings.HasSuffix(path, "credentials") {
		t.Errorf("path should end with 'credentials': %q", path)
	}
}

// --- defaultConfigPath ---

func TestDefaultConfigPath_Success(t *testing.T) {
	path, err := defaultConfigPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join("sentinelone", defaultConfigFileName)
	if !strings.HasSuffix(path, want) {
		t.Errorf("path = %q, want suffix %q", path, want)
	}
}

func TestDefaultConfigPath_UserConfigDirError(t *testing.T) {
	// Replace the injected function to simulate os.UserConfigDir failure.
	orig := userConfigDirFn
	userConfigDirFn = func() (string, error) { return "", errors.New("no home directory") }
	defer func() { userConfigDirFn = orig }()

	_, err := defaultConfigPath()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "resolve config directory") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewClientFromConfig_ConfigFilePathError(t *testing.T) {
	// Clear the env override so defaultConfigPath is called, then make it fail.
	t.Setenv(EnvConfig, "")

	orig := userConfigDirFn
	userConfigDirFn = func() (string, error) { return "", errors.New("no home directory") }
	defer func() { userConfigDirFn = orig }()

	_, err := NewClientFromConfig()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "resolve config directory") {
		t.Errorf("unexpected error: %v", err)
	}
}
