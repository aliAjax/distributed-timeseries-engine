package shard_router

func CanMoveMigration(from, to MigrationState) bool {
	edges := map[MigrationState]map[MigrationState]bool{
		MigrationQueued:    {MigrationCopying: true, MigrationFailed: true},
		MigrationCopying:   {MigrationCompleted: true, MigrationFailed: true},
		MigrationFailed:    {MigrationRetrying: true},
		MigrationRetrying:  {MigrationFailed: true},
		MigrationCompleted: {MigrationCopying: true},
	}
	return edges[from][to]
}
