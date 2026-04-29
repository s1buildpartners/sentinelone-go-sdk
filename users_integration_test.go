//go:build integration

package sentinelone_test

import (
	"context"
	"testing"
)

func TestIntegration_Users_List(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()

	users, pag, err := cli.Users.List(ctx, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(users) == 0 {
		t.Fatal("expected at least one user")
	}

	if pag == nil {
		t.Fatal("expected pagination metadata")
	}

	t.Logf("found %d user(s); first ID: %s", len(users), users[0].ID)
}

func TestIntegration_Users_Get(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()

	users, _, err := cli.Users.List(ctx, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(users) == 0 {
		t.Skip("no users accessible; skipping")
	}

	userID := users[0].ID

	user, err := cli.Users.Get(ctx, userID)
	if err != nil {
		t.Fatalf("Get(%s): %v", userID, err)
	}

	if user.ID != userID {
		t.Errorf("ID = %q, want %q", user.ID, userID)
	}

	if user.Email != nil {
		t.Logf("user: %s (%s)", *user.Email, user.Scope)
	}
}

func TestIntegration_Users_GetAPITokenDetails(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()

	users, _, err := cli.Users.List(ctx, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(users) == 0 {
		t.Skip("no users accessible; skipping")
	}

	detail, err := cli.Users.GetAPITokenDetails(ctx, users[0].ID)
	if err != nil {
		// Not all accounts have API token management enabled; treat as non-fatal.
		t.Logf("GetAPITokenDetails: %v (may not be enabled for this account)", err)
		return
	}

	t.Logf("token expires: %s", detail.ExpiresAt)
}
