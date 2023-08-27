package flow

import "github.com/hjwalt/runway/format"

func JsonTopic[K any, V any](topic string) Topic[K, V] {
	return GenericTopic[K, V](topic, format.Json[K](), format.Json[V]())
}
