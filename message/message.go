package message

import "time"

type Message[K any, V any] struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       K
	Value     V
	Headers   map[string][]Bytes
	Timestamp time.Time
}

type Bytes = []byte
