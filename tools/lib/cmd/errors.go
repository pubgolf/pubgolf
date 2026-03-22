package cmd

import (
	"log"
	"os"
	"strings"
)

// Exit codes for classifyAndExit.
const (
	// ExitCodeTestFailure indicates a code/test/lint failure (exit 1).
	ExitCodeTestFailure = 1
	// ExitCodeInfraFailure indicates an infrastructure failure (exit 2).
	ExitCodeInfraFailure = 2
)

// infraPatterns are known stderr substrings that indicate infrastructure failures
// rather than code/test problems.
var infraPatterns = []string{
	"address already in use",
	"cannot connect to the docker daemon",
	"connection refused",
	"permission denied",
	"no such host",
	"container exited with code 137", // OOM killed
}

// isInfraError checks whether an error message matches known infrastructure failure patterns.
func isInfraError(err error) bool {
	msg := strings.ToLower(err.Error())

	for _, p := range infraPatterns {
		if strings.Contains(msg, p) {
			return true
		}
	}

	return false
}

// classifyAndExit logs the error with classification and exits with the appropriate code.
// Exit code 1 = code/test failure. Exit code 2 = infrastructure failure.
func classifyAndExit(err error) {
	if err == nil {
		return
	}

	if isInfraError(err) {
		log.Printf("ERROR: Infrastructure failure\n  %s\n  This is not a code issue.", err)
		os.Exit(ExitCodeInfraFailure)
	}

	log.Printf("ERROR: %s", err)
	os.Exit(ExitCodeTestFailure)
}
