package server

import (
	"context"
	"crypto/rand"
	"database/sql"
	"math/big"
	"regexp"

	"github.com/escavelo/pubgolf/api/lib/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var uuidFormat *regexp.Regexp = regexp.MustCompile(
	"^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
var phoneNumberFormat *regexp.Regexp = regexp.MustCompile(
	"^\\+[1-9]\\d{1,14}$")

func validateAuthenticatedRequest(server *APIServer, ctx context.Context,
	eventKey *string) (*sql.Tx, string, string, error) {
	if isEmpty(eventKey) {
		return nil, "", "", invalidArgumentError()
	}

	authHeader, err := getAuthTokenFromHeader(ctx)
	if err != nil {
		return nil, "", "", err
	}

	tx, err := server.DB.Begin()
	if err != nil {
		return nil, "", "", temporaryServerError(err)
	}

	playerEventID, playerID, err := db.ValidateAuthToken(tx, &authHeader)
	if err != nil {
		tx.Rollback()
		return nil, "", "", temporaryServerError(err)
	}
	if playerEventID == "" || playerID == "" {
		tx.Rollback()
		return nil, "", "", insufficientPermissionsError()
	}

	eventID, err := db.GetEventID(tx, eventKey)
	if err != nil {
		tx.Rollback()
		return nil, "", "", temporaryServerError(err)
	}
	if eventID == "" {
		tx.Rollback()
		return nil, "", "", eventNotFoundError(eventKey)
	}

	if playerEventID != eventID {
		tx.Rollback()
		return nil, "", "", insufficientPermissionsError()
	}

	return tx, eventID, playerID, nil
}

func validateUnauthenticatedRequest(server *APIServer, eventKey *string) (
	*sql.Tx, string, error) {
	if isEmpty(eventKey) {
		return nil, "", invalidArgumentError()
	}

	tx, err := server.DB.Begin()
	if err != nil {
		return nil, "", temporaryServerError(err)
	}

	eventID, err := db.GetEventID(tx, eventKey)
	if err != nil {
		tx.Rollback()
		return nil, "", temporaryServerError(err)
	}
	if eventID == "" {
		tx.Rollback()
		return nil, "", eventNotFoundError(eventKey)
	}

	return tx, eventID, nil
}

func getAuthTokenFromHeader(ctx context.Context) (string, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.DataLoss,
			"Error reading 'authorization' header.")
	}

	authHeader, ok := metadata["authorization"]
	if !ok || len(authHeader) != 1 {
		return "", insufficientPermissionsError()
	}

	validFormat := uuidFormat.MatchString(authHeader[0])
	if !validFormat {
		return "", insufficientPermissionsError()
	}

	return authHeader[0], nil
}

func generateAuthCode() (uint32, error) {
	randNum, err := rand.Int(rand.Reader, big.NewInt(899999))
	if err != nil {
		return 0, err
	}
	randCode := randNum.Int64() + 100000
	return uint32(randCode), nil
}

func isEmpty(arg *string) bool {
	return *arg == ""
}

func invalidPhoneNumberFormat(phoneNumber *string) bool {
	isValid := phoneNumberFormat.MatchString(*phoneNumber)
	return !isValid
}

func invalidIDFormat(id *string) bool {
	isValid := uuidFormat.MatchString(*id)
	return !isValid
}

func invalidAuthCodeFormat(authCode uint32) bool {
	return authCode < 100000 || authCode > 999999
}
