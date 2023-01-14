package event

import "github.com/yohamta/donburi/features/events"

type Score struct {
    Value int
}

var ScoreEvent = events.NewEventType[Score]()
