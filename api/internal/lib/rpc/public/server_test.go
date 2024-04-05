package public

import (
	"os"
	"testing"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"
)

func TestMain(m *testing.M) {
	testguard.UnitTest()
	os.Exit(m.Run())
}

func makeTestServer(dao dao.QueryProvider) *Server {
	mockMessenger := new(sms.MockMessenger)

	return NewServer(dao, mockMessenger)
}
