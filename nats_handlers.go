package main

import (
	"encoding/json"
	"log"

	"github.com/nats-io/go-nats"
)

func (h *Handler) ContactHandler(msg *nats.Msg) {
	// Validate JSON
	updateMsg := UpdateMsg{}
	err := json.Unmarshal(msg.Data, &updateMsg)
	if err != nil {
		log.Println(err)
		return
	}

	contact := Contact{}
	err = json.Unmarshal([]byte(updateMsg.Data), &contact)
	if err != nil {
		log.Println(err)
		return
	}

	// Transmit
	for _, conn := range h.contactConnections {
		conn <- msg.Data
	}
}

func (h *Handler) ConversationHandler(msg *nats.Msg) {
	// Validate JSON
	updateMsg := UpdateMsg{}
	err := json.Unmarshal(msg.Data, &updateMsg)
	if err != nil {
		log.Println(err)
		return
	}

	conversation := Conversation{}
	err = json.Unmarshal([]byte(updateMsg.Data), &conversation)
	if err != nil {
		log.Println(err)
		return
	}

	// Transmit
	for _, conn := range h.conversationConnections {
		conn <- msg.Data
	}
}

func (h *Handler) UserHandler(msg *nats.Msg) {
	// Validate JSON
	updateMsg := UpdateMsg{}
	err := json.Unmarshal(msg.Data, &updateMsg)
	if err != nil {
		log.Println(err)
		return
	}

	user := User{}
	err = json.Unmarshal([]byte(updateMsg.Data), &user)
	if err != nil {
		log.Println(err)
		return
	}

	// Transmit
	for _, conn := range h.userConnections {
		conn <- msg.Data
	}
}

func (h *Handler) MemberHandler(msg *nats.Msg) {
	// Validate JSON
	updateMsg := UpdateMsg{}
	err := json.Unmarshal(msg.Data, &updateMsg)
	if err != nil {
		log.Println(err)
		return
	}

	member := Member{}
	err = json.Unmarshal([]byte(updateMsg.Data), &member)
	if err != nil {
		log.Println(err)
		return
	}

	// Get transmit channel
	if channels, ok := h.memberConnections[member.Conversation]; ok {
		// Transmit
		for _, conn := range channels {
			conn <- msg.Data
		}
	} else {
		log.Printf("member conversation %s not found\n", member.Conversation)
	}
}
