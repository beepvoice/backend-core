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
	"gopkg.in/guregu/null.v3"
)

func TestConversation(t *testing.T) {
	db := connect()
	defer db.Close()
	h := NewHandler(db, nil)
	r := NewRouter(h)

	users := setupConversationUsers(t, db, r)

	t.Run("Create", testCreateConversation(db, r, users))
	t.Run("Get", testGetConversations(db, r, users))
}

func setupConversationUsers(t *testing.T, db *sql.DB, router http.Handler) []User {

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

func testCreateConversation(db *sql.DB, router http.Handler, users []User) func(t *testing.T) {
	return func(t *testing.T) {

		// Test
		mockConversation := &Conversation{
			Title: null.StringFrom("Test Conversation 1"),
			DM:    false,
		}
		b, _ := json.Marshal(mockConversation)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/user/conversation", bytes.NewBuffer(b))
		claim, _ := json.Marshal(&RawClient{UserId: users[0].ID, ClientId: "test"})
		r.Header.Add("X-User-Claim", string(claim))

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)

		// Assert
		got, want := &Conversation{}, mockConversation
		json.NewDecoder(w.Body).Decode(&got)
		if got.DM != want.DM || got.Title.String != want.Title.String {
			t.Error("Wanted a Conversation with same Title, DM. Got something else")
		}

		assertDB(t, db, `SELECT * FROM "conversation" WHERE title = $1 AND dm = $2`, mockConversation.Title, mockConversation.DM)
		assertDB(t, db, `SELECT * FROM member WHERE "user" = $1 AND "conversation" = $2`, users[0].ID, got.ID)

	}
}

func testGetConversations(db *sql.DB, router http.Handler, users []User) func(t *testing.T) {
	return func(t *testing.T) {

		// Setup
		mockConversation := &Conversation{
			Title: null.StringFrom("Test Conversation 2"),
			DM:    false,
		}
		bs, _ := json.Marshal(mockConversation)

		ws := httptest.NewRecorder()
		rs := httptest.NewRequest("POST", "/user/conversation", bytes.NewBuffer(bs))
		claims, _ := json.Marshal(&RawClient{UserId: users[1].ID, ClientId: "test"})
		rs.Header.Add("X-User-Claim", string(claims))

		router.ServeHTTP(ws, rs)
		assertCode(t, ws, 200)

		// Test
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/user/conversation", nil)
		claim, _ := json.Marshal(&RawClient{UserId: users[1].ID, ClientId: "test"})
		r.Header.Add("X-User-Claim", string(claim))

		router.ServeHTTP(w, r)
		assertCode(t, w, 200)
		conversations := make([]Conversation, 1)
		json.NewDecoder(w.Body).Decode(&conversations)

		// Assert
		got, want := conversations[0], Conversation{}
		json.NewDecoder(ws.Body).Decode(&want)
		if diff := cmp.Diff(got, want); len(diff) != 0 {
			t.Error(diff)
		}

	}
}
