package e2e

import (
	"flag"
	"os"
)

var enableE2ETests = flag.Bool("e2e", false, "run e2e tests")

func GuardUnitTests() {
	flag.Parse()

	if *enableE2ETests {
		os.Exit(0)
	}
}

func GuardE2ETests() {
	flag.Parse()

	if !*enableE2ETests {
		os.Exit(0)
	}
}
