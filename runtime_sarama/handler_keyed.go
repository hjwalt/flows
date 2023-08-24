package runtime_sarama

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/metric"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/reflect"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
	"go.uber.org/zap"
)

// constructor
func NewKeyedHandler(configurations ...runtime.Configuration[*KeyedHandler]) ConsumerHandler {
	loop := &KeyedHandler{
		MaxBufferred: 1000,
		MaxPerKey:    1,
		MaxDelay:     100 * time.Millisecond,
	}
	for _, configuration := range configurations {
		loop = configuration(loop)
	}
	return loop
}

// configuration
func WithKeyedHandlerMaxBufferred(maxBuffered int64) runtime.Configuration[*KeyedHandler] {
	return func(cbl *KeyedHandler) *KeyedHandler {
		cbl.MaxBufferred = maxBuffered
		return cbl
	}
}

func WithKeyedHandlerMaxDelay(maxDelay time.Duration) runtime.Configuration[*KeyedHandler] {
	return func(cbl *KeyedHandler) *KeyedHandler {
		cbl.MaxDelay = maxDelay
		return cbl
	}
}

func WithKeyedHandlerMaxPerKey(maxPerKey int64) runtime.Configuration[*KeyedHandler] {
	return func(cbl *KeyedHandler) *KeyedHandler {
		cbl.MaxPerKey = maxPerKey
		return cbl
	}
}

func WithKeyedHandlerFunction(loopFunction stateless.BatchFunction) runtime.Configuration[*KeyedHandler] {
	return func(cbl *KeyedHandler) *KeyedHandler {
		cbl.F = loopFunction
		return cbl
	}
}

func WithKeyedHandlerKeyFunction(keyFunction stateful.PersistenceIdFunction[message.Bytes, message.Bytes]) runtime.Configuration[*KeyedHandler] {
	return func(cbl *KeyedHandler) *KeyedHandler {
		cbl.K = keyFunction
		return cbl
	}
}

func WithKeyedHandlerPrometheus() runtime.Configuration[*KeyedHandler] {
	return func(cbl *KeyedHandler) *KeyedHandler {
		cbl.metric = metric.PrometheusConsume()
		return cbl
	}
}

// implementation
type KeyedHandler struct {
	// batching configuration
	MaxBufferred int64
	MaxDelay     time.Duration
	MaxPerKey    int64
	F            stateless.BatchFunction
	K            stateful.PersistenceIdFunction[message.Bytes, message.Bytes]

	// metrics
	metric metric.Consume
}

func (h *KeyedHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	khs := &keyedHandlerState{
		id:             "claim-" + claim.Topic() + "-" + reflect.GetString(claim.Partition()),
		messageClaimed: make(chan *sarama.ConsumerMessage),
		countReached:   make(chan bool),
		keyCount:       structure.NewCountMap[string](),
	}
	khs.context, khs.cancel = context.WithCancel(context.Background())

	khs.reset(h.MaxDelay)
	go khs.claimLoop(session, claim)
	for {
		completed, err := h.ConsumeClaimIteration(session, claim, khs)
		if err != nil {
			khs.cancel()
			return err
		}

		if completed {
			return nil
		}
	}
}

func (kh *KeyedHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (kh *KeyedHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (kh *KeyedHandler) Start() error {
	return nil
}

func (kh *KeyedHandler) Stop() {
}

func (h *KeyedHandler) ConsumeClaimIteration(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, khs *keyedHandlerState) (bool, error) {

	select {
	case saramaMessage := <-khs.messageClaimed:
		logger.Info("read", zap.String("handler", khs.id), zap.String("topic", saramaMessage.Topic), zap.Int32("partition", saramaMessage.Partition), zap.Int64("offset", saramaMessage.Offset))

		flowMessage, flowMessageErr := FromConsumerMessage(saramaMessage)
		if flowMessageErr != nil {
			return true, flowMessageErr
		}

		keyForMessage, keyError := h.K(khs.context, flowMessage)
		if keyError != nil {
			return true, keyError
		}

		if khs.keyCount.Get(keyForMessage) >= h.MaxPerKey {
			// prevents updating the same key on the same batch
			if err := khs.consumeBatch("key", h.MaxDelay, session, h.F, h.metric); err != nil {
				return true, err
			}
		}

		khs.messages = append(khs.messages, flowMessage)
		khs.messageToCommit = saramaMessage
		khs.keyCount.Add(keyForMessage, 1)

		if len(khs.messages) >= int(h.MaxBufferred) {
			if err := khs.consumeBatch("batch", h.MaxDelay, session, h.F, h.metric); err != nil {
				return true, err
			}
		}
	case <-khs.timerReached:
		if err := khs.consumeBatch("timer", h.MaxDelay, session, h.F, h.metric); err != nil {
			return true, err
		}
	}
	if khs.context.Err() != nil {
		return true, nil
	}
	return false, nil
}

type keyedHandlerState struct {
	id string

	// batched data
	messageToCommit *sarama.ConsumerMessage
	messages        []message.Message[[]byte, []byte]
	keyCount        structure.CountMap[string]

	// batching mechanic
	countReached chan bool
	timerReached <-chan time.Time

	// channel termination mechanic
	messageClaimed chan *sarama.ConsumerMessage
	context        context.Context
	cancel         context.CancelFunc
}

func (khs *keyedHandlerState) reset(maxDelay time.Duration) {
	khs.messageToCommit = nil
	khs.messages = make([]message.Message[[]byte, []byte], 0)
	khs.timerReached = time.After(maxDelay)
	khs.keyCount.Clear()
}

func (khs *keyedHandlerState) claimLoop(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) {
	// the claim loop channel will stop when the consumer is stopped, triggering end, triggering consume loop end
	for saramaMessage := range claim.Messages() {
		khs.messageClaimed <- saramaMessage
		if khs.context.Err() != nil {
			return
		}
	}
	khs.cancel()
}

func (khs *keyedHandlerState) consumeBatch(
	source string,
	maxDelay time.Duration,
	session sarama.ConsumerGroupSession,
	fn stateless.BatchFunction,
	metric metric.Consume,
) error {
	if len(khs.messages) == 0 {
		khs.reset(maxDelay)
		return nil
	}
	logger.Info("handling", zap.String("source", source), zap.Int("batch", len(khs.messages)), zap.String("handler", khs.id), zap.String("topic", khs.messageToCommit.Topic), zap.Int32("partition", khs.messageToCommit.Partition), zap.Int64("offset", khs.messageToCommit.Offset))
	if _, err := fn(context.Background(), khs.messages); err != nil {
		return err
	}
	session.MarkMessage(khs.messageToCommit, "")
	logger.Info("commit", zap.String("topic", khs.messageToCommit.Topic), zap.Int32("partition", khs.messageToCommit.Partition), zap.Int64("offset", khs.messageToCommit.Offset))

	if metric != nil {
		metric.MessagesProcessedIncrement(khs.messageToCommit.Topic, khs.messageToCommit.Partition, int64(len(khs.messages)))
	}

	khs.reset(maxDelay)
	return nil
}
