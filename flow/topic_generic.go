package flow

import "github.com/hjwalt/runway/format"

type genericTopic[K any, V any] struct {
	topic       string
	keyFormat   format.Format[K]
	valueFormat format.Format[V]
}

func (t genericTopic[K, V]) Name() string {
	return t.topic
}

func (t genericTopic[K, V]) KeyFormat() format.Format[K] {
	return t.keyFormat
}

func (t genericTopic[K, V]) ValueFormat() format.Format[V] {
	return t.valueFormat
}

func GenericTopic[K any, V any](topic string, keyFormat format.Format[K], valueFormat format.Format[V]) Topic[K, V] {
	return &genericTopic[K, V]{
		topic:       topic,
		keyFormat:   keyFormat,
		valueFormat: valueFormat,
	}
}
