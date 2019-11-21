package server

import (
	"context"
	"log"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"

	"github.com/escavelo/pubgolf/api/lib/db"
	"github.com/escavelo/pubgolf/api/lib/sms"
)

func (server *APIServer) RegisterPlayer(ctx context.Context,
	req *pg.RegisterPlayerRequest) (*pg.RegisterPlayerReply, error) {
	if isEmpty(&req.Name) || invalidPhoneNumberFormat(&req.PhoneNumber) {
		return nil, invalidArgumentError()
	}

	tx, eventID, err := validateUnauthenticatedRequest(server, &req.EventKey)
	if err != nil {
		return nil, err
	}

	userExists, err := db.CheckPlayerExists(tx, &eventID, &req.PhoneNumber)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}
	if userExists {
		tx.Rollback()
		return nil, playerAlreadyExistsError(&req.EventKey, &req.PhoneNumber)
	}

	authCode, err := generateAuthCode()
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}

	err = db.CreatePlayer(tx, &eventID, &req.Name, req.League, &req.PhoneNumber,
		authCode)
	if err != nil {
		log.Printf("%s - %s", eventID, req)
		tx.Rollback()
		return nil, temporaryServerError(err)
	}

	err = sms.SendAuthCodeSms(&req.PhoneNumber, authCode)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}

	tx.Commit()
	return &pg.RegisterPlayerReply{}, nil
}

func (server *APIServer) RequestPlayerLogin(ctx context.Context,
	req *pg.RequestPlayerLoginRequest) (*pg.RequestPlayerLoginReply, error) {
	if invalidPhoneNumberFormat(&req.PhoneNumber) {
		return nil, invalidArgumentError()
	}

	tx, eventID, err := validateUnauthenticatedRequest(server, &req.EventKey)
	if err != nil {
		return nil, err
	}

	userExists, err := db.CheckPlayerExists(tx, &eventID, &req.PhoneNumber)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}
	if !userExists {
		// Player doesn't exist, so don't send an SMS, but return a success code to
		// avoid exposing what users do/don't exist.
		tx.Rollback()
		return &pg.RequestPlayerLoginReply{}, nil
	}

	authCode, err := generateAuthCode()
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}

	err = db.SetAuthCode(tx, &eventID, &req.PhoneNumber, authCode)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}

	err = sms.SendAuthCodeSms(&req.PhoneNumber, authCode)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}

	tx.Commit()
	return &pg.RequestPlayerLoginReply{}, nil
}

func (server *APIServer) PlayerLogin(ctx context.Context,
	req *pg.PlayerLoginRequest) (*pg.PlayerLoginReply, error) {
	if invalidPhoneNumberFormat(&req.PhoneNumber) ||
		invalidAuthCodeFormat(req.AuthCode) {
		return nil, invalidArgumentError()
	}

	tx, eventID, err := validateUnauthenticatedRequest(server, &req.EventKey)
	if err != nil {
		return nil, err
	}

	authCodeValid, err := db.ValidateAuthCode(tx, &eventID, &req.PhoneNumber,
		req.AuthCode)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}
	if !authCodeValid {
		tx.Rollback()
		return nil, invalidAuthError()
	}

	err = db.GenerateAuthToken(tx, &eventID, &req.PhoneNumber)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}

	authToken, err := db.GetAuthToken(tx, &eventID, &req.PhoneNumber)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}

	tx.Commit()
	return &pg.PlayerLoginReply{AuthToken: authToken}, nil
}
