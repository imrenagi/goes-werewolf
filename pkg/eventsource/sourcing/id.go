package sourcing

import "github.com/google/uuid"

type EventSourceId uuid.UUID

func NewEventSourceId() EventSourceId {
	return EventSourceId(uuid.New())
}
