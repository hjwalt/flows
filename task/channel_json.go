package task

import (
	"github.com/hjwalt/runway/format"
)

func JsonChannel[V any](channel string) Channel[V] {
	return GenericChannel(channel, format.Json[V]())
}
