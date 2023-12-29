package main

import (
	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/example/example_word_adapter"
	"github.com/hjwalt/flows/example/example_word_collect"
	"github.com/hjwalt/flows/example/example_word_count"
	"github.com/hjwalt/flows/example/example_word_join"
	"github.com/hjwalt/flows/example/example_word_materialise"
	"github.com/hjwalt/flows/example/example_word_remap"
	"github.com/hjwalt/runway/environment"
)

func main() {
	m := flows.NewMain()

	example_word_adapter.Register(m)
	example_word_collect.Register(m)
	example_word_count.Register(m)
	example_word_join.Register(m)
	example_word_materialise.Register(m)
	example_word_remap.Register(m)

	err := m.Start(environment.GetString("INSTANCE", flows.AllInstances))

	if err != nil {
		panic(err)
	}
}
