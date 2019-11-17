package server

import (
	"database/sql"

	_ "github.com/lib/pq"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

type APIServer struct {
	DB *sql.DB
	// Include a default implementation of all RPC methods, even if we don't get
	// around to defining it.
	pg.UnimplementedAPIServer
}
