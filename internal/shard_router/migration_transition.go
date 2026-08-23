package shard_router

func RetryMigration(current MigrationState, copySucceeded bool) MigrationState {
	if current != MigrationFailed && current != MigrationRetrying {
		return current
	}
	if copySucceeded {
		return MigrationCompleted
	}
	return MigrationRetrying
}
