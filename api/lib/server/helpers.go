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

var authTokenFormat *regexp.Regexp = regexp.MustCompile(
	"^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
var phoneNumberFormat *regexp.Regexp = regexp.MustCompile(
	"^\\+[1-9]\\d{1,14}$")

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

	validFormat := authTokenFormat.MatchString(authHeader[0])
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

func invalidAuthCodeFormat(authCode uint32) bool {
	return authCode < 100000 || authCode > 999999
}
