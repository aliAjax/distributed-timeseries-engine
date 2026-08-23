# Bug Reproduction

## What happens

Zero-value query and group options panic on their first map write. Typed-nil
policy and validator values also pass ordinary interface nil checks, causing a
nil dereference or retaining an unusable validator.

## How to trigger it

Create zero values for `QueryOptions` and `GroupOptions` and perform their first
write. Pass typed-nil policy and validator pointers through their interface
entry points and invoke the normal policy or validation path.

## Error

The runtime and assertion failures are:

```text
panic: assignment to entry in nil map
panic: runtime error: invalid memory address or nil pointer dereference
typed-nil validator retained: len=1
```
