package flows

import (
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

type RuntimeRegistrar interface {
	Register(ci inverse.Container)
}

func Runtimes(ci inverse.Container, rf func(inverse.Container) RuntimeRegistrar) []runtime.Runtime {
	r := rf(ci)

	r.Register(ci)

	return InjectedRuntimes(ci)
}

func Register(m Main, instance string, rf func(ci inverse.Container) RuntimeRegistrar) {
	m.Register(
		instance,
		func(ci inverse.Container) []runtime.Runtime {
			return Runtimes(ci, rf)
		},
	)
}
