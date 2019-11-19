package server

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func invalidArgumentError(request interface{}) error {
	errorMsg := "Missing or invalid argument in request: %s."
	return status.New(codes.InvalidArgument, fmt.Sprintf(errorMsg, request)).Err()
}
}

func invalidAuthError() error {
	errorMsg := "Invalid or expired authorization."
	return status.New(codes.Unauthenticated, errorMsg).Err()
}

func eventNotFoundError(eventKey *string) error {
	errorMsg := "No event found with key '%s'."
	return status.New(codes.NotFound, fmt.Sprintf(errorMsg, *eventKey)).Err()
}

func userAlreadyExistsError(eventKey *string, phoneNumber *string) error {
	errorMsg := "User already exists for event '%s' with phone number '%s'."
	return status.New(codes.AlreadyExists, fmt.Sprintf(errorMsg, *eventKey,
		*phoneNumber)).Err()
}

func temporaryServerError(err error) error {
	errorMsg := "Server error: %s."
	return status.New(codes.Unavailable, fmt.Sprintf(errorMsg, err)).Err()
}
