// +build integration

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	//"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUser(t *testing.T) {
	db := connect()
	defer db.Close()
	h := NewHandler(db)
	r := NewRouter(h)

	t.Run("Create", testCreateUser(db, r))
}

func testCreateUser(db *sql.DB, router http.Handler) func(t *testing.T) {
	return func(t *testing.T) {
		mockUser := &User{
			PhoneNumber: "+65 99999999",
			FirstName:   "Test",
			LastName:    "User",
		}
		b, _ := json.Marshal(mockUser)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/user", bytes.NewBuffer(b))

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		assertDB(t, db, `SELECT * FROM "user" WHERE phone_number = '+65 9999 9999' AND first_name = 'Test' AND last_name = 'User'`)
	}
}

/*
	got, want := new(User), created
	json.NewDecoder(w.Body).Decode(got)
	if diff := cmp.Diff(got, want); len(diff) != 0 {
		t.Error(diff)
	}
*/
