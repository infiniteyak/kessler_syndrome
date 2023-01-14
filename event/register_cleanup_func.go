package event

import (
	"github.com/yohamta/donburi/features/events"
)

type RegisterCleanupFunc struct {
    Function func()
}

var RegisterCleanupFuncEvent = events.NewEventType[RegisterCleanupFunc]()
