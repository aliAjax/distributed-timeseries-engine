# Bug Reproduction

## What happens

Concurrent block-store updates and snapshots race on internal maps. Returned
snapshots also retain aliases to internal maps and slices, so later mutations
change data that callers already captured.

## How to trigger it

Run concurrent writers and snapshot readers against the catalog, quarantine,
series, and health views with the Go race detector enabled. Mutate a returned
snapshot after each concurrent run and inspect the source view.

## Error

All four paths report a data race, followed by aliasing failures such as:

```text
WARNING: DATA RACE
catalog snapshot aliases internal entries
quarantine snapshot aliases internal reasons
series snapshot aliases internal slices
health snapshot aliases internal status
```
