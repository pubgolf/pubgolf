package config

import (
	"os"
	"testing"

	"github.com/pubgolf/pubgolf/api/internal/e2e"
)

func TestMain(m *testing.M) {
	e2e.GuardUnitTests()
	os.Exit(m.Run())
}
