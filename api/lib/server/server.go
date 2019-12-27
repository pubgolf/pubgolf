package server

import (
	"database/sql"

	// Needed to import the Postgres driver correctly.
	_ "github.com/lib/pq"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

// APIServer is a struct for passing global context, such as the databse handle.
type APIServer struct {
	DB *sql.DB
	// Include a default implementation of all RPC methods, even if we don't get around to defining it.
	pg.UnimplementedAPIServer
}
