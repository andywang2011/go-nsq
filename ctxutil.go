package nsq

import "context"

func ctxOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

type buildResult[T any] struct {
	v   T
	err error
}

// buildWithCtx runs a blocking constructor (go-nsq dials and queries lookupd
// without a context) and returns ctx.Err() as soon as ctx is cancelled. A value
// that completes after cancellation is released through cleanup so nothing leaks.
func buildWithCtx[T any](ctx context.Context, build func() (T, error), cleanup func(T)) (T, error) {
	var zero T
	if err := ctx.Err(); err != nil {
		return zero, err
	}
	ch := make(chan buildResult[T], 1)
	go func() {
		v, err := build()
		ch <- buildResult[T]{v: v, err: err}
	}()
	select {
	case r := <-ch:
		return r.v, r.err
	case <-ctx.Done():
		go func() {
			if r := <-ch; r.err == nil {
				cleanup(r.v)
			}
		}()
		return zero, ctx.Err()
	}
}
