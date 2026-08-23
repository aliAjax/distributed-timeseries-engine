# Bug Reproduction

## What happens

Replica result coordination can finish before workers start, leave result and
error senders blocked, and keep consumers alive after the result stream ends.

## How to trigger it

Exercise producer completion, worker waiting, cancellation during error
delivery, and consumer shutdown under the race detector. Bound every case with
a timeout so leaked goroutines and missing closes are observable.

## Error

The four lifecycle checks fail with:

```text
replica result channel was not closed
coordinator completed before workers started
canceled error sender leaked
consumer did not exit after channel close
```
