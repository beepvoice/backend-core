package main

import (
	"database/sql"
	"net/http/httptest"
	"testing"
)

func assertCode(t *testing.T, w *httptest.ResponseRecorder, code int) {
	if got, want := w.Code, code; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}

func assertDB(t *testing.T, db *sql.DB, query string, args ...interface{}) {
	rows, err := db.Query(query, args...)
	if err != nil {
		t.Errorf("Error during query %s: %s", query, err)
		return
	}
	defer rows.Close()
	if rows.Next() != true {
		t.Errorf("Want one result, found none for query %s", query)
	}
}
