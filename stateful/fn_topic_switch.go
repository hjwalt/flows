package stateful

import (
	"context"
	"fmt"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"go.uber.org/zap"
)

// constructor
func NewTopicSwitch(configurations ...runtime.Configuration[*TopicSwitch]) SingleFunction {
	singleFunction := &TopicSwitch{
		functions: make(map[string]SingleFunction),
	}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configuration
func WithTopicSwitchFunction(topic string, f SingleFunction) runtime.Configuration[*TopicSwitch] {
	return func(sts *TopicSwitch) *TopicSwitch {
		sts.functions[topic] = f
		return sts
	}
}

// implementation
type TopicSwitch struct {
	functions map[string]SingleFunction
}

func (r *TopicSwitch) Apply(c context.Context, m message.Message[message.Bytes, message.Bytes], s State[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], State[message.Bytes], error) {
	logger.Info("stateful switch", zap.String("topic", m.Topic))
	fn, fnExists := r.functions[m.Topic]
	if !fnExists {
		return make([]message.Message[[]byte, []byte], 0), s, TopicMissingError(m.Topic)
	}
	return fn(c, m, s)
}

func TopicMissingError(topic string) error {
	return fmt.Errorf("function for topic %s missing", topic)
}
