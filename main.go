package main

import (
	"database/sql"
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
	router.POST("/user/", h.CreateUser)
	router.GET("/user/", h.GetUsersByPhone)
	router.GET("/user/:user", h.GetUser)
	//router.PATCH("/user/:user", h.UpdateUser)
	// Conversations
	router.POST("/user/:user/conversation/", h.CreateConversation)
	router.GET("/user/:user/conversation/", h.GetConversations) // USER MEMBER CONVERSATION
	router.DELETE("/user/:user/conversation/:conversation", h.DeleteConversation)
	//router.GET("/user/:user/conversation/bymembers/", h.GetConversationsByMembers) // TODO
	router.GET("/user/:user/conversation/:conversation", h.GetConversation)      // USER MEMBER CONVERSATION
	router.PATCH("/user/:user/conversation/:conversation", h.UpdateConversation) // USER MEMBER CONVERSATION ADMIN=true -> update conversation title
	//router.DELETE("/user/:user/conversation/:conversation", h.DeleteConversation) // USER MEMBER CONVERSATION -> delete membership
	router.POST("/user/:user/conversation/:conversation/member/", h.CreateConversationMember) // USER MEMBER CONVERSATION ADMIN=true -> create new membership
	router.GET("/user/:user/conversation/:conversation/member/", h.GetConversationMembers)    // USER MEMBER CONVERSATION
	//router.DELETE("/user/:user/conversation/:conversation/member/:member", h.DeleteConversationMember) // USER MEMBER CONVERSATION ADMIN=true -> delete membership
	// Last heard
	//router.GET("/user/:user/lastheard/:conversation", h.GetLastheard)
	//router.PUT("/user/:user/lastheard/:conversation", h.SetLastheard)
	// Contacts
	router.POST("/user/:user/contact/", h.CreateContact)
	router.GET("/user/:user/contact/", h.GetContacts)
	//router.GET("/user/:user/contact/:contact", h.GetContact)
	//router.DELETE("/user/:user/contact/:contact", h.DeleteContact)
	//router.GET("/user/:user/contact/:contact/conversation/", h.GetContactConversations)

	log.Printf("starting server on %s", listen)
	log.Fatal(http.ListenAndServe(listen, router))
}
