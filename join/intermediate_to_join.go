package join

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

var ErrorIntermediateToJoinDeserialiseMessage = errors.New("error deserialising message")

// constructor
func NewIntermediateToJoinMap(configurations ...runtime.Configuration[*IntermediateToJoinMap]) stateless.SingleFunction {
	singleFunction := &IntermediateToJoinMap{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configuration
func WithIntermediateToJoinMapTransactionWrappedFunction(f stateless.SingleFunction) runtime.Configuration[*IntermediateToJoinMap] {
	return func(itjm *IntermediateToJoinMap) *IntermediateToJoinMap {
		itjm.transactionWrapped = f
		return itjm
	}
}

// implementation
type IntermediateToJoinMap struct {
	transactionWrapped stateless.SingleFunction
}

func (r *IntermediateToJoinMap) Apply(c context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	logger.Info("intermediate to join", zap.String("topic", m.Topic))
	messageDeserialized, messageDeserialisationError := IntermediateValueFormat.Unmarshal(m.Value)
	if messageDeserialisationError != nil {
		return make([]message.Message[[]byte, []byte], 0), errors.Join(ErrorIntermediateToJoinDeserialiseMessage, messageDeserialisationError)
	}

	messageDeserialized.Partition = m.Partition
	messageDeserialized.Offset = m.Offset

	return r.transactionWrapped(c, messageDeserialized)
}
