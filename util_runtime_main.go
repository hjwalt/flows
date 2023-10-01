package flows

import (
	"errors"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	AllInstances = "all_instances"
)

func NewMain() Main {
	return &main{
		allInstance: AllInstances,
		runtimes:    make(map[string]func(inverse.Container) []runtime.Runtime),
	}
}

func NewMainAllInstance(allInstanceName string) Main {
	return &main{
		allInstance: allInstanceName,
		runtimes:    make(map[string]func(inverse.Container) []runtime.Runtime),
	}
}

type Main interface {
	Runtimes(i string, r func(inverse.Container) []runtime.Runtime) error
	Prebuilt(i string, rf func(ci inverse.Container) Prebuilt) error
	Registrar(i string, rf func(ci inverse.Container)) error
	Start(i string) error
}

type main struct {
	allInstance string
	runtimes    map[string]func(inverse.Container) []runtime.Runtime
}

func (m *main) Runtimes(i string, r func(inverse.Container) []runtime.Runtime) error {
	if m == nil {
		return errors.New("main is missing")
	}
	m.runtimes[i] = r
	return nil
}

func (m *main) Prebuilt(i string, rf func(ci inverse.Container) Prebuilt) error {
	return m.Runtimes(
		i,
		func(ci inverse.Container) []runtime.Runtime {
			r := rf(ci)
			r.Register(ci)
			return InjectedRuntimes(ci)
		},
	)
}

func (m *main) Registrar(i string, rf func(ci inverse.Container)) error {
	return m.Runtimes(
		i,
		func(ci inverse.Container) []runtime.Runtime {
			rf(ci)
			return InjectedRuntimes(ci)
		},
	)
}

func (m *main) Start(i string) error {
	if m == nil {
		return errors.New("main is missing")
	}

	allRuntimes := []runtime.Runtime{}

	if i == m.allInstance {
		for _, rConstructor := range m.runtimes {
			rs := rConstructor(inverse.NewContainer())
			allRuntimes = append(allRuntimes, rs...)
		}

	} else {
		rConstructor, rExist := m.runtimes[i]
		if !rExist {
			return errors.New("instance is missing")
		}

		rs := rConstructor(inverse.NewContainer())
		allRuntimes = append(allRuntimes, rs...)
	}

	if len(allRuntimes) == 0 {
		return errors.New("no runtimes resolved")
	}

	r := &RuntimeFacade{
		Runtimes: allRuntimes,
	}

	return r.Start()
}
