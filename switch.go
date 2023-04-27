package flows

import (
	"errors"

	"github.com/hjwalt/flows/runtime"
)

func Main() *main {
	return &main{
		runtimes: make(map[string]runtime.Runtime),
	}
}

type main struct {
	runtimes map[string]runtime.Runtime
}

func (m *main) Register(i string, r runtime.Runtime) error {
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

	r, rExist := m.runtimes[i]
	if !rExist {
		return errors.New("instance is missing")
	}

	return r.Start()
}
