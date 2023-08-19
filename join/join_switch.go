package join

import (
	"context"
	"fmt"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
	"go.uber.org/zap"
)

// constructor
func NewJoinSwitch(configurations ...runtime.Configuration[*JoinSwitch]) stateless.BatchFunction {
	singleFunction := &JoinSwitch{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// // configuration
func WithJoinSwitchSourceTopicFunction(f stateless.BatchFunction) runtime.Configuration[*JoinSwitch] {
	return func(sts *JoinSwitch) *JoinSwitch {
		sts.sourceTopicHandlerFunction = f
		return sts
	}
}

func WithJoinSwitchIntermediateTopicFunction(topic string, f stateless.BatchFunction) runtime.Configuration[*JoinSwitch] {
	return func(sts *JoinSwitch) *JoinSwitch {
		sts.intermediateTopic = topic
		sts.intermediateTopicHandlerFunction = f
		return sts
	}
}

// implementation
type JoinSwitch struct {
	sourceTopicHandlerFunction       stateless.BatchFunction
	intermediateTopic                string
	intermediateTopicHandlerFunction stateless.BatchFunction
}

func (r *JoinSwitch) Apply(c context.Context, m []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {

	messageMultiMap := structure.NewMultiMap[string, message.Message[message.Bytes, message.Bytes]]()
	for _, mi := range m {
		messageMultiMap.Add(mi.Topic, mi)
	}

	resultMessages := []message.Message[message.Bytes, message.Bytes]{}
	for k, v := range messageMultiMap.GetAll() {
		logger.Info("join switch", zap.String("topic", k))

		currGroupHandler := r.sourceTopicHandlerFunction
		if k == r.intermediateTopic {
			currGroupHandler = r.intermediateTopicHandlerFunction
		}
		currGroupMessages, currGroupHandlerErr := currGroupHandler(c, v)
		if currGroupHandlerErr != nil {
			return make([]message.Message[[]byte, []byte], 0), currGroupHandlerErr
		}
		resultMessages = append(resultMessages, currGroupMessages...)
	}

	return resultMessages, nil
}

func TopicMissingError(topic string) error {
	return fmt.Errorf("function for topic %s missing", topic)
}
