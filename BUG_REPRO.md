# Bug Reproduction

## What happens

Repository operations turn storage failures into plain text. Callers can no
longer identify the storage sentinel for open, ingest, query, or seal errors.

## How to trigger it

Run the repository error checks for the four public operation wrappers. Each
check passes a distinct storage sentinel through its operation and then uses
the standard error-chain API to identify it.

## Error

The four observed failures are:

```text
open failure detached from storage sentinel
ingest failure detached from storage sentinel
query failure detached from storage sentinel
seal failure detached from storage sentinel
```
