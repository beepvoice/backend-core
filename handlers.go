package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	db *sql.DB
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse
	user := User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Validate
	phone, err := ParsePhone(user.PhoneNumber)
	if err != nil || len(user.FirstName) < 1 || len(user.LastName) < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user.PhoneNumber = phone // shouldn't be needed but makes life easier

	// Generate ID
	id := "u-" + RandomHex()
	user.ID = id

	// Log
	log.Print(user)

	// Insert
	var finalId string
	err = h.db.QueryRow(`
		INSERT INTO "user" (id, username, bio, profile_pic, first_name, last_name, phone_number)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT(phone_number)
			DO UPDATE SET phone_number=EXCLUDED.phone_number, username=$2, first_name=$5, last_name=$6
			RETURNING id
	`, user.ID, user.Username, user.Bio, user.ProfilePic, user.FirstName, user.LastName, user.PhoneNumber).Scan(&finalId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	user.ID = finalId

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUsersByPhone(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse
	phone, err := ParsePhone(r.FormValue("phone_number"))

	// Validate
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Response object
	users := make([]User, 0)

	// Select
	rows, err := h.db.Query(`
		SELECT id, username, bio, profile_pic, first_name, last_name FROM "user" WHERE phone_number = $1
	`, phone)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer rows.Close()

	// Scan
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Bio, &user.ProfilePic, &user.FirstName, &user.LastName); err != nil {
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

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse
	userID := p.ByName("user")

	// Response object
	user := User{}

	// Select
	err := h.db.QueryRow(`
		SELECT id, username, bio, profile_pic, first_name, last_name, phone_number FROM "user" WHERE id = $1
	`, userID).Scan(&user.ID, &user.Username, &user.Bio, &user.ProfilePic, &user.FirstName, &user.LastName, &user.PhoneNumber)

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
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUserByUsername(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse
	username := p.ByName("userername")

	// Response object
	user := User{}

	// Select
	err := h.db.QueryRow(`
		SELECT id, username, bio, profile_pic, first_name, last_name, phone_number FROM "user" WHERE username = $1
	`, username).Scan(&user.ID, &user.Username, &user.Bio, &user.ProfilePic, &user.FirstName, &user.LastName, &user.PhoneNumber)

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
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  // Parse
  userID := r.Context().Value("user").(string)
  user := User{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

  // Check for duplicate username
  var _id string
  err = h.db.QueryRow(`
    SELECT id FROM "user" WHERE "user".id <> $1 AND "user".username = $2
  `, userID, user.Username).Scan(&_id)
  if err == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
  } else if err != sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
  }

  // Update
  _, err = h.db.Exec(`
		UPDATE "user"
		SET 
      username = $2,
      bio = $3,
      profile_pic = $4,
      first_name = $5,
      last_name = $6
		WHERE id = $1
	`, userID, user.Username, user.Bio, user.ProfilePic, user.FirstName, user.LastName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.WriteHeader(200);
}

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
		INSERT INTO "conversation" (id, title, dm, picture) VALUES ($1, $2, $3, $4)
	`, conversation.ID, conversation.Title, conversation.DM, conversation.Picture)
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
		SELECT id, CASE
      WHEN dm THEN (SELECT CONCAT("user".first_name, ' ', "user".last_name) FROM "user", member WHERE "user".id <> $1 AND "user".id = member.user AND member.conversation = "conversation".id)
      ELSE title
    END AS title,
    picture
    FROM "conversation"
		INNER JOIN member
		ON member.conversation = "conversation".id AND member.user = $1
	`, userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer rows.Close()

	// Scan
	for rows.Next() {
		var id, title, picture string
		if err := rows.Scan(&id, &title, &picture); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
    conversations = append(conversations, Conversation{ID:id, Title:title, DM:false, Picture:picture})
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
		SELECT id, CASE 
      WHEN dm THEN (SELECT CONCAT("user".first_name, ' ', "user".last_name) FROM "user", member WHERE "user".id <> $1 AND "user".id = member.user AND member.conversation = "conversation".id)
      ELSE title
    END AS title,
    picture
    FROM "conversation"
		INNER JOIN member
		ON member.conversation = "conversation".id AND member.user = $1 AND member.conversation = $2
	`, userID, conversationID).Scan(&conversation.ID, &conversation.Title, &conversation.Picture)

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
	if len(conversation.Title) > 0 {
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

	w.WriteHeader(200);
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

	w.WriteHeader(200)
}

func (h *Handler) CreateConversationMember(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse
	userID := r.Context().Value("user").(string)
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

  // Check for valid conversation and prevent duplicate entries
  var test string
  err = h.db.QueryRow(`
    SELECT "conversation".id FROM "conversation", "member"
    WHERE 
      "conversation".id = $1
      AND (
        "conversation".dm = FALSE
        OR (SELECT 
          COUNT("member".user)
          FROM "member"
          WHERE "member".conversation = $1)
        <= 2)
      AND "member".conversation = "conversation".id
      AND "member".user <> $2
  `, conversationID, member.ID).Scan(&test)
	switch {
    case err == sql.ErrNoRows:
      http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
      return
    case err != nil:
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      log.Print(err)
      return
	}

  // Check user adding the user is in conversation
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

	// Insert
	_, err = h.db.Exec(`
		INSERT INTO member ("user", "conversation") VALUES ($2, $1)
	`, conversationID, member.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Respond
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(member)
	w.WriteHeader(200)
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
		var id, username, bio, profilePic, firstName, lastName, phoneNumber string
		if err := rows.Scan(&id, &username, &bio, &profilePic, &firstName, &lastName, &phoneNumber); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
    users = append(users, User{ID: id, Username: username, Bio: bio, ProfilePic: profilePic, FirstName: firstName, LastName: lastName, PhoneNumber: phoneNumber})
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

type PhoneNumber struct {
	PhoneNumber string `json:"phone_number"`
}
func (h *Handler) CreateContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse
	userID := r.Context().Value("user").(string)
	contact := PhoneNumber{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&contact)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Validate
	phone, err := ParsePhone(contact.PhoneNumber)
	if err != nil || len(contact.PhoneNumber) < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Generate ID (just in case)
	id := "u-" + RandomHex()

	// Create contact if not exists, returning the id regardless
	var contactId string
	err = h.db.QueryRow(`
		INSERT INTO "user" (id, username, bio, profile_pic, first_name, last_name, phone_number)
			VALUES ($1, '', '', '', '', '', $2)
			ON CONFLICT(phone_number)
			DO UPDATE SET phone_number=EXCLUDED.phone_number
			RETURNING id
	`, id, phone).Scan(&contactId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Insert
	_, err = h.db.Exec(`
		INSERT INTO contact ("user", contact) VALUES ($1, $2)
	`, userID, contactId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Respond
	w.WriteHeader(200)
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(contact)
}

func (h *Handler) GetContacts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse
	userID := r.Context().Value("user").(string)

	// Response object
	contacts := make([]User, 0)

	// Select
	rows, err := h.db.Query(`
		SELECT id, username, bio, profile_pic, first_name, last_name, phone_number FROM "user"
		INNER JOIN contact
		ON contact.contact = "user".id AND contact.user = $1
	`, userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer rows.Close()

	// Scan
	for rows.Next() {
		var id, username, bio, profilePic, firstName, lastName, phone string
		if err := rows.Scan(&id, &username, &bio, &profilePic, &firstName, &lastName, &phone); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		contacts = append(contacts, User{id, username, bio, profilePic, firstName, lastName, phone})
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db}
}
