package public

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// StartPlayerLogin creates a player if the phone number hasn't been seen before and triggers a verification SMS.
func (s *Server) StartPlayerLogin(ctx context.Context, req *connect.Request[apiv1.StartPlayerLoginRequest]) (*connect.Response[apiv1.StartPlayerLoginResponse], error) {
	num, err := models.NewPhoneNum(req.Msg.GetPhoneNumber())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	playerExists := false
	span := trace.SpanFromContext(ctx)

	_, err = s.dao.CreatePlayer(ctx, "", num)
	if err != nil {
		if !errors.Is(err, dao.ErrAlreadyCreated) {
			return nil, connect.NewError(connect.CodeUnknown, err)
		}

		// Player already exists. Do not return an error, but check and log if we've already verified the phone number.

		isVerified, err := s.dao.PhoneNumberIsVerified(ctx, num)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("check if phone number is verified: %w", err))
		}

		span.SetAttributes(attribute.Bool("player.phone_number_verified", isVerified))

		playerExists = true
	}

	span.SetAttributes(attribute.Bool("player.new_phone_number", !playerExists))

	err = s.mes.SendVerification(ctx, num)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	return connect.NewResponse(&apiv1.StartPlayerLoginResponse{}), nil
}
