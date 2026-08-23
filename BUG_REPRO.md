# Bug Reproduction

## What happens

Quality filtering, window extraction, label filtering, and sample partitioning
reuse their input slice storage. Building or mutating a derived result changes
the caller's original time-series data.

## How to trigger it

Keep a copy of each source slice, call the corresponding derivation function,
then mutate the returned slice and compare the source with its saved value.

## Error

The red run produced these failures:

```text
quality filter mutated source
window extraction corrupted the source point slice
label filter mutated source
partition aliases source
```
