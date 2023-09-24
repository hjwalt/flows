package stateful_deduplicate_offset

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/structure"
	"go.uber.org/zap"
)

func New(
	next stateful.SingleFunction,
) stateful.SingleFunction {
	f := fn{
		next: next,
	}
	return f.apply
}

type fn struct {
	next stateful.SingleFunction
}

func (r fn) apply(c context.Context, m flow.Message[structure.Bytes, structure.Bytes], inState stateful.State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], stateful.State[structure.Bytes], error) {

	s := stateful.SetDefault(inState)

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
