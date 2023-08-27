package flow

import (
	"time"

	"github.com/hjwalt/runway/structure"
)

type Message[K any, V any] struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       K
	Value     V
	Headers   map[string][]structure.Bytes
	Timestamp time.Time
}

func EmptySlice() []Message[structure.Bytes, structure.Bytes] {
	return []Message[structure.Bytes, structure.Bytes]{}
}
