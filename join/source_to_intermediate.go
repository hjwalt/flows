package join

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
	"go.uber.org/zap"
)

func NewSourceToIntermediateMap(configurations ...runtime.Configuration[*SourceToIntermediateMap]) stateless.SingleFunction {
	singleFunction := &SourceToIntermediateMap{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configuration
func WithSourceToIntermediateMapIntermediateTopic(intermediateTopic string) runtime.Configuration[*SourceToIntermediateMap] {
	return func(stim *SourceToIntermediateMap) *SourceToIntermediateMap {
		stim.intermediateTopic = intermediateTopic
		return stim
	}
}

func WithSourceToIntermediateMapPersistenceIdFunction(f stateful.PersistenceIdFunction[[]byte, []byte]) runtime.Configuration[*SourceToIntermediateMap] {
	return func(stim *SourceToIntermediateMap) *SourceToIntermediateMap {
		stim.persistenceId = f
		return stim
	}
}

// implementation
type SourceToIntermediateMap struct {
	persistenceId     stateful.PersistenceIdFunction[[]byte, []byte]
	intermediateTopic string
}

func (r *SourceToIntermediateMap) Apply(c context.Context, m flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
	logger.Info("source to intermediate", zap.String("topic", m.Topic))
	persistenceId, persistenceIdError := r.persistenceId(c, m)
	if persistenceIdError != nil {
		logger.ErrorErr("error getting persistence id", persistenceIdError)
		return make([]flow.Message[[]byte, []byte], 0), persistenceIdError
	}

	// To ensure changes are sequenced
	joinKey := &protobuf.JoinKey{
		PersistenceId: persistenceId,
	}
	joinKeyBytes, joinKeySerialisationErr := IntermediateKeyFormat.Marshal(joinKey)
	if joinKeySerialisationErr != nil {
		logger.ErrorErr("error serialising join key", joinKeySerialisationErr)
		return make([]flow.Message[[]byte, []byte], 0), joinKeySerialisationErr
	}

	// To keep all information about the source message
	joinValueBytes, joinValueSerialisationErr := IntermediateValueFormat.Marshal(m)
	if joinValueSerialisationErr != nil {
		logger.ErrorErr("error serialising join value", joinValueSerialisationErr)
		return make([]flow.Message[[]byte, []byte], 0), joinValueSerialisationErr
	}

	remappedMessage := flow.Message[[]byte, []byte]{
		Topic: r.intermediateTopic,
		Key:   joinKeyBytes,
		Value: joinValueBytes,
	}

	return []flow.Message[structure.Bytes, structure.Bytes]{remappedMessage}, nil
}

var IntermediateValueFormat = flow.Format()
var IntermediateKeyFormat = format.Protobuf[*protobuf.JoinKey]()

func IntermediateTopicKeyFunction(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes]) (string, error) {
	keyValue, keyError := IntermediateKeyFormat.Unmarshal(m.Key)
	if keyError != nil {
		return "", keyError
	}
	return keyValue.PersistenceId, nil
}
