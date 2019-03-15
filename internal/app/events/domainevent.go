package events

import "time"

// DomainEvent is base class for all domain events.
type DomainEvent struct {
	ID           string
	EventVersion int32
	Timestamp    time.Time
	EventName    string
}
