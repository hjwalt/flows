package stateless_error_handler_skip_list

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/structure"
)

func New(skipList []error) stateless.ErrorHandlerFunction {
	f := fn{skipList: skipList}
	return f.handle
}

func Default() stateless.ErrorHandlerFunction {
	return New(
		[]error{
			format.ErrFormatConversionMarshal,
			format.ErrFormatConversionUnmarshal,
		},
	)
}

type fn struct {
	skipList []error
}

func (r *fn) handle(c context.Context, m flow.Message[structure.Bytes, structure.Bytes], e error) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
	for _, skipErr := range r.skipList {
		if errors.Is(e, skipErr) {
			logger.WarnErr("stateless skipped error", e)
			return flow.EmptySlice(), nil
		}
	}
	return flow.EmptySlice(), e
}
