package main

import (
	"github.com/julienschmidt/httprouter"
)

func NewRouter(h *Handler) *httprouter.Router {
	router := httprouter.New()

	// Users
	router.POST("/user", h.CreateUser)
	router.GET("/user", h.GetUserByPhone)
	router.GET("/user/id/:user", h.GetUser)
	router.GET("/user/username/:username", h.GetUserByUsername)
	router.PATCH("/user", AuthMiddleware(h.UpdateUser))

	// Conversations
	router.POST("/user/conversation", AuthMiddleware(h.CreateConversation))
	router.GET("/user/conversation", AuthMiddleware(h.GetConversations)) // USER MEMBER CONVERSATION
	router.DELETE("/user/conversation/:conversation", AuthMiddleware(h.DeleteConversation))
	//router.GET("/user/:user/conversation/bymembers/", h.GetConversationsByMembers) // TODO
	router.GET("/user/conversation/:conversation", AuthMiddleware(h.GetConversation))      // USER MEMBER CONVERSATION
	router.PATCH("/user/conversation/:conversation", AuthMiddleware(h.UpdateConversation)) // USER MEMBER CONVERSATION ADMIN=true -> update conversation title
	//router.DELETE("/user/:user/conversation/:conversation", h.DeleteConversation) // USER MEMBER CONVERSATION -> delete membership
	router.POST("/user/conversation/:conversation/pin", AuthMiddleware(h.PinConversation))
	router.DELETE("/user/conversation/:conversation/pin", AuthMiddleware(h.UnpinConversation))
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

	// Subscribe
	router.GET("/user/subscribe/contact", AuthMiddleware(h.SubscribeContact))
	router.GET("/user/subscribe/conversation", AuthMiddleware(h.SubscribeConversation))
	router.GET("/user/subscribe", AuthMiddleware(h.SubscribeUser))
	router.GET("/user/subscribe/conversation/:conversation/member", AuthMiddleware(h.SubscribeMember))

	return router
}
