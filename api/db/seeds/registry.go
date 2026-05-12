// Package seeds provides named seed datasets for populating databases.
package seeds

import (
	"context"
	"database/sql"
	"sort"
)

// ExpectedCounts declares the expected row counts for a seed's event data.
// Used by the status command to detect drift from the baseline.
type ExpectedCounts struct {
	Stages              int
	Rules               int
	Players             int
	Scores              int
	AdjustmentTemplates int
}

// Seed defines a named seed dataset that targets a single event.
type Seed struct {
	Name     string
	EventKey string
	Expected ExpectedCounts
	Run      func(ctx context.Context, tx *sql.Tx) error
}

// Registry maps seed names to their definitions.
var Registry = map[string]Seed{
	"admin-test": AdminTestSeed,
}

// Names returns the sorted list of available seed names.
func Names() []string {
	names := make([]string, 0, len(Registry))
	for name := range Registry {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}
