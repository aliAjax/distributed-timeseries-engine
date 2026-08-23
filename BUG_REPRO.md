# Bug Reproduction

## What happens

A successful shard-migration retry remains in the intermediate `retrying`
state. The active view and audit mapping disagree about the same migration.

## How to trigger it

Start a migration, enter the retry path, and complete the retried operation.
Then inspect the stored state, allowed transition edges, active-migration view,
and audit status for that migration.

## Error

The observed contradictions are:

```text
successful retry state=retrying
retrying to completed edge rejected
retrying migration omitted from active view
retrying audit status=completed
```
