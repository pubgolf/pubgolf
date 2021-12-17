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

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	pubg "github.com/pubgolf/pubgolf/api/gen/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/lib/config"
	"github.com/pubgolf/pubgolf/api/lib/honeycomb"
	"github.com/pubgolf/pubgolf/api/lib/rpc"
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
func makeServer(cfg *config.App, serverImpl pubg.PubGolfServiceServer) *http.Server {
	// Construct gRPC server.
	grpcServer := grpc.NewServer()
	pubg.RegisterPubGolfServiceServer(grpcServer, serverImpl)

	// Configure gRPC-web wrapper.
	var options []grpcweb.Option
	if cfg.EnvName != config.DeployEnvProd {
		options = []grpcweb.Option{
			// Allow FE dev server to hit the API directly (not via reverse proxy) to allow client-only dev workflows to target staging or preview environments.
			grpcweb.WithOriginFunc(func(origin string) bool {
				return origin == "http://localhost:3000"
			}),
		}
	}
	wrappedGRPCServer := grpcweb.WrapServer(grpcServer, options...)

	// Mount routes.
	mux := http.NewServeMux()
	mux.HandleFunc("/health-check", healthCheck(cfg))
	mux.Handle("/rpc/", http.StripPrefix("/rpc", webGRPC(wrappedGRPCServer)))

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
		Handler: honeycomb.WrapMux(mux),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

// webGRPC passes the request through the gRPC-web proxy, which converts it into a standard gRPC request for the handler.
func webGRPC(wrappedGRPCServer *grpcweb.WrappedGrpcServer) http.HandlerFunc {
	return honeycomb.WrapHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedGRPCServer.ServeHTTP(w, r)
	})
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
