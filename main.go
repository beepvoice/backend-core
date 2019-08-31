package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

var listen string
var postgres string

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	listen = os.Getenv("LISTEN")
	postgres = os.Getenv("POSTGRES")

	// Open postgres
	log.Printf("connecting to postgres %s", postgres)
	db, err := sql.Open("postgres", postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Handler
	h := NewHandler(db)

	// Routes
	router := httprouter.New()
	// Users
	router.POST("/user", h.CreateUser)
	router.GET("/user", h.GetUserByPhone)
	router.GET("/user/id/:user", h.GetUser)
	router.GET("/user/username/:username", h.GetUserByUsername)
	router.PATCH("/user", h.UpdateUser)
	// Conversations
	router.POST("/user/conversation", AuthMiddleware(h.CreateConversation))
	router.GET("/user/conversation", AuthMiddleware(h.GetConversations)) // USER MEMBER CONVERSATION
	router.DELETE("/user/conversation/:conversation", AuthMiddleware(h.DeleteConversation))
	//router.GET("/user/:user/conversation/bymembers/", h.GetConversationsByMembers) // TODO
	router.GET("/user/conversation/:conversation", AuthMiddleware(h.GetConversation))      // USER MEMBER CONVERSATION
	router.PATCH("/user/conversation/:conversation", AuthMiddleware(h.UpdateConversation)) // USER MEMBER CONVERSATION ADMIN=true -> update conversation title
	//router.DELETE("/user/:user/conversation/:conversation", h.DeleteConversation) // USER MEMBER CONVERSATION -> delete membership
	router.POST("/user/conversation/:conversation/pin", AuthMiddleware(h.PinConversation))
	router.POST("/user/conversation/:conversation/member", AuthMiddleware(h.CreateConversationMember)) // USER MEMBER CONVERSATION ADMIN=true -> create new membership
	router.GET("/user/conversation/:conversation/member", AuthMiddleware(h.GetConversationMembers))    // USER MEMBER CONVERSATION
	//router.DELETE("/user/:user/conversation/:conversation/member/:member", h.DeleteConversationMember) // USER MEMBER CONVERSATION ADMIN=true -> delete membership
	// Last heard
	//router.GET("/user/:user/lastheard/:conversation", h.GetLastheard)
	//router.PUT("/user/:user/lastheard/:conversation", h.SetLastheard)
	// Contacts
	router.POST("/user/contact", AuthMiddleware(h.CreateContact))
	router.GET("/user/contact", AuthMiddleware(h.GetContacts))
	//router.GET("/user/:user/contact/:contact", h.GetContact)
	//router.DELETE("/user/:user/contact/:contact", h.DeleteContact)
	//router.GET("/user/:user/contact/:contact/conversation/", h.GetContactConversations)

	log.Printf("starting server on %s", listen)
	log.Fatal(http.ListenAndServe(listen, router))
}

type RawClient struct {
	UserId   string `json:"userid"`
	ClientId string `json:"clientid"`
}

func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ua := r.Header.Get("X-User-Claim")
		if ua == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var client RawClient
		err := json.Unmarshal([]byte(ua), &client)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		context := context.WithValue(r.Context(), "user", client.UserId)
		next(w, r.WithContext(context), p)
	}
}
