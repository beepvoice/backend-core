package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
