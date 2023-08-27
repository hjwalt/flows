package flow

import "github.com/hjwalt/runway/format"

type GenericTopic[K any, V any] struct {
	topic       string
	keyFormat   format.Format[K]
	valueFormat format.Format[V]
}

func (t GenericTopic[K, V]) Name() string {
	return t.topic
}

func (t GenericTopic[K, V]) KeyFormat() format.Format[K] {
	return t.keyFormat
}

func (t GenericTopic[K, V]) ValueFormat() format.Format[V] {
	return t.valueFormat
}

func Generic[K any, V any](topic string, keyFormat format.Format[K], valueFormat format.Format[V]) Topic[K, V] {
	return &GenericTopic[K, V]{
		topic:       topic,
		keyFormat:   keyFormat,
		valueFormat: valueFormat,
	}
}
