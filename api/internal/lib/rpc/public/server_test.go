package public

import (
	"testing"

	"go.uber.org/goleak"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"
)

func TestMain(m *testing.M) {
	testguard.UnitTest()
	goleak.VerifyTestMain(m,
		goleak.IgnoreTopFunction("github.com/hashicorp/golang-lru/v2/expirable.NewLRU[...].func1"),
		goleak.IgnoreTopFunction("net/http.(*http2clientConnReadLoop).run"),
	)
}

func makeTestServer(dao dao.QueryProvider) *Server {
	mockMessenger := new(sms.MockMessenger)

	return NewServer(dao, mockMessenger)
}
