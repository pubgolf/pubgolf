package config

import "fmt"

// DB holds the connection info for a Postgres database.
type DB struct {
	Host     string `default:"localhost"`
	User     string
	Name     string
	Password string
	Port     int    `default:"5432"`
	SSLMode  string `default:"disable" split_words:"true"`
}

// DataSourceName returns a string suitable for consumption by `sql.Open()`.
func (db DB) DataSourceName() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		db.User, db.Password, db.Name, db.Host, db.Port, db.SSLMode)
}
