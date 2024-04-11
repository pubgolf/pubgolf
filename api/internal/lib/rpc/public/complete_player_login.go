package public

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
)

var errInvalidAuthCode = errors.New("invalid auth code")

// CompletePlayerLogin creates a player if the phone number hasn't been seen before and triggers a verification SMS.
func (s *Server) CompletePlayerLogin(ctx context.Context, req *connect.Request[apiv1.CompletePlayerLoginRequest]) (*connect.Response[apiv1.CompletePlayerLoginResponse], error) {
	num, err := models.NewPhoneNum(req.Msg.GetPhoneNumber())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse phone number: %w", err))
	}

	authCode := req.Msg.GetAuthCode()
	if !sms.AuthCodeFormat.MatchString(authCode) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("incorrect auth code format: %w", errInvalidAuthCode))
	}

	valid, err := s.mes.CheckVerification(ctx, num, authCode)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("check auth code: %w", err))
	}

	if !valid {
		return nil, connect.NewError(connect.CodePermissionDenied, errInvalidAuthCode)
	}

	didVerify, err := s.dao.VerifyPhoneNumber(ctx, num)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("mark phone number as verified: %w", err))
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Bool("player.phone_number_verified", !didVerify))

	playerID, authToken, err := s.dao.GenerateAuthToken(ctx, num)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("generate auth token: %w", err))
	}

	player, err := s.dao.PlayerByID(ctx, playerID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("get player from DB: %w", err))
	}

	p, err := player.Proto()
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("convert player model to proto: %w", err))
	}

	return connect.NewResponse(&apiv1.CompletePlayerLoginResponse{
		Player:    p,
		AuthToken: authToken.String(),
	}), nil
}
