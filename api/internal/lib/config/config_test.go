package config

import (
	"testing"

	"go.uber.org/goleak"

	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"
)

func TestMain(m *testing.M) {
	testguard.UnitTest()
	goleak.VerifyTestMain(m,
		// Background cache eviction goroutine from expirable LRU cache used by config package.
		goleak.IgnoreTopFunction("github.com/hashicorp/golang-lru/v2/expirable.NewLRU[...].func1"),
	)
}
