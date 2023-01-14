package event

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type RemoveEntity struct {
    Entity *donburi.Entity
}

var RemoveEntityEvent = events.NewEventType[RemoveEntity]()
