package public

import (
	"context"

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
	_, err = s.dao.CreatePlayer(ctx, "", num)
	if err != nil {
		if err != dao.ErrAlreadyCreated {
			return nil, connect.NewError(connect.CodeUnknown, err)
		}
		playerExists = true
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Bool("login_request.player_exists", playerExists))

	err = s.mes.SendVerification(num)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	return connect.NewResponse(&apiv1.StartPlayerLoginResponse{}), nil
}
