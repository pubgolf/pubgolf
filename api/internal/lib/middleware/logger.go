package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

const (
	// eventKeyHeader is the header name to check for client-provided context on the currently selected event.
	eventKeyHeader = "X-Pubgolf-Eventkey"
	// playerIDHeader is the header name to check for client-provided context on the requesting player.
	playerIDHeader = "X-Pubgolf-Playerid"
	// deviceIDHeader is the header name to check for client-provided context on the requesting device.
	deviceIDHeader = "X-Pubgolf-Deviceid"
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

			ek := req.Header().Get(eventKeyHeader)
			if ek != "" {
				telemetry.AddRecursiveAttribute(&ctx, "client.ctx_header.event_key", ek)
			}

			pID := req.Header().Get(playerIDHeader)
			if pID != "" {
				telemetry.AddRecursiveAttribute(&ctx, "client.ctx_header.player_id", pID)
			}

			dID := req.Header().Get(deviceIDHeader)
			if dID != "" {
				telemetry.AddRecursiveAttribute(&ctx, "client.ctx_header.device_id", dID)
			}

			if err != nil {
				if connectErr := new(connect.Error); errors.As(err, &connectErr) {
					span.SetAttributes(attribute.String("error.message", connectErr.Message()))

					log.Printf("%s.%s(%s) completed in %q with code: \"%d %s\" error: %q\n", service, method, args, duration, connectErr.Code(), connectErr.Code(), connectErr.Message())
				} else {
					span.SetAttributes(attribute.String("error.message", err.Error()))

					log.Printf("%s.%s(%s) completed in %q with an unhandled error: %q\n", service, method, args, duration, err)
				}
			} else {
				log.Printf("%s.%s(%s) completed in %q successfully\n", service, method, args, duration)
			}

			return resp, err
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
