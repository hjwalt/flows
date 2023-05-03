package flows

import (
	"errors"

	"github.com/hjwalt/flows/runtime"
)

func Main() *main {
	return &main{
		runtimes: make(map[string]func() runtime.Runtime),
	}
}

type main struct {
	runtimes map[string]func() runtime.Runtime
}

func (m *main) Register(i string, r func() runtime.Runtime) error {
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

	rConstructor, rExist := m.runtimes[i]
	if !rExist {
		return errors.New("instance is missing")
	}

	r := rConstructor()

	return r.Start()
}
