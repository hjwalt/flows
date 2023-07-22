package topic

import "github.com/hjwalt/runway/format"

type Topic[K any, V any] interface {
	Topic() string
	KeyFormat() format.Format[K]
	ValueFormat() format.Format[V]
}
