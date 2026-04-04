package admin

import (
	"os"
	"testing"

	"github.com/pubgolf/pubgolf/api/internal/lib/blobstore"
	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"
)

func TestMain(m *testing.M) {
	testguard.UnitTest()
	os.Exit(m.Run())
}

func makeTestServer(q dao.QueryProvider) *Server {
	return NewServer(q, new(blobstore.MockBlobStore))
}
