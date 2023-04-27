package runtime_retry

import (
	"context"

	"github.com/hjwalt/runway/reflect"
)

type retryContextKey struct {
	key string
}

func countContextKey() any {
	return retryContextKey{key: "retry count"}
}

func SetTryCount(ctx context.Context, trycount int64) context.Context {
	return context.WithValue(ctx, countContextKey(), trycount)
}

func GetTryCount(ctx context.Context) int64 {
	return reflect.GetInt64(ctx.Value(countContextKey()))
}
