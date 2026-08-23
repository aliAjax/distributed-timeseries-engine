# Bug Reproduction

## What happens

Canceled or expired contexts do not stop WAL replay, append, scan, and copy
operations. Work continues across the cancellation boundary.

## How to trigger it

Pass an already canceled context to replay and short-deadline contexts to the
remaining WAL operations. Observe whether callbacks, visitors, and writes are
still invoked after cancellation.

## Error

The red run reported continued work with no returned error:

```text
Replay: err=<nil> called=1
Append: err=<nil>
Scan: err=<nil> visited=1
Copy: err=<nil> written=7 dst="payload"
```
