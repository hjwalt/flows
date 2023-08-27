package flow

import "github.com/hjwalt/runway/format"

func Json[K any, V any](topic string) Topic[K, V] {
	return Generic[K, V](topic, format.Json[K](), format.Json[V]())
}
