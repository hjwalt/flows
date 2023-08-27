package stateful

import (
	"context"
	"fmt"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
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

func (r *TopicSwitch) Apply(c context.Context, m flow.Message[structure.Bytes, structure.Bytes], s State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], State[structure.Bytes], error) {
	logger.Info("stateful switch", zap.String("topic", m.Topic))
	fn, fnExists := r.functions[m.Topic]
	if !fnExists {
		return make([]flow.Message[[]byte, []byte], 0), s, TopicMissingError(m.Topic)
	}
	return fn(c, m, s)
}

func TopicMissingError(topic string) error {
	return fmt.Errorf("function for topic %s missing", topic)
}
