package task

import "time"

type Message[V any] struct {
	Channel   string
	Value     V
	Headers   map[string]any
	Timestamp time.Time
}
