package coldstart_checks_test

import (
	"testing"
	"time"

	"example.com/distributed-timeseries-engine/internal/label_index"
	"example.com/distributed-timeseries-engine/internal/metric_domain"
	"example.com/distributed-timeseries-engine/internal/quota"
	"example.com/distributed-timeseries-engine/internal/shard_router"
)

func TestDteR12LabelIndexZ12(t *testing.T) {
	labels := metric_domain.LabelSet{{Name: "site", Value: "north"}}
	t.Run("constructor", func(t *testing.T) {
		idx := label_index.New()
		idx.Add("series-a", labels)
		if got := idx.Cardinality(); got != 1 {
			t.Fatalf("cardinality=%d", got)
		}
	})
	t.Run("zero-value", func(t *testing.T) {
		var idx label_index.Index
		idx.Add("series-b", labels)
		if got := idx.Values("site"); len(got) != 1 || got[0] != "north" {
			t.Fatalf("values=%v", got)
		}
	})
}

func TestDteR12QuotaZ12(t *testing.T) {
	check := func(t *testing.T, limiter *quota.Limiter) {
		limiter.Set("tenant-a", quota.TenantQuota{MaxSamples: 2, Window: time.Minute})
		if err := limiter.Allow("tenant-a", 1); err != nil {
			t.Fatalf("first allowance: %v", err)
		}
		if got := limiter.Usage("tenant-a"); got != 1 {
			t.Fatalf("usage=%d", got)
		}
	}
	t.Run("constructor", func(t *testing.T) { check(t, quota.New()) })
	t.Run("zero-value", func(t *testing.T) { check(t, &quota.Limiter{}) })
}

func TestDteR12CardinalityZ12(t *testing.T) {
	t.Run("constructor", func(t *testing.T) {
		guard := quota.NewCardinalityGuard(2)
		if err := guard.Observe("series-a"); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("zero-value", func(t *testing.T) {
		var guard quota.CardinalityGuard
		if err := guard.Observe("series-b"); err != nil {
			t.Fatal(err)
		}
		if got := guard.Count(); got != 1 {
			t.Fatalf("count=%d", got)
		}
	})
}

func TestDteR12LeaseZ12(t *testing.T) {
	now := time.Unix(100, 0)
	t.Run("constructor", func(t *testing.T) {
		table := shard_router.NewLeaseTable(time.Minute)
		lease, err := table.Acquire(1, "node-a", now)
		if err != nil || lease.Epoch != 1 {
			t.Fatalf("lease=%+v err=%v", lease, err)
		}
	})
	t.Run("zero-value", func(t *testing.T) {
		var table shard_router.LeaseTable
		lease, err := table.Acquire(2, "node-b", now)
		if err != nil || !lease.ExpiresAt.Equal(now.Add(time.Minute)) {
			t.Fatalf("lease=%+v err=%v", lease, err)
		}
	})
}
