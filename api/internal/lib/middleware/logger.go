package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NewLoggingInterceptor logs (and annotates OTel spans for) all gRPC calls.
func NewLoggingInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			resp, err := next(ctx, req)
			duration := time.Since(start)

			parts := strings.Split(req.Spec().Procedure, "/")
			service := parts[1]
			method := parts[2]
			args, _ := json.Marshal(req.Any())

			span := trace.SpanFromContext(ctx)
			span.SetAttributes(attribute.String("rpc.args", string(args)))

			if connectErr := new(connect.Error); errors.As(err, &connectErr) {
				span.SetAttributes(attribute.String("error.message", connectErr.Message()))

				log.Printf("%s.%s(%s) completed in %q with code: \"%d %s\" error: %q\n", service, method, args, duration, connectErr.Code(), connectErr.Code(), connectErr.Message())
			} else {
				log.Printf("%s.%s(%s) completed in %q successfully\n", service, method, args, duration)
			}

			return resp, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
