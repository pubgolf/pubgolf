package public

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/oklog/ulid/v2"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
)

var errInvalidAuthCode = errors.New("invalid auth code")

// CompletePlayerLogin creates a player if the phone number hasn't been seen before and triggers a verification SMS.
func (s *Server) CompletePlayerLogin(_ context.Context, req *connect.Request[apiv1.CompletePlayerLoginRequest]) (*connect.Response[apiv1.CompletePlayerLoginResponse], error) {
	num, err := models.NewPhoneNum(req.Msg.GetPhoneNumber())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	authCode := req.Msg.GetAuthCode()
	if !sms.AuthCodeFormat.MatchString(authCode) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("incorrect auth code format: %w", errInvalidAuthCode))
	}

	valid, err := s.mes.CheckVerification(num, authCode)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	if !valid {
		return nil, connect.NewError(connect.CodePermissionDenied, errInvalidAuthCode)
	}

	// TODO: Fetch player.

	return connect.NewResponse(&apiv1.CompletePlayerLoginResponse{
		Player: &apiv1.Player{
			Id: ulid.Make().String(),
			Data: &apiv1.PlayerData{
				Name: "Name of Player",
			},
		},
	}), nil
}
