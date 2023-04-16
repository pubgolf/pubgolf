package telemetry

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"strings"

	"github.com/XSAM/otelsql"
	"go.opentelemetry.io/otel/attribute"
)

// spanNameFormatterFn is a helper to convert an anonymous function to an `otelsql.SpanNameFormatter`.
type spanNameFormatterFn func(ctx context.Context, method otelsql.Method, query string) string

// Format calls the `spanNameFormatterFn` as a function.
func (s spanNameFormatterFn) Format(ctx context.Context, method otelsql.Method, query string) string {
	return s(ctx, method, query)
}

// parseDBCQueryName parses the SQLc comment of a SQL query to override the span name.
var parseDBCQueryName spanNameFormatterFn = spanNameFormatterFn(func(ctx context.Context, method otelsql.Method, query string) string {
	const queryPrefix = "-- name: "

	if strings.HasPrefix(query, queryPrefix) {
		qPl := len(queryPrefix)
		return "dbc." + query[qPl:qPl+strings.Index(query[qPl:], " ")]
	}

	return string(method)
})

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
