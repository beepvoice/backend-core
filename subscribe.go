package main

import (
  "fmt"
  "net/http"
  "time"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) SubscribeContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  Subscribe(h.contactConnections, w, r, p)
}

func (h *Handler) SubscribeConversation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  Subscribe(h.conversationConnections, w, r, p)
}

func (h *Handler) SubscribeUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  Subscribe(h.userConnections, w, r, p)
}

func (h *Handler) SubscribeMember(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  conversation := p.ByName("conversation")
  if _, ok := h.memberConnections[conversation]; !ok {
    h.memberConnections[conversation] = make(map[string]chan []byte)
  }

  Subscribe(h.memberConnections[conversation], w, r, p)
}

func Subscribe(channels map[string]chan []byte, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  flusher, ok := w.(http.Flusher)
  if !ok {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "text/event-stream")
  w.Header().Set("Cache-Control", "no-cache")
  w.Header().Set("Connection", "keep-alive")

  id := RandomHex()
  recv := make(chan []byte)
  channels[id] = recv

  // Refresh connection periodically
  resClosed := w.(http.CloseNotifier).CloseNotify()
  ticker := time.NewTicker(25 * time.Second)

  for {
    select {
      case msg := <-recv:
        fmt.Fprintf(w, "data: %s\n\n", msg)
        flusher.Flush()
      case <- ticker.C:
        w.Write([]byte(":\n\n"))
      case <- resClosed:
        ticker.Stop()
        delete(channels, id)
        return
    }
  }
}
