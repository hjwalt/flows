package example_word_multi

import (
	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/example/example_word_count"
	"github.com/hjwalt/flows/example/example_word_materialise"
	"github.com/hjwalt/runway/runtime"
)

const (
	Instance = "flow-word-multi"
)

func instance() runtime.Runtime {
	wordCount := example_word_count.Registrar()
	wordMaterialise := example_word_materialise.Registrar()

	wordMaterialise.RegisterRuntime()
	wordCount.Register()
	wordMaterialise.Register()

	return &flows.RuntimeFacade{
		Runtimes: flows.InjectedRuntimes(),
	}
}

func Register(m flows.Main) {
	m.Register(Instance, instance)
}
