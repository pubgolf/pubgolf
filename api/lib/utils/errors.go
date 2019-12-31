package utils

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InvalidArgumentError constructs a gRPC error to indicate that the client request was malformed.
func InvalidArgumentError() error {
	errorMsg := "Missing or invalid argument in request."
	return status.New(codes.InvalidArgument, errorMsg).Err()
}

// InsufficientPermissionsError constructs a gRPC error to indicate that the request lacked an authorized token for the
// requested action. Invalid tokens should use `InvalidAuthError()` instead.
func InsufficientPermissionsError() error {
	errorMsg := "Lacking sufficient authorization for request."
	return status.New(codes.PermissionDenied, errorMsg).Err()
}

// InvalidAuthError constructs a gRPC error to indicate that the auth token is not valid. Valid tokens that aren't
// authorized for a particular action should use `InsufficientPermissionsError()` instead.
func InvalidAuthError() error {
	errorMsg := "Invalid or expired authorization."
	return status.New(codes.Unauthenticated, errorMsg).Err()
}

// EventNotFoundError constructs a gRPC error to indicate that the provided event key was invalid.
func EventNotFoundError(eventKey *string) error {
	errorMsg := "No event found with key '%s'."
	return status.New(codes.NotFound, fmt.Sprintf(errorMsg, *eventKey)).Err()
}

// PlayerNotFoundError constructs a gRPC error to indicate that the provided player ID was invalid.
func PlayerNotFoundError(playerID *string) error {
	errorMsg := "No player found with ID '%s'."
	return status.New(codes.NotFound, fmt.Sprintf(errorMsg, *playerID)).Err()
}

// PlayerAlreadyExistsError constructs a gRPC error to indicate that a conflicting registration already exists. This
// error response is deprecated.
func PlayerAlreadyExistsError(eventKey *string, phoneNumber *string) error {
	errorMsg := "Player already exists for event '%s' with phone number '%s'."
	return status.New(codes.AlreadyExists, fmt.Sprintf(errorMsg, *eventKey,
		*phoneNumber)).Err()
}

// TemporaryServerError constructs a gRPC error to indicate that a miscellaneous server error has occurred, and that
// the request may be retried.
func TemporaryServerError(err error) error {
	errorMsg := "Server error: %s."
	return status.New(codes.Unavailable, fmt.Sprintf(errorMsg, err)).Err()
}
