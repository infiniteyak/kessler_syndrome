package event

import "github.com/yohamta/donburi/features/events"

type AsteroidsCountUpdate struct {
    Value int
}

var AsteroidsCountUpdateEvent = events.NewEventType[AsteroidsCountUpdate]()
