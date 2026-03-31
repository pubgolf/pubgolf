package middleware

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

// NewContextErrorInterceptor remaps Connect error codes when the underlying
// error wraps a context cancellation or deadline exceeded error. This allows
// handlers and guards to uniformly wrap errors as CodeUnavailable without
// needing to special-case context errors — the interceptor corrects the code
// to CodeCanceled or CodeDeadlineExceeded before the response reaches OTel
// and logging interceptors.
func NewContextErrorInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if err == nil {
				return resp, nil
			}

			var connectErr *connect.Error
			if !errors.As(err, &connectErr) {
				return resp, err
			}

			// Check DeadlineExceeded first: a deadline expiration also cancels
			// the context, so both context.Canceled and context.DeadlineExceeded
			// may be present.
			if errors.Is(err, context.DeadlineExceeded) {
				return nil, connect.NewError(connect.CodeDeadlineExceeded, connectErr.Unwrap())
			}

			if errors.Is(err, context.Canceled) {
				return nil, connect.NewError(connect.CodeCanceled, connectErr.Unwrap())
			}

			return resp, err
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
