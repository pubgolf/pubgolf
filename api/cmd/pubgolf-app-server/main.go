package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
	"github.com/pubgolf/pubgolf/api/internal/lib/honeycomb"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
	"github.com/pubgolf/pubgolf/api/internal/lib/rpc"
)

func main() {
	// Initialize app config.
	cfg, err := config.Init()
	guard(err, "parse env config")

	// Initialize monitoring.
	flush := honeycomb.Init(cfg)
	defer flush()

	// Initialize server.
	server := makeServer(cfg, &rpc.PubGolfServiceServer{})
	makeShutdownWatcher(server)

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

// makeServer initializes an HTTP server with settings and the router.
func makeServer(cfg *config.App, serverImpl apiv1connect.PubGolfServiceClient) *http.Server {
	// Construct gRPC server.

	rpcMux := http.NewServeMux()
	rpcMux.Handle(apiv1connect.NewPubGolfServiceHandler(serverImpl))

	// Mount routes.
	mux := http.NewServeMux()
	mux.HandleFunc("/health-check", healthCheck(cfg))
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
		Handler: h2c.NewHandler(honeycomb.WrapMux(mux), &http2.Server{}),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

// healthCheck returns a 200 if the app is online and able to process requests.
func healthCheck(cfg *config.App) http.HandlerFunc {
	return honeycomb.WrapHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Saluton mundo, from `%s`!", cfg.EnvName)
	})
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
