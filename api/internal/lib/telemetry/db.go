package telemetry

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"strings"

	"github.com/XSAM/otelsql"
	"go.opentelemetry.io/otel/attribute"
)

// parseDBCQueryName parses the SQLc comment of a SQL query to override the span name.
func parseDBCQueryName(ctx context.Context, method otelsql.Method, query string) string {
	const queryPrefix = "-- name: "

	if strings.HasPrefix(query, queryPrefix) {
		qPl := len(queryPrefix)
		return "dbc." + strings.Split(string(method), ".")[2] + "." + query[qPl:qPl+strings.Index(query[qPl:], " ")]
	}

	return string(method)
}

func parseDBCQueryAttributes(ctx context.Context, method otelsql.Method, query string, args []driver.NamedValue) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("sql.method", string(method)),
	}
}

// WrapDB returns an OTel-instrumented DB handle.
func WrapDB(db driver.Connector) *sql.DB {
	return otelsql.OpenDB(db,
		otelsql.WithSpanNameFormatter(parseDBCQueryName),
		otelsql.WithAttributesGetter(parseDBCQueryAttributes),
	)
}
