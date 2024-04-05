package testguard

import (
	"flag"
	"os"
)

var enableE2ETests = flag.Bool("e2e", false, "run e2e tests")

// UnitTest marks this test suite as unit-level and skips running if `-e2e=true` is passed to the test command.
func UnitTest() {
	flag.Parse()

	if *enableE2ETests {
		os.Exit(0)
	}
}

// E2ETest marks this test suite as e2e-level and skips running if `-e2e=true` is not passed to the test command.
func E2ETest() {
	flag.Parse()

	if !*enableE2ETests {
		os.Exit(0)
	}
}
