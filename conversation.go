package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) CreateConversation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse
	userID := r.Context().Value("user").(string)
	conversation := Conversation{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&conversation)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Generate ID
	id := "c-" + RandomHex()
	conversation.ID = id

	// Log
	log.Print(conversation)

	// Insert
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Conversation
	_, err1 := tx.Exec(`
		INSERT INTO "conversation" (id, title, picture) VALUES ($1, $2, $3)
	`, conversation.ID, conversation.Title, conversation.Picture)
	// First member
	_, err2 := tx.Exec(`
		INSERT INTO member ("user", "conversation") VALUES ($1, $2)
	`, userID, conversation.ID)
	if err1 != nil || err2 != nil {
		// likely 404...
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		log.Print(err1, err2)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Publish NATs
	if h.nc != nil {
		conversationString, err := json.Marshal(&conversation)
		if err == nil {
			updateMsg := UpdateMsg{
				Type: "add",
				Data: string(conversationString),
			}
			updateMsgString, err := json.Marshal(&updateMsg)
			if err == nil {
				h.nc.Publish("conversation", updateMsgString)
			}
		}
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse
	userID := r.Context().Value("user").(string)

	// Response object
	conversations := make([]Conversation, 0)

	// Select
	rows, err := h.db.Query(`
		SELECT "conversation".id, "conversation".title, "conversation".picture, member.pinned
		FROM "conversation", member
		WHERE member.conversation = "conversation".id AND member.user = $1
	`, userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer rows.Close()

	// Scan
	for rows.Next() {
		conversation := Conversation{}
		if err := rows.Scan(&conversation.ID, &conversation.Title, &conversation.Picture, &conversation.Pinned); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		conversations = append(conversations, conversation)
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

func (h *Handler) GetConversation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse
	userID := r.Context().Value("user").(string)
	conversationID := p.ByName("conversation")

	// Response object
	conversation := Conversation{}

	// Select
	err := h.db.QueryRow(`
		SELECT "conversation".id, "conversation".title, "conversation".picture, member.pinned
		FROM "conversation", member
		WHERE member.conversation = "conversation".id AND member.user = $1 AND member.conversation = $2
	`, userID, conversationID).Scan(&conversation.ID, &conversation.Title, &conversation.Picture, &conversation.Pinned)

	switch {
	case err == sql.ErrNoRows:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func (h *Handler) UpdateConversation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse
	userID := r.Context().Value("user").(string)
	conversationID := p.ByName("conversation")
	conversation := Conversation{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&conversation)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check
	var conversationID2 string
	err = h.db.QueryRow(`
		SELECT id FROM "conversation"
		INNER JOIN member
		ON member.conversation = "conversation".id AND member.user = $1 AND member.conversation = $2
	`, userID, conversationID).Scan(&conversationID2)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Update
	if conversation.Title.Valid {
		_, err = h.db.Exec(`
			UPDATE "conversation"
			SET title = $2, picture = $3
			WHERE id = $1
		`, conversationID, conversation.Title, conversation.Picture)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
	}

	// Publish NATs
	if h.nc != nil {
		conversationString, err := json.Marshal(&conversation)
		if err == nil {
			updateMsg := UpdateMsg{
				Type: "update",
				Data: string(conversationString),
			}
			updateMsgString, err := json.Marshal(&updateMsg)
			if err == nil {
				h.nc.Publish("conversation", updateMsgString)
			}
		}
	}

	w.WriteHeader(200)
}

func (h *Handler) DeleteConversation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userID := r.Context().Value("user").(string)
	conversationID := p.ByName("conversation")

	// Delete
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Check
	var conversationID2 string
	err = h.db.QueryRow(`
		SELECT id FROM "conversation"
		INNER JOIN member
		ON member.conversation = "conversation".id AND member.user = $1 AND member.conversation = $2
	`, userID, conversationID).Scan(&conversationID2)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Users in Conversation
	_, err1 := tx.Exec(`
		DELETE FROM "member" WHERE "conversation" = $1
	`, conversationID)
	// Conversation
	_, err2 := tx.Exec(`
		DELETE FROM "conversation" WHERE "id" = $1
	`, conversationID)

	if err1 != nil || err2 != nil {
		// likely 404...
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		log.Print(err1, err2)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Publish NATs
	if h.nc != nil {
		conversation := Conversation{
			ID: conversationID,
		}
		conversationString, err := json.Marshal(&conversation)
		if err == nil {
			updateMsg := UpdateMsg{
				Type: "delete",
				Data: string(conversationString),
			}
			updateMsgString, err := json.Marshal(&updateMsg)
			if err == nil {
				h.nc.Publish("conversation", updateMsgString)
			}
		}
	}

	w.WriteHeader(200)
}

func (h *Handler) CreateConversationMember(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse
	// We don't need the user ID here because when we first create a conversation, it should have no members
	// TODO: conversations should have conversation owners?
	//userID := r.Context().Value("user").(string)
	conversationID := p.ByName("conversation")
	member := User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&member)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Validate
	if len(member.ID) < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Log
	log.Print(member)

	// TODO: When we need stronger constraints, add some policy around existing conversations with a title set

	// Insert
	_, err = h.db.Exec(`
		INSERT INTO member ("user", "conversation") VALUES ($2, $1)
	`, conversationID, member.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Publish NATs
	if h.nc != nil {
		member := Member{
			User:         member.ID,
			Conversation: conversationID,
			Pinned:       false, // default
		}
		memberString, err := json.Marshal(&member)
		if err == nil {
			updateMsg := UpdateMsg{
				Type: "add",
				Data: string(memberString),
			}
			updateMsgString, err := json.Marshal(&updateMsg)
			if err == nil {
				h.nc.Publish("member", updateMsgString)
			}
		}
	}

	// Respond
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(member)
	w.Write([]byte(conversationID))
}

func (h *Handler) GetConversationMembers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse
	userID := r.Context().Value("user").(string)
	conversationID := p.ByName("conversation")

	// Response object
	users := make([]User, 0)

	// Select
	rows, err := h.db.Query(`
		SELECT "user".id, "user".username, "user".bio, "user".profile_pic, "user".first_name, "user".last_name, "user".phone_number FROM "user"
		INNER JOIN member m ON "user".id = m.user AND "user".id != $1
		INNER JOIN conversation ON "conversation".id = m.conversation
		INNER JOIN member
		ON member.conversation = "conversation".id AND member.user = $1 AND member.conversation = $2
	`, userID, conversationID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer rows.Close()

	// Scan
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Bio, &user.ProfilePic, &user.FirstName, &user.LastName, &user.PhoneNumber); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		users = append(users, user)
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) PinConversation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conversationID := p.ByName("conversation")
	userID := r.Context().Value("user").(string)

	// Check relation exists
	var exists int
	err := h.db.QueryRow(`SELECT 1 FROM member WHERE "user" = $1 AND "conversation" = $2`, userID, conversationID).Scan(&exists)
	if err == sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Update relation
	_, err = h.db.Exec(`UPDATE "member" SET "pinned" = TRUE WHERE "user" = $1 AND "conversation" = $2`, userID, conversationID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Publish NATs
	if h.nc != nil {
		member := Member{
			User:         userID,
			Conversation: conversationID,
			Pinned:       true,
		}
		memberString, err := json.Marshal(&member)
		if err == nil {
			updateMsg := UpdateMsg{
				Type: "update",
				Data: string(memberString),
			}
			updateMsgString, err := json.Marshal(&updateMsg)
			if err == nil {
				h.nc.Publish("member", updateMsgString)
			}
		}
	}

	w.WriteHeader(200)
}

func (h *Handler) UnpinConversation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conversationID := p.ByName("conversation")
	userID := r.Context().Value("user").(string)

	// Check relation exists
	var exists int
	err := h.db.QueryRow(`SELECT 1 FROM member WHERE "user" = $1 AND "conversation" = $2`, userID, conversationID).Scan(&exists)
	if err == sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Update relation
	_, err = h.db.Exec(`UPDATE "member" SET "pinned" = FALSE WHERE "user" = $1 AND "conversation" = $2`, userID, conversationID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Publish NATs
	if h.nc != nil {
		member := Member{
			User:         userID,
			Conversation: conversationID,
			Pinned:       false,
		}
		memberString, err := json.Marshal(&member)
		if err == nil {
			updateMsg := UpdateMsg{
				Type: "update",
				Data: string(memberString),
			}
			updateMsgString, err := json.Marshal(&updateMsg)
			if err == nil {
				h.nc.Publish("member", updateMsgString)
			}
		}
	}

	w.WriteHeader(200)
}
