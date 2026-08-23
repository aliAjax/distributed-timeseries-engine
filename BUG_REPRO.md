# Bug Reproduction

## What happens

HTTP request contexts lose cancellation, deadline, and tenant values at bridge
boundaries. A cached context can also carry the first request's value into the
next request.

## How to trigger it

Cancel a parent request context before dispatch, pass a deadline and tenant
value through the bridge, then issue a second request with a different value.
Observe what reaches the handler in each case.

## Error

The red run reported:

```text
err=<nil> called=true
deadline=0001-01-01 00:00:00 +0000 UTC ok=false
second request inherited value=first
tenant value=<nil>
```
