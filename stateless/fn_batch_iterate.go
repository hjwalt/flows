package stateless

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
)

// constructor
func NewBatchIterateFunction(configurations ...runtime.Configuration[*BatchIterateFunction]) BatchFunction {
	batchIterateFunction := &BatchIterateFunction{}
	for _, configuration := range configurations {
		batchIterateFunction = configuration(batchIterateFunction)
	}
	return batchIterateFunction.Apply
}

// configuration
func WithBatchIterateFunctionNextFunction(next SingleFunction) runtime.Configuration[*BatchIterateFunction] {
	return func(pbif *BatchIterateFunction) *BatchIterateFunction {
		pbif.next = next
		return pbif
	}
}

// implementation
type BatchIterateFunction struct {
	next SingleFunction
}

func (r *BatchIterateFunction) Apply(c context.Context, ms []flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
	res := make([]flow.Message[structure.Bytes, structure.Bytes], 0)

	for _, m := range ms {
		if mres, err := r.next(c, m); err != nil {
			return make([]flow.Message[structure.Bytes, structure.Bytes], 0), err
		} else {
			res = append(res, mres...)
		}
	}
	return res, nil
}
