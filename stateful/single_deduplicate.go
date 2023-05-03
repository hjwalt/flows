package stateful

import (
	"context"

	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

// constructor
func NewSingleStatefulDeduplicate(configurations ...runtime.Configuration[*SingleStatefulDeduplicate]) SingleFunction {
	singleFunction := &SingleStatefulDeduplicate{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configuration
func WithSingleStatefulDeduplicateNextFunction(next SingleFunction) runtime.Configuration[*SingleStatefulDeduplicate] {
	return func(st *SingleStatefulDeduplicate) *SingleStatefulDeduplicate {
		st.next = next
		return st
	}
}

// implementation
type SingleStatefulDeduplicate struct {
	next SingleFunction
}

func (r SingleStatefulDeduplicate) Apply(c context.Context, m message.Message[message.Bytes, message.Bytes], inState SingleState[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], SingleState[message.Bytes], error) {

	s := SetDefault(inState)

	// for edge case of offset zero
	storedOffset, offsetExists := s.Internal.GetV1().OffsetProgress[m.Partition]
	if !offsetExists {
		storedOffset = -1
	}

	// extract messages
	storedMessages := make([]message.Message[message.Bytes, message.Bytes], 0)
	for _, lastResultBytes := range s.Results.GetV1().Messages {
		messageConverted, convertErr := format.Convert(lastResultBytes, format.Bytes(), message.Format())
		if convertErr != nil {
			return make([]message.Message[message.Bytes, message.Bytes], 0), s, convertErr
		}
		storedMessages = append(storedMessages, messageConverted)
	}

	inputOffset := m.Offset

	if storedOffset > inputOffset {
		logger.Debug("skipping without output", zap.Int32("partition", m.Partition), zap.Int64("offset", m.Offset))
		return make([]message.Message[message.Bytes, message.Bytes], 0), s, nil
	}

	if storedOffset == inputOffset {
		logger.Debug("skipping with output", zap.Int32("partition", m.Partition), zap.Int64("offset", m.Offset))
		return storedMessages, s, nil
	}

	logger.Debug("acting", zap.Int32("partition", m.Partition), zap.Int64("offset", m.Offset))
	nm, ns, applyErr := r.next(c, m, s)

	if applyErr != nil {
		return make([]message.Message[message.Bytes, message.Bytes], 0), s, applyErr
	}

	ns.Internal.GetV1().OffsetProgress[m.Partition] = inputOffset
	ns.Results.GetV1().Messages = make([][]byte, 0)

	for _, msgs := range nm {
		messageConverted, convertErr := format.Convert(msgs, message.Format(), format.Bytes())
		if convertErr != nil {
			return make([]message.Message[message.Bytes, message.Bytes], 0), s, convertErr
		}
		ns.Results.GetV1().Messages = append(ns.Results.GetV1().Messages, messageConverted)
	}

	if storedOffset == -1 {
		logger.Info("offset data reset, prepending last results before output", zap.Int32("partition", m.Partition), zap.Int64("offset", m.Offset))
		nm = append(storedMessages, nm...)
	}

	return nm, ns, applyErr

}
