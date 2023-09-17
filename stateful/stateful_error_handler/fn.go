package stateful_error_handler

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/structure"
)

func New(
	next stateful.SingleFunction,
	handler stateful.ErrorHandlerFunction,
) stateful.SingleFunction {
	res := &fn{next: next, handler: handler}
	return res.apply
}

type fn struct {
	next    stateful.SingleFunction
	handler stateful.ErrorHandlerFunction
}

func (r *fn) apply(c context.Context, m flow.Message[structure.Bytes, structure.Bytes], s stateful.State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], stateful.State[structure.Bytes], error) {
	res, next, err := r.next(c, m, s)
	if err != nil {
		return r.handler(c, m, s, err)
	} else {
		return res, next, nil
	}
}
