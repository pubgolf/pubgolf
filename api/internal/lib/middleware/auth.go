package middleware

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
)

const tokenHeader = "X-PubGolf-Auth" //nolint:gosec

var errMissingAuthToken = errors.New("no auth token provided")

var errInvalidAuthToken = errors.New("invalid auth token")

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
					errMissingAuthToken,
				)
			}

			if req.Header().Get(tokenHeader) != cfg.AdminAuth.AdminServiceToken {
				return nil, connect.NewError(
					connect.CodePermissionDenied,
					errInvalidAuthToken,
				)
			}

			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
