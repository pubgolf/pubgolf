// Admin defines the admin API service for the game management UI.

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: api/v1/admin.proto

package apiv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

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
	// AdminServicePurgeAllCachesProcedure is the fully-qualified name of the AdminService's
	// PurgeAllCaches RPC.
	AdminServicePurgeAllCachesProcedure = "/api.v1.AdminService/PurgeAllCaches"
	// AdminServiceCreatePlayerProcedure is the fully-qualified name of the AdminService's CreatePlayer
	// RPC.
	AdminServiceCreatePlayerProcedure = "/api.v1.AdminService/CreatePlayer"
	// AdminServiceUpdatePlayerProcedure is the fully-qualified name of the AdminService's UpdatePlayer
	// RPC.
	AdminServiceUpdatePlayerProcedure = "/api.v1.AdminService/UpdatePlayer"
	// AdminServiceListPlayersProcedure is the fully-qualified name of the AdminService's ListPlayers
	// RPC.
	AdminServiceListPlayersProcedure = "/api.v1.AdminService/ListPlayers"
	// AdminServiceListVenuesProcedure is the fully-qualified name of the AdminService's ListVenues RPC.
	AdminServiceListVenuesProcedure = "/api.v1.AdminService/ListVenues"
	// AdminServiceUpdateStageProcedure is the fully-qualified name of the AdminService's UpdateStage
	// RPC.
	AdminServiceUpdateStageProcedure = "/api.v1.AdminService/UpdateStage"
	// AdminServiceListEventStagesProcedure is the fully-qualified name of the AdminService's
	// ListEventStages RPC.
	AdminServiceListEventStagesProcedure = "/api.v1.AdminService/ListEventStages"
	// AdminServiceCreateAdjustmentTemplateProcedure is the fully-qualified name of the AdminService's
	// CreateAdjustmentTemplate RPC.
	AdminServiceCreateAdjustmentTemplateProcedure = "/api.v1.AdminService/CreateAdjustmentTemplate"
	// AdminServiceUpdateAdjustmentTemplateProcedure is the fully-qualified name of the AdminService's
	// UpdateAdjustmentTemplate RPC.
	AdminServiceUpdateAdjustmentTemplateProcedure = "/api.v1.AdminService/UpdateAdjustmentTemplate"
	// AdminServiceListAdjustmentTemplatesProcedure is the fully-qualified name of the AdminService's
	// ListAdjustmentTemplates RPC.
	AdminServiceListAdjustmentTemplatesProcedure = "/api.v1.AdminService/ListAdjustmentTemplates"
	// AdminServiceCreateStageScoreProcedure is the fully-qualified name of the AdminService's
	// CreateStageScore RPC.
	AdminServiceCreateStageScoreProcedure = "/api.v1.AdminService/CreateStageScore"
	// AdminServiceUpdateStageScoreProcedure is the fully-qualified name of the AdminService's
	// UpdateStageScore RPC.
	AdminServiceUpdateStageScoreProcedure = "/api.v1.AdminService/UpdateStageScore"
	// AdminServiceListStageScoresProcedure is the fully-qualified name of the AdminService's
	// ListStageScores RPC.
	AdminServiceListStageScoresProcedure = "/api.v1.AdminService/ListStageScores"
	// AdminServiceDeleteStageScoreProcedure is the fully-qualified name of the AdminService's
	// DeleteStageScore RPC.
	AdminServiceDeleteStageScoreProcedure = "/api.v1.AdminService/DeleteStageScore"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	adminServiceServiceDescriptor                        = v1.File_api_v1_admin_proto.Services().ByName("AdminService")
	adminServicePurgeAllCachesMethodDescriptor           = adminServiceServiceDescriptor.Methods().ByName("PurgeAllCaches")
	adminServiceCreatePlayerMethodDescriptor             = adminServiceServiceDescriptor.Methods().ByName("CreatePlayer")
	adminServiceUpdatePlayerMethodDescriptor             = adminServiceServiceDescriptor.Methods().ByName("UpdatePlayer")
	adminServiceListPlayersMethodDescriptor              = adminServiceServiceDescriptor.Methods().ByName("ListPlayers")
	adminServiceListVenuesMethodDescriptor               = adminServiceServiceDescriptor.Methods().ByName("ListVenues")
	adminServiceUpdateStageMethodDescriptor              = adminServiceServiceDescriptor.Methods().ByName("UpdateStage")
	adminServiceListEventStagesMethodDescriptor          = adminServiceServiceDescriptor.Methods().ByName("ListEventStages")
	adminServiceCreateAdjustmentTemplateMethodDescriptor = adminServiceServiceDescriptor.Methods().ByName("CreateAdjustmentTemplate")
	adminServiceUpdateAdjustmentTemplateMethodDescriptor = adminServiceServiceDescriptor.Methods().ByName("UpdateAdjustmentTemplate")
	adminServiceListAdjustmentTemplatesMethodDescriptor  = adminServiceServiceDescriptor.Methods().ByName("ListAdjustmentTemplates")
	adminServiceCreateStageScoreMethodDescriptor         = adminServiceServiceDescriptor.Methods().ByName("CreateStageScore")
	adminServiceUpdateStageScoreMethodDescriptor         = adminServiceServiceDescriptor.Methods().ByName("UpdateStageScore")
	adminServiceListStageScoresMethodDescriptor          = adminServiceServiceDescriptor.Methods().ByName("ListStageScores")
	adminServiceDeleteStageScoreMethodDescriptor         = adminServiceServiceDescriptor.Methods().ByName("DeleteStageScore")
)

// AdminServiceClient is a client for the api.v1.AdminService service.
type AdminServiceClient interface {
	// PurgeAllCaches clears all caches in the API server. Used for testing.
	PurgeAllCaches(context.Context, *connect.Request[v1.PurgeAllCachesRequest]) (*connect.Response[v1.PurgeAllCachesResponse], error)
	// CreatePlayer creates a new player profile for a given event.
	CreatePlayer(context.Context, *connect.Request[v1.AdminServiceCreatePlayerRequest]) (*connect.Response[v1.AdminServiceCreatePlayerResponse], error)
	// UpdatePlayer modifies the player's profile and settings for a given event.
	UpdatePlayer(context.Context, *connect.Request[v1.UpdatePlayerRequest]) (*connect.Response[v1.UpdatePlayerResponse], error)
	// ListPlayers returns all players for a given event.
	ListPlayers(context.Context, *connect.Request[v1.ListPlayersRequest]) (*connect.Response[v1.ListPlayersResponse], error)
	// ListVenues returns all venues that a stage can link to.
	ListVenues(context.Context, *connect.Request[v1.ListVenuesRequest]) (*connect.Response[v1.ListVenuesResponse], error)
	// UpdateStage sets the details for a stage.
	UpdateStage(context.Context, *connect.Request[v1.UpdateStageRequest]) (*connect.Response[v1.UpdateStageResponse], error)
	// ListEventStages returns a full schedule for an event.
	ListEventStages(context.Context, *connect.Request[v1.ListEventStagesRequest]) (*connect.Response[v1.ListEventStagesResponse], error)
	// CreateAdjustmentTemplate creates an adjustment to surface in player score submission.
	CreateAdjustmentTemplate(context.Context, *connect.Request[v1.CreateAdjustmentTemplateRequest]) (*connect.Response[v1.CreateAdjustmentTemplateResponse], error)
	// CreateAdjustmentTemplate updates an adjustment template.
	UpdateAdjustmentTemplate(context.Context, *connect.Request[v1.UpdateAdjustmentTemplateRequest]) (*connect.Response[v1.UpdateAdjustmentTemplateResponse], error)
	// ListAdjustmentTemplates returns all adjustment templates for a given event.
	ListAdjustmentTemplates(context.Context, *connect.Request[v1.ListAdjustmentTemplatesRequest]) (*connect.Response[v1.ListAdjustmentTemplatesResponse], error)
	// CreateStageScore sets the score and adjustments for a given pair of player and stage IDs.
	CreateStageScore(context.Context, *connect.Request[v1.CreateStageScoreRequest]) (*connect.Response[v1.CreateStageScoreResponse], error)
	// CreateStageScore updates the score and adjustments for a player/stage pair, based on their IDs.
	UpdateStageScore(context.Context, *connect.Request[v1.UpdateStageScoreRequest]) (*connect.Response[v1.UpdateStageScoreResponse], error)
	// ListStageScores returns all sets of (scores, adjustments[]) for an event, ordered chronologically by event stage, then chronologically by score creation time.
	ListStageScores(context.Context, *connect.Request[v1.ListStageScoresRequest]) (*connect.Response[v1.ListStageScoresResponse], error)
	// DeleteStageScore removes all scoring data for a player/stage pair.
	DeleteStageScore(context.Context, *connect.Request[v1.DeleteStageScoreRequest]) (*connect.Response[v1.DeleteStageScoreResponse], error)
}

// NewAdminServiceClient constructs a client for the api.v1.AdminService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAdminServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) AdminServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &adminServiceClient{
		purgeAllCaches: connect.NewClient[v1.PurgeAllCachesRequest, v1.PurgeAllCachesResponse](
			httpClient,
			baseURL+AdminServicePurgeAllCachesProcedure,
			connect.WithSchema(adminServicePurgeAllCachesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createPlayer: connect.NewClient[v1.AdminServiceCreatePlayerRequest, v1.AdminServiceCreatePlayerResponse](
			httpClient,
			baseURL+AdminServiceCreatePlayerProcedure,
			connect.WithSchema(adminServiceCreatePlayerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updatePlayer: connect.NewClient[v1.UpdatePlayerRequest, v1.UpdatePlayerResponse](
			httpClient,
			baseURL+AdminServiceUpdatePlayerProcedure,
			connect.WithSchema(adminServiceUpdatePlayerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listPlayers: connect.NewClient[v1.ListPlayersRequest, v1.ListPlayersResponse](
			httpClient,
			baseURL+AdminServiceListPlayersProcedure,
			connect.WithSchema(adminServiceListPlayersMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listVenues: connect.NewClient[v1.ListVenuesRequest, v1.ListVenuesResponse](
			httpClient,
			baseURL+AdminServiceListVenuesProcedure,
			connect.WithSchema(adminServiceListVenuesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateStage: connect.NewClient[v1.UpdateStageRequest, v1.UpdateStageResponse](
			httpClient,
			baseURL+AdminServiceUpdateStageProcedure,
			connect.WithSchema(adminServiceUpdateStageMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listEventStages: connect.NewClient[v1.ListEventStagesRequest, v1.ListEventStagesResponse](
			httpClient,
			baseURL+AdminServiceListEventStagesProcedure,
			connect.WithSchema(adminServiceListEventStagesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createAdjustmentTemplate: connect.NewClient[v1.CreateAdjustmentTemplateRequest, v1.CreateAdjustmentTemplateResponse](
			httpClient,
			baseURL+AdminServiceCreateAdjustmentTemplateProcedure,
			connect.WithSchema(adminServiceCreateAdjustmentTemplateMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateAdjustmentTemplate: connect.NewClient[v1.UpdateAdjustmentTemplateRequest, v1.UpdateAdjustmentTemplateResponse](
			httpClient,
			baseURL+AdminServiceUpdateAdjustmentTemplateProcedure,
			connect.WithSchema(adminServiceUpdateAdjustmentTemplateMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listAdjustmentTemplates: connect.NewClient[v1.ListAdjustmentTemplatesRequest, v1.ListAdjustmentTemplatesResponse](
			httpClient,
			baseURL+AdminServiceListAdjustmentTemplatesProcedure,
			connect.WithSchema(adminServiceListAdjustmentTemplatesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createStageScore: connect.NewClient[v1.CreateStageScoreRequest, v1.CreateStageScoreResponse](
			httpClient,
			baseURL+AdminServiceCreateStageScoreProcedure,
			connect.WithSchema(adminServiceCreateStageScoreMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateStageScore: connect.NewClient[v1.UpdateStageScoreRequest, v1.UpdateStageScoreResponse](
			httpClient,
			baseURL+AdminServiceUpdateStageScoreProcedure,
			connect.WithSchema(adminServiceUpdateStageScoreMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listStageScores: connect.NewClient[v1.ListStageScoresRequest, v1.ListStageScoresResponse](
			httpClient,
			baseURL+AdminServiceListStageScoresProcedure,
			connect.WithSchema(adminServiceListStageScoresMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deleteStageScore: connect.NewClient[v1.DeleteStageScoreRequest, v1.DeleteStageScoreResponse](
			httpClient,
			baseURL+AdminServiceDeleteStageScoreProcedure,
			connect.WithSchema(adminServiceDeleteStageScoreMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// adminServiceClient implements AdminServiceClient.
type adminServiceClient struct {
	purgeAllCaches           *connect.Client[v1.PurgeAllCachesRequest, v1.PurgeAllCachesResponse]
	createPlayer             *connect.Client[v1.AdminServiceCreatePlayerRequest, v1.AdminServiceCreatePlayerResponse]
	updatePlayer             *connect.Client[v1.UpdatePlayerRequest, v1.UpdatePlayerResponse]
	listPlayers              *connect.Client[v1.ListPlayersRequest, v1.ListPlayersResponse]
	listVenues               *connect.Client[v1.ListVenuesRequest, v1.ListVenuesResponse]
	updateStage              *connect.Client[v1.UpdateStageRequest, v1.UpdateStageResponse]
	listEventStages          *connect.Client[v1.ListEventStagesRequest, v1.ListEventStagesResponse]
	createAdjustmentTemplate *connect.Client[v1.CreateAdjustmentTemplateRequest, v1.CreateAdjustmentTemplateResponse]
	updateAdjustmentTemplate *connect.Client[v1.UpdateAdjustmentTemplateRequest, v1.UpdateAdjustmentTemplateResponse]
	listAdjustmentTemplates  *connect.Client[v1.ListAdjustmentTemplatesRequest, v1.ListAdjustmentTemplatesResponse]
	createStageScore         *connect.Client[v1.CreateStageScoreRequest, v1.CreateStageScoreResponse]
	updateStageScore         *connect.Client[v1.UpdateStageScoreRequest, v1.UpdateStageScoreResponse]
	listStageScores          *connect.Client[v1.ListStageScoresRequest, v1.ListStageScoresResponse]
	deleteStageScore         *connect.Client[v1.DeleteStageScoreRequest, v1.DeleteStageScoreResponse]
}

// PurgeAllCaches calls api.v1.AdminService.PurgeAllCaches.
func (c *adminServiceClient) PurgeAllCaches(ctx context.Context, req *connect.Request[v1.PurgeAllCachesRequest]) (*connect.Response[v1.PurgeAllCachesResponse], error) {
	return c.purgeAllCaches.CallUnary(ctx, req)
}

// CreatePlayer calls api.v1.AdminService.CreatePlayer.
func (c *adminServiceClient) CreatePlayer(ctx context.Context, req *connect.Request[v1.AdminServiceCreatePlayerRequest]) (*connect.Response[v1.AdminServiceCreatePlayerResponse], error) {
	return c.createPlayer.CallUnary(ctx, req)
}

// UpdatePlayer calls api.v1.AdminService.UpdatePlayer.
func (c *adminServiceClient) UpdatePlayer(ctx context.Context, req *connect.Request[v1.UpdatePlayerRequest]) (*connect.Response[v1.UpdatePlayerResponse], error) {
	return c.updatePlayer.CallUnary(ctx, req)
}

// ListPlayers calls api.v1.AdminService.ListPlayers.
func (c *adminServiceClient) ListPlayers(ctx context.Context, req *connect.Request[v1.ListPlayersRequest]) (*connect.Response[v1.ListPlayersResponse], error) {
	return c.listPlayers.CallUnary(ctx, req)
}

// ListVenues calls api.v1.AdminService.ListVenues.
func (c *adminServiceClient) ListVenues(ctx context.Context, req *connect.Request[v1.ListVenuesRequest]) (*connect.Response[v1.ListVenuesResponse], error) {
	return c.listVenues.CallUnary(ctx, req)
}

// UpdateStage calls api.v1.AdminService.UpdateStage.
func (c *adminServiceClient) UpdateStage(ctx context.Context, req *connect.Request[v1.UpdateStageRequest]) (*connect.Response[v1.UpdateStageResponse], error) {
	return c.updateStage.CallUnary(ctx, req)
}

// ListEventStages calls api.v1.AdminService.ListEventStages.
func (c *adminServiceClient) ListEventStages(ctx context.Context, req *connect.Request[v1.ListEventStagesRequest]) (*connect.Response[v1.ListEventStagesResponse], error) {
	return c.listEventStages.CallUnary(ctx, req)
}

// CreateAdjustmentTemplate calls api.v1.AdminService.CreateAdjustmentTemplate.
func (c *adminServiceClient) CreateAdjustmentTemplate(ctx context.Context, req *connect.Request[v1.CreateAdjustmentTemplateRequest]) (*connect.Response[v1.CreateAdjustmentTemplateResponse], error) {
	return c.createAdjustmentTemplate.CallUnary(ctx, req)
}

// UpdateAdjustmentTemplate calls api.v1.AdminService.UpdateAdjustmentTemplate.
func (c *adminServiceClient) UpdateAdjustmentTemplate(ctx context.Context, req *connect.Request[v1.UpdateAdjustmentTemplateRequest]) (*connect.Response[v1.UpdateAdjustmentTemplateResponse], error) {
	return c.updateAdjustmentTemplate.CallUnary(ctx, req)
}

// ListAdjustmentTemplates calls api.v1.AdminService.ListAdjustmentTemplates.
func (c *adminServiceClient) ListAdjustmentTemplates(ctx context.Context, req *connect.Request[v1.ListAdjustmentTemplatesRequest]) (*connect.Response[v1.ListAdjustmentTemplatesResponse], error) {
	return c.listAdjustmentTemplates.CallUnary(ctx, req)
}

// CreateStageScore calls api.v1.AdminService.CreateStageScore.
func (c *adminServiceClient) CreateStageScore(ctx context.Context, req *connect.Request[v1.CreateStageScoreRequest]) (*connect.Response[v1.CreateStageScoreResponse], error) {
	return c.createStageScore.CallUnary(ctx, req)
}

// UpdateStageScore calls api.v1.AdminService.UpdateStageScore.
func (c *adminServiceClient) UpdateStageScore(ctx context.Context, req *connect.Request[v1.UpdateStageScoreRequest]) (*connect.Response[v1.UpdateStageScoreResponse], error) {
	return c.updateStageScore.CallUnary(ctx, req)
}

// ListStageScores calls api.v1.AdminService.ListStageScores.
func (c *adminServiceClient) ListStageScores(ctx context.Context, req *connect.Request[v1.ListStageScoresRequest]) (*connect.Response[v1.ListStageScoresResponse], error) {
	return c.listStageScores.CallUnary(ctx, req)
}

// DeleteStageScore calls api.v1.AdminService.DeleteStageScore.
func (c *adminServiceClient) DeleteStageScore(ctx context.Context, req *connect.Request[v1.DeleteStageScoreRequest]) (*connect.Response[v1.DeleteStageScoreResponse], error) {
	return c.deleteStageScore.CallUnary(ctx, req)
}

// AdminServiceHandler is an implementation of the api.v1.AdminService service.
type AdminServiceHandler interface {
	// PurgeAllCaches clears all caches in the API server. Used for testing.
	PurgeAllCaches(context.Context, *connect.Request[v1.PurgeAllCachesRequest]) (*connect.Response[v1.PurgeAllCachesResponse], error)
	// CreatePlayer creates a new player profile for a given event.
	CreatePlayer(context.Context, *connect.Request[v1.AdminServiceCreatePlayerRequest]) (*connect.Response[v1.AdminServiceCreatePlayerResponse], error)
	// UpdatePlayer modifies the player's profile and settings for a given event.
	UpdatePlayer(context.Context, *connect.Request[v1.UpdatePlayerRequest]) (*connect.Response[v1.UpdatePlayerResponse], error)
	// ListPlayers returns all players for a given event.
	ListPlayers(context.Context, *connect.Request[v1.ListPlayersRequest]) (*connect.Response[v1.ListPlayersResponse], error)
	// ListVenues returns all venues that a stage can link to.
	ListVenues(context.Context, *connect.Request[v1.ListVenuesRequest]) (*connect.Response[v1.ListVenuesResponse], error)
	// UpdateStage sets the details for a stage.
	UpdateStage(context.Context, *connect.Request[v1.UpdateStageRequest]) (*connect.Response[v1.UpdateStageResponse], error)
	// ListEventStages returns a full schedule for an event.
	ListEventStages(context.Context, *connect.Request[v1.ListEventStagesRequest]) (*connect.Response[v1.ListEventStagesResponse], error)
	// CreateAdjustmentTemplate creates an adjustment to surface in player score submission.
	CreateAdjustmentTemplate(context.Context, *connect.Request[v1.CreateAdjustmentTemplateRequest]) (*connect.Response[v1.CreateAdjustmentTemplateResponse], error)
	// CreateAdjustmentTemplate updates an adjustment template.
	UpdateAdjustmentTemplate(context.Context, *connect.Request[v1.UpdateAdjustmentTemplateRequest]) (*connect.Response[v1.UpdateAdjustmentTemplateResponse], error)
	// ListAdjustmentTemplates returns all adjustment templates for a given event.
	ListAdjustmentTemplates(context.Context, *connect.Request[v1.ListAdjustmentTemplatesRequest]) (*connect.Response[v1.ListAdjustmentTemplatesResponse], error)
	// CreateStageScore sets the score and adjustments for a given pair of player and stage IDs.
	CreateStageScore(context.Context, *connect.Request[v1.CreateStageScoreRequest]) (*connect.Response[v1.CreateStageScoreResponse], error)
	// CreateStageScore updates the score and adjustments for a player/stage pair, based on their IDs.
	UpdateStageScore(context.Context, *connect.Request[v1.UpdateStageScoreRequest]) (*connect.Response[v1.UpdateStageScoreResponse], error)
	// ListStageScores returns all sets of (scores, adjustments[]) for an event, ordered chronologically by event stage, then chronologically by score creation time.
	ListStageScores(context.Context, *connect.Request[v1.ListStageScoresRequest]) (*connect.Response[v1.ListStageScoresResponse], error)
	// DeleteStageScore removes all scoring data for a player/stage pair.
	DeleteStageScore(context.Context, *connect.Request[v1.DeleteStageScoreRequest]) (*connect.Response[v1.DeleteStageScoreResponse], error)
}

// NewAdminServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAdminServiceHandler(svc AdminServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	adminServicePurgeAllCachesHandler := connect.NewUnaryHandler(
		AdminServicePurgeAllCachesProcedure,
		svc.PurgeAllCaches,
		connect.WithSchema(adminServicePurgeAllCachesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceCreatePlayerHandler := connect.NewUnaryHandler(
		AdminServiceCreatePlayerProcedure,
		svc.CreatePlayer,
		connect.WithSchema(adminServiceCreatePlayerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceUpdatePlayerHandler := connect.NewUnaryHandler(
		AdminServiceUpdatePlayerProcedure,
		svc.UpdatePlayer,
		connect.WithSchema(adminServiceUpdatePlayerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceListPlayersHandler := connect.NewUnaryHandler(
		AdminServiceListPlayersProcedure,
		svc.ListPlayers,
		connect.WithSchema(adminServiceListPlayersMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceListVenuesHandler := connect.NewUnaryHandler(
		AdminServiceListVenuesProcedure,
		svc.ListVenues,
		connect.WithSchema(adminServiceListVenuesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceUpdateStageHandler := connect.NewUnaryHandler(
		AdminServiceUpdateStageProcedure,
		svc.UpdateStage,
		connect.WithSchema(adminServiceUpdateStageMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceListEventStagesHandler := connect.NewUnaryHandler(
		AdminServiceListEventStagesProcedure,
		svc.ListEventStages,
		connect.WithSchema(adminServiceListEventStagesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceCreateAdjustmentTemplateHandler := connect.NewUnaryHandler(
		AdminServiceCreateAdjustmentTemplateProcedure,
		svc.CreateAdjustmentTemplate,
		connect.WithSchema(adminServiceCreateAdjustmentTemplateMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceUpdateAdjustmentTemplateHandler := connect.NewUnaryHandler(
		AdminServiceUpdateAdjustmentTemplateProcedure,
		svc.UpdateAdjustmentTemplate,
		connect.WithSchema(adminServiceUpdateAdjustmentTemplateMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceListAdjustmentTemplatesHandler := connect.NewUnaryHandler(
		AdminServiceListAdjustmentTemplatesProcedure,
		svc.ListAdjustmentTemplates,
		connect.WithSchema(adminServiceListAdjustmentTemplatesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceCreateStageScoreHandler := connect.NewUnaryHandler(
		AdminServiceCreateStageScoreProcedure,
		svc.CreateStageScore,
		connect.WithSchema(adminServiceCreateStageScoreMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceUpdateStageScoreHandler := connect.NewUnaryHandler(
		AdminServiceUpdateStageScoreProcedure,
		svc.UpdateStageScore,
		connect.WithSchema(adminServiceUpdateStageScoreMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceListStageScoresHandler := connect.NewUnaryHandler(
		AdminServiceListStageScoresProcedure,
		svc.ListStageScores,
		connect.WithSchema(adminServiceListStageScoresMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminServiceDeleteStageScoreHandler := connect.NewUnaryHandler(
		AdminServiceDeleteStageScoreProcedure,
		svc.DeleteStageScore,
		connect.WithSchema(adminServiceDeleteStageScoreMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/api.v1.AdminService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AdminServicePurgeAllCachesProcedure:
			adminServicePurgeAllCachesHandler.ServeHTTP(w, r)
		case AdminServiceCreatePlayerProcedure:
			adminServiceCreatePlayerHandler.ServeHTTP(w, r)
		case AdminServiceUpdatePlayerProcedure:
			adminServiceUpdatePlayerHandler.ServeHTTP(w, r)
		case AdminServiceListPlayersProcedure:
			adminServiceListPlayersHandler.ServeHTTP(w, r)
		case AdminServiceListVenuesProcedure:
			adminServiceListVenuesHandler.ServeHTTP(w, r)
		case AdminServiceUpdateStageProcedure:
			adminServiceUpdateStageHandler.ServeHTTP(w, r)
		case AdminServiceListEventStagesProcedure:
			adminServiceListEventStagesHandler.ServeHTTP(w, r)
		case AdminServiceCreateAdjustmentTemplateProcedure:
			adminServiceCreateAdjustmentTemplateHandler.ServeHTTP(w, r)
		case AdminServiceUpdateAdjustmentTemplateProcedure:
			adminServiceUpdateAdjustmentTemplateHandler.ServeHTTP(w, r)
		case AdminServiceListAdjustmentTemplatesProcedure:
			adminServiceListAdjustmentTemplatesHandler.ServeHTTP(w, r)
		case AdminServiceCreateStageScoreProcedure:
			adminServiceCreateStageScoreHandler.ServeHTTP(w, r)
		case AdminServiceUpdateStageScoreProcedure:
			adminServiceUpdateStageScoreHandler.ServeHTTP(w, r)
		case AdminServiceListStageScoresProcedure:
			adminServiceListStageScoresHandler.ServeHTTP(w, r)
		case AdminServiceDeleteStageScoreProcedure:
			adminServiceDeleteStageScoreHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAdminServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAdminServiceHandler struct{}

func (UnimplementedAdminServiceHandler) PurgeAllCaches(context.Context, *connect.Request[v1.PurgeAllCachesRequest]) (*connect.Response[v1.PurgeAllCachesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.PurgeAllCaches is not implemented"))
}

func (UnimplementedAdminServiceHandler) CreatePlayer(context.Context, *connect.Request[v1.AdminServiceCreatePlayerRequest]) (*connect.Response[v1.AdminServiceCreatePlayerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.CreatePlayer is not implemented"))
}

func (UnimplementedAdminServiceHandler) UpdatePlayer(context.Context, *connect.Request[v1.UpdatePlayerRequest]) (*connect.Response[v1.UpdatePlayerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.UpdatePlayer is not implemented"))
}

func (UnimplementedAdminServiceHandler) ListPlayers(context.Context, *connect.Request[v1.ListPlayersRequest]) (*connect.Response[v1.ListPlayersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.ListPlayers is not implemented"))
}

func (UnimplementedAdminServiceHandler) ListVenues(context.Context, *connect.Request[v1.ListVenuesRequest]) (*connect.Response[v1.ListVenuesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.ListVenues is not implemented"))
}

func (UnimplementedAdminServiceHandler) UpdateStage(context.Context, *connect.Request[v1.UpdateStageRequest]) (*connect.Response[v1.UpdateStageResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.UpdateStage is not implemented"))
}

func (UnimplementedAdminServiceHandler) ListEventStages(context.Context, *connect.Request[v1.ListEventStagesRequest]) (*connect.Response[v1.ListEventStagesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.ListEventStages is not implemented"))
}

func (UnimplementedAdminServiceHandler) CreateAdjustmentTemplate(context.Context, *connect.Request[v1.CreateAdjustmentTemplateRequest]) (*connect.Response[v1.CreateAdjustmentTemplateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.CreateAdjustmentTemplate is not implemented"))
}

func (UnimplementedAdminServiceHandler) UpdateAdjustmentTemplate(context.Context, *connect.Request[v1.UpdateAdjustmentTemplateRequest]) (*connect.Response[v1.UpdateAdjustmentTemplateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.UpdateAdjustmentTemplate is not implemented"))
}

func (UnimplementedAdminServiceHandler) ListAdjustmentTemplates(context.Context, *connect.Request[v1.ListAdjustmentTemplatesRequest]) (*connect.Response[v1.ListAdjustmentTemplatesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.ListAdjustmentTemplates is not implemented"))
}

func (UnimplementedAdminServiceHandler) CreateStageScore(context.Context, *connect.Request[v1.CreateStageScoreRequest]) (*connect.Response[v1.CreateStageScoreResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.CreateStageScore is not implemented"))
}

func (UnimplementedAdminServiceHandler) UpdateStageScore(context.Context, *connect.Request[v1.UpdateStageScoreRequest]) (*connect.Response[v1.UpdateStageScoreResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.UpdateStageScore is not implemented"))
}

func (UnimplementedAdminServiceHandler) ListStageScores(context.Context, *connect.Request[v1.ListStageScoresRequest]) (*connect.Response[v1.ListStageScoresResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.ListStageScores is not implemented"))
}

func (UnimplementedAdminServiceHandler) DeleteStageScore(context.Context, *connect.Request[v1.DeleteStageScoreRequest]) (*connect.Response[v1.DeleteStageScoreResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v1.AdminService.DeleteStageScore is not implemented"))
}
