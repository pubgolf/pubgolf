package server

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func invalidArgumentError() error {
	errorMsg := "Missing or invalid argument in request."
	return status.New(codes.InvalidArgument, errorMsg).Err()
}

func insufficientPermissionsError() error {
	errorMsg := "Lacking sufficent authorization for request."
	return status.New(codes.PermissionDenied, errorMsg).Err()
}

func invalidAuthError() error {
	errorMsg := "Invalid or expired authorization."
	return status.New(codes.Unauthenticated, errorMsg).Err()
}

func eventNotFoundError(eventKey *string) error {
	errorMsg := "No event found with key '%s'."
	return status.New(codes.NotFound, fmt.Sprintf(errorMsg, *eventKey)).Err()
}

func playerNotFoundError(playerID *string) error {
	errorMsg := "No player found with ID '%s'."
	return status.New(codes.NotFound, fmt.Sprintf(errorMsg, *playerID)).Err()
}

func playerAlreadyExistsError(eventKey *string, phoneNumber *string) error {
	errorMsg := "Player already exists for event '%s' with phone number '%s'."
	return status.New(codes.AlreadyExists, fmt.Sprintf(errorMsg, *eventKey,
		*phoneNumber)).Err()
}

func temporaryServerError(err error) error {
	errorMsg := "Server error: %s."
	return status.New(codes.Unavailable, fmt.Sprintf(errorMsg, err)).Err()
}
