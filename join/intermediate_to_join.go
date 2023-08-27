package join

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
	"go.uber.org/zap"
)

var ErrorIntermediateToJoinDeserialiseMessage = errors.New("error deserialising message")

// constructor
func NewIntermediateToJoinMap(configurations ...runtime.Configuration[*IntermediateToJoinMap]) stateless.BatchFunction {
	singleFunction := &IntermediateToJoinMap{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configuration
func WithIntermediateToJoinMapTransactionWrappedFunction(f stateless.BatchFunction) runtime.Configuration[*IntermediateToJoinMap] {
	return func(itjm *IntermediateToJoinMap) *IntermediateToJoinMap {
		itjm.transactionWrapped = f
		return itjm
	}
}

// implementation
type IntermediateToJoinMap struct {
	transactionWrapped stateless.BatchFunction
}

func (r *IntermediateToJoinMap) Apply(c context.Context, ms []flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
	messagesToMap := make([]flow.Message[structure.Bytes, structure.Bytes], len(ms))

	for i, m := range ms {
		logger.Info("intermediate to join", zap.String("topic", m.Topic))
		messageDeserialized, messageDeserialisationError := IntermediateValueFormat.Unmarshal(m.Value)
		if messageDeserialisationError != nil {
			return make([]flow.Message[[]byte, []byte], 0), errors.Join(ErrorIntermediateToJoinDeserialiseMessage, messageDeserialisationError)
		}

		messageDeserialized.Partition = m.Partition
		messageDeserialized.Offset = m.Offset

		messagesToMap[i] = messageDeserialized
	}

	return r.transactionWrapped(c, messagesToMap)
}
