package task

import "github.com/hjwalt/runway/format"

func StringChannel(channel string) Channel[string] {
	return GenericChannel(channel, format.String())
}
