package public

import (
	"context"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// StartPlayerLogin creates a player if the phone number hasn't been seen before and triggers a verification SMS.
func (s *Server) StartPlayerLogin(_ context.Context, req *connect.Request[apiv1.StartPlayerLoginRequest]) (*connect.Response[apiv1.StartPlayerLoginResponse], error) {
	// TODO: Upsert player based on phone number.
	num, err := models.NewPhoneNum(req.Msg.GetPhoneNumber())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	err = s.mes.SendVerification(num)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	return connect.NewResponse(&apiv1.StartPlayerLoginResponse{}), nil
}
