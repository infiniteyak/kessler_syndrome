package event

import "github.com/yohamta/donburi/features/events"

type GameOver struct {}

var GameOverEvent = events.NewEventType[GameOver]()
