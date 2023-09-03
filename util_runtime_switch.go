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
		runtimes: make(map[string]func(inverse.Container) []runtime.Runtime),
	}
}

type Main interface {
	Register(i string, r func(inverse.Container) []runtime.Runtime) error
	Start(i string) error
}

type main struct {
	runtimes map[string]func(inverse.Container) []runtime.Runtime
}

func (m *main) Register(i string, r func(inverse.Container) []runtime.Runtime) error {
	if m == nil {
		return errors.New("main is missing")
	}
	m.runtimes[i] = r
	return nil
}

func (m *main) Start(i string) error {
	if m == nil {
		return errors.New("main is missing")
	}

	allRuntimes := []runtime.Runtime{}

	if i == AllInstances {
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
