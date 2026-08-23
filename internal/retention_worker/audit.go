package retention_worker

import (
	"fmt"
	"sync"
	"time"
)

type AuditEvent struct {
	Time                    time.Time
	Action, BlockID, Reason string
}
type AuditLog struct {
	mu     sync.Mutex
	events []AuditEvent
	max    int
}

func NewAuditLog(max int) *AuditLog {
	if max < 1 {
		max = 1000
	}
	return &AuditLog{max: max, events: make([]AuditEvent, 0, max)}
}
func (a *AuditLog) Record(action, id, reason string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.events) >= a.max {
		copy(a.events, a.events[1:])
		a.events = a.events[:a.max-1]
	}
	a.events = append(a.events, AuditEvent{Time: time.Now(), Action: action, BlockID: id, Reason: reason})
}
func (a *AuditLog) List() []AuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]AuditEvent(nil), a.events...)
}
func (a *AuditLog) String() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return fmt.Sprintf("retention audit events=%d", len(a.events))
}
