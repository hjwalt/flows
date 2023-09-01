package flows

import "github.com/hjwalt/runway/runtime"

type RuntimeRegistrar interface {
	Register()
	RegisterRuntime()
}

func Register(m Main, instance string, rf func() RuntimeRegistrar) {
	m.Register(
		instance,
		func() runtime.Runtime {
			r := rf()

			r.RegisterRuntime()
			r.Register()

			return &RuntimeFacade{
				Runtimes: InjectedRuntimes(),
			}
		},
	)
}
