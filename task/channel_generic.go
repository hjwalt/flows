package task

import "github.com/hjwalt/runway/format"

type genericChannel[V any] struct {
	channel     string
	valueFormat format.Format[V]
}

func (c genericChannel[V]) Name() string {
	return c.channel
}

func (c genericChannel[V]) ValueFormat() format.Format[V] {
	return c.valueFormat
}

func GenericChannel[V any](channel string, valueFormat format.Format[V]) Channel[V] {
	return genericChannel[V]{
		channel:     channel,
		valueFormat: valueFormat,
	}
}
