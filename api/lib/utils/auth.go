package utils

import (
	"context"
	"crypto/rand"
	"math/big"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GetAuthTokenFromHeader parses the auth code out of the request metadata ('authorization' header).
func GetAuthTokenFromHeader(ctx context.Context) (string, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.DataLoss, "Error reading 'authorization' header.")
	}

	authHeader, ok := metadata["authorization"]
	if !ok || len(authHeader) != 1 {
		return "", InsufficientPermissionsError()
	}

	validFormat := uuidFormat.MatchString(authHeader[0])
	if !validFormat {
		return "", InsufficientPermissionsError()
	}

	return authHeader[0], nil
}

// GenerateAuthCode generates a valid, cryptographically random auth code.
func GenerateAuthCode() (uint32, error) {
	randNum, err := rand.Int(rand.Reader, big.NewInt(899999))
	if err != nil {
		return 0, err
	}
	randCode := randNum.Int64() + 100000
	return uint32(randCode), nil
}
