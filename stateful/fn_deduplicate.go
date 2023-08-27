package stateful

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
	"go.uber.org/zap"
)

// constructor
func NewDeduplicate(configurations ...runtime.Configuration[*Deduplicate]) SingleFunction {
	singleFunction := &Deduplicate{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configuration
func WithDeduplicateNextFunction(next SingleFunction) runtime.Configuration[*Deduplicate] {
	return func(st *Deduplicate) *Deduplicate {
		st.next = next
		return st
	}
}

// implementation
type Deduplicate struct {
	next SingleFunction
}

func (r Deduplicate) Apply(c context.Context, m flow.Message[structure.Bytes, structure.Bytes], inState State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], State[structure.Bytes], error) {

	s := SetDefault(inState)

	// for edge case of offset zero
	storedOffset, offsetExists := s.Internal.GetV1().OffsetProgress[m.Partition]
	if !offsetExists {
		storedOffset = -1
	}

	// extract messages
	storedMessages := make([]flow.Message[structure.Bytes, structure.Bytes], 0)
	for _, lastResultBytes := range s.Results.GetV1().Messages {
		messageConverted, convertErr := format.Convert(lastResultBytes, format.Bytes(), flow.Format())
		if convertErr != nil {
			return make([]flow.Message[structure.Bytes, structure.Bytes], 0), s, convertErr
		}
		storedMessages = append(storedMessages, messageConverted)
	}

	inputOffset := m.Offset

	if storedOffset > inputOffset {
		logger.Debug("skipping without output", zap.Int32("partition", m.Partition), zap.Int64("offset", m.Offset))
		return make([]flow.Message[structure.Bytes, structure.Bytes], 0), s, nil
	}

	if storedOffset == inputOffset {
		logger.Debug("skipping with output", zap.Int32("partition", m.Partition), zap.Int64("offset", m.Offset))
		return storedMessages, s, nil
	}

	logger.Debug("acting", zap.Int32("partition", m.Partition), zap.Int64("offset", m.Offset))
	nm, ns, applyErr := r.next(c, m, s)

	if applyErr != nil {
		return make([]flow.Message[structure.Bytes, structure.Bytes], 0), s, applyErr
	}

	ns.Internal.GetV1().OffsetProgress[m.Partition] = inputOffset
	ns.Results.GetV1().Messages = make([][]byte, 0)

	for _, msgs := range nm {
		messageConverted, convertErr := format.Convert(msgs, flow.Format(), format.Bytes())
		if convertErr != nil {
			return make([]flow.Message[structure.Bytes, structure.Bytes], 0), s, convertErr
		}
		ns.Results.GetV1().Messages = append(ns.Results.GetV1().Messages, messageConverted)
	}

	if storedOffset == -1 {
		logger.Info("offset data reset, prepending last results before output", zap.Int32("partition", m.Partition), zap.Int64("offset", m.Offset))
		nm = append(storedMessages, nm...)
	}

	return nm, ns, applyErr

}
