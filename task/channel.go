package task

import "github.com/hjwalt/runway/format"

type Channel[V any] interface {
	Name() string
	ValueFormat() format.Format[V]
}
