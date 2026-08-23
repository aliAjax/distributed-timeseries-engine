# Bug Reproduction

## What happens

Recovery errors retain their human-readable text, but the wrapped error identity is lost across chunk decoding, WAL replay, block sealing, and repository recovery. Callers using `errors.Is` or `errors.As` therefore cannot recognize the outer recovery category and the underlying storage failure together.

## How to trigger it

Use the recovery scenarios for a short chunk, a corrupt WAL record, an empty block seal, and a repository recovery failure. The unmodified behavior fails all four identity checks; after the fix, each scenario preserves both the wrapper type and its sentinel errors.

## Observed failure messages

- `decode identity lost: decode chunk: short chunk`
- `replay identity lost: replay: corrupt wal record`
- `block identity lost: block write: empty block`
- `recovery identity lost: recovery: replay: corrupt wal record`
