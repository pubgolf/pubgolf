package config

import (
	"os"
	"testing"

	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"
)

func TestMain(m *testing.M) {
	testguard.UnitTest()
	os.Exit(m.Run())
}
