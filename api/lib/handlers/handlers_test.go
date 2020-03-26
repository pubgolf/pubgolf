package handlers_test

import (
	"context"
	"database/sql"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	. "github.com/pubgolf/pubgolf/api/lib/handlers"
)

var testDB *sql.DB

// makeTestUnauthenticatedRequestData initializes a RequestData object with a connection to the test DB and returns a
// test logger.
func makeTestUnauthenticatedRequestData() (rd *RequestData, logHook *test.Hook) {
	tx, err := testDB.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		log.Fatalf("unable to connect to DB: %s", err)
	}

	logger, hook := test.NewNullLogger()

	return &RequestData{
		Tx:  tx,
		Ctx: context.Background(),
		Log: log.NewEntry(logger),
	}, hook
}

// countLogEntries returns the number of log messages in a testing logger that match the specified message text.
func countLogEntries(logHook *test.Hook, logMessage string) int {
	lineCount := 0
	for _, logLine := range logHook.AllEntries() {
		if logLine.Message == logMessage {
			lineCount++
		}
	}
	return lineCount
}
