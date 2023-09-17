package stateless_error_handler

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/structure"
)

func New(
	next stateless.SingleFunction,
	handler stateless.ErrorHandlerFunction,
) stateless.SingleFunction {
	res := &fn{next: next, handler: handler}
	return res.apply
}

type fn struct {
	next    stateless.SingleFunction
	handler stateless.ErrorHandlerFunction
}

func (r *fn) apply(c context.Context, m flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
	res, err := r.next(c, m)
	if err != nil {
		return r.handler(c, m, err)
	} else {
		return res, nil
	}
}
