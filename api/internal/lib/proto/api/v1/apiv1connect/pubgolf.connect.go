// PubGolf defines the app-facing API service for the in-game apps.

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: api/v1/pubgolf.proto

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
	// PubGolfServiceName is the fully-qualified name of the PubGolfService service.
	PubGolfServiceName = "api.v1.PubGolfService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// PubGolfServiceClientVersionProcedure is the fully-qualified name of the PubGolfService's
	// ClientVersion RPC.
	PubGolfServiceClientVersionProcedure = "/api.v1.PubGolfService/ClientVersion"
	// PubGolfServiceCreatePlayerProcedure is the fully-qualified name of the PubGolfService's
	// CreatePlayer RPC.
	PubGolfServiceCreatePlayerProcedure = "/api.v1.PubGolfService/CreatePlayer"
	// PubGolfServiceStartPlayerLoginProcedure is the fully-qualified name of the PubGolfService's
	// StartPlayerLogin RPC.
	PubGolfServiceStartPlayerLoginProcedure = "/api.v1.PubGolfService/StartPlayerLogin"
	// PubGolfServiceCompletePlayerLoginProcedure is the fully-qualified name of the PubGolfService's
	// CompletePlayerLogin RPC.
	PubGolfServiceCompletePlayerLoginProcedure = "/api.v1.PubGolfService/CompletePlayerLogin"
	// PubGolfServiceGetMyPlayerProcedure is the fully-qualified name of the PubGolfService's
	// GetMyPlayer RPC.
	PubGolfServiceGetMyPlayerProcedure = "/api.v1.PubGolfService/GetMyPlayer"
	// PubGolfServiceGetScheduleProcedure is the fully-qualified name of the PubGolfService's
	// GetSchedule RPC.
	PubGolfServiceGetScheduleProcedure = "/api.v1.PubGolfService/GetSchedule"
	// PubGolfServiceGetVenueProcedure is the fully-qualified name of the PubGolfService's GetVenue RPC.
	PubGolfServiceGetVenueProcedure = "/api.v1.PubGolfService/GetVenue"
	// PubGolfServiceListContentItemsProcedure is the fully-qualified name of the PubGolfService's
	// ListContentItems RPC.
	PubGolfServiceListContentItemsProcedure = "/api.v1.PubGolfService/ListContentItems"
	// PubGolfServiceGetContentItemProcedure is the fully-qualified name of the PubGolfService's
	// GetContentItem RPC.
	PubGolfServiceGetContentItemProcedure = "/api.v1.PubGolfService/GetContentItem"
	// PubGolfServiceGetPlayerProcedure is the fully-qualified name of the PubGolfService's GetPlayer
	// RPC.
	PubGolfServiceGetPlayerProcedure = "/api.v1.PubGolfService/GetPlayer"
	// PubGolfServiceGetScoresForCategoryProcedure is the fully-qualified name of the PubGolfService's
	// GetScoresForCategory RPC.
	PubGolfServiceGetScoresForCategoryProcedure = "/api.v1.PubGolfService/GetScoresForCategory"
	// PubGolfServiceGetScoresForPlayerProcedure is the fully-qualified name of the PubGolfService's
	// GetScoresForPlayer RPC.
	PubGolfServiceGetScoresForPlayerProcedure = "/api.v1.PubGolfService/GetScoresForPlayer"
	// PubGolfServiceGetScoresForVenueProcedure is the fully-qualified name of the PubGolfService's
	// GetScoresForVenue RPC.
	PubGolfServiceGetScoresForVenueProcedure = "/api.v1.PubGolfService/GetScoresForVenue"
)

// PubGolfServiceClient is a client for the api.v1.PubGolfService service.
type PubGolfServiceClient interface {
	// ClientVersion (unauthenticated) indicates to the server that a client of a given version is attempting to connect, and allows the server to respond with a "soft" or "hard" upgrade notification.
	ClientVersion(context.Context, *connect_go.Request[v1.ClientVersionRequest]) (*connect_go.Response[v1.ClientVersionResponse], error)
	// CreatePlayer creates a new player profile for a given event.
	//
	// Deprecated: Use `StartPlayerLogin` RPC instead.
	//
	// Deprecated: do not use.
	CreatePlayer(context.Context, *connect_go.Request[v1.PubGolfServiceCreatePlayerRequest]) (*connect_go.Response[v1.PubGolfServiceCreatePlayerResponse], error)
	// StartPlayerLogin (unauthenticated) registers the player's contact info if the player doesn't exist, then sends an auth code.
	StartPlayerLogin(context.Context, *connect_go.Request[v1.StartPlayerLoginRequest]) (*connect_go.Response[v1.StartPlayerLoginResponse], error)
	// CompletePlayerLogin (unauthenticated) accepts an auth code and logs in the player, returning the data necessary to bootstrap a player's session in the app.
	CompletePlayerLogin(context.Context, *connect_go.Request[v1.CompletePlayerLoginRequest]) (*connect_go.Response[v1.CompletePlayerLoginResponse], error)
	// GetMyPlayer is an authenticated request that returns the same data as `CompletePlayerLogin()` if the player's auth token is still valid.
	GetMyPlayer(context.Context, *connect_go.Request[v1.GetMyPlayerRequest]) (*connect_go.Response[v1.GetMyPlayerResponse], error)
	// GetSchedule returns the list of visble venues, as well as the next venue transition time. It optionally accepts a data version to allow local caching.
	GetSchedule(context.Context, *connect_go.Request[v1.GetScheduleRequest]) (*connect_go.Response[v1.GetScheduleResponse], error)
	// GetVenue performs a bulk lookup of venue metadata by ID. IDs are scoped to an event key.
	GetVenue(context.Context, *connect_go.Request[v1.GetVenueRequest]) (*connect_go.Response[v1.GetVenueResponse], error)
	// ListContentItems
	ListContentItems(context.Context, *connect_go.Request[v1.ListContentItemsRequest]) (*connect_go.Response[v1.ListContentItemsResponse], error)
	// GetContentItem
	GetContentItem(context.Context, *connect_go.Request[v1.GetContentItemRequest]) (*connect_go.Response[v1.GetContentItemResponse], error)
	// GetPlayer
	GetPlayer(context.Context, *connect_go.Request[v1.GetPlayerRequest]) (*connect_go.Response[v1.GetPlayerResponse], error)
	// GetScoresForCategory
	GetScoresForCategory(context.Context, *connect_go.Request[v1.GetScoresForCategoryRequest]) (*connect_go.Response[v1.GetScoresForCategoryResponse], error)
	// GetScoresForPlayer
	GetScoresForPlayer(context.Context, *connect_go.Request[v1.GetScoresForPlayerRequest]) (*connect_go.Response[v1.GetScoresForPlayerResponse], error)
	// GetScoresForVenue
	GetScoresForVenue(context.Context, *connect_go.Request[v1.GetScoresForVenueRequest]) (*connect_go.Response[v1.GetScoresForVenueResponse], error)
}

// NewPubGolfServiceClient constructs a client for the api.v1.PubGolfService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewPubGolfServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) PubGolfServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &pubGolfServiceClient{
		clientVersion: connect_go.NewClient[v1.ClientVersionRequest, v1.ClientVersionResponse](
			httpClient,
			baseURL+PubGolfServiceClientVersionProcedure,
			opts...,
		),
		createPlayer: connect_go.NewClient[v1.PubGolfServiceCreatePlayerRequest, v1.PubGolfServiceCreatePlayerResponse](
			httpClient,
			baseURL+PubGolfServiceCreatePlayerProcedure,
			opts...,
		),
		startPlayerLogin: connect_go.NewClient[v1.StartPlayerLoginRequest, v1.StartPlayerLoginResponse](
			httpClient,
			baseURL+PubGolfServiceStartPlayerLoginProcedure,
			opts...,
		),
		completePlayerLogin: connect_go.NewClient[v1.CompletePlayerLoginRequest, v1.CompletePlayerLoginResponse](
			httpClient,
			baseURL+PubGolfServiceCompletePlayerLoginProcedure,
			opts...,
		),
		getMyPlayer: connect_go.NewClient[v1.GetMyPlayerRequest, v1.GetMyPlayerResponse](
			httpClient,
			baseURL+PubGolfServiceGetMyPlayerProcedure,
			opts...,
		),
		getSchedule: connect_go.NewClient[v1.GetScheduleRequest, v1.GetScheduleResponse](
			httpClient,
			baseURL+PubGolfServiceGetScheduleProcedure,
			opts...,
		),
		getVenue: connect_go.NewClient[v1.GetVenueRequest, v1.GetVenueResponse](
			httpClient,
			baseURL+PubGolfServiceGetVenueProcedure,
			opts...,
		),
		listContentItems: connect_go.NewClient[v1.ListContentItemsRequest, v1.ListContentItemsResponse](
			httpClient,
			baseURL+PubGolfServiceListContentItemsProcedure,
			opts...,
		),
		getContentItem: connect_go.NewClient[v1.GetContentItemRequest, v1.GetContentItemResponse](
			httpClient,
			baseURL+PubGolfServiceGetContentItemProcedure,
			opts...,
		),
		getPlayer: connect_go.NewClient[v1.GetPlayerRequest, v1.GetPlayerResponse](
			httpClient,
			baseURL+PubGolfServiceGetPlayerProcedure,
			opts...,
		),
		getScoresForCategory: connect_go.NewClient[v1.GetScoresForCategoryRequest, v1.GetScoresForCategoryResponse](
			httpClient,
			baseURL+PubGolfServiceGetScoresForCategoryProcedure,
			opts...,
		),
		getScoresForPlayer: connect_go.NewClient[v1.GetScoresForPlayerRequest, v1.GetScoresForPlayerResponse](
			httpClient,
			baseURL+PubGolfServiceGetScoresForPlayerProcedure,
			opts...,
		),
		getScoresForVenue: connect_go.NewClient[v1.GetScoresForVenueRequest, v1.GetScoresForVenueResponse](
			httpClient,
			baseURL+PubGolfServiceGetScoresForVenueProcedure,
			opts...,
		),
	}
}

// pubGolfServiceClient implements PubGolfServiceClient.
type pubGolfServiceClient struct {
	clientVersion        *connect_go.Client[v1.ClientVersionRequest, v1.ClientVersionResponse]
	createPlayer         *connect_go.Client[v1.PubGolfServiceCreatePlayerRequest, v1.PubGolfServiceCreatePlayerResponse]
	startPlayerLogin     *connect_go.Client[v1.StartPlayerLoginRequest, v1.StartPlayerLoginResponse]
	completePlayerLogin  *connect_go.Client[v1.CompletePlayerLoginRequest, v1.CompletePlayerLoginResponse]
	getMyPlayer          *connect_go.Client[v1.GetMyPlayerRequest, v1.GetMyPlayerResponse]
	getSchedule          *connect_go.Client[v1.GetScheduleRequest, v1.GetScheduleResponse]
	getVenue             *connect_go.Client[v1.GetVenueRequest, v1.GetVenueResponse]
	listContentItems     *connect_go.Client[v1.ListContentItemsRequest, v1.ListContentItemsResponse]
	getContentItem       *connect_go.Client[v1.GetContentItemRequest, v1.GetContentItemResponse]
	getPlayer            *connect_go.Client[v1.GetPlayerRequest, v1.GetPlayerResponse]
	getScoresForCategory *connect_go.Client[v1.GetScoresForCategoryRequest, v1.GetScoresForCategoryResponse]
	getScoresForPlayer   *connect_go.Client[v1.GetScoresForPlayerRequest, v1.GetScoresForPlayerResponse]
	getScoresForVenue    *connect_go.Client[v1.GetScoresForVenueRequest, v1.GetScoresForVenueResponse]
}

// ClientVersion calls api.v1.PubGolfService.ClientVersion.
func (c *pubGolfServiceClient) ClientVersion(ctx context.Context, req *connect_go.Request[v1.ClientVersionRequest]) (*connect_go.Response[v1.ClientVersionResponse], error) {
	return c.clientVersion.CallUnary(ctx, req)
}

// CreatePlayer calls api.v1.PubGolfService.CreatePlayer.
//
// Deprecated: do not use.
func (c *pubGolfServiceClient) CreatePlayer(ctx context.Context, req *connect_go.Request[v1.PubGolfServiceCreatePlayerRequest]) (*connect_go.Response[v1.PubGolfServiceCreatePlayerResponse], error) {
	return c.createPlayer.CallUnary(ctx, req)
}

// StartPlayerLogin calls api.v1.PubGolfService.StartPlayerLogin.
func (c *pubGolfServiceClient) StartPlayerLogin(ctx context.Context, req *connect_go.Request[v1.StartPlayerLoginRequest]) (*connect_go.Response[v1.StartPlayerLoginResponse], error) {
	return c.startPlayerLogin.CallUnary(ctx, req)
}

// CompletePlayerLogin calls api.v1.PubGolfService.CompletePlayerLogin.
func (c *pubGolfServiceClient) CompletePlayerLogin(ctx context.Context, req *connect_go.Request[v1.CompletePlayerLoginRequest]) (*connect_go.Response[v1.CompletePlayerLoginResponse], error) {
	return c.completePlayerLogin.CallUnary(ctx, req)
}

// GetMyPlayer calls api.v1.PubGolfService.GetMyPlayer.
func (c *pubGolfServiceClient) GetMyPlayer(ctx context.Context, req *connect_go.Request[v1.GetMyPlayerRequest]) (*connect_go.Response[v1.GetMyPlayerResponse], error) {
	return c.getMyPlayer.CallUnary(ctx, req)
}

// GetSchedule calls api.v1.PubGolfService.GetSchedule.
func (c *pubGolfServiceClient) GetSchedule(ctx context.Context, req *connect_go.Request[v1.GetScheduleRequest]) (*connect_go.Response[v1.GetScheduleResponse], error) {
	return c.getSchedule.CallUnary(ctx, req)
}

// GetVenue calls api.v1.PubGolfService.GetVenue.
func (c *pubGolfServiceClient) GetVenue(ctx context.Context, req *connect_go.Request[v1.GetVenueRequest]) (*connect_go.Response[v1.GetVenueResponse], error) {
	return c.getVenue.CallUnary(ctx, req)
}

// ListContentItems calls api.v1.PubGolfService.ListContentItems.
func (c *pubGolfServiceClient) ListContentItems(ctx context.Context, req *connect_go.Request[v1.ListContentItemsRequest]) (*connect_go.Response[v1.ListContentItemsResponse], error) {
	return c.listContentItems.CallUnary(ctx, req)
}

// GetContentItem calls api.v1.PubGolfService.GetContentItem.
func (c *pubGolfServiceClient) GetContentItem(ctx context.Context, req *connect_go.Request[v1.GetContentItemRequest]) (*connect_go.Response[v1.GetContentItemResponse], error) {
	return c.getContentItem.CallUnary(ctx, req)
}

// GetPlayer calls api.v1.PubGolfService.GetPlayer.
func (c *pubGolfServiceClient) GetPlayer(ctx context.Context, req *connect_go.Request[v1.GetPlayerRequest]) (*connect_go.Response[v1.GetPlayerResponse], error) {
	return c.getPlayer.CallUnary(ctx, req)
}

// GetScoresForCategory calls api.v1.PubGolfService.GetScoresForCategory.
func (c *pubGolfServiceClient) GetScoresForCategory(ctx context.Context, req *connect_go.Request[v1.GetScoresForCategoryRequest]) (*connect_go.Response[v1.GetScoresForCategoryResponse], error) {
	return c.getScoresForCategory.CallUnary(ctx, req)
}

// GetScoresForPlayer calls api.v1.PubGolfService.GetScoresForPlayer.
func (c *pubGolfServiceClient) GetScoresForPlayer(ctx context.Context, req *connect_go.Request[v1.GetScoresForPlayerRequest]) (*connect_go.Response[v1.GetScoresForPlayerResponse], error) {
	return c.getScoresForPlayer.CallUnary(ctx, req)
}

// GetScoresForVenue calls api.v1.PubGolfService.GetScoresForVenue.
func (c *pubGolfServiceClient) GetScoresForVenue(ctx context.Context, req *connect_go.Request[v1.GetScoresForVenueRequest]) (*connect_go.Response[v1.GetScoresForVenueResponse], error) {
	return c.getScoresForVenue.CallUnary(ctx, req)
}

// PubGolfServiceHandler is an implementation of the api.v1.PubGolfService service.
type PubGolfServiceHandler interface {
	// ClientVersion (unauthenticated) indicates to the server that a client of a given version is attempting to connect, and allows the server to respond with a "soft" or "hard" upgrade notification.
	ClientVersion(context.Context, *connect_go.Request[v1.ClientVersionRequest]) (*connect_go.Response[v1.ClientVersionResponse], error)
	// CreatePlayer creates a new player profile for a given event.
	//
	// Deprecated: Use `StartPlayerLogin` RPC instead.
	//
	// Deprecated: do not use.
	CreatePlayer(context.Context, *connect_go.Request[v1.PubGolfServiceCreatePlayerRequest]) (*connect_go.Response[v1.PubGolfServiceCreatePlayerResponse], error)
	// StartPlayerLogin (unauthenticated) registers the player's contact info if the player doesn't exist, then sends an auth code.
	StartPlayerLogin(context.Context, *connect_go.Request[v1.StartPlayerLoginRequest]) (*connect_go.Response[v1.StartPlayerLoginResponse], error)
	// CompletePlayerLogin (unauthenticated) accepts an auth code and logs in the player, returning the data necessary to bootstrap a player's session in the app.
	CompletePlayerLogin(context.Context, *connect_go.Request[v1.CompletePlayerLoginRequest]) (*connect_go.Response[v1.CompletePlayerLoginResponse], error)
	// GetMyPlayer is an authenticated request that returns the same data as `CompletePlayerLogin()` if the player's auth token is still valid.
	GetMyPlayer(context.Context, *connect_go.Request[v1.GetMyPlayerRequest]) (*connect_go.Response[v1.GetMyPlayerResponse], error)
	// GetSchedule returns the list of visble venues, as well as the next venue transition time. It optionally accepts a data version to allow local caching.
	GetSchedule(context.Context, *connect_go.Request[v1.GetScheduleRequest]) (*connect_go.Response[v1.GetScheduleResponse], error)
	// GetVenue performs a bulk lookup of venue metadata by ID. IDs are scoped to an event key.
	GetVenue(context.Context, *connect_go.Request[v1.GetVenueRequest]) (*connect_go.Response[v1.GetVenueResponse], error)
	// ListContentItems
	ListContentItems(context.Context, *connect_go.Request[v1.ListContentItemsRequest]) (*connect_go.Response[v1.ListContentItemsResponse], error)
	// GetContentItem
	GetContentItem(context.Context, *connect_go.Request[v1.GetContentItemRequest]) (*connect_go.Response[v1.GetContentItemResponse], error)
	// GetPlayer
	GetPlayer(context.Context, *connect_go.Request[v1.GetPlayerRequest]) (*connect_go.Response[v1.GetPlayerResponse], error)
	// GetScoresForCategory
	GetScoresForCategory(context.Context, *connect_go.Request[v1.GetScoresForCategoryRequest]) (*connect_go.Response[v1.GetScoresForCategoryResponse], error)
	// GetScoresForPlayer
	GetScoresForPlayer(context.Context, *connect_go.Request[v1.GetScoresForPlayerRequest]) (*connect_go.Response[v1.GetScoresForPlayerResponse], error)
	// GetScoresForVenue
	GetScoresForVenue(context.Context, *connect_go.Request[v1.GetScoresForVenueRequest]) (*connect_go.Response[v1.GetScoresForVenueResponse], error)
}

// NewPubGolfServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewPubGolfServiceHandler(svc PubGolfServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	pubGolfServiceClientVersionHandler := connect_go.NewUnaryHandler(
		PubGolfServiceClientVersionProcedure,
		svc.ClientVersion,
		opts...,
	)
	pubGolfServiceCreatePlayerHandler := connect_go.NewUnaryHandler(
		PubGolfServiceCreatePlayerProcedure,
		svc.CreatePlayer,
		opts...,
	)
	pubGolfServiceStartPlayerLoginHandler := connect_go.NewUnaryHandler(
		PubGolfServiceStartPlayerLoginProcedure,
		svc.StartPlayerLogin,
		opts...,
	)
	pubGolfServiceCompletePlayerLoginHandler := connect_go.NewUnaryHandler(
		PubGolfServiceCompletePlayerLoginProcedure,
		svc.CompletePlayerLogin,
		opts...,
	)
	pubGolfServiceGetMyPlayerHandler := connect_go.NewUnaryHandler(
		PubGolfServiceGetMyPlayerProcedure,
		svc.GetMyPlayer,
		opts...,
	)
	pubGolfServiceGetScheduleHandler := connect_go.NewUnaryHandler(
		PubGolfServiceGetScheduleProcedure,
		svc.GetSchedule,
		opts...,
	)
	pubGolfServiceGetVenueHandler := connect_go.NewUnaryHandler(
		PubGolfServiceGetVenueProcedure,
		svc.GetVenue,
		opts...,
	)
	pubGolfServiceListContentItemsHandler := connect_go.NewUnaryHandler(
		PubGolfServiceListContentItemsProcedure,
		svc.ListContentItems,
		opts...,
	)
	pubGolfServiceGetContentItemHandler := connect_go.NewUnaryHandler(
		PubGolfServiceGetContentItemProcedure,
		svc.GetContentItem,
		opts...,
	)
	pubGolfServiceGetPlayerHandler := connect_go.NewUnaryHandler(
		PubGolfServiceGetPlayerProcedure,
		svc.GetPlayer,
		opts...,
	)
	pubGolfServiceGetScoresForCategoryHandler := connect_go.NewUnaryHandler(
		PubGolfServiceGetScoresForCategoryProcedure,
		svc.GetScoresForCategory,
		opts...,
	)
	pubGolfServiceGetScoresForPlayerHandler := connect_go.NewUnaryHandler(
		PubGolfServiceGetScoresForPlayerProcedure,
		svc.GetScoresForPlayer,
		opts...,
	)
	pubGolfServiceGetScoresForVenueHandler := connect_go.NewUnaryHandler(
		PubGolfServiceGetScoresForVenueProcedure,
		svc.GetScoresForVenue,
		opts...,
	)
	return "/api.v1.PubGolfService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PubGolfServiceClientVersionProcedure:
			pubGolfServiceClientVersionHandler.ServeHTTP(w, r)
		case PubGolfServiceCreatePlayerProcedure:
			pubGolfServiceCreatePlayerHandler.ServeHTTP(w, r)
		case PubGolfServiceStartPlayerLoginProcedure:
			pubGolfServiceStartPlayerLoginHandler.ServeHTTP(w, r)
		case PubGolfServiceCompletePlayerLoginProcedure:
			pubGolfServiceCompletePlayerLoginHandler.ServeHTTP(w, r)
		case PubGolfServiceGetMyPlayerProcedure:
			pubGolfServiceGetMyPlayerHandler.ServeHTTP(w, r)
		case PubGolfServiceGetScheduleProcedure:
			pubGolfServiceGetScheduleHandler.ServeHTTP(w, r)
		case PubGolfServiceGetVenueProcedure:
			pubGolfServiceGetVenueHandler.ServeHTTP(w, r)
		case PubGolfServiceListContentItemsProcedure:
			pubGolfServiceListContentItemsHandler.ServeHTTP(w, r)
		case PubGolfServiceGetContentItemProcedure:
			pubGolfServiceGetContentItemHandler.ServeHTTP(w, r)
		case PubGolfServiceGetPlayerProcedure:
			pubGolfServiceGetPlayerHandler.ServeHTTP(w, r)
		case PubGolfServiceGetScoresForCategoryProcedure:
			pubGolfServiceGetScoresForCategoryHandler.ServeHTTP(w, r)
		case PubGolfServiceGetScoresForPlayerProcedure:
			pubGolfServiceGetScoresForPlayerHandler.ServeHTTP(w, r)
		case PubGolfServiceGetScoresForVenueProcedure:
			pubGolfServiceGetScoresForVenueHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedPubGolfServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedPubGolfServiceHandler struct{}

func (UnimplementedPubGolfServiceHandler) ClientVersion(context.Context, *connect_go.Request[v1.ClientVersionRequest]) (*connect_go.Response[v1.ClientVersionResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.ClientVersion is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) CreatePlayer(context.Context, *connect_go.Request[v1.PubGolfServiceCreatePlayerRequest]) (*connect_go.Response[v1.PubGolfServiceCreatePlayerResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.CreatePlayer is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) StartPlayerLogin(context.Context, *connect_go.Request[v1.StartPlayerLoginRequest]) (*connect_go.Response[v1.StartPlayerLoginResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.StartPlayerLogin is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) CompletePlayerLogin(context.Context, *connect_go.Request[v1.CompletePlayerLoginRequest]) (*connect_go.Response[v1.CompletePlayerLoginResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.CompletePlayerLogin is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) GetMyPlayer(context.Context, *connect_go.Request[v1.GetMyPlayerRequest]) (*connect_go.Response[v1.GetMyPlayerResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.GetMyPlayer is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) GetSchedule(context.Context, *connect_go.Request[v1.GetScheduleRequest]) (*connect_go.Response[v1.GetScheduleResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.GetSchedule is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) GetVenue(context.Context, *connect_go.Request[v1.GetVenueRequest]) (*connect_go.Response[v1.GetVenueResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.GetVenue is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) ListContentItems(context.Context, *connect_go.Request[v1.ListContentItemsRequest]) (*connect_go.Response[v1.ListContentItemsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.ListContentItems is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) GetContentItem(context.Context, *connect_go.Request[v1.GetContentItemRequest]) (*connect_go.Response[v1.GetContentItemResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.GetContentItem is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) GetPlayer(context.Context, *connect_go.Request[v1.GetPlayerRequest]) (*connect_go.Response[v1.GetPlayerResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.GetPlayer is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) GetScoresForCategory(context.Context, *connect_go.Request[v1.GetScoresForCategoryRequest]) (*connect_go.Response[v1.GetScoresForCategoryResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.GetScoresForCategory is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) GetScoresForPlayer(context.Context, *connect_go.Request[v1.GetScoresForPlayerRequest]) (*connect_go.Response[v1.GetScoresForPlayerResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.GetScoresForPlayer is not implemented"))
}

func (UnimplementedPubGolfServiceHandler) GetScoresForVenue(context.Context, *connect_go.Request[v1.GetScoresForVenueRequest]) (*connect_go.Response[v1.GetScoresForVenueResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("api.v1.PubGolfService.GetScoresForVenue is not implemented"))
}
