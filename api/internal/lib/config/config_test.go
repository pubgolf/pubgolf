package config

import (
	"testing"

	"go.uber.org/goleak"

	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"
)

func TestMain(m *testing.M) {
	testguard.UnitTest()
	goleak.VerifyTestMain(m)
}
