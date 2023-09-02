package example_word_multi

import (
	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/example/example_word_count"
	"github.com/hjwalt/flows/example/example_word_join"
	"github.com/hjwalt/flows/example/example_word_materialise"
	"github.com/hjwalt/runway/runtime"
)

const (
	Instance = "flow-word-multi"
)

func instance() runtime.Runtime {
	allRuntimes := []runtime.Runtime{}

	allRuntimes = append(allRuntimes, flows.Runtimes(example_word_count.Registrar)...)
	allRuntimes = append(allRuntimes, flows.Runtimes(example_word_materialise.Registrar)...)
	allRuntimes = append(allRuntimes, flows.Runtimes(example_word_join.Registrar)...)

	return &flows.RuntimeFacade{
		Runtimes: allRuntimes,
	}
}

func Register(m flows.Main) {
	m.Register(Instance, instance)
}
