package shard_router

func AuditMigrationState(state MigrationState) string {
	switch state {
	case MigrationQueued, MigrationCopying:
		return "active"
	case MigrationCompleted, MigrationRetrying:
		return "completed"
	case MigrationFailed:
		return "failed"
	default:
		return "unknown"
	}
}
