package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bufbuild/connect-go"
	otelconnect "github.com/bufbuild/connect-opentelemetry-go"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/pubgolf/pubgolf/api/internal/db"
	"github.com/pubgolf/pubgolf/api/internal/lib/config"
	"github.com/pubgolf/pubgolf/api/internal/lib/middleware"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
	"github.com/pubgolf/pubgolf/api/internal/lib/rpc"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

func main() {
	// Initialize app config.
	cfg, err := config.Init()
	guard(err, "parse env config")

	// Initialize telemetry.
	cleanupTelemetry, err := telemetry.Init(cfg)
	guard(err, "init otel")
	defer cleanupTelemetry()

	// Initialize server.
	dbConn := makeDB(cfg)
	server := makeServer(cfg, dbConn)
	makeShutdownWatcher(server)

	migrationFlag := flag.Bool("run-migrations", false, "run migrations and exit")
	flag.Parse()

	if *migrationFlag {
		err = db.RunMigrations(dbConn)
		guard(err, "run migrations")

		os.Exit(0)
	}

	// Run server.
	log.Printf("Listening on port %d...", cfg.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		guard(err, "listen and serve")
	}
	log.Println("Server stopped")
}

// guard logs and exits on error.
func guard(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err.Error())
	}
}

// makeDB instantiates a database connection, verifies ability to connect and initializes tracing/debugging tools as necessary.
func makeDB(cfg *config.App) *sql.DB {
	ctr, err := pq.NewConnector(cfg.AppDatabaseURL)
	guard(err, "open database connection")

	db := telemetry.WrapDB(ctr)

	err = db.Ping()
	guard(err, "ping database")

	return db
}

// makeServer initializes an HTTP server with settings and the router.
func makeServer(cfg *config.App, db *sql.DB) *http.Server {
	// Construct gRPC server.
	rpcMux := http.NewServeMux()
	rpcMux.Handle(apiv1connect.NewPubGolfServiceHandler(
		rpc.NewPubGolfServiceServer(db),
		connect.WithInterceptors(
			otelconnect.NewInterceptor(),
			middleware.NewLoggingInterceptor(),
		),
	))

	// Mount routes.
	mux := http.NewServeMux()
	mux.Handle("/health-check", healthCheck(cfg))
	mux.Handle("/rpc/", http.StripPrefix("/rpc", rpcMux))

	// Fallback to serving the built web-app assets, or the HMR server in the dev environment.
	if cfg.EnvName == config.DeployEnvDev {
		upstream, err := url.Parse("http://localhost:3000")
		guard(err, "parse upstream for web-app reverse proxy")
		mux.HandleFunc("/", httputil.NewSingleHostReverseProxy(upstream).ServeHTTP)
	} else {
		mux.HandleFunc("/", http.FileServer(http.Dir("./web-app/build")).ServeHTTP)
	}

	// Configure HTTP server.
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: h2c.NewHandler(mux, &http2.Server{}),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

// healthCheck returns a 200 if the app is online and able to process requests.
func healthCheck(cfg *config.App) http.Handler {
	return otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Saluton mundo, from `%s`!", cfg.EnvName)
	}), "health-check")
}

// makeShutdownWatcher spawns a child goroutine to gracefully close the server on an OS signal.
func makeShutdownWatcher(server *http.Server) {
	beginShutdown := make(chan os.Signal, 1)
	signal.Notify(beginShutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go shutdownWatcher(server, beginShutdown)
}

// shutdownWatcher watches the provided `beginShutdown` channel and begins a graceful shutdown of the provided `server`.
func shutdownWatcher(server *http.Server, beginShutdown <-chan os.Signal) {
	<-beginShutdown
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	guard(server.Shutdown(ctx), "call server shutdown command")
}
