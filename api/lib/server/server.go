package server

import (
	"context"
	"database/sql"

	// Needed to import the Postgres driver correctly.
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/escavelo/pubgolf/api/lib/db"
	"github.com/escavelo/pubgolf/api/lib/handlers"
	"github.com/escavelo/pubgolf/api/lib/utils"
	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

// APIServer is a struct for passing global context, such as the database handle.
type APIServer struct {
	DB  *sql.DB
	Log *log.Entry

	// Include a default implementation of all RPC methods, even if we don't get around to defining it.
	pg.UnimplementedAPIServer
}

func initializeRequestData(ctx context.Context, server *APIServer, methodName string) (*handlers.RequestData, error) {
	rd := handlers.RequestData{Ctx: ctx}

	rd.Log = server.Log.WithContext(ctx).WithField("request", log.Fields{
		"method": methodName,
	})

	tx, err := server.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return &rd, utils.TemporaryServerError(err)
	}
	rd.Tx = tx

	return &rd, nil
}

func validateAuthHeader(rd *handlers.RequestData) error {
	authHeader, err := utils.GetAuthTokenFromHeader(rd.Ctx)
	if err != nil {
		return err
	}

	eventID, playerID, err := db.ValidateAuthToken(rd.Tx, &authHeader)
	if err != nil {
		return utils.TemporaryServerError(err)
	}
	if eventID == "" || playerID == "" {
		return utils.InsufficientPermissionsError()
	}

	rd.EventID = eventID
	rd.PlayerID = playerID

	rd.Log = rd.Log.WithField("auth_token", log.Fields{
		"event_id":  eventID,
		"player_id": playerID,
	})

	return nil
}

func addResponseContextToLog(rd *handlers.RequestData, err error) {
	code := codes.OK
	msg := ""

	if err != nil {
		// Add raw error message as top-level field.
		rd.Log = rd.Log.WithError(err)

		if st, ok := status.FromError(err); ok {
			// Parse out an actual status code if we can.
			code = st.Code()
			msg = st.Message()
		} else {
			// Default to a generic "Unknown" error (which is what the gRPC server does if we throw an error without status
			// info).
			code = codes.Unknown
			msg = err.Error()
		}
	}

	responseDetails := log.Fields{
		"code":   code,
		"status": code.String(),
	}
	if msg != "" {
		responseDetails["msg"] = msg
	}

	rd.Log = rd.Log.WithField("response", responseDetails)

	// Log the event name as a standard message to make rollups easier.
	if err != nil {
		rd.Log.Error("gRPC Error")
	} else {
		rd.Log.Info("gRPC Success")
	}
}

func processUnauthenticatedRequest(ctx context.Context, server *APIServer, req interface{}, methodName string,
	handler func(*handlers.RequestData, interface{}) (interface{}, error)) (rep interface{}, err error) {
	// Set up logging and DB transactions.
	rd, err := initializeRequestData(ctx, server, methodName)
	// Flush log entry on return, using a closure to capture the named return vars.
	defer func() { addResponseContextToLog(rd, err) }()
	// Handle errors creating the DB transaction.
	if err != nil {
		return nil, err
	}

	// Call provided RPC handler, which contains the actual business logic.
	resp, err := handler(rd, req)
	if err != nil {
		rd.Tx.Rollback()
		return nil, err
	}

	rd.Tx.Commit()
	return resp, nil
}

func processAuthenticatedRequest(ctx context.Context, server *APIServer, req interface{}, methodName string,
	handler func(*handlers.RequestData, interface{}) (interface{}, error)) (rep interface{}, err error) {
	// Set up logging and DB transactions.
	rd, err := initializeRequestData(ctx, server, methodName)
	// Flush log entry on return, using a closure to capture the named return vars.
	defer func() { addResponseContextToLog(rd, err) }()
	// Handle errors creating the DB transaction.
	if err != nil {
		return nil, err
	}

	// Validate auth token in header.
	err = validateAuthHeader(rd)
	if err != nil {
		rd.Tx.Rollback()
		return nil, err
	}

	// Call provided RPC handler, which contains the actual business logic.
	resp, err := handler(rd, req)
	if err != nil {
		rd.Tx.Rollback()
		return nil, err
	}

	rd.Tx.Commit()
	return resp, nil
}
