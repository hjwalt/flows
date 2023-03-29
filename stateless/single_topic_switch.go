package stateless

import (
	"context"
	"fmt"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

// constructor
func NewSingleTopicSwitch(configurations ...runtime.Configuration[*SingleTopicSwitch]) SingleFunction {
	singleFunction := &SingleTopicSwitch{
		functions: make(map[string]SingleFunction),
	}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configuration
func WithSingleTopicSwitchStatelessSingleFunction(topic string, f SingleFunction) runtime.Configuration[*SingleTopicSwitch] {
	return func(sts *SingleTopicSwitch) *SingleTopicSwitch {
		sts.functions[topic] = f
		return sts
	}
}

// implementation
type SingleTopicSwitch struct {
	functions map[string]SingleFunction
}

func (r *SingleTopicSwitch) Apply(c context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	logger.Info("stateless switch", zap.String("topic", m.Topic))
	fn, fnExists := r.functions[m.Topic]
	if !fnExists {
		return make([]message.Message[[]byte, []byte], 0), TopicMissingError(m.Topic)
	}
	return fn(c, m)
}

func TopicMissingError(topic string) error {
	return fmt.Errorf("function for topic %s missing", topic)
}
