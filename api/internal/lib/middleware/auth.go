package middleware

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

type ctxKeyPlayerID struct{}

// tokenHeader is the header name to check for auth credentials.
const tokenHeader = "X-PubGolf-AuthToken" //nolint:gosec

var (
	errMissingAuthToken = errors.New("no auth token provided")
	errInvalidAuthToken = errors.New("invalid auth token")

	// nonAuthRPCs contains the public RPC endpoints which bypass the token auth.
	nonAuthRPCs = initRouteSet("api.v1.PubGolfService", []string{
		"ClientVersion",
		"StartPlayerLogin",
		"CompletePlayerLogin",
	})
)

// initRouteSet takes a list of RPC method names and a service name, returning a set for faster lookup.
func initRouteSet(service string, names []string) map[string]struct{} {
	m := make(map[string]struct{}, len(names))
	for _, n := range names {
		m[fmt.Sprintf("/%s/%s", service, n)] = struct{}{}
	}

	return m
}

// nonAuthRoute returns true if the route should be skipped for token auth checks.
func nonAuthRoute(req connect.AnyRequest) bool {
	_, match := nonAuthRPCs[req.Spec().Procedure]

	return match
}

// ContextWithPlayerID adds a playerID to a context. Used for testing RPCs without utilizing the full middleware.
func ContextWithPlayerID(ctx context.Context, playerID models.PlayerID) context.Context {
	return context.WithValue(ctx, ctxKeyPlayerID{}, playerID)
}

// PlayerID returns the playerID inferred from the request's auth token.
func PlayerID(ctx context.Context) (models.PlayerID, bool) {
	pID := ctx.Value(ctxKeyPlayerID{})
	playerID, ok := pID.(models.PlayerID)

	return playerID, ok
}

// NewAuthInterceptor checks for a valid auth token and adds the corresponding player ID to the context.
func NewAuthInterceptor(q dao.QueryProvider) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if nonAuthRoute(req) {
				return next(ctx, req)
			}

			t := req.Header().Get(tokenHeader)
			if t == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errMissingAuthToken)
			}

			token, err := models.AuthTokenFromString(t)
			if err != nil {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid auth token format: %w", err))
			}

			playerID, err := q.PlayerIDByAuthToken(ctx, token)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, connect.NewError(connect.CodePermissionDenied, errInvalidAuthToken)
				}

				return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("lookup auth token: %w", err))
			}

			return next(ContextWithPlayerID(ctx, playerID), req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}

// NewAdminAuthInterceptor guards against admin requests which don't contain a valid auth token.
func NewAdminAuthInterceptor(cfg *config.App) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if req.Header().Get(tokenHeader) == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errMissingAuthToken)
			}

			if req.Header().Get(tokenHeader) != cfg.AdminAuth.AdminServiceToken {
				return nil, connect.NewError(connect.CodePermissionDenied, errInvalidAuthToken)
			}

			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
