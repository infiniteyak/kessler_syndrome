package event

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type RegisterEntity struct {
    Entity *donburi.Entity
}

var RegisterEntityEvent = events.NewEventType[RegisterEntity]()
