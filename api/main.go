package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/escavelo/pubgolf/api/lib/server"
	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

const (
	defaultPort = "50051"
)

func loadEnv() {
	env := os.Getenv("PUBGOLF_ENV")
	if "" == env {
		env = "dev"
	}

	log.Printf("Running as PUBGOLF_ENV of %s", env)

	// In dev, we assume we're running in the monorepo, in which case the .env is at the repo root. Otherwise, we're
	// probably in the docker environment and can access values as real env vars.
	if strings.ToLower(env) == "dev" {
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
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

func initDb(connStr string) *sql.DB {
	log.Printf("Connecting to DB: %s", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func initServer(db *sql.DB) *grpc.Server {
	grpcServer := grpc.NewServer()
	pg.RegisterAPIServer(grpcServer, &server.APIServer{DB: db})
	return grpcServer
}

func bindServer(server *grpc.Server, port string) {
	log.Printf("Listening on port %s", port)
	portListener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Serving on port %s", port)
	if err := server.Serve(portListener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	log.SetPrefix("[API Server] ")
	log.Println("Starting API server...")

	loadEnv()

	db := initDb(getDbConnectionString())
	defer db.Close()

	server := initServer(db)
	port := os.Getenv("API_PORT")
	if port == "" {
		port = defaultPort
	}
	bindServer(server, port)
}
