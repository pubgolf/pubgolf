package public

import (
	"testing"

	"go.uber.org/goleak"

	"github.com/pubgolf/pubgolf/api/internal/lib/blobstore"
	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"
)

func TestMain(m *testing.M) {
	testguard.UnitTest()
	goleak.VerifyTestMain(m,
		// dao package initializes expirable LRU caches at package level; their eviction
		// goroutines run for the lifetime of the process and can't be stopped.
		goleak.IgnoreTopFunction("github.com/hashicorp/golang-lru/v2/expirable.NewLRU[...].func1"),
	)
}

func makeTestServer(dao dao.QueryProvider) *Server {
	mockMessenger := new(sms.MockMessenger)
	mockBlobStore := new(blobstore.MockBlobStore)

	return NewServer(dao, mockMessenger, mockBlobStore)
}
