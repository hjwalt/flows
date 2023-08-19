package stateless

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
	"go.uber.org/zap"
)

// constructor
func NewTopicSwitch(configurations ...runtime.Configuration[*TopicSwitch]) BatchFunction {
	fn := &TopicSwitch{
		functions: make(map[string]BatchFunction),
	}
	for _, configuration := range configurations {
		fn = configuration(fn)
	}
	return fn.Apply
}

// configuration
func WithTopicSwitchStatelessSingleFunction(topic string, f BatchFunction) runtime.Configuration[*TopicSwitch] {
	return func(sts *TopicSwitch) *TopicSwitch {
		sts.functions[topic] = f
		return sts
	}
}

// implementation
type TopicSwitch struct {
	functions map[string]BatchFunction
}

func (r *TopicSwitch) Apply(c context.Context, m []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {

	messageMultiMap := structure.NewMultiMap[string, message.Message[message.Bytes, message.Bytes]]()
	for _, mi := range m {
		messageMultiMap.Add(mi.Topic, mi)
	}

	resultMessages := []message.Message[message.Bytes, message.Bytes]{}
	for k, v := range messageMultiMap.GetAll() {
		logger.Info("join switch", zap.String("topic", k))

		fn, fnExists := r.functions[k]
		if !fnExists {
			return make([]message.Message[[]byte, []byte], 0), errors.Join(errors.New(k), ErrSwitchMissingTopic)
		}

		currGroupMessages, currGroupHandlerErr := fn(c, v)
		if currGroupHandlerErr != nil {
			return make([]message.Message[[]byte, []byte], 0), currGroupHandlerErr
		}
		resultMessages = append(resultMessages, currGroupMessages...)
	}

	return resultMessages, nil
}

var (
	ErrSwitchMissingTopic = errors.New("stateless switch missing topic")
)