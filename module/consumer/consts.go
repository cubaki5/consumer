package consumer

import "time"

const (
	ConsumerBuffer int           = 5
	ItemServeTime  time.Duration = 2 * time.Second
	PanicDuration  time.Duration = 30 * time.Second
	FoolServer     string        = "server is full"
)
