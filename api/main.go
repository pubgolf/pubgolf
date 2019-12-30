package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/escavelo/pubgolf/api/lib/server"
	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

const (
	defaultPort = "50051"
)

func loadEnv() string {
	env := os.Getenv("PUBGOLF_ENV")
	if "" == env {
		env = "dev"
	}

	// In dev, we assume we're running in the monorepo, in which case the .env is at the repo root. Otherwise, we're
	// probably in the docker environment and can access values as real env vars.
	if strings.ToLower(env) == "dev" {
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatal("error loading .env file")
		}
	}

	return env
}

func getDbConnectionString() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	if password != "" {
		password = ":" + password
	}
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	database := os.Getenv("DB_NAME")
	if database == "" {
		database = user
	}
	sslmode := os.Getenv("DB_SSL_MODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	return fmt.Sprintf("postgres://%s%s@%s:%s/%s?sslmode=%s", user, password,
		host, port, database, sslmode)
}

func initDb(logCtx *log.Entry, connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logCtx.Fatalf("unable to connect to DB: %s", err)
	}
	return db
}

func initServer(logCtx *log.Entry, db *sql.DB) *grpc.Server {
	grpcServer := grpc.NewServer()
	pg.RegisterAPIServer(grpcServer, &server.APIServer{LogCtx: logCtx, DB: db})
	return grpcServer
}

func bindServer(logCtx *log.Entry, server *grpc.Server, port string) {
	portListener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logCtx.Fatalf("failed to listen: %v", err)
	}

	logCtx.Info("Server Start")

	if err := server.Serve(portListener); err != nil {
		logCtx.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	env := loadEnv()
	port := os.Getenv("API_PORT")
	if port == "" {
		port = defaultPort
	}

	if env != "dev" {
		log.SetFormatter(&log.JSONFormatter{
			DataKey:     "data",
			PrettyPrint: false,
		})
	}
	serverLogCtx := log.WithField("server_info", log.Fields{
		"port": port,
		"env":  env,
	})

	db := initDb(serverLogCtx, getDbConnectionString())
	defer db.Close()

	server := initServer(serverLogCtx, db)
	go bindServer(serverLogCtx, server, port)

	<-shutdownChan
	server.GracefulStop()
	serverLogCtx.Info("Server Stop")
}
