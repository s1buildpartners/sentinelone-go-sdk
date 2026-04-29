//go:build integration

package sentinelone_test

import (
	"context"
	"net/http"
	"testing"

	s1 "github.com/s1buildpartners/sentinelone-go-sdk"
)

func TestIntegration_Accounts_List(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()

	accounts, pag, err := cli.Accounts.List(ctx, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(accounts) == 0 {
		t.Fatal("expected at least one account")
	}

	if pag == nil {
		t.Fatal("expected pagination metadata")
	}

	t.Logf("found %d account(s); first: %s (%s)", len(accounts), accounts[0].Name, accounts[0].ID)
}

func TestIntegration_Accounts_Get(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()
	accountID := firstAccountID(t, cli)

	account, err := cli.Accounts.Get(ctx, accountID)
	if err != nil {
		t.Fatalf("Get(%s): %v", accountID, err)
	}

	if account.ID != accountID {
		t.Errorf("ID = %q, want %q", account.ID, accountID)
	}

	t.Logf("account: %s (state=%s)", account.Name, account.State)
}

func TestIntegration_Accounts_GetPolicy(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()
	accountID := firstAccountID(t, cli)

	policy, err := cli.Accounts.GetPolicy(ctx, accountID)
	if err != nil {
		t.Fatalf("GetPolicy(%s): %v", accountID, err)
	}

	t.Logf("mitigation mode: %s", policy.MitigationMode)
}

func TestIntegration_Accounts_GetUninstallPasswordMetadata(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()
	accountID := firstAccountID(t, cli)

	meta, err := cli.Accounts.GetUninstallPasswordMetadata(ctx, accountID)
	if err != nil {
		if respErr, ok := s1.AsResponseError(err); ok && respErr.StatusCode == http.StatusForbidden {
			t.Skipf("uninstall password feature not available for account %s", accountID)
		}

		t.Fatalf("GetUninstallPasswordMetadata(%s): %v", accountID, err)
	}

	t.Logf("has uninstall password: %v", meta.HasPassword)
}
