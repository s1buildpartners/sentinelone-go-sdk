//go:build integration

package sentinelone_test

import (
	"context"
	"testing"

	s1 "github.com/s1buildpartners/sentinelone-go-sdk"
)

func TestIntegration_Agents_List(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()

	agents, pag, err := cli.Agents.List(ctx, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if pag == nil {
		t.Fatal("expected pagination metadata")
	}

	t.Logf("found %d agent(s) (total=%d)", len(agents), pag.TotalItems)
}

func TestIntegration_Agents_List_WithFilter(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()

	// Retrieve at most one agent per OS type to confirm filters are forwarded.
	params := &s1.ListAgentsParams{
		ListParams: s1.ListParams{Limit: s1.IntPtr(1)},
		OSTypes:    []string{"windows", "linux", "macos"},
	}

	agents, _, err := cli.Agents.List(ctx, params)
	if err != nil {
		t.Fatalf("List(osTypes): %v", err)
	}

	t.Logf("filtered list returned %d agent(s)", len(agents))
}

func TestIntegration_Agents_Count(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()

	total, err := cli.Agents.Count(ctx, nil)
	if err != nil {
		t.Fatalf("Count: %v", err)
	}

	t.Logf("total agents: %d", total)
}

func TestIntegration_Agents_Count_WithFilter(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()

	// Count only active agents; result must be <= the unfiltered total.
	activeParams := &s1.ListAgentsParams{IsActive: s1.BoolPtr(true)}

	active, err := cli.Agents.Count(ctx, activeParams)
	if err != nil {
		t.Fatalf("Count(isActive=true): %v", err)
	}

	total, err := cli.Agents.Count(ctx, nil)
	if err != nil {
		t.Fatalf("Count(all): %v", err)
	}

	if active > total {
		t.Errorf("active count %d exceeds total count %d", active, total)
	}

	t.Logf("active=%d total=%d", active, total)
}
