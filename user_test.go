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

func TestUser(t *testing.T) {
	db := connect()
	defer db.Close()
	h := NewHandler(db, nil)
	r := NewRouter(h)

	t.Run("Create", testCreateUser(db, r))
	t.Run("GetUserByPhone", testGetUserByPhone(db, r))
	t.Run("GetUser", testGetUser(db, r))
	t.Run("UpdateUser", testUpdateUser(db, r))
}

func testCreateUser(db *sql.DB, router http.Handler) func(t *testing.T) {
	return func(t *testing.T) {

		// Setup
		mockUser := User{
			PhoneNumber: "+65 99999999",
			FirstName:   "Test",
			LastName:    "User 1",
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

func testGetUserByPhone(db *sql.DB, router http.Handler) func(t *testing.T) {
	return func(t *testing.T) {

		// Setup
		mockUser := User{
			PhoneNumber: "+65 99999998",
			FirstName:   "Test",
			LastName:    "User 2",
		}
		b, _ := json.Marshal(mockUser)

		ws := httptest.NewRecorder()
		rs := httptest.NewRequest("POST", "/user", bytes.NewBuffer(b))
		router.ServeHTTP(ws, rs)

		createdUser := new(User)
		json.NewDecoder(ws.Body).Decode(createdUser)

		// Test
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/user?phone_number=%2B6599999998", nil)

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		// Assert
		got, want := new(User), createdUser
		json.NewDecoder(w.Body).Decode(got)
		if diff := cmp.Diff(got, want); len(diff) != 0 {
			t.Error(diff)
		}

	}
}

func testGetUser(db *sql.DB, router http.Handler) func(t *testing.T) {
	return func(t *testing.T) {

		// Setup
		mockUser := User{
			PhoneNumber: "+65 99999997",
			FirstName:   "User",
			LastName:    "Test 2",
		}
		b, _ := json.Marshal(mockUser)

		ws := httptest.NewRecorder()
		rs := httptest.NewRequest("POST", "/user", bytes.NewBuffer(b))
		router.ServeHTTP(ws, rs)

		createdUser := new(User)
		json.NewDecoder(ws.Body).Decode(createdUser)

		// Test
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/user/id/"+createdUser.ID, nil)

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		// Assert
		got, want := new(User), createdUser
		json.NewDecoder(w.Body).Decode(got)
		if diff := cmp.Diff(got, want); len(diff) != 0 {
			t.Error(diff)
		}

	}
}

func testUpdateUser(db *sql.DB, router http.Handler) func(t *testing.T) {
	return func(t *testing.T) {

		// Setup
		mockUser := User{
			PhoneNumber: "+65 99999996",
			FirstName:   "User",
			LastName:    "Test 3",
		}
		bs, _ := json.Marshal(mockUser)

		ws := httptest.NewRecorder()
		rs := httptest.NewRequest("POST", "/user", bytes.NewBuffer(bs))
		router.ServeHTTP(ws, rs)

		createdUser := new(User)
		json.NewDecoder(ws.Body).Decode(createdUser)

		// Test
		b := []byte(`{"first_name": "Ambrose", "last_name": "Chua"}`)
		updateUser := new(User)
		json.NewDecoder(bytes.NewBuffer(b)).Decode(updateUser)
		updatedUser := createdUser
		updatedUser.FirstName = updateUser.FirstName
		updatedUser.LastName = updateUser.LastName

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PATCH", "/user", bytes.NewBuffer(b))
		claim, _ := json.Marshal(&RawClient{UserId: createdUser.ID, ClientId: "test"})
		r.Header.Add("X-User-Claim", string(claim))

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		// Assert
		wt := httptest.NewRecorder()
		rt := httptest.NewRequest("GET", "/user/id/"+createdUser.ID, nil)

		router.ServeHTTP(wt, rt)

		got, want := new(User), updatedUser
		json.NewDecoder(wt.Body).Decode(got)
		if diff := cmp.Diff(got, want); len(diff) != 0 {
			t.Error(diff)
		}

	}
}
