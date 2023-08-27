package flow

import "github.com/hjwalt/runway/format"

func String(topic string) Topic[string, string] {
	return Generic(topic, format.String(), format.String())
}
