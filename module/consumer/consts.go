package consumer

import "time"

const (
	ConsumerBuffer int    = 5
	ItemServeTime         = 2 * time.Second
	PanicDuration         = 10 * time.Second
	FoolServer     string = "server is full"
)
