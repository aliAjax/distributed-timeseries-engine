package shard_router

func AuditMigrationState(state MigrationState) string {
	switch state {
	case MigrationQueued, MigrationCopying, MigrationRetrying:
		return "active"
	case MigrationCompleted:
		return "completed"
	case MigrationFailed:
		return "failed"
	default:
		return "unknown"
	}
}
