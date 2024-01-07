package task

import (
	"github.com/hjwalt/runway/format"
	"google.golang.org/protobuf/proto"
)

func ProtobufChannel[V proto.Message](channel string) Channel[V] {
	return GenericChannel(channel, format.Protobuf[V]())
}
