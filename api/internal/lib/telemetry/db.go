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
func parseDBCQueryName(_ context.Context, method otelsql.Method, query string) string {
	const queryPrefix = "-- name: "

	if strings.HasPrefix(query, queryPrefix) {
		qPl := len(queryPrefix)

		return "dbc." + strings.Split(string(method), ".")[2] + "." + query[qPl:qPl+strings.Index(query[qPl:], " ")]
	}

	return string(method)
}

func parseDBCQueryAttributes(_ context.Context, method otelsql.Method, _ string, _ []driver.NamedValue) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("sql.method", string(method)),
	}
}

// WrapDB wraps a database connector with OpenTelemetry instrumentation.
// statsProvider is called on each span to attach live pool stats; pass nil to omit.
func WrapDB(db driver.Connector, statsProvider func() sql.DBStats) *sql.DB {
	return otelsql.OpenDB(db,
		otelsql.WithSpanNameFormatter(parseDBCQueryName),
		otelsql.WithAttributesGetter(func(ctx context.Context, method otelsql.Method, query string, args []driver.NamedValue) []attribute.KeyValue {
			attrs := parseDBCQueryAttributes(ctx, method, query, args)

			if statsProvider != nil {
				s := statsProvider()
				attrs = append(attrs,
					attribute.Int("db.pool.open", s.OpenConnections),
					attribute.Int("db.pool.in_use", s.InUse),
					attribute.Int("db.pool.idle", s.Idle),
					attribute.Int64("db.pool.wait_count", s.WaitCount),
					attribute.Int64("db.pool.wait_duration_ms", s.WaitDuration.Milliseconds()),
				)
			}

			return attrs
		}),
	)
}
