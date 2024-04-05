package public

import (
	"os"
	"testing"

	"github.com/pubgolf/pubgolf/api/internal/e2e"
	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
)

func TestMain(m *testing.M) {
	e2e.GuardUnitTests()
	os.Exit(m.Run())
}

func makeTestServer(dao dao.QueryProvider) *Server {
	mockMessenger := new(sms.MockMessenger)

	return NewServer(dao, mockMessenger)
}
