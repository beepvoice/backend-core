package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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

func (h *Handler) GetUserByPhone(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse
	phone, err := ParsePhone(r.FormValue("phone_number"))

	// Validate
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Response object
	user := User{}

	// Select
	err = h.db.QueryRow(`
		SELECT id, username, bio, profile_pic, first_name, last_name, phone_number FROM "user" WHERE phone_number = $1
	`, phone).Scan(&user.ID, &user.Username, &user.Bio, &user.ProfilePic, &user.FirstName, &user.LastName, &user.PhoneNumber)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
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
	username := p.ByName("username")

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

	w.WriteHeader(200)
}
