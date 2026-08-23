package shard_router

import "testing"

func TestDteR07RetryZ07(t *testing.T) {
	if state := RetryMigration(MigrationFailed, true); state != MigrationCompleted {
		t.Fatalf("successful retry state=%s", state)
	}
	if state := RetryMigration(MigrationFailed, false); state != MigrationRetrying {
		t.Fatalf("failed retry state=%s", state)
	}
}

func TestDteR07EdgesZ07(t *testing.T) {
	if !CanMoveMigration(MigrationRetrying, MigrationCompleted) {
		t.Fatal("retrying to completed edge rejected")
	}
	if CanMoveMigration(MigrationCompleted, MigrationCopying) {
		t.Fatal("completed migration allowed to resume copying")
	}
}

func TestDteR07ViewZ07(t *testing.T) {
	active := ActiveMigrationStates()
	if !active[MigrationRetrying] {
		t.Fatal("retrying migration omitted from active view")
	}
	if active[MigrationCompleted] {
		t.Fatal("completed migration included in active view")
	}
}

func TestDteR07AuditZ07(t *testing.T) {
	if status := AuditMigrationState(MigrationRetrying); status != "active" {
		t.Fatalf("retrying audit status=%s", status)
	}
	if status := AuditMigrationState(MigrationCompleted); status != "completed" {
		t.Fatalf("completed audit status=%s", status)
	}
}
