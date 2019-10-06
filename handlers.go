package main

import (
	"database/sql"

  "github.com/nats-io/go-nats"
)

type Handler struct {
	db *sql.DB
  nc *nats.Conn

  contactConnections        map[string]chan []byte
  conversationConnections   map[string]chan []byte
  userConnections           map[string]chan []byte
  memberConnections         map[string]map[string]chan []byte
}

func NewHandler(db *sql.DB, nc *nats.Conn) *Handler {
  contactConnections := make(map[string]chan []byte)
  conversationConnections := make(map[string]chan []byte)
  userConnections := make(map[string]chan []byte)
  memberConnections := make(map[string]map[string]chan []byte)

  h := &Handler{
    db,
    nc,
    contactConnections,
    conversationConnections,
    userConnections,
    memberConnections,
  }

  if nc != nil {
    nc.Subscribe("contacts", h.ContactHandler)
    nc.Subscribe("conversations", h.ConversationHandler)
    nc.Subscribe("users", h.UserHandler)
    nc.Subscribe("members", h.MemberHandler)
  }

  return h
}
