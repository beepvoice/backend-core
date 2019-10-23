// +build integration

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestContact(t *testing.T) {
	db := connect()
	defer db.Close()
	h := NewHandler(db, nil)
	r := NewRouter(h)

	users := setupContactUsers(t, db, r)

	t.Run("Create", testCreateContact(db, r, users))
	t.Run("Get", testGetContacts(db, r, users))
}

func setupContactUsers(t *testing.T, db *sql.DB, router http.Handler) []User {

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

		got := User{}
		json.NewDecoder(w.Body).Decode(&got)

		resultUsers = append(resultUsers, got)
	}

	return resultUsers

}

func testCreateContact(db *sql.DB, router http.Handler, users []User) func(t *testing.T) {
	return func(t *testing.T) {

		// Setup
		mockUser := &User{
			PhoneNumber: "+65 9999 1001",
			FirstName:   "ContactOwner",
			LastName:    "User",
		}
		bs, _ := json.Marshal(mockUser)

		ws := httptest.NewRecorder()
		rs := httptest.NewRequest("POST", "/user", bytes.NewBuffer(bs))
		router.ServeHTTP(ws, rs)

		createdUser := new(User)
		json.NewDecoder(ws.Body).Decode(createdUser)

		// Test
		b := []byte(`{"phone_number": "` + users[0].PhoneNumber + `"}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/user/contact", bytes.NewBuffer(b))
		claim, _ := json.Marshal(&RawClient{UserId: createdUser.ID, ClientId: "test"})
		r.Header.Add("X-User-Claim", string(claim))

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		// Assert
		got, want := User{}, users[0]
		json.NewDecoder(w.Body).Decode(&got)
		if diff := cmp.Diff(got, want); len(diff) != 0 {
			t.Error(diff)
		}

		assertDB(t, db, `SELECT * FROM contact WHERE "user" = $1 AND contact = $2`, createdUser.ID, users[0].ID)

	}
}

func testGetContacts(db *sql.DB, router http.Handler, users []User) func(t *testing.T) {
	return func(t *testing.T) {

		// Setup
		mockUser := &User{
			PhoneNumber: "+65 9999 1002",
			FirstName:   "ContactOwner",
			LastName:    "User",
		}
		bs, _ := json.Marshal(mockUser)

		ws := httptest.NewRecorder()
		rs := httptest.NewRequest("POST", "/user", bytes.NewBuffer(bs))
		router.ServeHTTP(ws, rs)

		createdUser := new(User)
		json.NewDecoder(ws.Body).Decode(createdUser)

		b := []byte(`{"phone_number": "` + users[0].PhoneNumber + `"}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/user/contact", bytes.NewBuffer(b))
		claim, _ := json.Marshal(&RawClient{UserId: createdUser.ID, ClientId: "test"})
		r.Header.Add("X-User-Claim", string(claim))

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		b = []byte(`{"phone_number": "` + users[1].PhoneNumber + `"}`)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/user/contact", bytes.NewBuffer(b))
		r.Header.Add("X-User-Claim", string(claim))

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		b = []byte(`{"phone_number": "` + users[2].PhoneNumber + `"}`)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/user/contact", bytes.NewBuffer(b))
		r.Header.Add("X-User-Claim", string(claim))

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		// Test
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "/user/contact", nil)
		r.Header.Add("X-User-Claim", string(claim))

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		// Assert
		got, want := []User{}, users
		json.NewDecoder(w.Body).Decode(&got)
		if diff := cmp.Diff(got, want); len(diff) != 0 {
			t.Error(diff)
		}

	}
}
