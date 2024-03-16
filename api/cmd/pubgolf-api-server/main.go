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
	"github.com/go-chi/chi/v5"
	chim "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/pubgolf/pubgolf/api/internal/db"
	"github.com/pubgolf/pubgolf/api/internal/lib/config"
	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/middleware"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
	"github.com/pubgolf/pubgolf/api/internal/lib/rpc/admin"
	"github.com/pubgolf/pubgolf/api/internal/lib/rpc/public"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
	"github.com/pubgolf/pubgolf/api/internal/lib/webapi"
)

func main() {
	// Initialize app config.
	cfg, err := config.Init()
	guard(err, "parse env config")

	// Initialize telemetry.
	cleanupTelemetry, err := telemetry.Init(cfg)
	guard(err, "init otel")
	defer cleanupTelemetry()

	ctx, bootSpan := otel.Tracer("").Start(context.Background(), "ServerBoot", trace.WithSpanKind(trace.SpanKindInternal))

	// Initialize DB.
	dbConn := makeDB(ctx, cfg)

	// Run migrations and exit if migrator instance.
	migrationFlag := flag.Bool("run-migrations", false, "run migrations and exit")
	flag.Parse()
	if *migrationFlag {
		bootSpan.SetAttributes(attribute.String("service.type", "migrator"))
		log.Println("Migrator instance: starting database migrations...")

		err = db.RunMigrations(dbConn)
		if err != nil {
			bootSpan.SetStatus(codes.Error, fmt.Sprintf("run migrations: %v", err))
		}

		log.Println("Migrator instance: completed migrations and shutting down...")
		bootSpan.End()
		return
	}
	bootSpan.SetAttributes(attribute.String("service.type", "server"))

	// Initialize server.
	dao, err := dao.New(ctx, dbConn, cfg.EnvName == config.DeployEnvDev)
	guard(err, "init DAO")
	server := makeServer(cfg, dao)
	makeShutdownWatcher(server)

	// Run server.
	bootSpan.End()
	log.Printf("Listening on port %d...", cfg.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		guard(err, "listen and serve")
	}
	log.Println("Server stopped")
}

// guard logs and exits on error.
func guard(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %v", msg, err.Error())
	}
}

// makeDB instantiates a database connection, verifies ability to connect and initializes tracing/debugging tools as necessary.
func makeDB(ctx context.Context, cfg *config.App) *sql.DB {
	conConfig, err := pgx.ParseConfig(cfg.AppDatabaseURL)
	guard(err, "parse database config")

	db := telemetry.WrapDB(stdlib.GetConnector(*conConfig))

	if cfg.EnvName == config.DeployEnvDev {
		err = db.PingContext(ctx)
		guard(err, "ping database")
	}

	return db
}

// makeServer initializes an HTTP server with settings and the router.
func makeServer(cfg *config.App, dao dao.QueryProvider) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.ChiMiddleware(r)...)

	// Mount routes.
	r.Get("/health-check", healthCheck(cfg))
	r.Get("/robots.txt", robots(cfg))
	r.Route("/web-api/", webapi.Router(cfg))
	r.Route("/rpc/", func(r chi.Router) {
		r.Use(
			chim.NoCache,
			httprate.Limit(10, 1*time.Second, httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
				return r.Header.Get("X-PubGolf-User-ID"), nil
			})),
		)

		rpcMux := http.NewServeMux()
		rpcMux.Handle(apiv1connect.NewPubGolfServiceHandler(public.NewServer(dao),
			connect.WithInterceptors(middleware.ConnectInterceptors()...),
		))
		rpcMux.Handle(apiv1connect.NewAdminServiceHandler(admin.NewServer(dao),
			connect.WithInterceptors(middleware.ConnectInterceptors()...),
			connect.WithInterceptors(middleware.NewAdminAuthInterceptor(cfg)),
		))
		r.Mount("/", http.StripPrefix("/rpc", rpcMux))
	})

	// Reverse proxy the web-app's static deployment.
	if cfg.WebAppUpstreamHost != "" {
		upstream, err := url.Parse(cfg.WebAppUpstreamHost)
		guard(err, "parse upstream for web-app reverse proxy")
		r.HandleFunc("/*", httputil.NewSingleHostReverseProxy(upstream).ServeHTTP)
	} else {
		r.HandleFunc("/*", http.FileServer(http.Dir("./web-app-content/")).ServeHTTP)
	}

	// Configure HTTP server.
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: h2c.NewHandler(r, &http2.Server{}),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
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
