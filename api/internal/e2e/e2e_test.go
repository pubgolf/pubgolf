package e2e

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

var cleanup func()

func TestMain(m *testing.M) {
	testguard.E2ETest()

	cleanup = runAPIServer()
	ret := m.Run()

	cleanup()
	os.Exit(ret)
}

func runAPIServer() func() {
	server := exec.Command(
		"doppler", "run",
		"--project", "pubgolf-api-server",
		"--config", "e2e",
		"--",
		"go", "run", "../../cmd/pubgolf-api-server",
	)

	serverLog := log.New(log.Writer(), "[API Server] ", log.LstdFlags)

	server.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	server.Stdout = serverLog.Writer()
	server.Stderr = serverLog.Writer()
	server.Stdin = os.Stdin

	guard(server.Start(), "start API server")

	serverStarted := false

	for range 60 {
		res, _ := http.Get("http://localhost:3100/health-check") //nolint:noctx
		if res != nil {
			res.Body.Close()

			if res.StatusCode == http.StatusOK {
				serverStarted = true

				break
			}
		}

		time.Sleep(1 * time.Second)
	}

	if !serverStarted {
		log.Panicln("API server startup timed out")
	}

	return func() {
		pgid, err := syscall.Getpgid(server.Process.Pid)
		guard(err, "get process group ID for API server")
		guard(syscall.Kill(-pgid, syscall.SIGINT), "send SIGINT to API server")
	}
}

// guard logs and exits on error.
func guard(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %v", msg, err.Error())
	}
}
