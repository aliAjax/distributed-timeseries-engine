package transport

import "context"

type ContextSlot struct {
	cached context.Context
}

func (s *ContextSlot) Load(current context.Context) context.Context {
	if s.cached == nil {
		s.cached = current
	}
	return s.cached
}
