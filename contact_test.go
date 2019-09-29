// +build integration

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	// "github.com/google/go-cmp/cmp"
)

func TestContact(t *testing.T) {
	db := connect()
	defer db.Close()
	h := NewHandler(db)
	r := NewRouter(h)

	users := setupUsers(t, db, r)

	t.Run("Create", testCreateContact(db, r, users))
}

func setupUsers(t *testing.T, db *sql.DB, router http.Handler) {

	users := []User{
		User{
			PhoneNumber: "+65 9999 0001",
			FirstName:   "Contact 1",
			LastName:    "User",
		},
		User{
			PhoneNumber: "+65 9999 0002",
			FirstName:   "Contact 2",
			LastName:    "User",
		},
		User{
			PhoneNumber: "+65 9999 0003",
			FirstName:   "Contact 3",
			LastName:    "User",
		},
	}

	resultUsers := []User{}

	for _, user := range users {
		b, _ := json.Marshal(user)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/user", bytes.NewBuffer(b))

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		got := new(User)
		json.NewDecoder(w.Body).Decode(got)

		resultUsers = append(resultUsers, got)
	}

	return resultUsers

}

func testCreateContact(db *sql.DB, router http.Handler) func(t *testing.T) {
	return func(t *testing.T) {

		// Setup
		mockUser := &User{
			PhoneNumber: "+65 9999 1001",
			FirstName:   "ContactOwner",
			LastName:    "User",
		}
		b, _ := json.Marshal(mockUser)

		// Test
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/user", bytes.NewBuffer(b))

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		// Assert
		got, want := new(User), mockUser
		wantPhone, _ := ParsePhone(want.PhoneNumber)
		json.NewDecoder(w.Body).Decode(got)
		if got.FirstName != want.FirstName || got.LastName != want.LastName || got.PhoneNumber != wantPhone {
			t.Error("Wanted a User with same FirstName, LastName, PhoneNumber. Got something else")
		}

		assertDB(t, db, `SELECT * FROM "user" WHERE phone_number = '+65 9999 9999' AND first_name = 'Test' AND last_name = 'User 1'`)

	}
}
