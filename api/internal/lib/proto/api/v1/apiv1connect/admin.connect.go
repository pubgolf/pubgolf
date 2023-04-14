// Admin defines the admin API service for the game management UI.

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: api/v1/admin.proto

package apiv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// AdminServiceName is the fully-qualified name of the AdminService service.
	AdminServiceName = "api.v1.AdminService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AdminServiceCreatePlayerProcedure is the fully-qualified name of the AdminService's CreatePlayer
	// RPC.
	AdminServiceCreatePlayerProcedure = "/api.v1.AdminService/CreatePlayer"
	// AdminServiceListPlayersProcedure is the fully-qualified name of the AdminService's ListPlayers
	// RPC.
	AdminServiceListPlayersProcedure = "/api.v1.AdminService/ListPlayers"
)

// AdminServiceClient is a client for the api.v1.AdminService service.
type AdminServiceClient interface {
	// CreatePlayer creates a new player profile for a given event.
	CreatePlayer(context.Context, *connect_go.Request[v1.CreatePlayerRequest]) (*connect_go.Response[v1.CreatePlayerResponse], error)
	// ListPlayers returns all players for a given event.
	ListPlayers(context.Context, *connect_go.Request[v1.ListPlayersRequest]) (*connect_go.Response[v1.ListPlayersResponse], error)
}

// NewAdminServiceClient constructs a client for the api.v1.AdminService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAdminServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) AdminServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &adminServiceClient{
		createPlayer: connect_go.NewClient[v1.CreatePlayerRequest, v1.CreatePlayerResponse](
			httpClient,
			baseURL+AdminServiceCreatePlayerProcedure,
			opts...,
		),
		listPlayers: connect_go.NewClient[v1.ListPlayersRequest, v1.ListPlayersResponse](
			httpClient,
			baseURL+AdminServiceListPlayersProcedure,
			opts...,
		),
	}
}

// adminServiceClient implements AdminServiceClient.
type adminServiceClient struct {
	createPlayer *connect_go.Client[v1.CreatePlayerRequest, v1.CreatePlayerResponse]
	listPlayers  *connect_go.Client[v1.ListPlayersRequest, v1.ListPlayersResponse]
}

// CreatePlayer calls api.v1.AdminService.CreatePlayer.
func (c *adminServiceClient) CreatePlayer(ctx context.Context, req *connect_go.Request[v1.CreatePlayerRequest]) (*connect_go.Response[v1.CreatePlayerResponse], error) {
	return c.createPlayer.CallUnary(ctx, req)
}

// ListPlayers calls api.v1.AdminService.ListPlayers.
func (c *adminServiceClient) ListPlayers(ctx context.Context, req *connect_go.Request[v1.ListPlayersRequest]) (*connect_go.Response[v1.ListPlayersResponse], error) {
	return c.listPlayers.CallUnary(ctx, req)
}

// AdminServiceHandler is an implementation of the api.v1.AdminService service.
type AdminServiceHandler interface {
	// CreatePlayer creates a new player profile for a given event.
	CreatePlayer(context.Context, *connect_go.Request[v1.CreatePlayerRequest]) (*connect_go.Response[v1.CreatePlayerResponse], error)
	// ListPlayers returns all players for a given event.
	ListPlayers(context.Context, *connect_go.Request[v1.ListPlayersRequest]) (*connect_go.Response[v1.ListPlayersResponse], error)
}

// NewAdminServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAdminServiceHandler(svc AdminServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(AdminServiceCreatePlayerProcedure, connect_go.NewUnaryHandler(
		AdminServiceCreatePlayerProcedure,
		svc.CreatePlayer,
		opts...,
	))
	mux.Handle(AdminServiceListPlayersProcedure, connect_go.NewUnaryHandler(
		AdminServiceListPlayersProcedure,
		svc.ListPlayers,
		opts...,
	))
	return "/api.v1.AdminService/", mux
}

// UnimplementedAdminServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAdminServiceHandler struct{}

func (UnimplementedAdminServiceHandler) CreatePlayer(context.Context, *connect_go.Request[v1.CreatePlayerRequest]) (*connect_go.Response[v1.CreatePlayerResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.AdminService.CreatePlayer is not implemented"))
}

func (UnimplementedAdminServiceHandler) ListPlayers(context.Context, *connect_go.Request[v1.ListPlayersRequest]) (*connect_go.Response[v1.ListPlayersResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.AdminService.ListPlayers is not implemented"))
}
