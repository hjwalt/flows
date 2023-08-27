package flow

import (
	"github.com/hjwalt/runway/format"
	"google.golang.org/protobuf/proto"
)

func Protobuf[K proto.Message, V proto.Message](topic string) Topic[K, V] {
	return Generic[K, V](topic, format.Protobuf[K](), format.Protobuf[V]())
}
