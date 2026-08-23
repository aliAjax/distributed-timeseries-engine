# Bug Reproduction

## What happens

Four cold-start components panic when they are embedded as zero values and
receive their first mutation. Their regular constructors work, but the zero
value paths leave internal maps nil.

## How to trigger it

Create zero values for the following types and invoke the listed first-write
operation:

- `label_index.Index.Add`
- `quota.Limiter.Set`
- `quota.CardinalityGuard.Observe`
- `shard_router.LeaseTable.Acquire`

Each operation reaches a map assignment before the receiver has initialized
its internal map.

## Error

All four paths fail with the same runtime error:

```text
panic: assignment to entry in nil map [recovered]
    panic: assignment to entry in nil map
```

The observed production frames are `internal/label_index/index.go:22`,
`internal/quota/quota.go:27`, `internal/quota/cardinality.go:27`, and
`internal/shard_router/lease.go:35`.
