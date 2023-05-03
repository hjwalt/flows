package runtime_retry

import (
	"context"

	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/runway/reflect"
)

func SetTryCount(ctx context.Context, trycount int64) context.Context {
	return context.WithValue(ctx, runtime.Context("RetryCount"), trycount)
}

func GetTryCount(ctx context.Context) int64 {
	return reflect.GetInt64(ctx.Value(runtime.Context("RetryCount")))
}
