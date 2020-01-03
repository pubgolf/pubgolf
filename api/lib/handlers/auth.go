package handlers

import (
	"github.com/escavelo/pubgolf/api/lib/db"
	"github.com/escavelo/pubgolf/api/lib/sms"
	"github.com/escavelo/pubgolf/api/lib/utils"
	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

// RegisterPlayer adds a new player to an event in an unconfirmed state and sends an auth code SMS to the provided
// phone number.
func RegisterPlayer(rd *RequestData, req *pg.RegisterPlayerRequest) (*pg.RegisterPlayerReply,
	error) {
	if utils.IsEmpty(&req.Name) || utils.InvalidPhoneNumberFormat(&req.PhoneNumber) {
		return nil, utils.InvalidArgumentError()
	}

	eventID, err := db.GetEventID(rd.Tx, &req.EventKey)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}
	if eventID == "" {
		return nil, utils.EventNotFoundError(&req.EventKey)
	}

	playerExists, err := db.CheckPlayerExistsByPhoneNumber(rd.Tx, &eventID, &req.PhoneNumber)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}
	if playerExists {
		return nil, utils.PlayerAlreadyExistsError(&req.EventKey, &req.PhoneNumber)
	}

	authCode, err := utils.GenerateAuthCode()
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	err = db.CreatePlayer(rd.Tx, &eventID, &req.Name, req.League, &req.PhoneNumber, authCode)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	err = sms.SendAuthCodeSms(&req.PhoneNumber, authCode)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	return &pg.RegisterPlayerReply{}, nil
}

// RequestPlayerLogin sends an auth code via SMS to the user matching the provided event key and phone number, if one
// exists. A non-existent player will not trigger an error response, in order to prevent mining of user phone numbers.
func RequestPlayerLogin(rd *RequestData, req *pg.RequestPlayerLoginRequest) (
	*pg.RequestPlayerLoginReply, error) {
	if utils.InvalidPhoneNumberFormat(&req.PhoneNumber) {
		return nil, utils.InvalidArgumentError()
	}

	eventID, err := db.GetEventID(rd.Tx, &req.EventKey)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}
	if eventID == "" {
		return nil, utils.EventNotFoundError(&req.EventKey)
	}

	playerExists, err := db.CheckPlayerExistsByPhoneNumber(rd.Tx, &eventID, &req.PhoneNumber)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}
	if !playerExists {
		// Player doesn't exist, so don't send an SMS, but return a success code to
		// avoid exposing what players do/don't exist.
		return &pg.RequestPlayerLoginReply{}, nil
	}

	authCode, err := utils.GenerateAuthCode()
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	err = db.SetAuthCode(rd.Tx, &eventID, &req.PhoneNumber, authCode)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	err = sms.SendAuthCodeSms(&req.PhoneNumber, authCode)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	return &pg.RequestPlayerLoginReply{}, nil
}

// PlayerLogin accepts and validates an auth code, returning an auth token for use in authenticated API calls.
func PlayerLogin(rd *RequestData, req *pg.PlayerLoginRequest) (*pg.PlayerLoginReply, error) {
	if utils.InvalidPhoneNumberFormat(&req.PhoneNumber) || utils.InvalidAuthCodeFormat(req.AuthCode) {
		return nil, utils.InvalidArgumentError()
	}

	eventID, err := db.GetEventID(rd.Tx, &req.EventKey)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}
	if eventID == "" {
		return nil, utils.EventNotFoundError(&req.EventKey)
	}

	authCodeValid, err := db.ValidateAuthCode(rd.Tx, &eventID, &req.PhoneNumber, req.AuthCode)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}
	if !authCodeValid {
		return nil, utils.InvalidAuthError()
	}

	err = db.GenerateAuthToken(rd.Tx, &eventID, &req.PhoneNumber)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	authToken, err := db.GetAuthToken(rd.Tx, &eventID, &req.PhoneNumber)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	return &pg.PlayerLoginReply{AuthToken: authToken}, nil
}
