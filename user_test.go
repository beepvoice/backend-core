// +build integration

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/google/go-cmp/cmp"
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
	t.Run("GetUserByPhone", testGetUserByPhone(db, r))
	t.Run("GetUser", testGetUser(db, r))
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

func testGetUserByPhone(db *sql.DB, router http.Handler) func(t *testing.T) {
	return func(t *testing.T) {
		mockUser := &User{
			PhoneNumber: "+65 9999 9999",
			FirstName:   "Test",
			LastName:    "User",
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/user?phone_number=%2B6599999999", nil)

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		got, want := new(User), mockUser
		json.NewDecoder(w.Body).Decode(got)
		if got.FirstName != want.FirstName || got.LastName != want.LastName || got.PhoneNumber != want.PhoneNumber {
			t.Errorf("Want user %v, got %v", want, got)
		}
	}
}

func testGetUser(db *sql.DB, router http.Handler) func(t *testing.T) {
	return func(t *testing.T) {
		mockUser := &User{
			PhoneNumber: "+65 88888888",
			FirstName:   "User",
			LastName:    "Test",
		}
		cb, _ := json.Marshal(mockUser)

		cw := httptest.NewRecorder()
		cr := httptest.NewRequest("POST", "/user", bytes.NewBuffer(cb))

		router.ServeHTTP(cw, cr)

		createdUser := new(User)
		json.NewDecoder(cw.Body).Decode(createdUser)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/user/id/"+createdUser.ID, nil)

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		got, want := new(User), createdUser
		json.NewDecoder(w.Body).Decode(got)
		if diff := cmp.Diff(got, want); len(diff) != 0 {
			t.Error(diff)
		}
	}
}
