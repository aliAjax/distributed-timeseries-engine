# Bug Reproduction

## What happens

Retention batches keep resources active until the whole loop returns. Cleanup,
rollback, and close failures can also replace or detach the primary operation
error, so callers lose the original failure identity.

## How to trigger it

Process at least two retention resources and assert that the first is closed
before the second starts. Then inject independent primary, cleanup, rollback,
and close errors and inspect the returned error chain.

## Error

The red run reported:

```text
resource still active before b
combined error lost identity: cleanup failed
rolledBack=true err=rollback failed
close result lost error identity: close failed
```
