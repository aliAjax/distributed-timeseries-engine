package shard_router

func ActiveMigrationStates() map[MigrationState]bool {
	return map[MigrationState]bool{
		MigrationQueued:    true,
		MigrationCopying:   true,
		MigrationCompleted: true,
	}
}
