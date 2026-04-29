//go:build integration

package sentinelone_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	s1 "github.com/s1buildpartners/sentinelone-go-sdk"
)

// TestMain loads a .env file from the current directory before running any
// integration tests.  Variables already set in the process environment are
// never overwritten, so real env vars always take priority over the file.
func TestMain(m *testing.M) {
	loadDotEnv()
	os.Exit(m.Run())
}

// loadDotEnv reads ".env" from the working directory and exports any
// variables it contains that are not already present in the environment.
//
// Supported syntax:
//
//	KEY=value          plain value
//	KEY="value"        double-quoted value
//	KEY='value'        single-quoted value
//	export KEY=value   optional export prefix
//	# comment          ignored
//	                   blank lines ignored
func loadDotEnv() {
	data, err := os.ReadFile(".env")
	if errors.Is(err, os.ErrNotExist) {
		return
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: read .env: %v\n", err)
		return
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip optional "export " prefix used in shell-compatible .env files.
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))

		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			continue
		}

		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])

		// Strip a matching pair of surrounding quotes.
		if len(val) >= 2 {
			first, last := val[0], val[len(val)-1]
			if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		// Existing env vars take priority over the file.
		if _, alreadySet := os.LookupEnv(key); !alreadySet {
			_ = os.Setenv(key, val)
		}
	}
}

// uniqueName returns a resource name unlikely to collide across test runs.
// The "go-sdk-inttest-" prefix makes integration-test objects easy to
// identify in the console and safe to bulk-delete if cleanup fails.
func uniqueName() string {
	return fmt.Sprintf("go-sdk-inttest-%d", time.Now().UnixNano())
}

// integrationClient returns a Client configured via SENTINELONE_URL and
// SENTINELONE_TOKEN.  The test is skipped when either variable is unset.
func integrationClient(t *testing.T) *s1.Client {
	t.Helper()

	cli, err := s1.NewClientFromEnv()
	if err != nil {
		t.Skip("skipping integration test (credentials not configured):", err)
	}

	return cli
}

// firstAccountID lists accounts and returns the ID of the first result.
// The test is skipped when no accounts are accessible.
func firstAccountID(t *testing.T, cli *s1.Client) string {
	t.Helper()

	accounts, _, err := cli.Accounts.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("list accounts: %v", err)
	}

	if len(accounts) == 0 {
		t.Skip("no accounts accessible; skipping")
	}

	return accounts[0].ID
}
