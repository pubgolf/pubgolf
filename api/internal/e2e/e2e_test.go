package e2e

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/fatih/color"

	"github.com/pubgolf/pubgolf/api/internal/lib/dbtest"
	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"
)

var (
	sharedTestDB        *sql.DB
	dopplerFallbackFlag = flag.Bool("doppler-local", false, "use fallback file for doppler to avoid a network call")
	noTelemetry         = flag.Bool("disable-telemetry", false, "do not initialize OTel or send telemetry to Honeycomb")
)

type LogWriter struct {
	logger *log.Logger
}

func NewPrefixLogWriter(p string) *LogWriter {
	l := log.New(log.Writer(), fmt.Sprintf("[%s] ", p), log.Flags())
	l.SetFlags(0)

	return NewLogWriter(l)
}

func NewLogWriter(l *log.Logger) *LogWriter {
	return &LogWriter{
		logger: l,
	}
}

func (lw LogWriter) Write(p []byte) (int, error) {
	lw.logger.Print(string(p))

	return len(p), nil
}

func TestMain(m *testing.M) {
	testguard.E2ETest()

	color.NoColor = false

	dbURL, dbCleanupFn := dbtest.NewURL("pubgolf-e2e")

	var err error
	sharedTestDB, err = sql.Open("pgx", dbURL)
	guard(err, "open DB connection")

	runAPIMigrator(dbURL)

	serverCleanup := runAPIServer(dbURL)

	ret := m.Run()

	serverCleanup()
	dbCleanupFn()
	os.Exit(ret)
}

func runAPIMigrator(dbURL string) {
	args := []string{
		"run",
		"--project", "pubgolf-api-server",
		"--config", "e2e",
		"--preserve-env",
	}

	if *dopplerFallbackFlag {
		args = append(args, []string{
			"--no-check-version",
			"--fallback-only",
		}...)
	}

	args = append(args, []string{
		"--",
		"go", "run", "../../cmd/pubgolf-api-server", "--run-migrations",
	}...)

	if *noTelemetry {
		args = append(args, []string{
			"--disable-telemetry",
		}...)
	}

	migrator := exec.Command("doppler", args...)

	migrator.Env = append(os.Environ(), "PUBGOLF_APP_DATABASE_URL="+dbURL)

	migratorLog := NewPrefixLogWriter(color.RedString("Migrator"))
	migrator.Stdout = migratorLog
	migrator.Stderr = migratorLog
	migrator.Stdin = os.Stdin

	guard(migrator.Run(), "run API migrator")
}

func runAPIServer(dbURL string) func() {
	args := []string{
		"run",
		"--project", "pubgolf-api-server",
		"--config", "e2e",
		"--preserve-env",
	}

	if *dopplerFallbackFlag {
		args = append(args, []string{
			"--no-check-version",
			"--fallback-only",
		}...)
	}

	args = append(args, []string{
		"--",
		"go", "run", "../../cmd/pubgolf-api-server",
	}...)

	if *noTelemetry {
		args = append(args, []string{
			"--disable-telemetry",
		}...)
	}

	server := exec.Command("doppler", args...)

	server.Env = append(os.Environ(), "PUBGOLF_APP_DATABASE_URL="+dbURL)
	server.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	serverLog := NewPrefixLogWriter(color.BlueString("API Server"))
	server.Stdout = serverLog
	server.Stderr = serverLog
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

// requestWithAuth makes an RPC call with an auth header.
func requestWithAuth[T any](msg *T, token string) *connect.Request[T] {
	req := connect.NewRequest(msg)
	req.Header().Set("X-Pubgolf-Authtoken", token)

	return req
}

// requestWithAdminAuth makes an RPC call with an admin auth header.
func requestWithAdminAuth[T any](msg *T) *connect.Request[T] {
	req := connect.NewRequest(msg)
	req.Header().Set("X-Pubgolf-Authtoken", "admin-api-token-value")

	return req
}
