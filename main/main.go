package main

import (
	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/example"
	"github.com/hjwalt/runway/environment"
)

func main() {
	m := flows.NewMain()

	m.Register("word-count", example.WordCount)
	m.Register("word-collect", example.WordCollect)
	m.Register("word-remap", example.WordRemap)
	m.Register("word-join", example.WordJoin)
	m.Register("word-materialise", example.WordMaterialise)

	err := m.Start(environment.GetString("INSTANCE", "word-remap"))

	if err != nil {
		panic(err)
	}
}
