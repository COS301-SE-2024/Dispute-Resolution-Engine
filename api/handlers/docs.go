package handlers

import "net/http"

func (h handler) getDocs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}