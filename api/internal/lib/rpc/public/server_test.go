package public

import (
	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
)

func makeTestServer(dao dao.QueryProvider) *Server {
	mockMessenger := new(sms.MockMessenger)

	return NewServer(dao, mockMessenger)
}
