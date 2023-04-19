package middleware

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"github.com/pubgolf/pubgolf/api/internal/lib/config"
)

const tokenHeader = "X-PubGolf-Auth"

// NewAdminAuthInterceptor guards against requests which don't contain a valid auth token.
func NewAdminAuthInterceptor(cfg *config.App) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Header().Get(tokenHeader) == "" {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("no auth token provided"),
				)
			}

			if req.Header().Get(tokenHeader) != cfg.AdminAuth.AdminServiceToken {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("invalid auth token"),
				)
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
