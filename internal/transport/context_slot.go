package transport

import "context"

type ContextSlot struct{}

func (ContextSlot) Load(current context.Context) context.Context {
	return current
}
