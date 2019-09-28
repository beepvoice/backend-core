package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
		contacts = append(contacts, User{id, &username, bio, profilePic, firstName, lastName, phone})
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}
