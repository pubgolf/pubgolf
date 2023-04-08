package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
)

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

			if connectErr := new(connect.Error); errors.As(err, &connectErr) {
				log.Printf("%s.%s(%s) completed in %q with code: \"%d %s\" error: %q\n", service, method, args, duration, connectErr.Code(), connectErr.Code(), connectErr.Message())
			} else {
				log.Printf("%s.%s(%s) completed in %q successfully\n", service, method, args, duration)
			}

			return resp, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
