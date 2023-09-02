package flows

import (
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

type RuntimeRegistrar interface {
	Register()
	RegisterRuntime()
	Inverse() inverse.Container
}

func Runtimes(rf func() RuntimeRegistrar) []runtime.Runtime {
	r := rf()

	r.RegisterRuntime()
	r.Register()

	return InjectedRuntimes(r.Inverse())
}

func Register(m Main, instance string, rf func() RuntimeRegistrar) {
	m.Register(
		instance,
		func() runtime.Runtime {
			r := rf()

			r.RegisterRuntime()
			r.Register()

			return &RuntimeFacade{
				Runtimes: InjectedRuntimes(r.Inverse()),
			}
		},
	)
}
