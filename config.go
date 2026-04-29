package sentinelone

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// EnvURL is the environment variable name for the management console base URL.
	EnvURL = "SENTINELONE_URL"
	// EnvToken is the environment variable name for the API token.
	EnvToken = "SENTINELONE_TOKEN"
	// EnvConfig is the environment variable name for an explicit config file path.
	EnvConfig = "SENTINELONE_CONFIG"
	// EnvProfile is the environment variable name for the default profile name.
	EnvProfile = "SENTINELONE_PROFILE"

	defaultProfileName    = "default"
	defaultConfigFileName = "credentials"

	configKeyURL   = "url"
	configKeyToken = "token"
)

// Sentinel errors returned by credential-loading functions.
var (
	// ErrEnvVarMissing is returned when a required environment variable is not set.
	ErrEnvVarMissing = errors.New("environment variable is not set")
	// ErrProfileNotFound is returned when the requested profile does not exist in the credentials file.
	ErrProfileNotFound = errors.New("profile not found")
	// ErrMissingCredential is returned when a profile is present but is missing a required field.
	ErrMissingCredential = errors.New("missing credential")
)

// userConfigDirFn wraps os.UserConfigDir so tests can substitute an alternative.
//
//nolint:gochecknoglobals
var userConfigDirFn = os.UserConfigDir

// LoadOption is a sealed option accepted by [NewClientFromConfig] and
// [NewClientFromProfile].  Both [ClientOption] and [ConfigOption] satisfy it.
//
// The unexported marker method prevents external types from satisfying the
// interface, keeping the option set closed to this package.
type LoadOption interface {
	applyLoadOption()
}

// credentialConfig collects the credential-loading settings extracted from
// [LoadOption] values by [applyLoadOptions].
type credentialConfig struct {
	profile    string
	configFile string
}

// ConfigOption configures credential-loading behaviour such as the profile
// name and credentials file path.  It satisfies [LoadOption] and can be
// mixed with [ClientOption] values in [NewClientFromConfig] and
// [NewClientFromProfile].
type ConfigOption func(*credentialConfig)

func (ConfigOption) applyLoadOption() {}

// WithProfile sets the profile name to load from the credentials file.
// If not specified, [EnvProfile] (SENTINELONE_PROFILE) is checked, then
// "default" is used as a final fallback.
func WithProfile(name string) ConfigOption {
	return func(c *credentialConfig) { c.profile = name }
}

// WithConfigFile sets an explicit path to the credentials file.  When set,
// the [EnvConfig] environment variable and the default platform path are both
// bypassed.
func WithConfigFile(path string) ConfigOption {
	return func(c *credentialConfig) { c.configFile = path }
}

// Profile holds the base URL and API token for a named SentinelOne management
// tenant.  Values are populated by [NewClientFromConfig] or [NewClientFromProfile].
type Profile struct {
	URL   string
	Token string
}

// applyLoadOptions partitions a mixed slice of [LoadOption] values into
// credential-loading settings and [ClientOption] values.
func applyLoadOptions(opts []LoadOption) (credentialConfig, []ClientOption) {
	var cfg credentialConfig

	var clientOpts []ClientOption

	for _, opt := range opts {
		switch typedOpt := opt.(type) {
		case ConfigOption:
			typedOpt(&cfg)
		case ClientOption:
			clientOpts = append(clientOpts, typedOpt)
		}
	}

	return cfg, clientOpts
}

// NewClientFromEnv creates a Client using credentials from environment variables.
//
// [EnvURL] (SENTINELONE_URL) must be set to the management console base URL
// and [EnvToken] (SENTINELONE_TOKEN) must be set to a valid API token.
// Both variables must be present; an error is returned if either is empty.
func NewClientFromEnv(opts ...ClientOption) (*Client, error) {
	rawURL := os.Getenv(EnvURL)

	if rawURL == "" {
		return nil, fmt.Errorf("%s %s: %w", errTag, EnvURL, ErrEnvVarMissing)
	}

	token := os.Getenv(EnvToken)

	if token == "" {
		return nil, fmt.Errorf("%s %s: %w", errTag, EnvToken, ErrEnvVarMissing)
	}

	return NewClient(rawURL, token, opts...), nil
}

// NewClientFromConfig creates a Client by loading credentials from the
// credentials file under a named profile.
//
// Options:
//   - [WithProfile]: profile name to load; defaults to "default" when omitted.
//   - [WithConfigFile]: explicit path to the credentials file; overrides
//     [EnvConfig] and the default platform path.
//   - Any [ClientOption] (e.g. [WithRateLimiting], [WithTimeout]).
//
// The default file path is platform-specific:
//   - Linux/BSD: $XDG_CONFIG_HOME/sentinelone/credentials (or ~/.config/…)
//   - macOS:     ~/Library/Application Support/sentinelone/credentials
//   - Windows:   %AppData%\SentinelOne\credentials
//
// The file uses an INI-style format:
//
//	# lines beginning with '#' or ';' are comments
//	[default]
//	url   = https://tenant.sentinelone.net
//	token = your-api-token
//
//	[production]
//	url   = https://prod.sentinelone.net
//	token = prod-api-token
//
// Both '=' and ':' are accepted as key-value separators.
func NewClientFromConfig(opts ...LoadOption) (*Client, error) {
	cfg, clientOpts := applyLoadOptions(opts)

	prof, err := loadProfile(cfg.profile, cfg.configFile)
	if err != nil {
		return nil, err
	}

	return NewClient(prof.URL, prof.Token, clientOpts...), nil
}

// NewClientFromProfile creates a Client using a layered credential lookup.
//
// Priority order:
//  1. [EnvURL] and [EnvToken] environment variables — if both are set they
//     are used directly and the config file is not read.
//  2. Config file — credentials are loaded for the requested profile.
//     If [WithProfile] is not provided, [EnvProfile] (SENTINELONE_PROFILE) is
//     checked first; "default" is used as a final fallback.
//
// Options:
//   - [WithProfile]: profile name to load from the credentials file.
//   - [WithConfigFile]: explicit path to the credentials file.
//   - Any [ClientOption] (e.g. [WithRateLimiting], [WithTimeout]).
//
// This is the recommended constructor for applications that want to support
// both environment-variable-based (CI/containers) and file-based (developer
// workstation) credential management without code changes.
func NewClientFromProfile(opts ...LoadOption) (*Client, error) {
	rawURL := os.Getenv(EnvURL)
	token := os.Getenv(EnvToken)

	cfg, clientOpts := applyLoadOptions(opts)

	if rawURL != "" && token != "" {
		return NewClient(rawURL, token, clientOpts...), nil
	}

	if cfg.profile == "" {
		if p := os.Getenv(EnvProfile); p != "" {
			cfg.profile = p
		}
	}

	prof, err := loadProfile(cfg.profile, cfg.configFile)
	if err != nil {
		return nil, err
	}

	return NewClient(prof.URL, prof.Token, clientOpts...), nil
}

// defaultConfigPath returns the OS-appropriate default credentials file path.
func defaultConfigPath() (string, error) {
	dir, err := userConfigDirFn()
	if err != nil {
		return "", fmt.Errorf("%s resolve config directory: %w", errTag, err)
	}

	return filepath.Join(dir, "sentinelone", defaultConfigFileName), nil
}

// configFilePath returns the path to the credentials file, honouring
// the SENTINELONE_CONFIG environment variable override.
func configFilePath() (string, error) {
	if p := os.Getenv(EnvConfig); p != "" {
		return p, nil
	}

	return defaultConfigPath()
}

// loadProfile reads the named profile from the credentials file and validates
// that both URL and token are present.  If configFile is non-empty it is used
// directly; otherwise [configFilePath] is called to determine the path.
func loadProfile(name, configFile string) (Profile, error) {
	if name == "" {
		name = defaultProfileName
	}

	var path string

	if configFile != "" {
		path = configFile
	} else {
		filePath, err := configFilePath()
		if err != nil {
			return Profile{}, err
		}

		path = filePath
	}

	profiles, err := parseConfigFile(path)
	if err != nil {
		return Profile{}, err
	}

	prof, ok := profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("%s profile %q in %s: %w", errTag, name, path, ErrProfileNotFound)
	}

	if prof.URL == "" {
		return Profile{}, fmt.Errorf("%s profile %q missing url in %s: %w", errTag, name, path, ErrMissingCredential)
	}

	if prof.Token == "" {
		return Profile{}, fmt.Errorf("%s profile %q missing token in %s: %w", errTag, name, path, ErrMissingCredential)
	}

	return prof, nil
}

// parseConfigFile parses an INI-style credentials file into a map of profile
// names to [Profile] values.
//
// Format rules:
//   - Lines starting with '#' or ';' are comments and are ignored.
//   - A section header "[profile-name]" begins a new profile.
//   - Key-value pairs use '=' or ':' as separators; keys are lower-cased.
//   - Lines before the first section header are ignored.
//   - Unrecognised keys are silently ignored.
func parseConfigFile(path string) (map[string]Profile, error) {
	file, err := os.Open(path) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("%s open config file: %w", errTag, err)
	}

	defer func() { _ = file.Close() }()

	profiles := make(map[string]Profile)

	var current string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || line[0] == '#' || line[0] == ';' {
			continue
		}

		if line[0] == '[' && line[len(line)-1] == ']' {
			current = line[1 : len(line)-1]
			if _, exists := profiles[current]; !exists {
				profiles[current] = Profile{}
			}

			continue
		}

		if current == "" {
			continue
		}

		sep := strings.IndexAny(line, "=:")
		if sep < 0 {
			continue
		}

		key := strings.TrimSpace(line[:sep])
		val := strings.TrimSpace(line[sep+1:])
		prof := profiles[current]

		switch strings.ToLower(key) {
		case configKeyURL:
			prof.URL = val
		case configKeyToken:
			prof.Token = val
		}

		profiles[current] = prof
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		return nil, fmt.Errorf("%s read config file: %w", errTag, scanErr)
	}

	return profiles, nil
}
