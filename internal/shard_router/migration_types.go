package shard_router

type MigrationState string

const (
	MigrationQueued    MigrationState = "queued"
	MigrationCopying   MigrationState = "copying"
	MigrationRetrying  MigrationState = "retrying"
	MigrationCompleted MigrationState = "completed"
	MigrationFailed    MigrationState = "failed"
)
