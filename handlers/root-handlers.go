package handlers

import (
	"log"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("TODOER API"))
	if err != nil {
		log.Fatal(err)
	}
}
