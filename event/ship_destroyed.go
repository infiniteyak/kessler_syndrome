package event

import "github.com/yohamta/donburi/features/events"

type ShipDestroyed struct { }

var ShipDestroyedEvent = events.NewEventType[ShipDestroyed]()
